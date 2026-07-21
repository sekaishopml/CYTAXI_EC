import type { RouteData, LatLng } from "./types";

export function drawRoute(
  map: google.maps.Map,
  route: RouteData,
  color = "#3b82f6",
  opacity = 0.8
): google.maps.Polyline {
  const path = decodePolyline(route.polyline);

  const polyline = new google.maps.Polyline({
    path,
    map,
    strokeColor: color,
    strokeOpacity: opacity,
    strokeWeight: 4,
    zIndex: 10,
  });

  return polyline;
}

export function drawAnimatedRoute(
  map: google.maps.Map,
  route: RouteData,
  color = "#3b82f6",
  durationMs = 1500
): google.maps.Polyline {
  const path = decodePolyline(route.polyline);
  const totalPts = path.length;
  let currentPt = 0;

  const polyline = new google.maps.Polyline({
    path: [],
    map,
    strokeColor: color,
    strokeOpacity: 0.8,
    strokeWeight: 4,
    zIndex: 10,
  });

  const step = Math.max(1, Math.floor(totalPts / 60));

  function animate() {
    if (currentPt >= totalPts) return;
    const end = Math.min(currentPt + step, totalPts);
    polyline.setPath(path.slice(0, end));
    currentPt = end;
    requestAnimationFrame(() => animate());
  }

  setTimeout(animate, 100);
  return polyline;
}

export function updateRoutePath(
  polyline: google.maps.Polyline,
  newPolyline: string
): void {
  const path = decodePolyline(newPolyline);
  polyline.setPath(path);
}

export function fitMapToRoute(
  map: google.maps.Map,
  route: RouteData
): void {
  const path = decodePolyline(route.polyline);
  const bounds = new google.maps.LatLngBounds();
  path.forEach((pt) => bounds.extend(pt));
  map.fitBounds(bounds, { top: 80, right: 20, bottom: 300, left: 20 });
}

// ─── Decode Google polyline ──────────────────────────────
export function decodePolyline(encoded: string): google.maps.LatLng[] {
  if (!encoded) return [];
  const points: google.maps.LatLng[] = [];
  let index = 0;
  let lat = 0;
  let lng = 0;

  while (index < encoded.length) {
    let shift = 0;
    let result = 0;
    let byte: number;

    do {
      byte = encoded.charCodeAt(index++) - 63;
      result |= (byte & 0x1f) << shift;
      shift += 5;
    } while (byte >= 0x20);

    const deltaLat = result & 1 ? ~(result >> 1) : result >> 1;
    lat += deltaLat;

    shift = 0;
    result = 0;

    do {
      byte = encoded.charCodeAt(index++) - 63;
      result |= (byte & 0x1f) << shift;
      shift += 5;
    } while (byte >= 0x20);

    const deltaLng = result & 1 ? ~(result >> 1) : result >> 1;
    lng += deltaLng;

    points.push(new google.maps.LatLng(lat / 1e5, lng / 1e5));
  }

  return points;
}

export function clearOverlays(...overlays: (google.maps.MVCObject | null)[]): void {
  overlays.forEach((o) => {
    if (o && "setMap" in o) {
      (o as any).setMap(null);
    }
  });
}
