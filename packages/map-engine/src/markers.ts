import type { LatLng, DriverMarkerData, MarkerAnimationConfig } from "./types";

let _markerIdCounter = 0;

export function createDriverMarker(
  map: google.maps.Map,
  data: DriverMarkerData,
  config?: MarkerAnimationConfig
): google.maps.Marker {
  const heading = data.heading ?? 0;
  const marker = new google.maps.Marker({
    position: { lat: data.lat, lng: data.lng },
    map,
    icon: driverMarkerIcon(data.status || "idle", heading),
    zIndex: 100,
    title: data.name,
    animation: google.maps.Animation.DROP,
    optimized: false,
  });

  const id = `drv_${++_markerIdCounter}`;
  (marker as any)._cyId = id;

  if (config?.pulse) {
    startPulse(map, marker, config.pulseColor || "#3b82f6", config.pulseDuration || 1.5);
  }

  return marker;
}

export function updateDriverMarker(
  marker: google.maps.Marker,
  data: Partial<DriverMarkerData>,
  animate = true
): void {
  if (data.lat !== undefined && data.lng !== undefined) {
    const pos = marker.getPosition()!;
    const newPos = new google.maps.LatLng(data.lat, data.lng);
    if (animate) {
      smoothMove(marker, newPos, 1000);
    } else {
      marker.setPosition(newPos);
    }
  }
  if (data.heading !== undefined) {
    marker.setIcon(driverMarkerIcon(data.status || "driving", data.heading));
  }
  if (data.status !== undefined) {
    marker.setIcon(driverMarkerIcon(data.status, data.heading ?? 0));
  }
}

export function createPickupPulse(
  map: google.maps.Map,
  position: LatLng
): google.maps.Circle {
  return new google.maps.Circle({
    map,
    center: position,
    radius: 40,
    fillColor: "#3b82f6",
    fillOpacity: 0.12,
    strokeColor: "#60a5fa",
    strokeOpacity: 0.35,
    strokeWeight: 1.5,
    zIndex: 50,
  });
}

export function createDestinationMarker(
  map: google.maps.Map,
  position: LatLng
): google.maps.Marker {
  return new google.maps.Marker({
    position,
    map,
    icon: destinationIcon(),
    zIndex: 50,
    optimized: false,
  });
}

function startPulse(map: google.maps.Map, marker: google.maps.Marker, color: string, duration: number): void {
  let scale = 1;
  let growing = true;

  const interval = setInterval(() => {
    if (!marker.getMap()) { clearInterval(interval); return; }

    if (growing) {
      scale += 0.02;
      if (scale >= 1.3) growing = false;
    } else {
      scale -= 0.02;
      if (scale <= 0.95) growing = false;
      if (scale <= 0.85) growing = true;
    }

    marker.setIcon(driverMarkerIcon("arriving", 0, scale));
  }, duration * 16);

  (marker as any)._pulseInterval = interval;
}

function smoothMove(marker: google.maps.Marker, target: google.maps.LatLng, durationMs: number): void {
  const start = marker.getPosition()!;
  const startTime = performance.now();

  function animate(time: number) {
    const elapsed = time - startTime;
    const progress = Math.min(elapsed / durationMs, 1);
    const eased = 1 - Math.pow(1 - progress, 3);

    const lat = start.lat() + (target.lat() - start.lat()) * eased;
    const lng = start.lng() + (target.lng() - start.lng()) * eased;

    marker.setPosition(new google.maps.LatLng(lat, lng));

    if (progress < 1) {
      requestAnimationFrame(animate);
    }
  }

  requestAnimationFrame(animate);
}

function driverMarkerIcon(
  status: string,
  heading: number,
  scale = 1
): google.maps.Icon {
  const colors: Record<string, string> = {
    idle: "#9ea5a0",
    searching: "#448aff",
    arriving: "#ffc107",
    waiting: "#60a5fa",
    driving: "#3b82f6",
  };

  const color = colors[status] || "#3b82f6";
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${24 * scale}" height="${24 * scale}" viewBox="0 0 24 24">
    <circle cx="12" cy="12" r="10" fill="${color}" stroke="white" stroke-width="2"/>
    <path d="M7 10l5-5 5 5" fill="none" stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    <path d="M12 5v14" stroke="white" stroke-width="1.5" stroke-linecap="round"/>
  </svg>`;

  const scaledSize = 24 * scale;
  return {
    url: `data:image/svg+xml;base64,${btoa(svg)}`,
    scaledSize: new google.maps.Size(scaledSize, scaledSize),
    anchor: new google.maps.Point(scaledSize / 2, scaledSize / 2),
    labelOrigin: new google.maps.Point(scaledSize / 2, scaledSize / 2),
  };
}

function destinationIcon(): google.maps.Icon {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
    <rect x="2" y="2" width="20" height="20" rx="4" fill="#3b82f6" stroke="white" stroke-width="2"/>
    <path d="M12 8l4 4-4 4M8 12h8" stroke="white" stroke-width="1.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>`;

  return {
    url: `data:image/svg+xml;base64,${btoa(svg)}`,
    scaledSize: new google.maps.Size(24, 24),
    anchor: new google.maps.Point(12, 12),
  };
}
