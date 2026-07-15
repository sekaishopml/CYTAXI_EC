"use client";
import React, { useState, useCallback, useRef } from "react";
import { TripState, Place, FareBreakdown, TrackingUpdate } from "@/types";
import { searchPlaces, calculateRoute, estimateFare, requestTrip } from "@/services/api";
import { subscribeToTrip } from "@/services/tracking";

import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

export function useTripFlow() {
  const [state, setState] = useState<TripState>("idle");
  const [name] = useState("");
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
  const [paymentMethod] = useState<"cash" | "card">("cash");
  const sheetRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);

  const handleCenterChange = useCallback((data: { lat: number; lng: number; address: string }) => {
    setPickupAddress(data.address);
    setPickupCoords({ lat: data.lat, lng: data.lng });
  }, []);

  const doSearchDest = useCallback(async (q: string) => {
    if (q.length < 3) { setDestSuggestions([]); return; }
    const results = await searchPlaces(q);
    setDestSuggestions(results);
  }, []);

  const selectDest = useCallback((place: Place) => {
    setDest(place);
    setDestQuery(place.name);
    setDestSuggestions([]);
  }, []);

  const handleConfirm = useCallback(async () => {
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
        phone: "0000000000", passenger_name: name || "Usuario",
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
    } catch (e) { setLoading(false); }
  }, [pickupCoords, pickupAddress, dest, fare, name]);

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

  const reset = useCallback(() => {
    setState("idle"); setDest(null); setDestQuery(""); setRoute(null); setFare(null);
    setDriver(null); setTripId(""); setTracking(null); setEta(0);
  }, []);

  const formProps = { name, setName: () => {}, destQuery, setDestQuery: (q: string) => { setDestQuery(q); doSearchDest(q); }, destSuggestions, dest, onSearch: doSearchDest, onSelect: selectDest, onConfirm: handleConfirm, loading, pickupAddress };
  const confirmProps = { pickup: { name: pickupAddress, address: pickupAddress, lat: pickupCoords?.lat || 0, lng: pickupCoords?.lng || 0 }, dest: dest!, route, fare, onConfirm: handleRequestTrip, onBack: () => setState("idle"), loading };
  const trackingProps = { state, driver, eta, route, tracking, onStart: startTracking, paymentMethod, pickup: null, dest };
  const completedProps = { fare, driver, pickup: null, dest, onNewTrip: reset, paymentMethod };

  return {
    state, pickupAddress, pickupCoords, dest, route, fare, driver, eta, tracking, loading, tripId,
    sheetRef, contentRef, formProps, confirmProps, trackingProps, completedProps,
    handleCenterChange, setState,
  };
}
