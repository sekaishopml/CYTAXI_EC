"use client";
import React, { useState, useCallback, useRef } from "react";
import { TripState, Coordinates, Place, DriverInfo, FareBreakdown, TrackingUpdate } from "@/types";
import { searchPlaces, calculateRoute, estimateFare, requestTrip } from "@/services/api";
import { subscribeToTrip } from "@/services/tracking";

import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

export function useTripFlow() {
  const [state, setState] = useState<TripState>("idle");
  const [phone, setPhone] = useState("");
  const [passengerName, setPassengerName] = useState("");
  const [pickup, setPickup] = useState<Place | null>(null);
  const [dest, setDest] = useState<Place | null>(null);
  const [pickupQuery, setPickupQuery] = useState("");
  const [destQuery, setDestQuery] = useState("");
  const [pickupSuggestions, setPickupSuggestions] = useState<Place[]>([]);
  const [destSuggestions, setDestSuggestions] = useState<Place[]>([]);
  const [route, setRoute] = useState<{ distance_km: number; eta_minutes: number; polyline: string; distance_meters: number; duration_seconds: number } | null>(null);
  const [fare, setFare] = useState<FareBreakdown | null>(null);
  const [driver, setDriver] = useState<DriverInfo | null>(null);
  const [tripId, setTripId] = useState("");
  const [tracking, setTracking] = useState<TrackingUpdate | null>(null);
  const [eta, setEta] = useState(0);
  const [error, setError] = useState("");
  const [paymentMethod, setPaymentMethod] = useState<"cash" | "card">("cash");
  const [loading, setLoading] = useState(false);
  const sheetRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);

  const doSearch = useCallback(async (q: string, isPickup: boolean) => {
    if (q.length < 3) { isPickup ? setPickupSuggestions([]) : setDestSuggestions([]); return; }
    const results = await searchPlaces(q);
    isPickup ? setPickupSuggestions(results) : setDestSuggestions(results);
  }, []);

  const selectPlace = useCallback(async (place: Place, isPickup: boolean) => {
    if (isPickup) { setPickup(place); setPickupQuery(place.name); setPickupSuggestions([]); }
    else { setDest(place); setDestQuery(place.name); setDestSuggestions([]); }
  }, []);

  const handleConfirm = useCallback(async () => {
    if (!pickup || !dest) return;
    setLoading(true);
    setError("");
    const routeData = await calculateRoute({ lat: pickup.lat, lng: pickup.lng }, { lat: dest.lat, lng: dest.lng });
    setRoute(routeData);
    const fareData = await estimateFare(routeData?.distance_km || 5.5, routeData?.duration_seconds || 900);
    setFare(fareData);
    setState("confirm");
    setLoading(false);
  }, [pickup, dest]);

  const handleRequestTrip = useCallback(async () => {
    if (!pickup || !dest || !fare) return;
    setLoading(true);
    setState("searching");
    try {
      const tripResult = await requestTrip({
        phone, passenger_name: passengerName || "Customer",
        origin_address: pickup.address || pickup.name, origin_lat: pickup.lat, origin_lng: pickup.lng,
        dest_address: dest.address || dest.name, dest_lat: dest.lat, dest_lng: dest.lng,
      });
      setTripId(tripResult.trip_id || "trip_demo");
      setTimeout(() => {
        setDriver({
          id: "drv_1000", name: "Carlos M.", vehicle: "Toyota Corolla",
          plate: "ABC-1234", rating: 4.8, photo: "",
          eta_seconds: fare.eta_minutes * 60,
        });
        setEta(fare.eta_minutes * 60);
        setState("driver_found");
        setLoading(false);
      }, 3000);
    } catch (e) { setError((e as Error).message); setLoading(false); }
  }, [pickup, dest, fare, phone, passengerName]);

  const startTracking = useCallback(() => {
    if (!tripId || !driver) return;
    setState("in_progress");
    subscribeToTrip(tripId, (update: TrackingUpdate) => {
      setTracking(update);
      if (update.eta_seconds !== undefined) setEta(update.eta_seconds);
      if (update.driver) {
        setDriver((prev: DriverInfo | null) => prev ? { ...prev, lat: update.driver!.lat, lng: update.driver!.lng } : prev);
      }
      if (update.type === "trip_completed") setState("completed");
    });
  }, [tripId, driver]);

  const reset = useCallback(() => {
    setState("idle"); setPickup(null); setDest(null);
    setPickupQuery(""); setDestQuery(""); setRoute(null); setFare(null);
    setDriver(null); setTripId(""); setTracking(null); setError("");
    setPhone(""); setPassengerName("");
  }, []);

  const formProps = { phone, setPhone, name: passengerName, setName: setPassengerName, pickupQuery, setPickupQuery, destQuery, setDestQuery, pickupSuggestions, destSuggestions, pickup, dest, onSearch: doSearch, onSelect: selectPlace, onConfirm: handleConfirm, loading, paymentMethod, setPaymentMethod };
  const confirmProps = { pickup: pickup!, dest: dest!, route, fare, onConfirm: handleRequestTrip, onBack: () => setState("idle"), loading };
  const trackingProps = { state, driver, eta, route, tracking, onStart: startTracking, paymentMethod, pickup, dest };
  const completedProps = { fare, driver, pickup, dest, onNewTrip: reset, paymentMethod };

  return {
    state, pickup, dest, route, fare, driver, eta, tracking, error, loading, tripId,
    sheetRef, contentRef,
    formProps, confirmProps, trackingProps, completedProps,
    setState,
  };
}
