const ALLOWED_TYPES = ["street_address","route","premise","subpremise","establishment","park","airport","university","school","parking","bus_station"];
const BLOCKED_CATEGORIES = ["water","river","lake","sea","ocean","reservoir","bay","stream","canal","pond","wetland","forest","heath","moor","grassland","farmland","meadow","scrub","beach","dune","quarry","landfill","construction","brownfield","greenfield","industrial","military"];

async function nominatimReverse(lat: number, lng: number): Promise<{valid: boolean; address: string | null}> {
  try {
    const res = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&addressdetails=1`,
      { headers: { "User-Agent": "CYTAXI/1.0" }, signal: AbortSignal.timeout(3000) }
    );
    if (!res.ok) return { valid: true, address: null };
    const data = await res.json();
    if (!data || data.error) return { valid: false, address: null };
    const addr = data.address || {};
    const category = data.category || "";
    const type = data.type || "";
    if (BLOCKED_CATEGORIES.some(c => category.includes(c) || type.includes(c))) {
      return { valid: false, address: null };
    }
    if (addr?.road || addr?.pedestrian || addr?.footway || addr?.cycleway || addr?.residential || addr?.commercial || addr?.industrial || addr?.neighbourhood) {
      return { valid: true, address: data.display_name || null };
    }
    if (addr?.city || addr?.town || addr?.village || addr?.municipality) {
      if (addr?.road) return { valid: true, address: data.display_name || null };
      return { valid: false, address: null };
    }
    return { valid: false, address: null };
  } catch {
    return { valid: true, address: null };
  }
}

export async function validateGoogleCoords(lat: number, lng: number): Promise<{valid: boolean; address: string | null}> {
  const google = (globalThis as any).google;
  if (google?.maps?.Geocoder) {
    return new Promise(resolve => {
      const timer = setTimeout(() => resolve({ valid: true, address: null }), 3000);
      try {
        const geocoder = new google.maps.Geocoder();
        geocoder.geocode({ location: { lat, lng } }, (results: any, status: string) => {
          clearTimeout(timer);
          if (status !== "OK" || !results?.length) { resolve({ valid: false, address: null }); return; }
          const r = results[0];
          const types = r.types || [];
          const isGood = types.some((t: string) => ALLOWED_TYPES.includes(t));
          resolve({ valid: isGood, address: isGood ? r.formatted_address : null });
        });
      } catch { clearTimeout(timer); resolve({ valid: true, address: null }); }
    });
  }
  return nominatimReverse(lat, lng);
}
