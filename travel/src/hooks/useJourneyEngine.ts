"use client";
import { useState, useCallback, useRef, useEffect, useMemo } from "react";
import {
  createJourneyEngine, JourneyEngine,
  RideState, RideEvent,
  stateAnimations, getStateConfig,
} from "@cytaxi/ride-machine";
import { getGlobalBus } from "@cytaxi/events";
import type { Place, FareBreakdown, TrackingUpdate, DriverInfo, RoutePayload } from "@/types";
import { searchPlaces, calculateRoute, estimateFare, requestTrip } from "@/services/api";
import { subscribeToTrip } from "@/services/tracking";
import {
  saveSession, loadSession, clearSession, isSessionValid,
} from "@/services/state-recovery";
import {
  enqueueAction, getPendingActions, dequeueAction,
  isOnline, onOnline, onOffline,
} from "@/services/offline-queue";
import { trackJourneyEvent, trackStateDuration, trackError, trackLatency } from "@/services/telemetry";
import { DEMO_CONFIG } from "@/services/demo";

const PRE_TRIP_STATES: RideState[] = ["travel_home", "pickup_select", "input", "confirm"];

export function useJourneyEngine() {
  const engineRef = useRef<JourneyEngine>(createJourneyEngine("travel_home"));
  const [state, setState] = useState<RideState>("travel_home");
  const [prevState, setPrevState] = useState<RideState | null>(null);
  const [direction, setDirection] = useState<"forward" | "back">("forward");
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [online, setOnline] = useState(true);

  const [pickupAddress, setPickupAddress] = useState("");
  const [pickupCoords, setPickupCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [dest, setDest] = useState<Place | null>(null);
  const [destQuery, setDestQuery] = useState("");
  const [destSuggestions, setDestSuggestions] = useState<Place[]>([]);
  const [route, setRoute] = useState<RoutePayload | null>(null);
  const [fare, setFare] = useState<FareBreakdown | null>(null);
  const [driver, setDriver] = useState<DriverInfo | null>(null);
  const [tripId, setTripId] = useState("");
  const [tracking, setTracking] = useState<TrackingUpdate | null>(null);
  const [eta, setEta] = useState(0);
  const [loading, setLoading] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState<"cash" | "card">("cash");
  const [vehicleType, setVehicleType] = useState("standard");
  const [note, setNote] = useState("");
  const [coupon, setCoupon] = useState("");
  const [scheduledAt, setScheduledAt] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [noDrivers, setNoDrivers] = useState(false);

  const sheetRef = useRef<HTMLDivElement>(null);
  const searchTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const animTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const stateTimer = useRef<number>(Date.now());
  const searchTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const searchSubRef = useRef<(() => void) | null>(null);
  const bus = useMemo(() => getGlobalBus(), []);

  useEffect(() => {
    setOnline(isOnline());
    const off1 = onOnline(() => setOnline(true));
    const off2 = onOffline(() => setOnline(false));
    return () => { off1(); off2(); };
  }, []);

  useEffect(() => {
    const engine = engineRef.current;
    const unsub = engine.onTransition((newState, prev, event) => {
      const now = Date.now();
      trackStateDuration(prev || "", newState, now - stateTimer.current);
      stateTimer.current = now;
      trackJourneyEvent("STATE_TRANSITION", { state: newState, prevState: prev });
      setPrevState(prev);
      setDirection(engine.snapshot.direction);
      setIsTransitioning(true);
      setState(newState);
      setError(null);

      if (animTimer.current) clearTimeout(animTimer.current);
      animTimer.current = setTimeout(() => {
        engine.endTransition();
        setIsTransitioning(false);
      }, stateAnimations[newState]?.duration || 350);
    });
    return () => { unsub(); if (animTimer.current) clearTimeout(animTimer.current); };
  }, []);

  const send = useCallback((event: RideEvent) => {
    engineRef.current.send(event);
  }, []);

  const goTo = useCallback((next: RideState) => {
    engineRef.current.goTo(next);
  }, []);

  const handleCenterChange = useCallback((data: { lat: number; lng: number; address: string }) => {
    setPickupAddress(data.address);
    setPickupCoords({ lat: data.lat, lng: data.lng });
    bus.emit("LOCATION_DETECTED", { ...data, source: "map_drag" });
  }, [bus]);

  const handleMapDestChange = useCallback((data: { lat: number; lng: number; address: string }) => {
    const place: Place = {
      name: data.address.split(",")[0] || "Destino seleccionado",
      address: data.address,
      lat: data.lat,
      lng: data.lng,
    };
    setDest(place);
    setDestQuery(data.address);
    setDestSuggestions([]);
    setRoute(null);
    setFare(null);
  }, []);

  const doSearchDest = useCallback(async (q: string) => {
    if (q.length < 3) { setDestSuggestions([]); return; }
    if (searchTimer.current) clearTimeout(searchTimer.current);
    searchTimer.current = setTimeout(async () => {
      const t0 = performance.now();
      const results = await searchPlaces(q);
      trackLatency("search_places", performance.now() - t0);
      setDestSuggestions(results);
    }, 400);
  }, []);

  const selectDest = useCallback((place: Place) => {
    setDest(place);
    setDestQuery(place.name);
    setDestSuggestions([]);
    setRoute(null);
    setFare(null);
  }, []);

  const handleClearDest = useCallback(() => {
    setDest(null);
    setDestQuery("");
    setDestSuggestions([]);
    setRoute(null);
    setFare(null);
  }, []);

  const handleConfirmPickup = useCallback(() => {
    engineRef.current.send("SELECT_PICKUP");
  }, []);

  const handleBackToPickup = useCallback(() => {
    setRoute(null);
    setFare(null);
    engineRef.current.goTo("pickup_select");
  }, []);

  const handleConfirmDest = useCallback(async () => {
    if (!pickupCoords || !dest) return;
    setLoading(true);
    setError(null);
    const t0 = performance.now();

    const routeData = await calculateRoute(pickupCoords, { lat: dest.lat, lng: dest.lng });
    trackLatency("calculate_route", performance.now() - t0);

    const routePayload: RoutePayload = {
      distance_km: routeData?.distance_km || 0,
      distance_meters: routeData?.distance_meters || 0,
      duration_seconds: routeData?.duration_seconds || 0,
      eta_minutes: routeData?.eta_minutes || 0,
      polyline: routeData?.polyline || "",
      pickup: { lat: pickupCoords.lat, lng: pickupCoords.lng },
      dest: { lat: dest.lat, lng: dest.lng },
    };
    setRoute(routePayload);
    bus.emit("ROUTE_CALCULATED", routePayload);

    const t1 = performance.now();
    const fareData = await estimateFare(routePayload.distance_km, routePayload.duration_seconds);
    trackLatency("estimate_fare", performance.now() - t1);

    const farePayload: FareBreakdown = {
      base: fareData.base,
      distance: fareData.distance,
      time: fareData.time,
      subtotal: fareData.subtotal,
      total: fareData.total,
      currency: fareData.currency || "USD",
      distance_km: routePayload.distance_km,
      eta_minutes: routePayload.eta_minutes,
      pricing_model: "standard",
    };
    setFare(farePayload);
    bus.emit("FARE_ESTIMATED", farePayload);
    engineRef.current.send("SELECT_DEST");
    setLoading(false);
  }, [pickupCoords, dest, bus]);

  const handleRequestTrip = useCallback(async () => {
    if (!pickupCoords || !dest || !fare) return;
    setLoading(true);
    setError(null);
    setNoDrivers(false);

    if (!isOnline()) {
      enqueueAction("TRIP_REQUEST", { pickupCoords, dest, fare });
      (bus as any).emit("OFFLINE_QUEUED", { message: "Viaje encolado — se procesará cuando tengas conexión" });
      setError("Sin conexión. Tu solicitud se procesará automáticamente cuando recuperes señal.");
      setLoading(false);
      return;
    }

    bus.emit("SEARCHING_DRIVERS", {});
    engineRef.current.send("REQUEST");
    const t0 = performance.now();
    try {
      const tripResult = await requestTrip({
        phone: DEMO_CONFIG.passenger.phone, passenger_name: DEMO_CONFIG.passenger.name,
        origin_address: pickupAddress, origin_lat: pickupCoords.lat, origin_lng: pickupCoords.lng,
        dest_address: dest.address || dest.name, dest_lat: dest.lat, dest_lng: dest.lng,
      });
      trackLatency("request_trip", performance.now() - t0);
      const tid = tripResult.trip_id || "trip_demo";
      setTripId(tid);

      // Suscripción SSE temprana para recibir asignación de conductor
      searchSubRef.current?.();
      searchSubRef.current = subscribeToTrip(tid, (update) => {
        if (update.type === "driver_assigned" || update.driver) {
          setNoDrivers(false);
          const d = update.driver || DEMO_CONFIG.driver;
          const driverData: DriverInfo = {
            ...d as any,
            lat: update.driver?.lat ?? pickupCoords.lat - 0.01,
            lng: update.driver?.lng ?? pickupCoords.lng - 0.01,
            eta_seconds: update.eta_seconds ?? fare.eta_minutes * 60,
          };
          setDriver(driverData);
          setEta(update.eta_seconds ?? fare.eta_minutes * 60);
          bus.emit("DRIVER_FOUND", driverData);
          bus.emit("DRIVER_ASSIGNED", driverData);
          engineRef.current.send("DRIVER_FOUND");
          setLoading(false);
          if (searchTimeoutRef.current) clearTimeout(searchTimeoutRef.current);
          searchSubRef.current?.();
          searchSubRef.current = null;
        }
      }, () => {});

      if (DEMO_CONFIG.enabled) {
        if (DEMO_CONFIG.simulateNoDrivers) {
          searchTimeoutRef.current = setTimeout(() => {
            setNoDrivers(true);
            setLoading(false);
          }, DEMO_CONFIG.matchingDelay + 2000);
        } else {
          setTimeout(() => {
            const driverData: DriverInfo = {
              ...DEMO_CONFIG.driver,
              lat: pickupCoords.lat - 0.01, lng: pickupCoords.lng - 0.01,
              eta_seconds: fare.eta_minutes * 60,
            };
            setNoDrivers(false);
            setDriver(driverData);
            setEta(fare.eta_minutes * 60);
            bus.emit("DRIVER_FOUND", driverData);
            bus.emit("DRIVER_ASSIGNED", driverData);
            engineRef.current.send("DRIVER_FOUND");
            setLoading(false);
            if (searchTimeoutRef.current) clearTimeout(searchTimeoutRef.current);
            searchSubRef.current?.();
            searchSubRef.current = null;
          }, DEMO_CONFIG.matchingDelay);
        }
      } else {
        // Modo real: esperar SSE, con timeout
        searchTimeoutRef.current = setTimeout(() => {
          setNoDrivers(true);
          setLoading(false);
          searchSubRef.current?.();
          searchSubRef.current = null;
        }, DEMO_CONFIG.searchTimeout);
      }
    } catch (err: any) {
      setLoading(false);
      trackError("TRIP_REQUEST_FAILED", err.message);
      bus.emit("ERROR_OCCURRED", {
        code: "TRIP_REQUEST_FAILED", message: "No se pudo solicitar el viaje",
        context: {}, timestamp: new Date().toISOString(),
      });
      setError("No se pudo solicitar el viaje. Verifica tu conexión e intenta de nuevo.");
      engineRef.current.goTo("input");
    }
  }, [pickupCoords, pickupAddress, dest, fare, bus]);

  const handleAcceptDriver = useCallback(() => {
    if (!driver) return;
    bus.emit("DRIVER_ACCEPTED", driver);
    engineRef.current.send("DRIVER_ACCEPTED");
  }, [driver, bus]);

  const handleRetrySearch = useCallback(() => {
    setNoDrivers(false);
    setError(null);
    setTripId("");
    handleRequestTrip();
  }, [handleRequestTrip]);

  const handleRejectDriver = useCallback(() => {
    if (driver) bus.emit("DRIVER_REJECTED", { driver_id: driver.id });
    setDriver(null);
    setTripId("");
    engineRef.current.goTo("input");
  }, [driver, bus]);

  const handleArriveAtPickup = useCallback(() => {
    if (!driver) return;
    bus.emit("DRIVER_ARRIVING", driver);
    engineRef.current.send("DRIVER_ARRIVING");
  }, [driver, bus]);

  const handleTripStart = useCallback(() => {
    if (!driver || !tripId) return;
    bus.emit("PASSENGER_BOARDED", { driver_id: driver.id, trip_id: tripId });
    bus.emit("TRIP_STARTED", {
      trip_id: tripId, status: "in_progress",
      origin: { lat: 0, lng: 0, address: pickupAddress },
      destination: { lat: dest?.lat || 0, lng: dest?.lng || 0, address: dest?.address || "" },
      created_at: new Date().toISOString(),
    });
    engineRef.current.send("TRIP_START");
  }, [driver, tripId, pickupAddress, dest, bus]);

  const handleDestinationArriving = useCallback(() => {
    bus.emit("DESTINATION_ARRIVED", {
      trip_id: tripId, status: "destination",
      origin: { lat: 0, lng: 0, address: pickupAddress },
      destination: { lat: dest?.lat || 0, lng: dest?.lng || 0, address: dest?.address || "" },
      created_at: new Date().toISOString(),
    });
    engineRef.current.send("TRIP_ARRIVED");
  }, [tripId, pickupAddress, dest, bus]);

  const handleTripComplete = useCallback(() => {
    bus.emit("TRIP_COMPLETED", {
      trip_id: tripId, status: "completed",
      origin: { lat: 0, lng: 0, address: pickupAddress },
      destination: { lat: dest?.lat || 0, lng: dest?.lng || 0, address: dest?.address || "" },
      created_at: new Date().toISOString(),
    });
    engineRef.current.send("TRIP_COMPLETE");
  }, [tripId, pickupAddress, dest, bus]);

  const handlePaymentDone = useCallback(() => {
    if (!fare) return;
    bus.emit("PAYMENT_CONFIRMED", {
      trip_id: tripId, amount: fare.total, currency: fare.currency,
      method: paymentMethod, status: "completed", transaction_id: `txn_${Date.now()}`,
      paid_at: new Date().toISOString(),
    });
    engineRef.current.send("PAYMENT_DONE");
  }, [fare, tripId, paymentMethod, bus]);

  const handleRatingDone = useCallback((score: number) => {
    if (!driver) return;
    bus.emit("RATING_SUBMITTED", {
      trip_id: tripId, from: "passenger", to: driver.id,
      score, categories: { punctuality: score, service: score, comfort: score },
    });
    engineRef.current.send("RATING_DONE");
  }, [driver, tripId, bus]);

  const startTracking = useCallback(() => {
    if (!tripId || !driver) return;
    const unsubscribe = subscribeToTrip(
      tripId,
      (update: TrackingUpdate) => {
        setTracking(update);
        if (update.eta_seconds !== undefined) setEta(update.eta_seconds);
        if (update.driver) {
          setDriver((prev) => prev ? { ...prev, lat: update.driver!.lat, lng: update.driver!.lng } : prev);
        }
        if (update.type === "trip_started") handleTripStart();
        if (update.type === "destination_arriving") handleDestinationArriving();
        if (update.type === "trip_completed") handleTripComplete();
      },
      (err) => {
        trackError("SSE_ERROR", err);
        setError(err === "max_retries"
          ? "No se pudo mantener la conexión con el conductor. Recarga la página."
          : "Error de conexión. Reintentando...");
      },
    );
    return unsubscribe;
  }, [tripId, driver, handleTripStart, handleDestinationArriving, handleTripComplete]);

  const unsubscribeRef = useRef<(() => void) | null>(null);
  useEffect(() => {
    if (tripId && driver && state === "driver_found") {
      unsubscribeRef.current = startTracking() ?? null;
    }
    return () => { unsubscribeRef.current?.(); };
  }, [tripId, driver, state, startTracking]);

  const reset = useCallback(() => {
    trackJourneyEvent("JOURNEY_RESET");
    bus.emit("JOURNEY_RESET", {});
    engineRef.current.reset();
    setDest(null); setDestQuery(""); setRoute(null); setFare(null);
    setDriver(null); setTripId(""); setTracking(null); setEta(0);
    setNote(""); setCoupon(""); setScheduledAt(null);
    setError(null); setNoDrivers(false);
    clearSession();
    if (searchTimeoutRef.current) clearTimeout(searchTimeoutRef.current);
    searchSubRef.current?.();
    unsubscribeRef.current?.();
  }, [bus]);

  const handleCancelTrip = useCallback(() => {
    trackJourneyEvent("TRIP_CANCELLED");
    bus.emit("TRIP_CANCELLED", {
      trip_id: tripId || "unknown", status: "cancelled",
      origin: { lat: 0, lng: 0, address: pickupAddress },
      destination: { lat: dest?.lat || 0, lng: dest?.lng || 0, address: dest?.address || "" },
      created_at: new Date().toISOString(),
      cancelled_by: "passenger",
    });
    setDriver(null); setTripId(""); setTracking(null); setEta(0);
    setDest(null); setDestQuery(""); setRoute(null); setFare(null);
    setNote(""); setCoupon(""); setScheduledAt(null);
    setError(null); setNoDrivers(false);
    engineRef.current.reset();
    clearSession();
    if (searchTimeoutRef.current) clearTimeout(searchTimeoutRef.current);
    searchSubRef.current?.();
    unsubscribeRef.current?.();
  }, [pickupAddress, dest, tripId, bus]);

  const dismissError = useCallback(() => setError(null), []);

  useEffect(() => {
    if (!state || !PRE_TRIP_STATES.includes(state)) return;
  }, [state, dest, destQuery, pickupAddress, pickupCoords, route, fare, vehicleType, note, paymentMethod, scheduledAt]);

  useEffect(() => {
    if (!online) return;
    const pending = getPendingActions();
    for (const action of pending) {
      if (action.type === "TRIP_REQUEST") {
        enqueueAction("TRIP_REQUEST", action.payload);
        dequeueAction(action.id);
        handleRequestTrip();
      }
    }
  }, [online, handleRequestTrip]);

  const animConfig = stateAnimations[state];

  return {
    state, prevState, direction, isTransitioning, animConfig,
    pickupAddress, pickupCoords, dest, destQuery, destSuggestions,
    route, fare, driver, eta, tracking, loading, tripId,
    paymentMethod, vehicleType, note, coupon, scheduledAt,
    online, error, noDrivers,
    sheetRef,

    setPaymentMethod, setVehicleType, setNote, setCoupon, setScheduledAt,
    setDest, setDestQuery, setRoute, setFare,

    send, goTo,
    handleCenterChange, handleMapDestChange,
    doSearchDest, selectDest, handleClearDest,
    handleConfirmPickup, handleBackToPickup,
    handleConfirmDest, handleRequestTrip, handleRetrySearch,
    handleAcceptDriver, handleRejectDriver,
    handleArriveAtPickup, handleTripStart,
    handleDestinationArriving, handleTripComplete,
    handlePaymentDone, handleRatingDone,
    startTracking, reset, handleCancelTrip,
    dismissError,
  };
}
