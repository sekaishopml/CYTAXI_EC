"use client";
import { useEffect, useRef } from "react";

declare global { interface Window { google: any; } }

interface MapViewProps {
  onCenterChange?: (c: { lat: number; lng: number }) => void;
  onMapSelect?: (c: { lat: number; lng: number }) => void;
  onMapReady?: () => void;
  interactive?: boolean;
}

const MAP_STYLE: any[] = [
  { elementType: "geometry", stylers: [{ color: "#1c1c1e" }] },
  { elementType: "labels.text.fill", stylers: [{ color: "#a8a8ae" }] },
  { elementType: "labels.text.stroke", stylers: [{ color: "#1c1c1e" }] },
  { featureType: "administrative.land_parcel", elementType: "labels.text.fill", stylers: [{ color: "#828288" }] },
  { featureType: "landscape.man_made", elementType: "geometry", stylers: [{ color: "#242426" }] },
  { featureType: "landscape.natural", elementType: "geometry", stylers: [{ color: "#202022" }] },
  { featureType: "poi", elementType: "geometry", stylers: [{ color: "#2a2a30" }] },
  { featureType: "poi", elementType: "labels.text.fill", stylers: [{ color: "#c0c0c6" }] },
  { featureType: "poi", elementType: "labels.icon", stylers: [{ saturation: 40, lightness: -20, gamma: 0.9 }] },
  { featureType: "poi.business", elementType: "labels.text.fill", stylers: [{ color: "#f4a460" }] },
  { featureType: "poi.medical", elementType: "geometry", stylers: [{ color: "#3d2828" }] },
  { featureType: "poi.medical", elementType: "labels.text.fill", stylers: [{ color: "#f07070" }] },
  { featureType: "poi.park", elementType: "geometry", stylers: [{ color: "#1e3024" }] },
  { featureType: "poi.park", elementType: "labels.text.fill", stylers: [{ color: "#80c080" }] },
  { featureType: "poi.sports_complex", elementType: "geometry", stylers: [{ color: "#2a3a30" }] },
  { featureType: "poi.sports_complex", elementType: "labels.text.fill", stylers: [{ color: "#70c090" }] },
  { featureType: "poi.attraction", elementType: "labels.text.fill", stylers: [{ color: "#c0a0e0" }] },
  { featureType: "poi.school", elementType: "labels.text.fill", stylers: [{ color: "#70a0d0" }] },
  { featureType: "road", elementType: "geometry", stylers: [{ color: "#3a3a3e" }] },
  { featureType: "road", elementType: "labels.text.fill", stylers: [{ color: "#bcbcc2" }] },
  { featureType: "road.arterial", elementType: "geometry", stylers: [{ color: "#46464a" }] },
  { featureType: "road.highway", elementType: "geometry", stylers: [{ color: "#4e4638" }] },
  { featureType: "road.highway", elementType: "labels.text.fill", stylers: [{ color: "#e0d4b0" }] },
  { featureType: "road.local", elementType: "labels.text.fill", stylers: [{ color: "#a0a0a6" }] },
  { featureType: "transit.line", elementType: "geometry", stylers: [{ color: "#303034" }] },
  { featureType: "transit.station", elementType: "geometry", stylers: [{ color: "#363638" }] },
  { featureType: "transit", elementType: "labels.text.fill", stylers: [{ color: "#a0a0a8" }] },
  { featureType: "transit.station", elementType: "labels.text.fill", stylers: [{ color: "#d0c060" }] },
  { featureType: "water", elementType: "geometry", stylers: [{ color: "#141c26" }] },
  { featureType: "water", elementType: "labels.text.fill", stylers: [{ color: "#607888" }] },
];

const DEFAULT_LNG = -79.8893;
const DEFAULT_LAT = -2.1894;

export function MapView({ onCenterChange, onMapSelect, onMapReady, interactive = true }: MapViewProps) {
  const mapRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const initialized = useRef(false);
  const onCenterChangeRef = useRef(onCenterChange);
  const onMapSelectRef = useRef(onMapSelect);
  const onMapReadyRef = useRef(onMapReady);
  const listenersRef = useRef<any[]>([]);
  onCenterChangeRef.current = onCenterChange;
  onMapSelectRef.current = onMapSelect;
  onMapReadyRef.current = onMapReady;

  useEffect(() => {
    if (typeof window === "undefined" || initialized.current) return;

    const cleanListeners = () => {
      listenersRef.current.forEach(l => { try { l.remove(); } catch {} });
      listenersRef.current = [];
    };

    const tryInit = () => {
      if (!window.google?.maps) return false;
      initialized.current = true;

      const map = new window.google.maps.Map(containerRef.current!, {
        center: { lat: DEFAULT_LAT, lng: DEFAULT_LNG },
        zoom: 15, disableDefaultUI: true, zoomControl: false, styles: MAP_STYLE,
        gestureHandling: "greedy", mapTypeId: "roadmap", mapTypeControl: false,
        streetViewControl: false, fullscreenControl: false, rotateControl: false,
      });

      mapRef.current = map;
      (window as any).__cymap = map;
      onMapReadyRef.current?.();

      listenersRef.current = [
        map.addListener("center_changed", () => {
          const c = map.getCenter();
          onCenterChangeRef.current?.({ lat: c.lat(), lng: c.lng() });
        }),
        map.addListener("click", () => {
          const c = map.getCenter();
          onMapSelectRef.current?.({ lat: c.lat(), lng: c.lng() });
        }),
      ];

      return true;
    };

    if (tryInit()) return;
    let attempts = 0;
    const interval = setInterval(() => {
      attempts++;
      if (tryInit() || attempts >= 15) clearInterval(interval);
    }, 200);

    return () => { clearInterval(interval); cleanListeners(); initialized.current = false; };
  }, []);

  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    map.setOptions({ gestureHandling: interactive ? "greedy" : "none" });
  }, [interactive]);

  return <div ref={containerRef} style={{ width: "100%", height: "100%", background: "#1c1c1e" }} />;
}
