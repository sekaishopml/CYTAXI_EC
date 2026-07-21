"use client";
import { useEffect, useRef, useState } from "react";

interface MapPreviewProps {
  onClick: () => void;
  shape?: "rounded" | "p-top";
  noShadow?: boolean;
}

const GUAYAQUIL = { lat: -2.170997, lng: -79.922359 };

export function MapPreview({ onClick, shape = "rounded", noShadow }: MapPreviewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<google.maps.Map | null>(null);
  const markerRef = useRef<google.maps.Marker | null>(null);
  const [center, setCenter] = useState(GUAYAQUIL);
  const [ready, setReady] = useState(false);
  const [live, setLive] = useState(false);
  const [hover, setHover] = useState(false);

  useEffect(() => {
    if (typeof window === "undefined") return;
    if (!navigator.geolocation) { setReady(true); return; }
    navigator.geolocation.getCurrentPosition(
      (pos) => { setCenter({ lat: pos.coords.latitude, lng: pos.coords.longitude }); setReady(true); },
      () => { setReady(true); },
      { timeout: 5000, enableHighAccuracy: false },
    );
  }, []);

  useEffect(() => {
    if (!ready || !containerRef.current || mapRef.current) return;

    const tryInit = () => {
      if (!window.google?.maps || !containerRef.current) return false;
      const map = new google.maps.Map(containerRef.current, {
        center,
        zoom: 15,
        disableDefaultUI: true,
        zoomControl: false,
        mapTypeControl: false,
        streetViewControl: false,
        fullscreenControl: false,
        gestureHandling: "none",
        keyboardShortcuts: false,
        styles: [
          { featureType: "poi", stylers: [{ visibility: "off" }] },
          { featureType: "transit", stylers: [{ visibility: "off" }] },
        ],
      });

      markerRef.current = new google.maps.Marker({
        position: center,
        map,
        animation: google.maps.Animation.DROP,
        icon: {
          path: google.maps.SymbolPath.CIRCLE,
          scale: 10,
          fillColor: "#3b82f6",
          fillOpacity: 1,
          strokeColor: "#ffffff",
          strokeWeight: 3,
        },
      });

      mapRef.current = map;
      setLive(true);
      return true;
    };

    if (tryInit()) return;

    let attempts = 0;
    const interval = setInterval(() => {
      attempts++;
      if (tryInit() || attempts >= 10) clearInterval(interval);
    }, 500);

    return () => clearInterval(interval);
  }, [ready, center]);

  return (
      <div
        role="button"
        tabIndex={0}
        onClick={onClick}
        onKeyDown={e => { if (e.key === "Enter" || e.key === " ") onClick(); }}
        onMouseEnter={() => setHover(true)}
        onMouseLeave={() => setHover(false)}
        aria-label="Abrir mapa para iniciar viaje"
        style={{
          width: "100%",
          height: 180,
          borderRadius: shape === "p-top" ? "24px 24px 0 0" : 24,
          clipPath: shape === "p-top" ? "polygon(24px 0, calc(100% - 24px) 0, 100% 24px, 100% 180px, 88% 179.9px, 76% 179.4px, 64% 178.4px, 52% 176.9px, 40% 174.8px, 28% 172.1px, 16% 168.7px, 0 163px, 0 24px, 24px 0)" : undefined,
          cursor: "pointer",
          position: "relative",
          overflow: "hidden",
          boxShadow: noShadow ? undefined : (hover ? "0 12px 40px rgba(0,0,0,0.2)" : "0 8px 32px rgba(0,0,0,0.15)"),
        transform: hover ? "scale(1.01)" : "scale(1)",
        transition: "transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.25s ease",
      }}
    >
      <div ref={containerRef} style={{ width: "100%", height: "100%" }} />
      {!live && (
        <div style={{
          position: "absolute", inset: 0,
          background: "linear-gradient(135deg, #e8edf2 0%, #dde3ea 100%)",
          display: "flex", alignItems: "center", justifyContent: "center",
          zIndex: 1,
        }}>
          <div style={{
            width: 40, height: 40, borderRadius: "50%",
            border: "3px solid rgba(37,99,235,0.15)",
            borderTopColor: "#3b82f6",
            animation: "spin 0.8s linear infinite",
          }} />
          <style>{`@keyframes spin{to{transform:rotate(360deg)}}`}</style>
        </div>
      )}
      <div style={{
        position: "absolute", inset: 0,
        background: "linear-gradient(180deg, transparent 0%, rgba(0,0,0,0.12) 100%)",
        pointerEvents: "none",
      }} />
      <div style={{
        position: "absolute", bottom: 24, left: 14,
        display: "flex", alignItems: "center", gap: 6,
        background: "rgba(255,255,255,0.92)",
        backdropFilter: "blur(12px) saturate(180%)",
        WebkitBackdropFilter: "blur(12px) saturate(180%)",
        borderRadius: 20,
        padding: "6px 12px 6px 8px",
        boxShadow: "0 4px 16px rgba(0,0,0,0.18)",
        pointerEvents: "none",
        maxWidth: "calc(100% - 28px)",
      }}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#2563eb" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z" />
          <circle cx="12" cy="10" r="3" />
        </svg>
        <span style={{
          fontSize: 11, fontWeight: 600,
          color: "#121212", fontFamily: "'Inter', sans-serif",
          whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis",
        }}>
          {ready && (center.lat !== GUAYAQUIL.lat || center.lng !== GUAYAQUIL.lng)
            ? "Tu ubicación · Toca para comenzar"
            : "Guayaquil · Toca para comenzar"}
        </span>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="rgba(0,0,0,0.35)" strokeWidth="2.5" strokeLinecap="round"><path d="M9 18l6-6-6-6"/></svg>
      </div>
    </div>
  );
}
