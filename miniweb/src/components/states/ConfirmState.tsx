"use client";
import { Place, FareBreakdown } from "@/types";

const G = "#006c49";

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
      <h2 style={{ fontSize: 20, fontWeight: 600, marginBottom: 14 }}>Confirma tu viaje</h2>

      <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 14, marginBottom: 14 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 6 }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: G }}>location_on</span>
          <div><p style={{ fontSize: 14, fontWeight: 600, margin: 0 }}>{pickup.name}</p><p style={{ fontSize: 12, color: "#3c4a42", margin: 0 }}>{pickup.address}</p></div>
        </div>
        <div style={{ width: 2, height: 20, background: "#d9dadc", marginLeft: 9 }} />
        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: "#276ef1" }}>trip</span>
          <div><p style={{ fontSize: 14, fontWeight: 600, margin: 0 }}>{dest.name}</p><p style={{ fontSize: 12, color: "#3c4a42", margin: 0 }}>{dest.address}</p></div>
        </div>
      </div>

      {route && (
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 14 }}>
          <div style={{ background: "#f6f6f6", borderRadius: 12, padding: 12, textAlign: "center" }}>
            <p style={{ fontSize: 20, fontWeight: 700, color: G, margin: 0 }}>{route.distance_km.toFixed(1)} km</p>
            <p style={{ fontSize: 11, color: "#3c4a42", margin: "2px 0 0" }}>Distancia</p>
          </div>
          <div style={{ background: "#f6f6f6", borderRadius: 12, padding: 12, textAlign: "center" }}>
            <p style={{ fontSize: 20, fontWeight: 700, margin: 0 }}>{route.eta_minutes} min</p>
            <p style={{ fontSize: 11, color: "#3c4a42", margin: "2px 0 0" }}>Tiempo estimado</p>
          </div>
        </div>
      )}

      {fare && (
        <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 14, marginBottom: 14 }}>
          <p style={{ fontSize: 12, fontWeight: 600, marginBottom: 8, color: "#3c4a42" }}>Detalle del precio</p>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}><span style={{ fontSize: 13, color: "#3c4a42" }}>Tarifa base</span><span style={{ fontSize: 13 }}>${fare.base.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}><span style={{ fontSize: 13, color: "#3c4a42" }}>Distancia</span><span style={{ fontSize: 13 }}>${fare.distance.toFixed(2)}</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}><span style={{ fontSize: 13, color: "#3c4a42" }}>Tiempo</span><span style={{ fontSize: 13 }}>${fare.time.toFixed(2)}</span></div>
          <div style={{ height: 1, background: "#d9dadc", margin: "6px 0" }} />
          <div style={{ display: "flex", justifyContent: "space-between" }}><span style={{ fontSize: 16, fontWeight: 700 }}>Total</span><span style={{ fontSize: 16, fontWeight: 700, color: G }}>${fare.total.toFixed(2)}</span></div>
        </div>
      )}

      <div style={{ display: "flex", gap: 10 }}>
        <button onClick={onBack} style={{ flex: 1, padding: 14, borderRadius: 12, fontSize: 15, fontWeight: 600, border: "none", cursor: "pointer", fontFamily: "Inter", background: "#f6f6f6", color: "#191c1e" }}>Atrás</button>
        <button onClick={onConfirm} disabled={loading} style={{ flex: 1, padding: 14, borderRadius: 12, fontSize: 15, fontWeight: 600, border: "none", cursor: "pointer", fontFamily: "Inter", background: G, color: "#fff", opacity: loading ? 0.4 : 1 }}>
          {loading ? "Procesando..." : "Solicitar viaje"}
        </button>
      </div>
    </div>
  );
}
