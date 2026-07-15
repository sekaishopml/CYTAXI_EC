"use client";
import { useEffect, useRef, useState } from "react";
import { Coordinates } from "@/types";

declare global { interface Window { L: any; } }

interface MapViewProps {
  pickupCoords: Coordinates | null;
  destCoords: Coordinates | null;
  driverCoords: Coordinates | null;
  polyline: string | null;
  fitBounds: boolean;
  onLocate?: (coords: Coordinates) => void;
}

export function MapView({ pickupCoords, destCoords, driverCoords, polyline, fitBounds, onLocate }: MapViewProps) {
  const mapRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const pickupMarkerRef = useRef<any>(null);
  const destMarkerRef = useRef<any>(null);
  const driverMarkerRef = useRef<any>(null);
  const routeRef = useRef<any>(null);
  const initialized = useRef(false);
  const [mapReady, setMapReady] = useState(false);

  useEffect(() => {
    if (typeof window === "undefined" || !window.L || initialized.current) return;
    const L = window.L;
    const map = L.map(containerRef.current!, {
      center: [-2.1894, -79.8893],
      zoom: 14,
      zoomControl: false,
      attributionControl: false,
    });
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", { maxZoom: 19 }).addTo(map);
    mapRef.current = map;
    initialized.current = true;
    setMapReady(true);

    // Locate user
    map.locate({ setView: false, maxZoom: 16 });
    map.on("locationfound", (e: any) => {
      if (onLocate && !pickupCoords) {
        onLocate({ lat: e.latlng.lat, lng: e.latlng.lng });
      }
    });

    return () => { map.remove(); initialized.current = false; };
  }, []);

  // Pickup marker
  useEffect(() => {
    if (!mapRef.current || !pickupCoords) return;
    const L = window.L;
    if (pickupMarkerRef.current) mapRef.current.removeLayer(pickupMarkerRef.current);
    const icon = L.divIcon({ html: '<div style="width:24px;height:24px;background:#000;border:3px solid #fff;border-radius:50%;box-shadow:0 2px 8px rgba(0,0,0,.3)"></div>', className: "", iconSize: [24, 24], iconAnchor: [12, 12] });
    pickupMarkerRef.current = L.marker([pickupCoords.lat, pickupCoords.lng], { icon }).addTo(mapRef.current);
  }, [pickupCoords]);

  // Dest marker
  useEffect(() => {
    if (!mapRef.current || !destCoords) return;
    const L = window.L;
    if (destMarkerRef.current) mapRef.current.removeLayer(destMarkerRef.current);
    const icon = L.divIcon({ html: '<div style="width:24px;height:24px;background:#276ef1;border:3px solid #fff;border-radius:50%;box-shadow:0 2px 8px rgba(0,0,0,.3)"></div>', className: "", iconSize: [24, 24], iconAnchor: [12, 12] });
    destMarkerRef.current = L.marker([destCoords.lat, destCoords.lng], { icon }).addTo(mapRef.current);
  }, [destCoords]);

  // Driver marker
  useEffect(() => {
    if (!mapRef.current || !driverCoords) return;
    const L = window.L;
    if (driverMarkerRef.current) mapRef.current.removeLayer(driverMarkerRef.current);
    const icon = L.divIcon({ html: '<div style="width:36px;height:36px;background:#000;border:3px solid #fff;border-radius:50%;box-shadow:0 2px 12px rgba(0,0,0,.4);display:flex;align-items:center;justify-content:center;color:#fff;font-size:16px">🚗</div>', className: "", iconSize: [36, 36], iconAnchor: [18, 18] });
    driverMarkerRef.current = L.marker([driverCoords.lat, driverCoords.lng], { icon, zIndexOffset: 100 }).addTo(mapRef.current);
  }, [driverCoords]);

  // Polyline
  useEffect(() => {
    if (!mapRef.current || !polyline) return;
    const L = window.L;
    if (routeRef.current) mapRef.current.removeLayer(routeRef.current);
    try {
      const coords = decodePolyline(polyline);
      if (coords.length > 0) {
        routeRef.current = L.polyline(coords, { color: "#276ef1", weight: 5, opacity: 0.9, lineCap: "round" }).addTo(mapRef.current);
        mapRef.current.fitBounds(routeRef.current.getBounds().pad(0.2));
      }
    } catch (e) {
      if (pickupCoords && destCoords) {
        routeRef.current = L.polyline([[pickupCoords.lat, pickupCoords.lng], [destCoords.lat, destCoords.lng]], { color: "#276ef1", weight: 4, opacity: 0.5, dashArray: "10, 10" }).addTo(mapRef.current);
      }
    }
  }, [polyline, pickupCoords, destCoords]);

  return (
    <>
      <div ref={containerRef} style={{ width: "100%", height: "100%", position: "absolute", top: 0, left: 0 }} />
      {mapReady && (
        <button className="locate-btn" onClick={() => { mapRef.current?.locate({ setView: true, maxZoom: 16 }); }}
          aria-label="My location">📍</button>
      )}
    </>
  );
}

function decodePolyline(str: string, precision = 5): [number, number][] {
  let index = 0, lat = 0, lng = 0;
  const coords: [number, number][] = [];
  while (index < str.length) {
    let shift = 0, result = 0, byte: number;
    do { byte = str.charCodeAt(index++) - 63; result |= (byte & 0x1f) << shift; shift += 5; } while (byte >= 0x20);
    lat += (result & 1) ? ~(result >> 1) : (result >> 1);
    shift = 0; result = 0;
    do { byte = str.charCodeAt(index++) - 63; result |= (byte & 0x1f) << shift; shift += 5; } while (byte >= 0x20);
    lng += (result & 1) ? ~(result >> 1) : (result >> 1);
    coords.push([lat / Math.pow(10, precision), lng / Math.pow(10, precision)]);
  }
  return coords;
}
