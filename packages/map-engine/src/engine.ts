import type { LatLng, RouteData } from "./types";
import {
  createDriverMarker, updateDriverMarker,
  createPickupPulse, createDestinationMarker,
} from "./markers";
import {
  drawRoute, drawAnimatedRoute, updateRoutePath,
  fitMapToRoute, clearOverlays,
} from "./routes";

interface MapEngineConfig {
  map: google.maps.Map;
  onCenterChange?: (c: LatLng) => void;
  onMapClick?: (c: LatLng) => void;
}

export class MapEngine {
  private map: google.maps.Map;
  private originMarker: google.maps.Marker | null = null;
  private destMarker: google.maps.Marker | null = null;
  private originPulse: google.maps.Circle | null = null;
  private polyline: google.maps.Polyline | null = null;
  private driverMarker: google.maps.Marker | null = null;
  private listeners: google.maps.MapsEventListener[] = [];
  private _onCenterChange: ((c: { lat: number; lng: number }) => void) | undefined;
  private _onMapClick: ((c: { lat: number; lng: number }) => void) | undefined;
  private _ready = false;

  constructor(config: MapEngineConfig) {
    this.map = config.map;

    this._onCenterChange = config.onCenterChange;
    this._onMapClick = config.onMapClick;

    this.listeners.push(
      this.map.addListener("center_changed", () => {
        if (!this._ready || !this._onCenterChange) return;
        const c = this.map.getCenter();
        if (c) this._onCenterChange({ lat: c.lat(), lng: c.lng() });
      }),
    );

    this.listeners.push(
      this.map.addListener("click", (e: google.maps.MapMouseEvent) => {
        if (!e.latLng || !this._ready || !this._onMapClick) return;
        this._onMapClick({ lat: e.latLng.lat(), lng: e.latLng.lng() });
      }),
    );

    this._ready = true;
  }

  get isReady() { return this._ready; }

  setCenter(lat: number, lng: number, zoom?: number) {
    this.map.setCenter({ lat, lng });
    if (zoom !== undefined) this.map.setZoom(zoom);
  }

  panTo(lat: number, lng: number) {
    this.map.panTo({ lat, lng });
  }

  setZoom(zoom: number) {
    this.map.setZoom(zoom);
  }

  setInteractive(interactive: boolean) {
    this.map.setOptions({
      gestureHandling: interactive ? "greedy" : "none",
    });
  }

  drawOrigin(position: LatLng, animate = true) {
    this.clearOrigin();
    this.originPulse = createPickupPulse(this.map, position);
    this.originMarker = new google.maps.Marker({
      position,
      map: this.map,
      icon: {
        url: `data:image/svg+xml;base64,${btoa(
          `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 28 28">
            <circle cx="14" cy="14" r="12" fill="#3b82f6" stroke="white" stroke-width="2.5"/>
            <circle cx="14" cy="14" r="4" fill="white" opacity="0.9"/>
          </svg>`
        )}`,
        scaledSize: new google.maps.Size(28, 28),
        anchor: new google.maps.Point(14, 14),
      },
      zIndex: 60,
      animation: animate ? google.maps.Animation.DROP : undefined,
    });
  }

  drawDestination(position: LatLng, animate = true) {
    this.clearDestination();
    this.destMarker = createDestinationMarker(this.map, position);
  }

  drawRoute(route: RouteData, animated = true) {
    this.clearRoute();
    if (animated) {
      this.polyline = drawAnimatedRoute(this.map, route, "#3b82f6", 1200);
    } else {
      this.polyline = drawRoute(this.map, route, "#3b82f6");
    }
  }

  updateRoute(polyline: string) {
    if (this.polyline) {
      updateRoutePath(this.polyline, polyline);
    }
  }

  fitToMarkers(padding?: { top?: number; bottom?: number; left?: number; right?: number }) {
    const bounds = new google.maps.LatLngBounds();
    let hasPoint = false;

    if (this.originMarker) {
      bounds.extend(this.originMarker.getPosition()!);
      hasPoint = true;
    }
    if (this.destMarker) {
      bounds.extend(this.destMarker.getPosition()!);
      hasPoint = true;
    }

    if (hasPoint) {
      this.map.fitBounds(bounds, {
        top: padding?.top ?? 120,
        bottom: padding?.bottom ?? 320,
        left: padding?.left ?? 40,
        right: padding?.right ?? 40,
      });
    }
  }

  showDriver(data: { lat: number; lng: number; heading?: number; status?: "idle" | "searching" | "arriving" | "waiting" | "driving"; name?: string }) {
    if (this.driverMarker) {
      updateDriverMarker(this.driverMarker, data, true);
    } else {
      this.driverMarker = createDriverMarker(this.map, {
        id: "driver",
        lat: data.lat,
        lng: data.lng,
        heading: data.heading,
        name: data.name,
        status: data.status || "arriving",
      });
    }
  }

  clearOrigin() {
    if (this.originPulse) { this.originPulse.setMap(null); this.originPulse = null; }
    if (this.originMarker) { this.originMarker.setMap(null); this.originMarker = null; }
  }

  clearDestination() {
    if (this.destMarker) { this.destMarker.setMap(null); this.destMarker = null; }
  }

  clearRoute() {
    if (this.polyline) { this.polyline.setMap(null); this.polyline = null; }
  }

  clearDriver() {
    if (this.driverMarker) { this.driverMarker.setMap(null); this.driverMarker = null; }
  }

  clearAll() {
    this.clearOrigin();
    this.clearDestination();
    this.clearRoute();
    this.clearDriver();
  }

  destroy() {
    this.clearAll();
    this.listeners.forEach(l => l.remove());
    this.listeners = [];
    this._ready = false;
  }
}
