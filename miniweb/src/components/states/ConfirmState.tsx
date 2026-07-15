"use client";
import { Place, FareBreakdown } from "@/types";

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
      <div style={{ background: "var(--uk-input-bg)", borderRadius: 16, padding: 16, marginBottom: 16 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 14, marginBottom: 8 }}>
          <div className="dot dot-pickup" />
          <div><p className="text-body-md" style={{ fontWeight: 600 }}>{pickup.name}</p><p className="text-label-sm text-muted">{pickup.address}</p></div>
        </div>
        <div style={{ width: 2, height: 24, background: "var(--uk-outline)", marginLeft: 4 }} />
        <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
          <div className="dot dot-dest" />
          <div><p className="text-body-md" style={{ fontWeight: 600 }}>{dest.name}</p><p className="text-label-sm text-muted">{dest.address}</p></div>
        </div>
      </div>

      {/* Metrics */}
      {route && (
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 16 }}>
          <div style={{ background: "var(--uk-input-bg)", borderRadius: 14, padding: 14, textAlign: "center" }}>
            <p className="text-title" style={{ fontSize: 22 }}>{route.distance_km.toFixed(1)} km</p>
            <p className="text-label-sm text-muted">Distance</p>
          </div>
          <div style={{ background: "var(--uk-input-bg)", borderRadius: 14, padding: 14, textAlign: "center" }}>
            <p className="text-title" style={{ fontSize: 22 }}>{route.eta_minutes} min</p>
            <p className="text-label-sm text-muted">ETA</p>
          </div>
        </div>
      )}

      {/* Fare breakdown */}
      {fare && (
        <div style={{ background: "var(--uk-input-bg)", borderRadius: 16, padding: 16, marginBottom: 16 }}>
          <h3 className="text-label-md" style={{ marginBottom: 10 }}>Fare breakdown</h3>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Base fare</span><span>${fare.base.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Distance ({fare.distance_km.toFixed(1)} km)</span><span>${fare.distance.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Time ({fare.eta_minutes} min)</span><span>${fare.time.toFixed(2)}</span></div>
          <div className="divider" />
          <div style={{ display: "flex", justifyContent: "space-between" }} className="text-title"><span>Total</span><span>${fare.total.toFixed(2)}</span></div>
        </div>
      )}

      <div style={{ display: "flex", gap: 12 }}>
        <button onClick={onBack} className="btn btn-secondary" style={{ flex: 1 }}>Back</button>
        <button onClick={onConfirm} disabled={loading} className="btn btn-primary" style={{ flex: 1 }}>
          {loading ? "Processing..." : "Request Trip"}
        </button>
      </div>
    </div>
  );
}
