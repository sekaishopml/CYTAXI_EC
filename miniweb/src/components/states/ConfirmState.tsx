"use client";
import { useState } from "react";
import { Place, FareBreakdown } from "@/types";
import { colors, radius } from "@cytaxi/design-tokens";

interface ConfirmStateProps {
  pickup: Place;
  dest: Place;
  route: { distance_km: number; eta_minutes: number; polyline?: string } | null;
  fare: FareBreakdown | null;
  onConfirm: () => void;
  onBack: () => void;
  loading: boolean;
  paymentMethod?: "cash" | "card";
  onPaymentChange?: (m: "cash" | "card") => void;
  vehicleType: string;
  onVehicleChange?: (v: string) => void;
  note: string;
  onNoteChange?: (n: string) => void;
  coupon: string;
  onCouponChange?: (c: string) => void;
  scheduledAt: string | null;
  onScheduleChange?: (d: string | null) => void;
}

const VEHICLES = [
  { id: "standard", label: "Standard", icon: "🚗", price: 1 },
  { id: "xl", label: "XL", icon: "🚐", price: 1.5 },
  { id: "premium", label: "Premium", icon: "🚙", price: 2 },
];

const smallBtn: React.CSSProperties = {
  fontFamily: "'Inter', sans-serif", cursor: "pointer",
  fontSize: 11, fontWeight: 600, border: "1px solid rgba(0,0,0,0.07)",
  background: "transparent", color: colors.textSecondary,
  padding: "5px 10px", borderRadius: 6,
  transition: "background 0.15s ease, color 0.15s ease",
};

export function ConfirmState({
  pickup, dest, route, fare,
  onConfirm, onBack, loading,
  paymentMethod = "cash", onPaymentChange,
  vehicleType, onVehicleChange,
  note, onNoteChange,
  coupon, onCouponChange,
  scheduledAt, onScheduleChange,
}: ConfirmStateProps) {
  const [showNotes, setShowNotes] = useState(!!note);
  const [showCoupon, setShowCoupon] = useState(!!coupon);

  const multiplier = VEHICLES.find(v => v.id === vehicleType)?.price || 1;
  const adjustedTotal = fare ? (fare.total * multiplier) : 0;

  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "10px 16px 12px", gap: 8 }}>
      <p style={{ margin: 0, fontSize: 16, fontWeight: 600, fontFamily: "'Space Grotesk', sans-serif", color: colors.textPrimary, letterSpacing: "-0.01em" }}>
        Confirma tu viaje
      </p>

      <div style={{
        background: colors.surface.paperLight,
        borderRadius: radius.md, border: "1px solid rgba(0,0,0,0.06)", padding: 10,
      }}>
        <div style={{ display: "flex", alignItems: "flex-start", gap: 10 }}>
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 3, marginTop: 1 }}>
            <div style={{ width: 20, height: 20, borderRadius: "50%", background: colors.cobalt, color: "#fff", fontSize: 9, fontWeight: 700, display: "flex", alignItems: "center", justifyContent: "center", fontFamily: "'Inter', sans-serif" }}>A</div>
            <div style={{ width: 1.5, height: 18, background: "rgba(0,0,0,0.1)" }} />
            <div style={{ width: 20, height: 20, borderRadius: "50%", background: colors.cobaltLight, color: "#fff", fontSize: 9, fontWeight: 700, display: "flex", alignItems: "center", justifyContent: "center", fontFamily: "'Inter', sans-serif" }}>B</div>
          </div>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ marginBottom: 6 }}>
              <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>ORIGEN</p>
              <p style={{ fontSize: 13, fontWeight: 600, color: colors.textPrimary, margin: "1px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>{pickup.name}</p>
              <p style={{ fontSize: 11, color: colors.textMuted, margin: "1px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>{pickup.address}</p>
            </div>
            <div>
              <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>DESTINO</p>
              <p style={{ fontSize: 13, fontWeight: 600, color: colors.textPrimary, margin: "1px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>{dest.name}</p>
              <p style={{ fontSize: 11, color: colors.textMuted, margin: "1px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>{dest.address}</p>
            </div>
          </div>
        </div>
        <div style={{ display: "flex", gap: 6, marginTop: 8 }}>
          <button type="button" onClick={onBack} style={smallBtn} onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(0,0,0,0.03)"; }} onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}>Editar origen</button>
          <button type="button" onClick={onBack} style={smallBtn} onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(0,0,0,0.03)"; }} onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}>Editar destino</button>
        </div>
      </div>

      <div style={{ display: "flex", gap: 6 }}>
        <div style={{ flex: 1, background: colors.surface.paperLight, borderRadius: radius.md, padding: 8, border: "1px solid rgba(0,0,0,0.06)" }}>
          <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: "0 0 5px", letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>Vehículo</p>
          <div style={{ display: "flex", gap: 4 }}>
            {VEHICLES.map(v => {
              const active = vehicleType === v.id;
              return (
                <button type="button" key={v.id} onClick={() => onVehicleChange?.(v.id)}
                  aria-label={`Vehículo ${v.label}`} aria-pressed={active}
                  style={{
                    flex: 1, padding: "5px 0", borderRadius: radius.sm, fontSize: 10, fontWeight: 600,
                    border: active ? `1.5px solid ${colors.cobalt}` : "1px solid rgba(0,0,0,0.07)",
                    cursor: "pointer", fontFamily: "'Inter', sans-serif",
                    background: active ? colors.cobaltBg : "transparent",
                    color: active ? colors.cobalt : colors.textSecondary,
                    transition: "background 0.15s ease, color 0.15s ease, border-color 0.15s ease",
                  }}>
                  <div style={{ fontSize: 14 }}>{v.icon}</div>
                  <div>{v.label}</div>
                </button>
              );
            })}
          </div>
        </div>
        {route && (
          <div style={{ display: "flex", flexDirection: "column", gap: 4 }}>
            <div style={{ flex: 1, background: colors.surface.paperLight, borderRadius: radius.sm, padding: "8px 12px", textAlign: "center", border: "1px solid rgba(0,0,0,0.06)", display: "flex", flexDirection: "column", justifyContent: "center" }}>
              <p style={{ fontSize: 18, fontWeight: 700, color: colors.cobalt, margin: 0, fontVariantNumeric: "tabular-nums", fontFamily: "'Space Grotesk', sans-serif" }}>{route.distance_km.toFixed(1)}</p>
              <p style={{ fontSize: 8, color: colors.textMuted, margin: 0, textTransform: "uppercase", letterSpacing: "0.06em", fontWeight: 500, fontFamily: "'JetBrains Mono', monospace" }}>km</p>
            </div>
            <div style={{ flex: 1, background: colors.surface.paperLight, borderRadius: radius.sm, padding: "8px 12px", textAlign: "center", border: "1px solid rgba(0,0,0,0.06)", display: "flex", flexDirection: "column", justifyContent: "center" }}>
              <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: 0, fontVariantNumeric: "tabular-nums", fontFamily: "'Space Grotesk', sans-serif" }}>{route.eta_minutes}</p>
              <p style={{ fontSize: 8, color: colors.textMuted, margin: 0, textTransform: "uppercase", letterSpacing: "0.06em", fontWeight: 500, fontFamily: "'JetBrains Mono', monospace" }}>min</p>
            </div>
          </div>
        )}
      </div>

      {fare && (
        <div style={{ background: colors.surface.paperLight, borderRadius: radius.md, padding: 10, border: "1px solid rgba(0,0,0,0.06)" }}>
          <div style={{ display: "flex", gap: 6 }}>
            <div style={{ flex: 1 }}>
              {[
                ["Base", fare.base * multiplier],
                ["Dist.", fare.distance * multiplier],
                ["Tiempo", fare.time * multiplier],
              ].map(([label, val]) => (
                <div key={label as string} style={{ display: "flex", justifyContent: "space-between", padding: "2px 0" }}>
                  <span style={{ fontSize: 11, color: colors.textSecondary, fontFamily: "'Inter', sans-serif" }}>{label}</span>
                  <span style={{ fontSize: 11, fontWeight: 500, color: colors.textPrimary, fontFamily: "'Inter', sans-serif", fontVariantNumeric: "tabular-nums" }}>${(val as number).toFixed(2)}</span>
                </div>
              ))}
              <div style={{ height: 1, background: "rgba(0,0,0,0.05)", margin: "4px 0" }} />
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span style={{ fontSize: 14, fontWeight: 700, color: colors.textPrimary, fontFamily: "'Space Grotesk', sans-serif" }}>Total</span>
                <span style={{ fontSize: 14, fontWeight: 700, color: colors.cobalt, fontFamily: "'Space Grotesk', sans-serif", fontVariantNumeric: "tabular-nums" }}>${adjustedTotal.toFixed(2)}</span>
              </div>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: 4, justifyContent: "center" }}>
              {(["cash", "card"] as const).map(m => {
                const active = paymentMethod === m;
                return (
                  <button type="button" key={m} onClick={() => onPaymentChange?.(m)}
                    aria-label={`Pagar con ${m === "cash" ? "efectivo" : "tarjeta"}`} aria-pressed={active}
                    style={{
                      padding: "5px 10px", borderRadius: radius.sm, fontSize: 11, fontWeight: 600,
                      fontFamily: "'Inter', sans-serif", whiteSpace: "nowrap",
                      border: active ? `1.5px solid ${colors.cobalt}` : "1px solid rgba(0,0,0,0.07)",
                      cursor: "pointer",
                      background: active ? colors.cobaltBg : "transparent",
                      color: active ? colors.cobalt : colors.textMuted,
                      transition: "background 0.15s ease, color 0.15s ease, border-color 0.15s ease",
                    }}>
                    {m === "cash" ? "💵 Efectivo" : "💳 Tarjeta"}
                  </button>
                );
              })}
            </div>
          </div>
        </div>
      )}

      <div style={{ display: "flex", gap: 6, flexWrap: "wrap" }}>
        {!showCoupon ? (
          <button type="button" onClick={() => setShowCoupon(true)} style={{ ...smallBtn, border: "1px dashed rgba(0,0,0,0.12)", color: colors.info }}>+ Cupón</button>
        ) : (
          <div style={{ display: "flex", gap: 4, alignItems: "center", flex: 1 }}>
            <input aria-label="Código de cupón" value={coupon} onChange={e => onCouponChange?.(e.target.value)}
              placeholder="Código"
              style={{ flex: 1, padding: "5px 8px", borderRadius: radius.sm, fontSize: 11, fontFamily: "'Inter', sans-serif",
                border: "1px solid rgba(0,0,0,0.08)", outline: "none", background: colors.surface.paperLight, minWidth: 60 }} />
            <button type="button" onClick={() => { onCouponChange?.(""); setShowCoupon(false); }}
              style={{ ...smallBtn, background: colors.cobaltBg, color: colors.cobalt, border: "none", padding: "5px 8px" }}>
              OK
            </button>
          </div>
        )}
        {!showNotes ? (
          <button type="button" onClick={() => setShowNotes(true)} style={{ ...smallBtn, border: "1px dashed rgba(0,0,0,0.12)" }}>+ Nota</button>
        ) : (
          <div style={{ flex: 1 }}>
            <input aria-label="Observación para el conductor" value={note} onChange={e => onNoteChange?.(e.target.value)}
              placeholder="Ej: Puerta principal"
              style={{ width: "100%", padding: "5px 8px", borderRadius: radius.sm, fontSize: 11, fontFamily: "'Inter', sans-serif",
                border: "1px solid rgba(0,0,0,0.08)", outline: "none", background: colors.surface.paperLight, boxSizing: "border-box" as const }} />
          </div>
        )}
        {!scheduledAt ? (
          <button type="button" onClick={() => onScheduleChange?.(new Date(Date.now() + 3600000).toISOString())}
            style={{ ...smallBtn, border: "1px dashed rgba(0,0,0,0.12)" }}>🕐 Programar</button>
        ) : (
          <div style={{ display: "flex", gap: 4, alignItems: "center" }}>
            <span style={{ fontSize: 11, color: colors.textSecondary, fontFamily: "'Inter', sans-serif" }}>🕐 {new Date(scheduledAt).toLocaleString("es-EC", { dateStyle: "short", timeStyle: "short" })}</span>
            <button type="button" onClick={() => onScheduleChange?.(null)}
              style={{ padding: "3px 6px", borderRadius: radius.sm, fontSize: 10, fontWeight: 600, border: "none", cursor: "pointer", background: "transparent", color: colors.danger, fontFamily: "'Inter', sans-serif" }}>
              X
            </button>
          </div>
        )}
      </div>

      <div style={{ display: "flex", gap: 8, marginTop: 4 }}>
        <button type="button" onClick={onBack} aria-label="Volver"
          style={{
            flex: 1, height: 44, borderRadius: radius.sm, fontSize: 14, fontWeight: 600,
            border: "1px solid rgba(0,0,0,0.1)", cursor: "pointer",
            fontFamily: "'Inter', sans-serif", background: "transparent", color: colors.textPrimary,
            transition: "background 0.15s ease",
          }}
          onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(0,0,0,0.03)"; }}
          onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}>
          Atrás
        </button>
        <button type="button" onClick={onConfirm} disabled={loading}
          aria-label="Solicitar viaje"
          style={{
            flex: 2, height: 44, background: loading ? "#8e96a0" : colors.cobalt, color: "#fff",
            borderRadius: radius.sm, fontSize: 14, fontWeight: 600, fontFamily: "'Inter', sans-serif",
            border: "none", cursor: loading ? "not-allowed" : "pointer",
            transition: "background 0.15s ease, opacity 0.15s ease",
            opacity: loading ? 0.6 : 1,
          }}>
          {loading ? (
            <span style={{ display: "inline-flex", alignItems: "center", gap: 8 }}>
              <span style={{ width: 16, height: 16, border: "2px solid rgba(255,255,255,0.3)", borderTopColor: "#fff", borderRadius: "50%", animation: "spin 0.7s linear infinite" }} />
              Procesando...
            </span>
          ) : "Solicitar viaje"}
        </button>
      </div>
    </div>
  );
}
