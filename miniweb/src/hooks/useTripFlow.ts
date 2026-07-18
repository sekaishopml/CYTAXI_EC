"use client";
import React, { useState, useCallback, useRef, useEffect } from "react";
import { TripState, Place, FareBreakdown, TrackingUpdate } from "@/types";
import { searchPlaces, calculateRoute, estimateFare, requestTrip } from "@/services/api";
import { subscribeToTrip } from "@/services/tracking";

import { PickUpStep } from "@/components/states/PickUpStep";
import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

export function useTripFlow() {
  const [state, setState] = useState<TripState>("pickup_select");
  const [pickupAddress, setPickupAddress] = useState("");
  const [pickupCoords, setPickupCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [dest, setDest] = useState<Place | null>(null);
  const [destQuery, setDestQuery] = useState("");
  const [destSuggestions, setDestSuggestions] = useState<Place[]>([]);
  const [route, setRoute] = useState<{ distance_km: number; eta_minutes: number; polyline: string; distance_meters: number; duration_seconds: number } | null>(null);
  const [fare, setFare] = useState<FareBreakdown | null>(null);
  const [driver, setDriver] = useState<any>(null);
  const [tripId, setTripId] = useState("");
  const [tracking, setTracking] = useState<TrackingUpdate | null>(null);
  const [eta, setEta] = useState(0);
  const [loading, setLoading] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState<"cash" | "card">("cash");
  const sheetRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const searchTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleCenterChange = useCallback((data: { lat: number; lng: number; address: string }) => {
    setPickupAddress(data.address);
    setPickupCoords({ lat: data.lat, lng: data.lng });
  }, []);

  // Setea el destino al arrastrar el mapa en estado input
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
  }, []);

  const doSearchDest = useCallback(async (q: string) => {
    if (q.length < 3) { setDestSuggestions([]); return; }
    if (searchTimer.current) clearTimeout(searchTimer.current);
    searchTimer.current = setTimeout(async () => {
      const results = await searchPlaces(q);
      setDestSuggestions(results);
    }, 400);
  }, []);

  const selectDest = useCallback((place: Place) => {
    setDest(place);
    setDestQuery(place.name);
    setDestSuggestions([]);
  }, []);

  const handleClearDest = useCallback(() => {
    setDest(null);
    setDestQuery("");
    setDestSuggestions([]);
    setRoute(null);
    setFare(null);
  }, []);

  const handleConfirmPickup = useCallback(() => {
    setState("input");
  }, []);

  const handleBackPickup = useCallback(() => {
    setDest(null);
    setDestQuery("");
    setDestSuggestions([]);
    setState("pickup_select");
  }, []);

  const handleConfirmDest = useCallback(async () => {
    if (!pickupCoords || !dest) return;
    setLoading(true);
    const routeData = await calculateRoute(pickupCoords, { lat: dest.lat, lng: dest.lng });
    setRoute(routeData);
    const fareData = await estimateFare(routeData?.distance_km || 5.5, routeData?.duration_seconds || 900);
    setFare(fareData);
    setState("confirm");
    setLoading(false);
  }, [pickupCoords, dest]);

  const handleRequestTrip = useCallback(async () => {
    if (!pickupCoords || !dest || !fare) return;
    setLoading(true);
    setState("searching");
    try {
      const tripResult = await requestTrip({
        phone: "0000000000", passenger_name: "Usuario",
        origin_address: pickupAddress, origin_lat: pickupCoords.lat, origin_lng: pickupCoords.lng,
        dest_address: dest.address || dest.name, dest_lat: dest.lat, dest_lng: dest.lng,
      });
      setTripId(tripResult.trip_id || "trip_demo");
      setTimeout(() => {
        setDriver({ id: "drv_1000", name: "Carlos M.", vehicle: "Toyota Corolla", plate: "ABC-1234", rating: 4.8, photo: "", eta_seconds: fare.eta_minutes * 60 });
        setEta(fare.eta_minutes * 60);
        setState("driver_found");
        setLoading(false);
      }, 3000);
    } catch (e) { setLoading(false); setState("input"); }
  }, [pickupCoords, pickupAddress, dest, fare]);

  const startTracking = useCallback(() => {
    if (!tripId || !driver) return;
    setState("in_progress");
    subscribeToTrip(tripId, (update: TrackingUpdate) => {
      setTracking(update);
      if (update.eta_seconds !== undefined) setEta(update.eta_seconds);
      if (update.driver) setDriver((prev: any) => prev ? { ...prev, lat: update.driver!.lat, lng: update.driver!.lng } : prev);
      if (update.type === "trip_completed") setState("completed");
    });
  }, [tripId, driver]);

  // Auto-start tracking when driver is assigned
  useEffect(() => {
    if (tripId && driver && state === "driver_found") {
      startTracking();
    }
  }, [tripId, driver, state, startTracking]);

  const reset = useCallback(() => {
    setState("pickup_select");
    setDest(null); setDestQuery(""); setRoute(null); setFare(null);
    setDriver(null); setTripId(""); setTracking(null); setEta(0);
    localStorage.removeItem("cytaxi_session");
  }, []);

  const handleCancelTrip = useCallback(() => {
    setState("pickup_select");
    setDriver(null); setTripId(""); setTracking(null); setEta(0);
    setDest(null); setDestQuery(""); setRoute(null); setFare(null);
    localStorage.removeItem("cytaxi_session");
  }, []);

  const handleRejectDriver = useCallback(() => {
    setDriver(null);
    setTripId("");
    setState("input");
  }, []);

  // Session persistence — con versión para migraciones futuras
  const SESSION_VER = 2;
  useEffect(() => {
    try {
      const saved = localStorage.getItem("cytaxi_session");
      if (saved) {
        const s = JSON.parse(saved);
        if (s.v !== SESSION_VER) { localStorage.removeItem("cytaxi_session"); return; }
        if (s.state === "pickup_select" && s.dest) { localStorage.removeItem("cytaxi_session"); return; }
        if (s.state === "confirm" && !s.route) { localStorage.removeItem("cytaxi_session"); return; }
        if (s.state && ["pickup_select", "input", "confirm"].includes(s.state)) {
          setState(s.state);
          if (s.state === "input" && s.dest) setDest(s.dest);
          if (s.state === "input" && s.destQuery) setDestQuery(s.destQuery);
          if (s.pickupAddress) setPickupAddress(s.pickupAddress);
          if (s.pickupCoords) setPickupCoords(s.pickupCoords);
          if (s.route) setRoute(s.route);
          if (s.fare) setFare(s.fare);
        } else { localStorage.removeItem("cytaxi_session"); }
      }
    } catch (e) { console.warn("Session restore failed", e); localStorage.removeItem("cytaxi_session"); }
  }, []);

  useEffect(() => {
    if (!state || state === "searching" || state === "driver_found" || state === "in_progress") return;
    try {
      localStorage.setItem("cytaxi_session", JSON.stringify({
        v: SESSION_VER, state, dest, destQuery, pickupAddress, pickupCoords, route, fare,
      }));
    } catch (e) { console.warn("Session save failed", e); }
  }, [state, dest, destQuery, pickupAddress, pickupCoords, route, fare]);

  const pickupStepProps = { onConfirm: handleConfirmPickup, address: pickupAddress, loading };
  const destStepProps = { destQuery, setDestQuery: setDestQuery, destSuggestions, dest, onSearch: doSearchDest, onSelect: selectDest, onConfirm: handleConfirmDest, loading, pickupAddress, onBack: handleBackPickup, onClearDest: handleClearDest };
  const confirmProps = { pickup: { name: pickupAddress.split(",")[0] || pickupAddress, address: pickupAddress, lat: pickupCoords?.lat || 0, lng: pickupCoords?.lng || 0 }, dest: dest!, route, fare, onConfirm: handleRequestTrip, onBack: () => setState("input"), loading, paymentMethod, onPaymentChange: setPaymentMethod };
  const trackingProps = { state, driver, eta, route, tracking, onStart: startTracking, paymentMethod, pickup: null, dest, onCancel: handleCancelTrip, onRejectDriver: handleRejectDriver };
  const completedProps = { fare, driver, pickup: null, dest, onNewTrip: reset, paymentMethod };

  return {
    state, pickupAddress, pickupCoords, dest, destQuery, route, fare, driver, eta, tracking, loading, tripId,
    sheetRef, contentRef, pickupStepProps, destStepProps, confirmProps, trackingProps, completedProps,
    handleCenterChange, handleMapDestChange, setState,
    setDest, setDestQuery, setRoute, setFare,
  };
}
