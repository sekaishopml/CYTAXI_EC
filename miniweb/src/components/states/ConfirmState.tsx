"use client";
import { Place, FareBreakdown } from "@/types";

const STITCH_GREEN = "#006c49";

interface ConfirmStateProps {
  pickup: Place; dest: Place;
  route: { distance_km: number; eta_minutes: number } | null;
  fare: FareBreakdown | null;
  onConfirm: () => void; onBack: () => void;
  loading: boolean;
}

export function ConfirmState({ pickup, dest, route, fare, onConfirm, onBack, loading }: ConfirmStateProps) {
  return (
    <div style={{ padding: "12px 20px 24px" }}>
      <div className="sheet-handle" />
      <h2 className="text-headline-mobile" style={{ marginBottom: 16 }}>Confirm your trip</h2>

      {/* Route info */}
      <div style={{ background: "#f6f6f6", borderRadius: 16, padding: 16, marginBottom: 16 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 14, marginBottom: 8 }}>
          <span className="material-symbols-outlined" style={{ fontSize: 20, color: STITCH_GREEN }}>location_on</span>
          <div><p className="text-body-md" style={{ fontWeight: 600 }}>{pickup.name}</p><p className="text-label-sm text-muted">{pickup.address}</p></div>
        </div>
        <div style={{ width: 2, height: 24, background: "#d9dadc", marginLeft: 10 }} />
        <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
          <span className="material-symbols-outlined" style={{ fontSize: 20, color: "#276ef1" }}>trip</span>
          <div><p className="text-body-md" style={{ fontWeight: 600 }}>{dest.name}</p><p className="text-label-sm text-muted">{dest.address}</p></div>
        </div>
      </div>

      {/* Metrics */}
      {route && (
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 16 }}>
          <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 14, textAlign: "center" }}>
            <p className="text-title" style={{ fontSize: 22, color: STITCH_GREEN }}>{route.distance_km.toFixed(1)} km</p>
            <p className="text-label-sm text-muted">Distance</p>
          </div>
          <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 14, textAlign: "center" }}>
            <p className="text-title" style={{ fontSize: 22 }}>{route.eta_minutes} min</p>
            <p className="text-label-sm text-muted">ETA</p>
          </div>
        </div>
      )}

      {/* Fare breakdown */}
      {fare && (
        <div style={{ background: "#f6f6f6", borderRadius: 16, padding: 16, marginBottom: 16 }}>
          <h3 className="text-label-md" style={{ marginBottom: 10 }}>Fare breakdown</h3>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Base fare</span><span>${fare.base.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Distance ({fare.distance_km.toFixed(1)} km)</span><span>${fare.distance.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Time ({fare.eta_minutes} min)</span><span>${fare.time.toFixed(2)}</span></div>
          <div className="divider" />
          <div style={{ display: "flex", justifyContent: "space-between" }} className="text-title"><span>Total</span><span style={{ color: STITCH_GREEN }}>${fare.total.toFixed(2)}</span></div>
        </div>
      )}

      <div style={{ display: "flex", gap: 12 }}>
        <button onClick={onBack} style={{ flex: 1, padding: 14, borderRadius: 14, fontSize: 16, fontWeight: 600, border: "none", cursor: "pointer", fontFamily: "Inter", background: "#f6f6f6", color: "#191c1e" }}>Back</button>
        <button onClick={onConfirm} disabled={loading} style={{ flex: 1, padding: 14, borderRadius: 14, fontSize: 16, fontWeight: 600, border: "none", cursor: "pointer", fontFamily: "Inter", transition: "opacity 0.2s",
          background: STITCH_GREEN, color: "#fff", opacity: loading ? 0.4 : 1 }}>
          {loading ? "Processing..." : "Request Trip"}
        </button>
      </div>
    </div>
  );
}
