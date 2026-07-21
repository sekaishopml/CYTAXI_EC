"use client";
import { useEffect, useRef, useCallback, memo } from "react";
import "leaflet/dist/leaflet.css";
import L from "leaflet";
import { MapEngine } from "@cytaxi/map-engine";
import type { RoutePayload } from "@/types";
import { RideState } from "@cytaxi/ride-machine";

interface MapControllerProps {
  state: RideState;
  interactive: boolean;
  showPin: boolean;
  pickupCoords: { lat: number; lng: number } | null;
  destCoords: { lat: number; lng: number } | null;
  route: RoutePayload | null;
  driverLat?: number;
  driverLng?: number;
  driverHeading?: number;
  onCenterChange: (c: { lat: number; lng: number }) => void;
  onMapClick: (c: { lat: number; lng: number }) => void;
  onMapReady: () => void;
}

const DEFAULT_CENTER = { lat: -2.1894, lng: -79.8893 };

function decodePolyline(encoded: string): [number, number][] {
  if (!encoded) return [];
  const points: [number, number][] = [];
  let index = 0, lat = 0, lng = 0;
  while (index < encoded.length) {
    let shift = 0, result = 0, byte: number;
    do {
      byte = encoded.charCodeAt(index++) - 63;
      result |= (byte & 0x1f) << shift;
      shift += 5;
    } while (byte >= 0x20);
    const deltaLat = result & 1 ? ~(result >> 1) : result >> 1;
    lat += deltaLat;
    shift = 0; result = 0;
    do {
      byte = encoded.charCodeAt(index++) - 63;
      result |= (byte & 0x1f) << shift;
      shift += 5;
    } while (byte >= 0x20);
    const deltaLng = result & 1 ? ~(result >> 1) : result >> 1;
    lng += deltaLng;
    points.push([lat / 1e5, lng / 1e5]);
  }
  return points;
}

function pickupHtml(color: string, label: string): string {
  return `<div style="position:relative;width:32px;height:32px">
    <div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);width:50px;height:50px;border-radius:50%;background:${color}22;border:2px solid ${color}44;animation:lfPulse 1.8s ease-in-out infinite"></div>
    <div style="position:absolute;top:0;left:0;width:32px;height:32px;border-radius:50%;background:${color};border:3px solid #fff;box-shadow:0 2px 8px rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center;color:#fff;font-size:14px;font-weight:700;font-family:Inter,sans-serif">${label}</div>
  </div>`;
}

function destinationHtml(): string {
  return `<div style="width:32px;height:32px;border-radius:6px;background:#3b82f6;border:3px solid #fff;box-shadow:0 2px 8px rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center">
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 8l4 4-4 4M8 12h8"/>
    </svg>
  </div>`;
}

function driverHtml(): string {
  return `<div style="width:34px;height:34px;border-radius:50% 50% 50% 0;background:#3b82f6;border:2.5px solid #fff;box-shadow:0 2px 6px rgba(0,0,0,.35);transform:rotate(-45deg);display:flex;align-items:center;justify-content:center"><div style="width:10px;height:10px;border-radius:50%;background:#fff"></div></div>`;
}

function MapController_({
  state, interactive, showPin,
  pickupCoords, destCoords, route,
  driverLat, driverLng, driverHeading,
  onCenterChange, onMapClick, onMapReady,
}: MapControllerProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const engineRef = useRef<MapEngine | null>(null);
  const leafletRef = useRef<{
    map: L.Map;
    pickup?: L.Marker;
    pickupCircle?: L.Circle;
    dest?: L.Marker;
    route?: L.Polyline;
    driver?: L.Marker;
  } | null>(null);
  const useLeaflet = useRef(false);
  const initialized = useRef(false);

  const onCenterRef = useRef(onCenterChange);
  const onClickRef = useRef(onMapClick);
  const onReadyRef = useRef(onMapReady);
  const animRouteRef = useRef<number | null>(null);
  onCenterRef.current = onCenterChange;
  onClickRef.current = onMapClick;
  onReadyRef.current = onMapReady;

  const initLeaflet = useCallback(() => {
    if (!containerRef.current || leafletRef.current) return;
    const map = L.map(containerRef.current, {
      center: [DEFAULT_CENTER.lat, DEFAULT_CENTER.lng],
      zoom: 15,
      zoomControl: false,
      attributionControl: false,
    });
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 19,
    }).addTo(map);
    map.on("moveend", () => {
      const c = map.getCenter();
      onCenterRef.current?.({ lat: c.lat, lng: c.lng });
    });
    map.on("click", (e: L.LeafletMouseEvent) => {
      onClickRef.current?.({ lat: e.latlng.lat, lng: e.latlng.lng });
    });
    leafletRef.current = { map };
    useLeaflet.current = true;
    (window as any).__cymap = {
      setCenter: (c: { lat: number; lng: number }, zoom?: number) => {
        if (zoom) map.setView([c.lat, c.lng], zoom);
        else map.setView([c.lat, c.lng]);
      },
      setZoom: (z: number) => map.setZoom(z),
      panTo: (c: { lat: number; lng: number }) => map.panTo([c.lat, c.lng]),
      getCenter: () => { const c = map.getCenter(); return { lat: c.lat, lng: c.lng }; },
    };
    // Inject pulse animation for Leaflet markers
    if (typeof document !== "undefined" && !document.getElementById("lf-pulse-style")) {
      const style = document.createElement("style");
      style.id = "lf-pulse-style";
      style.textContent = `@keyframes lfPulse{0%,100%{transform:translate(-50%,-50%) scale(0.85);opacity:0.6}50%{transform:translate(-50%,-50%) scale(1.15);opacity:0.2}}`;
      document.head.appendChild(style);
    }
    onReadyRef.current?.();
  }, []);

  useEffect(() => {
    if (typeof window === "undefined" || initialized.current) return;

    const tryInitGoogle = () => {
      if (!window.google?.maps) return false;
      const map = new window.google.maps.Map(containerRef.current!, {
        center: DEFAULT_CENTER,
        zoom: 15,
        disableDefaultUI: true,
        zoomControl: false,
        mapTypeId: "roadmap",
        mapTypeControl: false,
        streetViewControl: false,
        fullscreenControl: false,
        rotateControl: false,
        gestureHandling: interactive ? "greedy" : "none",
      });
      const engine = new MapEngine({
        map,
        onCenterChange: (c) => onCenterRef.current?.(c),
        onMapClick: (c) => onClickRef.current?.(c),
      });
      engineRef.current = engine;
      (window as any).__cymap = map;
      initialized.current = true;
      onReadyRef.current?.();
      return true;
    };

    if (tryInitGoogle()) return;
    let attempts = 0;
    const interval = setInterval(() => {
      attempts++;
      if (tryInitGoogle() || attempts >= 6) {
        clearInterval(interval);
        if (!initialized.current && attempts >= 6) {
          initLeaflet();
          initialized.current = true;
        }
      }
    }, 300);

    return () => { clearInterval(interval); };
  }, [interactive, initLeaflet]);

  useEffect(() => {
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (g) g.setInteractive(interactive);
    if (lf) { if (interactive) lf.map.dragging.enable(); else lf.map.dragging.disable(); }
  }, [interactive]);

  // ── ORIGEN (pickup A) ──
  useEffect(() => {
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (!g && !lf) return;

    // No dibujar marker permanente durante pickup_select (solo pin overlay)
    if (state === "pickup_select") {
      if (g) g.clearOrigin();
      if (lf) {
        if (lf.pickup) { lf.pickup.remove(); lf.pickup = undefined; }
        if (lf.pickupCircle) { lf.pickupCircle.remove(); lf.pickupCircle = undefined; }
      }
      return;
    }

    if (g) {
      if (pickupCoords) {
        g.drawOrigin(pickupCoords, true);
      } else {
        g.clearOrigin();
      }
    }
    if (lf) {
      if (pickupCoords) {
        if (lf.pickup) {
          lf.pickup.setLatLng([pickupCoords.lat, pickupCoords.lng]);
          lf.pickupCircle?.setLatLng([pickupCoords.lat, pickupCoords.lng]);
        } else {
          lf.pickupCircle = L.circle([pickupCoords.lat, pickupCoords.lng], {
            radius: 25,
            color: "#3b82f6",
            fillColor: "#3b82f6",
            fillOpacity: 0.08,
            weight: 1.5,
            opacity: 0.25,
          }).addTo(lf.map);
          lf.pickup = L.marker([pickupCoords.lat, pickupCoords.lng], {
            icon: L.divIcon({ className: "", html: pickupHtml("#3b82f6", "A"), iconSize: [32, 32], iconAnchor: [16, 16] }),
          }).addTo(lf.map);
        }
      }
      if (!pickupCoords) {
        if (lf.pickup) { lf.pickup.remove(); lf.pickup = undefined; }
        if (lf.pickupCircle) { lf.pickupCircle.remove(); lf.pickupCircle = undefined; }
      }
    }
  }, [pickupCoords, state]);

  // ── DESTINO (B) ──
  useEffect(() => {
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (!g && !lf) return;

    // No dibujar marker permanente durante selección de destino (solo pin overlay B)
    if (state === "input" && !destCoords) {
      if (g) g.clearDestination();
      if (lf && lf.dest) { lf.dest.remove(); lf.dest = undefined; }
      return;
    }

    if (g) {
      if (destCoords) {
        g.drawDestination(destCoords, true);
      } else {
        g.clearDestination();
      }
    }
    if (lf) {
      if (destCoords) {
        if (lf.dest) {
          lf.dest.setLatLng([destCoords.lat, destCoords.lng]);
        } else {
          lf.dest = L.marker([destCoords.lat, destCoords.lng], {
            icon: L.divIcon({ className: "", html: destinationHtml(), iconSize: [32, 32], iconAnchor: [16, 16] }),
          }).addTo(lf.map);
        }
      }
      if (!destCoords && lf.dest) { lf.dest.remove(); lf.dest = undefined; }
    }
  }, [destCoords, state]);

  // ── FIT TO MARKERS ──
  useEffect(() => {
    if (state === "pickup_select") return;
    if (!pickupCoords || !destCoords) return;
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (g) {
      g.fitToMarkers({ top: 120, bottom: 320 });
    }
    if (lf) {
      const b = L.latLngBounds([pickupCoords.lat, pickupCoords.lng], [destCoords.lat, destCoords.lng]);
      lf.map.fitBounds(b, { paddingTopLeft: [40, 120], paddingBottomRight: [40, 320] });
    }
  }, [pickupCoords, destCoords, route, state]);

  // ── ROUTE ──
  useEffect(() => {
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (g) {
      if (route && route.polyline && pickupCoords && destCoords) {
        g.drawRoute({ polyline: route.polyline, distance_km: route.distance_km, duration_seconds: route.duration_seconds }, true);
        g.fitToMarkers({ top: 120, bottom: 320 });
      }
      if (!route) { g.clearRoute(); }
    }
    if (lf) {
      if (route && route.polyline && pickupCoords && destCoords) {
        lf.route?.remove();
        if (animRouteRef.current) cancelAnimationFrame(animRouteRef.current);

        const pts = decodePolyline(route.polyline);
        if (pts.length === 0) return;

        const DRAW_DURATION = 900;
        const drawnLine = L.polyline([], {
          color: "#3b82f6", weight: 5, opacity: 0.9,
          lineCap: "round", lineJoin: "round",
        }).addTo(lf.map);

        const trailLine = L.polyline(pts, {
          color: "#3b82f6", weight: 3, opacity: 0.12,
          dashArray: "6 8", lineCap: "round",
        }).addTo(lf.map);

        lf.route = drawnLine;

        const startTime = performance.now();
        const drawFrame = (now: number) => {
          const elapsed = now - startTime;
          const progress = Math.min(elapsed / DRAW_DURATION, 1);
          const eased = 1 - Math.pow(1 - progress, 3);
          const count = Math.max(1, Math.round(eased * pts.length));
          drawnLine.setLatLngs(pts.slice(0, count));
          if (progress < 1) {
            animRouteRef.current = requestAnimationFrame(drawFrame);
          } else {
            trailLine.remove();
          }
        };
        animRouteRef.current = requestAnimationFrame(drawFrame);

        if (pickupCoords && destCoords) {
          const allPts = [pickupCoords, ...pts.map(p => ({ lat: p[0], lng: p[1] })), destCoords];
          const b = L.latLngBounds(allPts.map(p => [p.lat, p.lng]));
          lf.map.fitBounds(b, { paddingTopLeft: [40, 120], paddingBottomRight: [40, 320] });
        }
      }
      if (!route && lf.route) {
        if (animRouteRef.current) cancelAnimationFrame(animRouteRef.current);
        lf.route.remove(); lf.route = undefined;
      }
    }
  }, [route, pickupCoords, destCoords]);

  // ── DRIVER ──
  useEffect(() => {
    const g = engineRef.current;
    const lf = leafletRef.current;
    if (g && driverLat !== undefined && driverLng !== undefined) {
      const status: "arriving" | "waiting" | "driving" =
        state === "arriving" ? "arriving" : state === "arrived" ? "waiting" : "driving";
      g.showDriver({ lat: driverLat, lng: driverLng, heading: driverHeading || 0, status });
    }
    if (lf && driverLat !== undefined && driverLng !== undefined) {
      const pos: [number, number] = [driverLat, driverLng];
      if (lf.driver) lf.driver.setLatLng(pos);
      else lf.driver = L.marker(pos, {
        icon: L.divIcon({ className: "", html: driverHtml(), iconSize: [34, 34], iconAnchor: [17, 17] }),
        zIndexOffset: 1000,
      }).addTo(lf.map);
    }
  }, [driverLat, driverLng, driverHeading, state]);

  return (
    <div
      ref={containerRef}
      style={{ width: "100%", height: "100%", background: "#1c1c1e" }}
    />
  );
}

export const MapController = memo(MapController_, (prev, next) => {
  return (
    prev.state === next.state &&
    prev.interactive === next.interactive &&
    prev.showPin === next.showPin &&
    prev.pickupCoords?.lat === next.pickupCoords?.lat &&
    prev.pickupCoords?.lng === next.pickupCoords?.lng &&
    prev.destCoords?.lat === next.destCoords?.lat &&
    prev.destCoords?.lng === next.destCoords?.lng &&
    prev.route === next.route &&
    prev.driverLat === next.driverLat &&
    prev.driverLng === next.driverLng
  );
});
