"use client";
import { useEffect, useRef, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import dynamic from "next/dynamic";
import { AnimatePresence, motion } from "framer-motion";
import { useJourneyEngine } from "@/hooks/useJourneyEngine";
import { BottomSheet } from "@/components/bottom-sheet/BottomSheet";
import { TripTimeline } from "@/features/trip/ui/TripTimeline";
import { PickUpStep } from "@/components/states/PickUpStep";
import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { ArrivingState } from "@/components/states/ArrivingState";
import { DestinationState } from "@/components/states/DestinationState";
import { PaymentState } from "@/components/states/PaymentState";
import { RatingState } from "@/components/states/RatingState";
import { CompletedState } from "@/components/states/CompletedState";
import { HomeState } from "@/components/states/HomeState";
import { getStateConfig, RideState } from "@cytaxi/ride-machine";
import { ErrorBoundary } from "@/components/shared/ErrorBoundary";
import { colors, shadows, zIndex } from "@cytaxi/design-tokens";

const API_URL = typeof window !== "undefined" ? `${window.location.protocol}//${window.location.host}/api/v1` : "";

const MapController = dynamic(() => import("@/components/map/MapController").then(m => ({ default: m.MapController })), { ssr: false });

const DEFAULT_LNG = -79.8893;
const DEFAULT_LAT = -2.1894;

const MAP_STATES: RideState[] = ["pickup_select", "input", "confirm", "searching", "driver_found", "arriving", "arrived", "in_progress", "destination"];
const INTERACTIVE_STATES: RideState[] = ["pickup_select", "input"];

async function fetchJSON(url: string, opts?: RequestInit, timeoutMs = 5000) {
  const ctrl = new AbortController();
  const timer = setTimeout(() => ctrl.abort(), timeoutMs);
  try {
    return await fetch(url, { ...opts, signal: ctrl.signal });
  } finally { clearTimeout(timer); }
}

export default function HomePage() {
  const flow = useJourneyEngine();
  const [keyboardOpen, setKeyboardOpen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const dragTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const moveTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [centerAddr, setCenterAddr] = useState("Detectando ubicación...");
  const [keyboardH, setKeyboardH] = useState(0);

  const showMap = MAP_STATES.includes(flow.state);
  const showPin = flow.state === "pickup_select" || (flow.state === "input" && !flow.dest);
  const mapInteractive = INTERACTIVE_STATES.includes(flow.state);
  const config = getStateConfig(flow.state);
  const router = useRouter();

  const reverseGeocode = useCallback(async (lat: number, lng: number) => {
    if (moveTimer.current) clearTimeout(moveTimer.current);
    moveTimer.current = setTimeout(async () => {
      try {
        const res = await fetchJSON(`${API_URL}/geo/reverse?lat=${lat}&lng=${lng}`);
        if (res.ok) {
          const data = await res.json();
          const addr = data?.FormattedAddress || data?.formatted_address || `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
          const hasStreet = data?.street || data?.road;
          if (!hasStreet) { setCenterAddr(`⚠️ Sin acceso: ${addr}`); return; }
          setCenterAddr(addr);
          if (flow.state === "pickup_select") {
            flow.handleCenterChange({ lat, lng, address: addr });
          }
        }
      } catch {
        setCenterAddr(`${lat.toFixed(4)}, ${lng.toFixed(4)}`);
        if (flow.state === "pickup_select") {
          flow.handleCenterChange({ lat, lng, address: `${lat.toFixed(4)}, ${lng.toFixed(4)}` });
        }
      }
    }, 300);
  }, [flow]);

  const handleMapCenterChange = useCallback((c: { lat: number; lng: number }) => {
    setIsDragging(true);
    if (dragTimer.current) clearTimeout(dragTimer.current);
    dragTimer.current = setTimeout(() => setIsDragging(false), 300);
    reverseGeocode(c.lat, c.lng);
  }, [reverseGeocode]);

  const validateAndSelectDest = useCallback(async (lat: number, lng: number) => {
    try {
      const res = await fetchJSON(`${API_URL}/geo/reverse?lat=${lat}&lng=${lng}`);
      let addr = `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
      if (res.ok) {
        const data = await res.json();
        addr = data?.FormattedAddress || data?.formatted_address || addr;
        const hasStreet = data?.street || data?.road;
        if (!hasStreet) {
          setCenterAddr("⚠️ Sin acceso vial — elige otra ubicación");
          setTimeout(() => { if (flow.state === "input" && !flow.dest) setCenterAddr("Selecciona un destino en el mapa"); }, 2500);
          return;
        }
      }
      flow.handleMapDestChange({ lat, lng, address: addr });
      setCenterAddr(addr);
    } catch {
      const addr = `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
      flow.handleMapDestChange({ lat, lng, address: addr });
      setCenterAddr(addr);
    }
  }, [flow]);

  const handleMapSelect = useCallback((c: { lat: number; lng: number }) => {
    if (flow.state === "input") validateAndSelectDest(c.lat, c.lng);
    else if (flow.state === "pickup_select") reverseGeocode(c.lat, c.lng);
  }, [flow, validateAndSelectDest, reverseGeocode]);

  const handleMapReady = useCallback(() => {
    setCenterAddr("Buscando tu ubicación...");
    if (!navigator.geolocation) { reverseGeocode(DEFAULT_LAT, DEFAULT_LNG); return; }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude: lat, longitude: lng } = pos.coords;
        const map = (window as any).__cymap;
        if (map) { map.setCenter({ lat, lng }); map.setZoom(16); }
        reverseGeocode(lat, lng);
      },
      () => { setCenterAddr("Ubicación predeterminada"); reverseGeocode(DEFAULT_LAT, DEFAULT_LNG); },
      { enableHighAccuracy: true, timeout: 8000 }
    );
  }, [reverseGeocode]);

  const handleGPS = useCallback(() => {
    if (!navigator.geolocation) return;
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude: lat, longitude: lng } = pos.coords;
        const map = (window as any).__cymap;
        if (map) map.panTo({ lat, lng });
        reverseGeocode(lat, lng);
      },
      () => {},
      { enableHighAccuracy: true, timeout: 5000 }
    );
  }, [reverseGeocode]);

  useEffect(() => {
    const vv = window.visualViewport;
    if (!vv) return;
    const handleResize = () => {
      const diff = window.innerHeight - vv.height;
      if (diff > 100) {
        setKeyboardOpen(true);
        setKeyboardH(diff);
        requestAnimationFrame(() => document.querySelector('input:focus')?.scrollIntoView({ behavior: "smooth", block: "center" }));
      } else {
        setKeyboardOpen(false);
        setKeyboardH(0);
      }
    };
    vv.addEventListener("resize", handleResize);
    return () => vv.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    if (showPin) setCenterAddr("Detectando ubicación...");
  }, [showPin]);

  const navigateTo = useCallback((nextState: RideState, dir: "forward" | "back" = "forward") => {
    if (nextState === flow.state) return;
    flow.goTo(nextState);
  }, [flow]);

  const renderState = () => {
    switch (flow.state) {
      case "pickup_select":
        return <PickUpStep
          onConfirm={() => flow.handleConfirmPickup()}
          address={flow.pickupAddress}
          loading={flow.loading}
        />;
      case "input":
        return (
          <FormState
            destQuery={flow.destQuery}
            setDestQuery={flow.setDestQuery}
            destSuggestions={flow.destSuggestions}
            dest={flow.dest}
            onSearch={flow.doSearchDest}
            onSelect={flow.selectDest}
            onConfirm={async () => {
              if (!flow.dest) {
                const map = (window as any).__cymap;
                if (map) {
                  const c = map.getCenter();
                  const lat = typeof c.lat === "function" ? c.lat() : c.lat;
                  const lng = typeof c.lng === "function" ? c.lng() : c.lng;
                  await validateAndSelectDest(lat, lng);
                  if (flow.dest) flow.handleConfirmDest();
                }
              } else { flow.handleConfirmDest(); }
            }}
            onBack={() => flow.handleBackToPickup()}
            onClearDest={flow.handleClearDest}
            loading={flow.loading}
            pickupAddress={flow.pickupAddress}
          />
        );
      case "confirm":
        return (
          <ConfirmState
            pickup={{
              name: flow.pickupAddress.split(",")[0] || flow.pickupAddress,
              address: flow.pickupAddress,
              lat: flow.pickupCoords?.lat || 0,
              lng: flow.pickupCoords?.lng || 0,
            }}
            dest={flow.dest!}
            route={flow.route}
            fare={flow.fare}
            onConfirm={flow.handleRequestTrip}
            onBack={() => flow.send("CONFIRM")}
            loading={flow.loading}
            paymentMethod={flow.paymentMethod}
            onPaymentChange={flow.setPaymentMethod}
            vehicleType={flow.vehicleType}
            onVehicleChange={flow.setVehicleType}
            note={flow.note}
            onNoteChange={flow.setNote}
            coupon={flow.coupon}
            onCouponChange={flow.setCoupon}
            scheduledAt={flow.scheduledAt}
            onScheduleChange={flow.setScheduledAt}
          />
        );
      case "searching":
      case "driver_found":
        return (
          <TrackingState
            state={flow.state}
            driver={flow.driver}
            eta={flow.eta}
            route={flow.route}
            paymentMethod={flow.paymentMethod}
            pickup={null}
            dest={flow.dest}
            onCancel={flow.handleCancelTrip}
            onRejectDriver={flow.handleRejectDriver}
            noDrivers={flow.noDrivers}
            onRetry={flow.handleRetrySearch}
          />
        );
      case "arriving":
        return (
          <ArrivingState
            driver={flow.driver}
            eta={flow.eta}
            route={flow.route}
            arrived={false}
            onArrived={() => flow.handleArriveAtPickup()}
            onCancel={flow.handleCancelTrip}
          />
        );
      case "arrived":
        return (
          <ArrivingState
            driver={flow.driver}
            eta={flow.eta}
            route={flow.route}
            arrived={true}
            onArrived={() => flow.handleTripStart()}
          />
        );
      case "in_progress":
        return (
          <TrackingState
            state={flow.state}
            driver={flow.driver}
            eta={flow.eta}
            route={flow.route}
            paymentMethod={flow.paymentMethod}
            pickup={null}
            dest={flow.dest}
            onCancel={flow.handleCancelTrip}
            noDrivers={flow.noDrivers}
            onRetry={flow.handleRetrySearch}
          />
        );
      case "destination":
        return (
          <DestinationState
            dest={flow.dest}
            driver={flow.driver}
            onComplete={() => flow.handleTripComplete()}
          />
        );
      case "payment":
        return (
          <PaymentState
            fare={flow.fare}
            method={flow.paymentMethod || "cash"}
            onDone={() => flow.handlePaymentDone()}
          />
        );
      case "rating":
        return (
          <RatingState
            driver={flow.driver}
            onDone={(score) => flow.handleRatingDone(score)}
          />
        );
      case "completed":
        return (
          <CompletedState
            fare={flow.fare}
            driver={flow.driver}
            pickup={null}
            dest={flow.dest}
            onNewTrip={() => flow.reset()}
            paymentMethod={flow.paymentMethod || "cash"}
          />
        );
      case "travel_home":
        return null; // rendered outside BottomSheet
      default:
        return null;
    }
  };

  return (
    <ErrorBoundary>
    <div role="application" aria-label="CYTAXI" style={{
      width: "100vw", height: "100dvh", overflow: "hidden",
      display: "flex", flexDirection: "column",
      position: "fixed", inset: 0, touchAction: "none",
    }}>
      {flow.state === "travel_home" ? (
        <HomeState onStartTrip={() => flow.send("START_TRIP")} />
      ) : (<>
        {showMap && (
          <div style={{
            position: "absolute", inset: 0, zIndex: 1,
            minHeight: 0,
            touchAction: mapInteractive ? "auto" : "none",
          }}>
            <MapController
              state={flow.state}
              interactive={mapInteractive}
              showPin={showPin}
              pickupCoords={flow.pickupCoords}
              destCoords={flow.dest ? { lat: flow.dest.lat, lng: flow.dest.lng } : null}
              route={flow.route}
              driverLat={flow.driver?.lat}
              driverLng={flow.driver?.lng}
              driverHeading={0}
              onCenterChange={handleMapCenterChange}
              onMapClick={handleMapSelect}
              onMapReady={handleMapReady}
            />

            {/* Map pin overlay — punto A/B animado */}
            <AnimatePresence>
              {showPin && (
                <>
                  {/* Wrapper: posiciona el pin centrado en el mapa. CSS translate para centrar, Motion para animar. */}
                  <div style={{
                    position: "absolute", top: "50%", left: "50%",
                    zIndex: 700, pointerEvents: "none",
                    transform: "translate(-50%, -100%)",
                    filter: "drop-shadow(0 4px 8px rgba(0,0,0,0.3))",
                  }}>
                    <motion.div
                      key={`pin-${flow.state}`}
                      initial={{ opacity: 0, scale: 0.3, y: 10 }}
                      animate={{
                        opacity: 1,
                        scale: isDragging ? 1.15 : 1,
                        y: isDragging ? -5 : 0,
                      }}
                      exit={{ opacity: 0, scale: 0.3, y: 10 }}
                      transition={{ type: "spring", stiffness: 480, damping: 22, mass: 0.6 }}
                    >
                      <svg width="32" height="44" viewBox="0 0 34 44" fill="none">
                        <defs>
                          {flow.state === "input" ? (
                            <linearGradient id="pinGradB" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                              <stop offset="0%" stopColor="#93c5fd" /><stop offset="50%" stopColor="#60a5fa" /><stop offset="100%" stopColor="#3b82f6" />
                            </linearGradient>
                          ) : (
                            <linearGradient id="pinGradA" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                              <stop offset="0%" stopColor="#93c5fd" /><stop offset="50%" stopColor="#3b82f6" /><stop offset="100%" stopColor="#2563eb" />
                            </linearGradient>
                          )}
                        </defs>
                        <path d="M17 0 C26 0 34 8 34 16 C34 26 26 33 20 41 C18.5 42.5 15.5 42.5 14 41 C8 33 0 26 0 16 C0 8 8 0 17 0 Z"
                          fill={flow.state === "input" ? "url(#pinGradB)" : "url(#pinGradA)"} />
                        <circle cx="17" cy="15" r="7" fill="white" opacity="0.93" />
                        <circle cx="17" cy="15" r="3.5" fill={flow.state === "input" ? "#3b82f6" : "#2563eb"} />
                      </svg>
                    </motion.div>
                  </div>

                  {/* Label above pin */}
                  <motion.div
                    key={`label-${flow.state}`}
                    initial={{ opacity: 0, y: 6 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -6, scale: 0.95 }}
                    transition={{ type: "spring", stiffness: 300, damping: 25, delay: 0.1 }}
                    style={{
                      position: "absolute",
                      top: "calc(50% - 82px)",
                      left: "50%",
                      transform: "translateX(-50%)",
                      zIndex: 710,
                      background: "rgba(255,255,255,0.92)",
                      backdropFilter: "blur(16px) saturate(180%)",
                      WebkitBackdropFilter: "blur(16px) saturate(180%)",
                      borderRadius: 10,
                      padding: "6px 14px",
                      boxShadow: "0 2px 12px rgba(0,0,0,0.1), 0 0 0 1px rgba(0,0,0,0.04)",
                      fontSize: 12, fontWeight: 500,
                      fontFamily: "'Inter', sans-serif",
                      color: "#121212",
                      whiteSpace: "nowrap",
                      maxWidth: 220,
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                      pointerEvents: "none",
                      lineHeight: "1.3",
                    }}
                  >
                    {flow.state === "pickup_select" ? (
                      <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
                        <span style={{
                          width: 16, height: 16, borderRadius: "50%",
                          background: "#2563eb", color: "#fff", fontSize: 9, fontWeight: 700,
                          display: "inline-flex", alignItems: "center", justifyContent: "center",
                          flexShrink: 0,
                        }}>A</span>
                        {centerAddr}
                      </span>
                    ) : (
                      <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
                        <span style={{
                          width: 16, height: 16, borderRadius: "50%",
                          background: "#3b82f6", color: "#fff", fontSize: 9, fontWeight: 700,
                          display: "inline-flex", alignItems: "center", justifyContent: "center",
                          flexShrink: 0,
                        }}>B</span>
                        {centerAddr}
                      </span>
                    )}
                  </motion.div>
                </>
              )}
            </AnimatePresence>

            {/* GPS button */}
            <button type="button" onClick={handleGPS} aria-label="Mi ubicación"
              style={{
                position: "absolute", right: 12, bottom: 14, zIndex: 690,
                width: 36, height: 36, borderRadius: 10,
                background: "rgba(255,255,255,0.95)", border: "none",
                boxShadow: shadows.button, cursor: "pointer",
                display: "flex", alignItems: "center", justifyContent: "center",
                color: colors.cobalt,
              }}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="3" /><path d="M12 2v4m0 12v4m10-10h-4M6 12H2" />
              </svg>
            </button>
          </div>
        )}

        {/* Offline banner */}
        {!flow.online && (
          <div role="alert" aria-live="assertive" style={{
            position: "fixed", top: 0, left: 0, right: 0, zIndex: zIndex.max,
            background: "#ea580c", color: "#fff", padding: "10px 16px",
            fontSize: 13, fontWeight: 600, fontFamily: "'Inter', sans-serif",
            textAlign: "center",
          }}>
            Sin conexión — las solicitudes se encolarán automáticamente
          </div>
        )}

        {/* Error banner */}
        {flow.error && (
          <div role="alert" aria-live="assertive" style={{
            position: "fixed",
            top: flow.online ? 0 : 40, left: 0, right: 0, zIndex: zIndex.max,
            background: "#dc2626", color: "#fff", padding: "10px 16px",
            fontSize: 13, fontWeight: 500, fontFamily: "'Inter', sans-serif",
            display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12,
          }}>
            <span>{flow.error}</span>
            <button type="button" onClick={flow.dismissError} aria-label="Cerrar mensaje de error"
              style={{ background: "none", border: "none", color: "#fff", cursor: "pointer", fontSize: 18, padding: "0 4px" }}>
              ✕
            </button>
          </div>
        )}

        {/* Timeline bar */}
        <TripTimeline state={flow.state} />

        {/* Bottom sheet */}
        <BottomSheet
          state={flow.state}
          prevState={flow.prevState}
          direction={flow.direction}
          isTransitioning={flow.isTransitioning}
          sheetRef={flow.sheetRef}
          showNavbar={config.showNavbar}
          keyboardOpen={keyboardOpen}
          keyboardH={keyboardH}
        >
          {renderState()}
        </BottomSheet>
      </>)}

      {/* Navbar */}
      {!keyboardOpen && flow.state !== "travel_home" && (
        <nav style={{
          position: "fixed", bottom: "max(16px, env(safe-area-inset-bottom, 16px))",
          left: "50%", transform: "translateX(-50%)",
          width: "calc(100% - 40px)", maxWidth: 400,
          zIndex: zIndex.navbar,
          background: colors.surface.paper,
          borderRadius: 28,
          boxShadow: "0 -4px 24px rgba(0,0,0,0.12), 0 0 0 1px rgba(0,0,0,0.04)",
          display: "flex", justifyContent: "space-evenly", alignItems: "center",
          padding: "8px 8px",
        }}>
          {([
            { icon: "M3 12l9-9 9 9", label: "Inicio", href: "/", active: ["travel_home", "pickup_select", "input", "confirm"].includes(flow.state) },
            { icon: "M3 17l6-6 4 4 8-8", label: "Actividad", href: "#", active: false },
            { icon: "M5 17h14M5 17a2 2 0 01-2-2V7a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2M5 17l-1 4h12l-1-4", label: "Viajes", href: "/", active: ["searching", "driver_found", "arriving", "arrived", "in_progress", "destination", "payment", "rating", "completed"].includes(flow.state) },
            { icon: "M2 7h20v10H2zm2 0V5a2 2 0 012-2h12a2 2 0 012 2v2", label: "Billetera", href: "#", active: false },
            { icon: "M12 12c2.7 0 5-2.3 5-5s-2.3-5-5-5-5 2.3-5 5 2.3 5 5 5zm0 2c-3.3 0-10 1.7-10 5v3h20v-3c0-3.3-6.7-5-10-5z", label: "Perfil", href: "/profile", active: false },
          ]).map(tab => (
            <a key={tab.label} href={tab.href} onClick={e => { if (tab.href === "#") e.preventDefault(); else router.push(tab.href); }} style={{
              display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center",
              padding: "6px 10px", borderRadius: 24,
              fontSize: 10, fontWeight: 600, fontFamily: "'Inter', sans-serif",
              letterSpacing: "0.01em",
              color: tab.active ? colors.cobalt : "#9ea5a0",
              background: tab.active ? colors.cobaltBg : "transparent",
            }}>
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke={tab.active ? colors.cobalt : "#9ea5a0"}
                strokeWidth={tab.active ? 2.2 : 1.8}
                strokeLinecap="round" strokeLinejoin="round"
                style={{ marginBottom: 3, transition: "stroke 0.25s", transform: tab.active ? "scale(1.08)" : "scale(1)" }}>
                <path d={tab.icon} />
              </svg>
              <span>{tab.label}</span>
            </a>
          ))}
        </nav>
      )}
    </div>
    </ErrorBoundary>
  );
}
