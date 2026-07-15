"use client";
import { useEffect, useRef, useState, useCallback } from "react";

declare global { interface Window { L: any; } }

const API_URL = typeof window !== "undefined"
  ? `${window.location.protocol}//${window.location.host}/api/v1`
  : "http://64.176.219.221/api/v1";

interface MapViewProps {
  onCenterChange?: (coords: { lat: number; lng: number; address: string }) => void;
  onMapReady?: () => void;
  showPin?: boolean;
}

export function MapView({ onCenterChange, onMapReady, showPin = true }: MapViewProps) {
  const mapRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const initialized = useRef(false);
  const [mapReady, setMapReady] = useState(false);
  const [centerAddr, setCenterAddr] = useState("Detectando ubicación...");

  const reverseGeocode = useCallback(async (lat: number, lng: number) => {
    try {
      const res = await fetch(`${API_URL}/geo/reverse?lat=${lat}&lng=${lng}`);
      if (res.ok) {
        const data = await res.json();
        const addr = data?.FormattedAddress || data?.formatted_address || `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
        setCenterAddr(addr);
        onCenterChange?.({ lat, lng, address: addr });
      }
    } catch {
      setCenterAddr(`${lat.toFixed(4)}, ${lng.toFixed(4)}`);
      onCenterChange?.({ lat, lng, address: `${lat.toFixed(4)}, ${lng.toFixed(4)}` });
    }
  }, [onCenterChange]);

  useEffect(() => {
    if (typeof window === "undefined" || !window.L || initialized.current) return;
    const L = window.L;
    const map = L.map(containerRef.current!, { center: [-2.1894, -79.8893], zoom: 15, zoomControl: false, attributionControl: false });
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", { maxZoom: 20 }).addTo(map);
    map.on("moveend", () => { const c = map.getCenter(); reverseGeocode(c.lat, c.lng); });
    mapRef.current = map;
    initialized.current = true;
    setMapReady(true);
    onMapReady?.();
    map.locate({ setView: true, maxZoom: 16 });
    map.on("locationfound", (e: any) => { map.setView([e.latlng.lat, e.latlng.lng], 16); reverseGeocode(e.latlng.lat, e.latlng.lng); });
    map.on("locationerror", () => { reverseGeocode(-2.1894, -79.8893); });
    return () => { map.remove(); initialized.current = false; };
  }, []);

  return (
    <>
      <div ref={containerRef} style={{ width: "100%", height: "100%" }} />
      {mapReady && showPin && (
        <>
          <div style={{ position: "absolute", top: "50%", left: "50%", zIndex: 10, transform: "translate(-50%, -100%)", pointerEvents: "none", display: "flex", flexDirection: "column", alignItems: "center" }}>
            <div style={{ width: 60, height: 60, borderRadius: "50%", background: "rgba(0,108,73,0.15)", position: "absolute", top: 0, left: "50%", transform: "translate(-50%, -50%)", animation: "pulse 2s cubic-bezier(0.215, 0.61, 0.355, 1) infinite" }} />
            <svg width="32" height="42" viewBox="0 0 32 42" fill="none" style={{ filter: "drop-shadow(0 2px 6px rgba(0,0,0,0.3))" }}>
              <path d="M16 0C7.16 0 0 7.16 0 16c0 12 16 26 16 26s16-14 16-26C32 7.16 24.84 0 16 0z" fill="#006c49" />
              <circle cx="16" cy="15" r="8" fill="white" />
            </svg>
          </div>
          <div style={{ position: "absolute", top: "calc(50% - 56px)", left: "50%", transform: "translateX(-50%)", zIndex: 11, background: "rgba(255,255,255,0.95)", backdropFilter: "blur(12px)", padding: "6px 16px", borderRadius: 12, boxShadow: "0 2px 12px rgba(0,0,0,0.12)", fontSize: 13, fontWeight: 500, fontFamily: "Inter", color: "#191c1e", whiteSpace: "nowrap", maxWidth: 260, overflow: "hidden", textOverflow: "ellipsis", border: "1px solid rgba(0,0,0,0.06)", pointerEvents: "none" }}>
            📍 {centerAddr}
          </div>
        </>
      )}
      {mapReady && (
        <button style={{ position: "fixed", right: 16, bottom: "calc(50dvh - 50px)", zIndex: 9, width: 44, height: 44, borderRadius: 14, background: "var(--uk-surface-container-lowest)", border: "none", boxShadow: "0 2px 8px rgba(0,0,0,0.08)", cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 20 }} onClick={() => { mapRef.current?.locate({ setView: true, maxZoom: 16 }); }} aria-label="Mi ubicación">📍</button>
      )}
    </>
  );
}
