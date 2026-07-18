"use client";
import { Place, FareBreakdown } from "@/types";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";

interface ConfirmStateProps {
  pickup: Place; dest: Place;
  route: { distance_km: number; eta_minutes: number } | null;
  fare: FareBreakdown | null;
  onConfirm: () => void; onBack: () => void;
  loading: boolean;
  paymentMethod?: "cash" | "card";
  onPaymentChange?: (m: "cash" | "card") => void;
}

export function ConfirmState({ pickup, dest, route, fare, onConfirm, onBack, loading, paymentMethod = "cash", onPaymentChange }: ConfirmStateProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "16px 20px 14px", gap: 14 }}>
      <p style={{ fontSize: 19, fontWeight: 700, color: T1, margin: 0, letterSpacing: "-0.02em" }}>Confirma tu viaje</p>

      {/* Route card */}
      <div style={{
        background: "rgba(255,255,255,0.88)", backdropFilter: "blur(20px)", borderRadius: 14,
        boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 6px 20px rgba(0,0,0,0.04)",
        border: "1px solid rgba(0,0,0,0.04)", padding: 14,
      }}>
        <div style={{ display: "flex", alignItems: "flex-start", gap: 12 }}>
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 4, marginTop: 2 }}>
            <div style={{ width: 10, height: 10, borderRadius: "50%", background: G, flexShrink: 0 }} />
            <div style={{ width: 1, height: 24, background: "rgba(0,0,0,0.1)", marginLeft: 4.5 }} />
            <div style={{ width: 10, height: 10, borderRadius: "50%", background: "#448aff", flexShrink: 0 }} />
          </div>
          <div style={{ flex: 1 }}>
            <div style={{ marginBottom: 20 }}>
              <p style={{ fontSize: 10, fontWeight: 600, color: "#9a9a9a", margin: 0, letterSpacing: "0.04em", textTransform: "uppercase" }}>Recogida</p>
              <p style={{ fontSize: 14, fontWeight: 600, color: T1, margin: "2px 0 0" }}>{pickup.name}</p>
              <p style={{ fontSize: 12, color: "#8a8a8a", margin: "1px 0 0" }}>{pickup.address}</p>
            </div>
            <div>
              <p style={{ fontSize: 10, fontWeight: 600, color: "#9a9a9a", margin: 0, letterSpacing: "0.04em", textTransform: "uppercase" }}>Destino</p>
              <p style={{ fontSize: 14, fontWeight: 600, color: T1, margin: "2px 0 0" }}>{dest.name}</p>
              <p style={{ fontSize: 12, color: "#8a8a8a", margin: "1px 0 0" }}>{dest.address}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Distance & ETA */}
      {route && (
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10 }}>
          <div style={{ background: "rgba(255,255,255,0.7)", borderRadius: 12, padding: "14px 12px", textAlign: "center" }}>
            <p style={{ fontSize: 22, fontWeight: 700, color: G, margin: 0, letterSpacing: "-0.02em" }}>{route.distance_km.toFixed(1)}</p>
            <p style={{ fontSize: 11, color: "#8a8a8a", margin: "2px 0 0", textTransform: "uppercase", letterSpacing: "0.04em", fontWeight: 600 }}>Kilómetros</p>
          </div>
          <div style={{ background: "rgba(255,255,255,0.7)", borderRadius: 12, padding: "14px 12px", textAlign: "center" }}>
            <p style={{ fontSize: 22, fontWeight: 700, color: T1, margin: 0, letterSpacing: "-0.02em" }}>{route.eta_minutes}</p>
            <p style={{ fontSize: 11, color: "#8a8a8a", margin: "2px 0 0", textTransform: "uppercase", letterSpacing: "0.04em", fontWeight: 600 }}>Minutos</p>
          </div>
        </div>
      )}

      {/* Fare breakdown */}
      {fare && (
        <div style={{ background: "rgba(255,255,255,0.7)", borderRadius: 14, padding: 14 }}>
          <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", margin: "0 0 10px", textTransform: "uppercase", letterSpacing: "0.04em" }}>Detalle del viaje</p>
          {[
            ["Tarifa base", fare.base],
            ["Distancia", fare.distance],
            ["Tiempo", fare.time],
          ].map(([label, val]) => (
            <div key={label as string} style={{ display: "flex", justifyContent: "space-between", padding: "3px 0" }}>
              <span style={{ fontSize: 13, color: T2 }}>{label}</span>
              <span style={{ fontSize: 13, fontWeight: 500 }}>${(val as number).toFixed(2)}</span>
            </div>
          ))}
          <div style={{ height: 1, background: "rgba(0,0,0,0.08)", margin: "8px 0" }} />
          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <span style={{ fontSize: 17, fontWeight: 700 }}>Total</span>
            <span style={{ fontSize: 17, fontWeight: 700, color: G }}>${fare.total.toFixed(2)}</span>
          </div>

          {/* Payment toggle */}
          <div style={{ display: "flex", gap: 8, marginTop: 10 }}>
            {(["cash", "card"] as const).map(m => (
              <button key={m} onClick={() => onPaymentChange?.(m)}
                style={{ flex: 1, padding: "10px 0", borderRadius: 10, fontSize: 13, fontWeight: 600, fontFamily: "Inter", border: paymentMethod === m ? `1.5px solid ${G}` : "1px solid rgba(0,0,0,0.1)", cursor: "pointer", background: paymentMethod === m ? `${G}0A` : "transparent", color: paymentMethod === m ? G : "#8a8a8a", transition: "all 0.15s" }}>
                {m === "cash" ? "💵 Efectivo" : "💳 Tarjeta"}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Actions */}
      <div style={{ display: "flex", gap: 10 }}>
        <button onClick={onBack} style={{ flex: 1, height: 52, borderRadius: 14, fontSize: 15, fontWeight: 600, border: "1px solid rgba(0,0,0,0.08)", cursor: "pointer", fontFamily: "Inter", background: "transparent", color: T1, transition: "all 0.15s" }}>
          Atrás
        </button>
        <button onClick={onConfirm} disabled={loading} style={{ flex: 2, height: 52, background: G, color: "#fff", borderRadius: 14, fontSize: 15, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: "pointer", boxShadow: "0 4px 16px rgba(0,108,73,0.25)", transition: "all 0.15s", opacity: loading ? 0.5 : 1 }}>
          {loading ? "Procesando..." : "Solicitar viaje"}
        </button>
      </div>
    </div>
  );
}
