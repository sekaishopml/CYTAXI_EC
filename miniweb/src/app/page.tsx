"use client";
import { useEffect, useRef, useState, useCallback } from "react";
import dynamic from "next/dynamic";
import { gsap } from "@/services/gsap";
import { useTripFlow } from "@/hooks/useTripFlow";
import { TripState } from "@/types";
import { PickUpStep } from "@/components/states/PickUpStep";
import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

const API_URL = typeof window !== "undefined" ? `${window.location.protocol}//${window.location.host}/api/v1` : "";

const MapView = dynamic(() => import("@/components/MapView").then(m => ({ default: m.MapView })), { ssr: false });

const DEFAULT_LNG = -79.8893;
const DEFAULT_LAT = -2.1894;

async function fetchJSON(url: string, opts?: RequestInit, timeoutMs = 5000) {
  const ctrl = new AbortController();
  const timer = setTimeout(() => ctrl.abort(), timeoutMs);
  try {
    const res = await fetch(url, { ...opts, signal: ctrl.signal });
    return res;
  } finally { clearTimeout(timer); }
}

export default function HomePage() {
  const flow = useTripFlow();
  const prevState = useRef(flow.state);
  const initial = useRef(true);
  const [keyboardOpen, setKeyboardOpen] = useState(false);
  const [transitioning, setTransitioning] = useState(false);
  const navDir = useRef<"forward" | "back">("forward");
  const [isDragging, setIsDragging] = useState(false);
  const dragTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const moveTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [centerAddr, setCenterAddr] = useState("Detectando ubicación...");
  const directionsRendererRef = useRef<any>(null);
  const routeMarkersRef = useRef<any[]>([]);
  const routeReqId = useRef(0);

  // Validación con Google Maps Geocoder para detectar ríos/mar/bosques
  const validateGoogle = useCallback(async (lat: number, lng: number): Promise<{valid: boolean; address: string | null}> => {
    const google = (window as any).google;
    if (!google?.maps?.Geocoder) return { valid: true, address: null };
    return new Promise(resolve => {
      const timer = setTimeout(() => resolve({ valid: true, address: null }), 3000);
      try {
        const geocoder = new google.maps.Geocoder();
        geocoder.geocode({ location: { lat, lng } }, (results: any, status: string) => {
          clearTimeout(timer);
          if (status !== "OK" || !results?.length) { resolve({ valid: false, address: null }); return; }
          const r = results[0];
          const types = r.types || [];
          const ALLOWED = ["street_address","route","premise","subpremise","establishment","park","airport","university","school","parking","bus_station"];
          const isGood = types.some((t: string) => ALLOWED.includes(t));
          resolve({ valid: isGood, address: isGood ? r.formatted_address : null });
        });
      } catch { clearTimeout(timer); resolve({ valid: true, address: null }); }
    });
  }, []);

  const reverseGeocode = useCallback(async (lat: number, lng: number) => {
    if (moveTimer.current) clearTimeout(moveTimer.current);
    moveTimer.current = setTimeout(async () => {
      try {
        const res = await fetchJSON(`${API_URL}/geo/reverse?lat=${lat}&lng=${lng}`);
        if (res.ok) {
          const data = await res.json();
          const addr = data?.FormattedAddress || data?.formatted_address || `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
          const hasStreet = data?.street || data?.road;
          const valid = await validateGoogle(lat, lng);
          if (!valid.valid) {
            setCenterAddr("⚠️ Sin acceso — río, mar o zona sin calles");
            return;
          }
          setCenterAddr(hasStreet ? addr : `⚠️ Sin acceso: ${addr}`);
          if (flow.state !== "input" && hasStreet) {
            flow.handleCenterChange({ lat, lng, address: addr });
          }
        }
      } catch {
        setCenterAddr(`${lat.toFixed(4)}, ${lng.toFixed(4)}`);
        if (flow.state !== "input") {
          flow.handleCenterChange({ lat, lng, address: `${lat.toFixed(4)}, ${lng.toFixed(4)}` });
        }
      }
    }, 300);
  }, [flow, validateGoogle]);

  const handleMapCenterChange = useCallback((c: { lat: number; lng: number }) => {
    setIsDragging(true);
    if (dragTimer.current) clearTimeout(dragTimer.current);
    dragTimer.current = setTimeout(() => setIsDragging(false), 300);
    reverseGeocode(c.lat, c.lng);
  }, [reverseGeocode]);

  const validateAndSelectDest = useCallback(async (lat: number, lng: number) => {
    try {
      const valid = await validateGoogle(lat, lng);
      if (!valid.valid) {
        setCenterAddr("⚠️ Sin acceso — río, mar o zona sin calles");
        setTimeout(() => { if (flow.state === "input" && !flow.dest) setCenterAddr("Selecciona un destino en el mapa"); }, 2500);
        return;
      }
      const res = await fetchJSON(`${API_URL}/geo/reverse?lat=${lat}&lng=${lng}`);
      if (!res.ok) return;
      const data = await res.json();
      const addr = data?.FormattedAddress || data?.formatted_address || valid.address;
      const hasStreet = data?.street || data?.road;
      if (!hasStreet || !addr) {
        setCenterAddr("⚠️ Sin acceso vial — elige otra ubicación");
        setTimeout(() => { if (flow.state === "input" && !flow.dest) setCenterAddr("Selecciona un destino en el mapa"); }, 2500);
        return;
      }
      flow.handleMapDestChange({ lat, lng, address: addr });
      setCenterAddr(addr);
    } catch {}
  }, [flow, validateGoogle]);

  const handleMapSelect = useCallback((c: { lat: number; lng: number }) => {
    if (flow.state === "input") {
      validateAndSelectDest(c.lat, c.lng);
    } else if (flow.state === "pickup_select") {
      reverseGeocode(c.lat, c.lng);
    }
  }, [flow, validateAndSelectDest, reverseGeocode]);

  const handleMapReady = useCallback(() => {
    setCenterAddr("Buscando tu ubicación...");
    if (!navigator.geolocation) {
      reverseGeocode(DEFAULT_LAT, DEFAULT_LNG);
      return;
    }
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
      if (diff > 100) { setKeyboardOpen(true); requestAnimationFrame(() => document.querySelector('input:focus')?.scrollIntoView({ behavior: "smooth", block: "center" })); }
      else setKeyboardOpen(false);
    };
    vv.addEventListener("resize", handleResize);
    return () => vv.removeEventListener("resize", handleResize);
  }, []);

  const navigateTo = useCallback((nextState: TripState, dir: "forward" | "back" = "forward") => {
    if (transitioning || nextState === flow.state) return;
    navDir.current = dir;
    const content = flow.contentRef.current;
    if (!content) { flow.setState(nextState); return; }
    setTransitioning(true);
    gsap.to(content, { opacity: 0, y: dir === "back" ? 24 : -24, scale: 0.97, duration: 0.18, ease: "power3.in", onComplete: () => {
      // Limpiar estado al navegar a pickup_select (evita datos stale)
      if (nextState === "pickup_select") {
        flow.setDest(null);
        flow.setDestQuery("");
        flow.setRoute(null);
        flow.setFare(null);
      }
      flow.setState(nextState);
    } });
  }, [flow, transitioning]);

  useEffect(() => {
    const content = flow.contentRef.current;
    if (!content) return;
    if (initial.current) { initial.current = false; prevState.current = flow.state; gsap.fromTo(content, { opacity: 0, y: 20, scale: 0.97 }, { opacity: 1, y: 0, scale: 1, duration: 0.4, ease: "expo.out" }); return; }
    if (prevState.current === flow.state) return;
    if (transitioning) {
      gsap.fromTo(content, { opacity: 0, y: navDir.current === "back" ? -24 : 24, scale: 0.97 }, { opacity: 1, y: 0, scale: 1, duration: 0.35, ease: "expo.out", onComplete: () => setTransitioning(false) });
    } else {
      const t = ["searching", "driver_found", "in_progress"].includes(flow.state);
      if (t) gsap.fromTo(content, { opacity: 0, y: 30, scale: 0.95 }, { opacity: 1, y: 0, scale: 1, duration: 0.5, ease: "back.out(1.4)" });
      else if (flow.state === "completed") gsap.fromTo(content, { opacity: 0, scale: 0.92 }, { opacity: 1, scale: 1, duration: 0.6, ease: "back.out(1.4)" });
      else gsap.fromTo(content, { opacity: 0, y: 16, scale: 0.97 }, { opacity: 1, y: 0, scale: 1, duration: 0.35, ease: "expo.out" });
      setTransitioning(false);
    }
    prevState.current = flow.state;
  }, [flow.state]);

  // GSAP timeline premium: mapa → pin → panel → nav (staggered entry)
  useEffect(() => {
    const tl = gsap.timeline({ defaults: { ease: "expo.out" } });
    tl.fromTo(flow.sheetRef.current, { opacity: 0, y: 50 }, { opacity: 1, y: 0, duration: 0.7, delay: 0.1 })
      .fromTo("#cytaxi-pin", { opacity: 0, scale: 0, y: -20 }, { opacity: 1, scale: 1, y: 0, duration: 0.5, ease: "back.out(1.7)" }, "-=0.3")
      .fromTo("#cytaxi-label", { opacity: 0, y: -10 }, { opacity: 1, y: 0, duration: 0.4 }, "-=0.2")
      .fromTo("nav", { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.5 }, "-=0.1");
    return () => { tl.kill(); };
  }, []);

  const showPin = (flow.state === "pickup_select" || flow.state === "input") && !(flow.pickupCoords && flow.dest);
  const mapInteractive = flow.state === "pickup_select" || flow.state === "input";
  const showMap = flow.state !== "completed";

  const [pinEnter, setPinEnter] = useState(false);
  useEffect(() => {
    if (showPin) setPinEnter(true);
    const t = setTimeout(() => setPinEnter(false), 700);
    return () => clearTimeout(t);
  }, [showPin]);

  // Markers inmediatos para punto A (pickup) y punto B (destino)
  useEffect(() => {
    routeMarkersRef.current.forEach(m => m.setMap(null));
    routeMarkersRef.current = [];
    const map = (window as any).__cymap;
    const google = (window as any).google;
    if (!map || !google?.maps) return;

    const markers: any[] = [];

    // Punto A: pickup (green square) — en input o confirm
    if (flow.pickupCoords && (flow.state === "input" || flow.state === "confirm")) {
      markers.push(new google.maps.Marker({
        position: new google.maps.LatLng(flow.pickupCoords.lat, flow.pickupCoords.lng),
        map, zIndex: 200,
        icon: { url: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='18' height='18' viewBox='0 0 16 16'%3E%3Crect x='0' y='0' width='16' height='16' rx='3' fill='%23ffffff'/%3E%3Crect x='2.5' y='2.5' width='11' height='11' rx='1.5' fill='%2300a152'/%3E%3C/svg%3E", anchor: new google.maps.Point(8, 8) },
      }));
    }

    // Punto B: destino (blue triangle) — cuando está seleccionado
    if (flow.dest && (flow.state === "input" || flow.state === "confirm")) {
      markers.push(new google.maps.Marker({
        position: new google.maps.LatLng(flow.dest.lat, flow.dest.lng),
        map, zIndex: 200,
        icon: { url: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='18' height='20' viewBox='0 0 16 18'%3E%3Cpath d='M8 0 L16 18 L0 18 Z' fill='%23ffffff' stroke-linejoin='round'/%3E%3Cpath d='M8 2.5 L13.5 15.5 L2.5 15.5 Z' fill='%23448aff' stroke-linejoin='round'/%3E%3C/svg%3E", anchor: new google.maps.Point(8, 9) },
      }));
    }

    routeMarkersRef.current = markers;
  }, [flow.pickupCoords, flow.dest, flow.state]);

  // Ruta con Google Maps — gradiente suave + flash blanco pulsante (GSAP)
  const flashAnimRef = useRef<any>(null);
  useEffect(() => {
    const map = (window as any).__cymap;
    const google = (window as any).google;

    if (directionsRendererRef.current) { directionsRendererRef.current.setMap(null); directionsRendererRef.current = null; }
    if ((window as any).__gradientPolyline) { (window as any).__gradientPolyline.forEach((p: any) => p.setMap(null)); (window as any).__gradientPolyline = null; }
    if ((window as any).__flashLine) { (window as any).__flashLine.setMap(null); (window as any).__flashLine = null; }
    if (flashAnimRef.current) { flashAnimRef.current.kill(); flashAnimRef.current = null; }

    if (!map || !google?.maps || !flow.pickupCoords || !flow.dest || (flow.state !== "input" && flow.state !== "confirm")) return;

    const reqId = ++routeReqId.current;
    const ds = new google.maps.DirectionsService();
    const dr = new google.maps.DirectionsRenderer({
      map, polylineOptions: { strokeColor: "#ffffff", strokeWeight: 4, strokeOpacity: 0.3 },
      suppressMarkers: true, suppressInfoWindows: true, preserveViewport: false, draggable: false,
    });
    directionsRendererRef.current = dr;

    function gradientColor(t: number) {
      const r = Math.round(0 + t * 68);
      const g = Math.round(161 - t * 23);
      const b = Math.round(82 + t * 173);
      return `rgb(${r},${g},${b})`;
    }

    function drawRouteAndFlash(status: string, result: any, pCoords: any, dCoords: any) {
      let path: any[];
      if (status === "OK" && result?.routes?.length) {
        path = result.routes[0].overview_path;
        dr.setDirections(result);
        map.fitBounds(result.routes[0].bounds, { padding: 40 });
      } else {
        path = [new google.maps.LatLng(pCoords.lat, pCoords.lng), new google.maps.LatLng(dCoords.lat, dCoords.lng)];
      }

      // Gradiente suave (baja intensidad)
      const SEG = 10, SL = Math.floor(path.length / SEG);
      const polys: any[] = [];
      for (let i = 0; i < SEG; i++) {
        const start = i * SL;
        const end = (i === SEG - 1) ? path.length : (i + 1) * SL + 1;
        const seg = path.slice(start, Math.min(end, path.length));
        if (seg.length < 2) continue;
        const t = i / (SEG - 1);
        polys.push(new google.maps.Polyline({ path: seg, strokeColor: gradientColor(t), strokeWeight: 3, strokeOpacity: 0.35, map, zIndex: 50 }));
      }
      (window as any).__gradientPolyline = polys;

      // Flash elegante: cometa con estela (brillante → medio → tenue)
      const flashLine = new google.maps.Polyline({
        path, strokeColor: "#ffffff", strokeWeight: 4, strokeOpacity: 0.2, map, zIndex: 60,
        icons: [
          { icon: { path: "M 0 -4 0 4", strokeOpacity: 0.9, scale: 3.5 }, offset: "0%", repeat: "60px" },
          { icon: { path: "M 0 -3 0 3", strokeOpacity: 0.4, scale: 2 }, offset: "5%", repeat: "60px" },
          { icon: { path: "M 0 -2 0 2", strokeOpacity: 0.15, scale: 1 }, offset: "10%", repeat: "60px" },
        ],
      });
      (window as any).__flashLine = flashLine;

      const animObj = { offset: 0 };
      flashAnimRef.current = gsap.to(animObj, {
        offset: 100, duration: 6, ease: "power1.inOut", repeat: -1,
        onUpdate: () => {
          const o = animObj.offset + "%";
          flashLine.setOptions({
            icons: [
              { icon: { path: "M 0 -4 0 4", strokeOpacity: 0.9, scale: 3.5 }, offset: o, repeat: "60px" },
              { icon: { path: "M 0 -3 0 3", strokeOpacity: 0.4, scale: 2 }, offset: (animObj.offset + 5) + "%", repeat: "60px" },
              { icon: { path: "M 0 -2 0 2", strokeOpacity: 0.15, scale: 1 }, offset: (animObj.offset + 10) + "%", repeat: "60px" },
            ],
          });
        },
      });
    }

    ds.route({
      origin: new google.maps.LatLng(flow.pickupCoords.lat, flow.pickupCoords.lng),
      destination: new google.maps.LatLng(flow.dest.lat, flow.dest.lng),
      travelMode: google.maps.TravelMode.DRIVING, provideRouteAlternatives: false, avoidFerries: true,
    } as any, (result: any, status: any) => {
      if (reqId !== routeReqId.current) return;
      drawRouteAndFlash(status, result, flow.pickupCoords!, flow.dest!);
    });

    return () => {
      if (directionsRendererRef.current) { directionsRendererRef.current.setMap(null); directionsRendererRef.current = null; }
      if ((window as any).__gradientPolyline) { (window as any).__gradientPolyline.forEach((p: any) => p.setMap(null)); (window as any).__gradientPolyline = null; }
      if ((window as any).__flashLine) { (window as any).__flashLine.setMap(null); (window as any).__flashLine = null; }
      if (flashAnimRef.current) { flashAnimRef.current.kill(); flashAnimRef.current = null; }
    };
  }, [flow.pickupCoords, flow.dest, flow.state]);

  const wrappedProps = {
    pickupStepProps: { ...flow.pickupStepProps, onConfirm: () => {
      if (!flow.pickupCoords) { setCenterAddr("Esperando ubicación GPS..."); return; }
      navigateTo("input", "forward");
    } },
    destStepProps: {
      ...flow.destStepProps,
      onBack: () => navigateTo("pickup_select", "back"),
      onConfirm: async () => {
        if (flow.state === "input" && !flow.dest) {
          const map = (window as any).__cymap;
          if (map) {
            const c = map.getCenter();
            await validateAndSelectDest(c.lat(), c.lng());
            if (flow.dest) flow.destStepProps.onConfirm();
          }
        } else if (flow.state === "input" && flow.dest) {
          flow.destStepProps.onConfirm();
        }
      },
    },
    confirmProps: { ...flow.confirmProps, onBack: () => navigateTo("input", "back") },
    trackingProps: flow.trackingProps,
    completedProps: { ...flow.completedProps, onNewTrip: () => navigateTo("pickup_select", "back") },
  };

  return (
    <div style={{ width: "100vw", height: "100dvh", overflow: "hidden", display: "flex", flexDirection: "column", position: "fixed", inset: 0, touchAction: "none" }}>
      {showMap && (
        <div style={{ flex: 1, position: "relative", minHeight: 180, touchAction: mapInteractive ? "auto" : "none" }}>
          <MapView onCenterChange={handleMapCenterChange} onMapSelect={handleMapSelect} onMapReady={handleMapReady} interactive={mapInteractive} />
          {showPin && (<>
          <div id="cytaxi-pin" className={pinEnter ? "pin-enter" : ""} style={{
            position: "absolute", top: "50%", left: "50%", zIndex: 700,
            transform: pinEnter ? undefined : `translate(-50%, -100%) scale(${isDragging ? 1.08 : 1}) translateY(${isDragging ? -6 : 0}px)`,
            pointerEvents: "none",
            transition: pinEnter ? undefined : (isDragging ? "none" : "transform 0.45s cubic-bezier(0.34, 1.56, 0.64, 1)") }}>
            <svg width="30" height="40" viewBox="0 0 34 44" fill="none">
              <defs>
                <linearGradient id="pinGradGreen" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                  <stop offset="0%" stopColor="#00e676" /><stop offset="50%" stopColor="#00a152" /><stop offset="100%" stopColor="#006c49" />
                </linearGradient>
                <linearGradient id="pinGradBlue" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                  <stop offset="0%" stopColor="#82b1ff" /><stop offset="50%" stopColor="#448aff" /><stop offset="100%" stopColor="#1565c0" />
                </linearGradient>
              </defs>
              <path d="M17 0 C26 0 34 8 34 16 C34 26 26 33 20 41 C18.5 42.5 15.5 42.5 14 41 C8 33 0 26 0 16 C0 8 8 0 17 0 Z" fill={flow.state === "input" ? "url(#pinGradBlue)" : "url(#pinGradGreen)"} />
              <path d="M17 4 C24 4 30 10 30 16 C30 24 23 30 19 37 C18 39 16 39 15 37 C11 30 4 24 4 16 C4 10 10 4 17 4 Z" fill="rgba(255,255,255,0.15)" />
              <circle cx="17" cy="15" r="7" fill="white" opacity="0.93" />
              <circle cx="17" cy="15" r="3" fill={flow.state === "input" ? "#1565c0" : "#006c49"} />
            </svg>
          </div>
          <div id="cytaxi-label" className={pinEnter ? "label-enter" : ""} style={{
            position: "absolute", top: "calc(50% - 80px)", left: "50%", transform: pinEnter ? undefined : "translateX(-50%)",
            zIndex: 710, background: "rgba(255,255,255,0.35)", backdropFilter: "blur(28px) saturate(500%) brightness(1.1)", WebkitBackdropFilter: "blur(28px) saturate(500%) brightness(1.1)", borderRadius: 8, padding: "5px 14px", boxShadow: "0 8px 40px rgba(0,0,0,0.15), inset 0 1px 0 rgba(255,255,255,0.9), inset 0 -1px 0 rgba(0,0,0,0.04), 0 0 0 0.5px rgba(255,255,255,0.4)",
            fontSize: 12, fontWeight: 500, fontFamily: "Inter", color: "#121212", whiteSpace: "nowrap", maxWidth: 220, overflow: "hidden", textOverflow: "ellipsis", pointerEvents: "none", lineHeight: "1.3", transition: "opacity 0.25s" }}>
            {centerAddr}
          </div>
          </>)}
          <button onClick={handleGPS} aria-label="Mi ubicación" onMouseDown={(e) => { e.currentTarget.style.transform = "scale(0.88)"; e.currentTarget.style.boxShadow = "0 1px 4px rgba(0,0,0,0.1)"; }} onMouseUp={(e) => { e.currentTarget.style.transform = "scale(1)"; e.currentTarget.style.boxShadow = "0 2px 8px rgba(0,0,0,0.08)"; }} onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; e.currentTarget.style.boxShadow = "0 2px 8px rgba(0,0,0,0.08)"; }}
            style={{ position: "absolute", right: 12, bottom: 14, zIndex: 690, width: 36, height: 36, borderRadius: 10, background: "rgba(255,255,255,0.95)", border: "none", boxShadow: "0 2px 8px rgba(0,0,0,0.08)", cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center", transition: "transform 0.12s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.12s ease" }}>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#00a152" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="3" /><path d="M12 2v4m0 12v4m10-10h-4M6 12H2" /></svg>
          </button>
        </div>
      )}
      <div ref={flow.sheetRef} style={{ width: "100%", flex: showMap ? "0 0 auto" : 1, overflow: "hidden", background: "var(--uk-background)", paddingBottom: keyboardOpen ? 0 : 56 }}>
        <div ref={flow.contentRef} style={{ opacity: 0, overflowY: keyboardOpen && flow.state === "input" ? "auto" : "hidden", WebkitOverflowScrolling: "touch" }}>
          {flow.state === "pickup_select" && <PickUpStep {...wrappedProps.pickupStepProps} />}
          {flow.state === "input" && <FormState {...wrappedProps.destStepProps} />}
          {flow.state === "confirm" && <ConfirmState {...wrappedProps.confirmProps} />}
          {["searching", "driver_found", "in_progress"].includes(flow.state) && <TrackingState {...wrappedProps.trackingProps} />}
          {flow.state === "completed" && <CompletedState {...wrappedProps.completedProps} />}
        </div>
      </div>
      {!keyboardOpen && (
        <nav style={{ position: "fixed", bottom: 0, width: "100%", zIndex: 50, background: "rgba(255,255,255,0.88)", backdropFilter: "blur(20px) saturate(180%)", borderTop: "1px solid rgba(0,0,0,0.04)", display: "flex", justifyContent: "space-around", alignItems: "center", padding: "6px 0 max(8px, env(safe-area-inset-bottom))" }}>
          {([
            { icon: "M3 12l9-9 9 9", label: "Inicio", active: ["pickup_select", "input", "confirm"].includes(flow.state) },
            { icon: "M3 17l6-6 4 4 8-8", label: "Actividad", active: false },
            { icon: "M5 17h14M5 17a2 2 0 01-2-2V7a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2M5 17l-1 4h12l-1-4", label: "Viajes", active: ["searching", "driver_found", "in_progress", "completed"].includes(flow.state) },
            { icon: "M2 7h20v10H2zm2 0V5a2 2 0 012-2h12a2 2 0 012 2v2", label: "Billetera", active: false },
            { icon: "M12 12c2.7 0 5-2.3 5-5s-2.3-5-5-5-5 2.3-5 5 2.3 5 5 5zm0 2c-3.3 0-10 1.7-10 5v3h20v-3c0-3.3-6.7-5-10-5z", label: "Perfil", active: false },
          ]).map(tab => (
            <a key={tab.label} href="#" onClick={e => e.preventDefault()} style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", padding: "6px 16px", borderRadius: 16, fontSize: 10, fontWeight: 600, fontFamily: "Inter", transition: "all 0.25s cubic-bezier(0.4, 0, 0.2, 1)", letterSpacing: "0.01em", color: tab.active ? "#006c49" : "#9ea5a0", background: tab.active ? "rgba(0,108,73,0.08)" : "transparent" }}>
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke={tab.active ? "#006c49" : "#9ea5a0"} strokeWidth={tab.active ? 2.2 : 1.8} strokeLinecap="round" strokeLinejoin="round" style={{ marginBottom: 3, transition: "stroke 0.25s", transform: tab.active ? "scale(1.08)" : "scale(1)" }}>
                <path d={tab.icon} />
              </svg>
              <span>{tab.label}</span>
            </a>
          ))}
        </nav>
      )}
    </div>
  );
}
