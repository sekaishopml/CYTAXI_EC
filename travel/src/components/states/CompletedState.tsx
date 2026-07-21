"use client";
import { FareBreakdown, DriverInfo, Place } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

export function CompletedState({ fare, pickup, dest, onNewTrip, paymentMethod }: {
  fare: FareBreakdown | null;
  driver: DriverInfo | null;
  pickup: Place | null;
  dest: Place | null;
  onNewTrip: () => void;
  paymentMethod: "cash" | "card";
}) {
  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "24px 20px 16px", gap: 16, textAlign: "center" }}>
      <div style={{
        width: 72, height: 72, borderRadius: "50%",
        background: colors.cobaltBg,
        display: "flex", alignItems: "center", justifyContent: "center",
        margin: "0 auto",
      }}>
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke={colors.cobalt} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M20 6L9 17l-5-5"/>
        </svg>
      </div>
      <div>
        <p style={{ fontSize: 20, fontWeight: 700, color: colors.textPrimary, margin: 0, letterSpacing: "-0.02em" }}>Viaje completado</p>
        <p style={{ fontSize: 13, color: colors.textMuted, margin: "4px 0 0" }}>Gracias por viajar con CYTAXI</p>
      </div>

      {pickup && dest && (
        <div style={{
          background: colors.surface.paperLight, borderRadius: 10,
          boxShadow: shadows.card, border: "1px solid rgba(0,0,0,0.06)",
          padding: 14, textAlign: "left",
        }}>
          <div style={{ display: "flex", alignItems: "flex-start", gap: 12 }}>
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 4, marginTop: 2 }}>
              <div style={{ width: 8, height: 8, borderRadius: "50%", background: colors.cobalt }} />
              <div style={{ width: 1.5, height: 20, background: "rgba(0,0,0,0.08)" }} />
              <div style={{ width: 8, height: 8, borderRadius: "50%", background: colors.cobaltLight }} />
            </div>
            <div style={{ flex: 1 }}>
              <p style={{ fontSize: 10, fontWeight: 600, color: colors.textMuted, margin: 0, letterSpacing: "0.05em", textTransform: "uppercase" as const }}>RECOGIDA</p>
              <p style={{ fontSize: 14, fontWeight: 600, color: colors.textPrimary, margin: "2px 0 0" }}>{pickup.name}</p>
              <div style={{ height: 12 }} />
              <p style={{ fontSize: 10, fontWeight: 600, color: colors.textMuted, margin: 0, letterSpacing: "0.05em", textTransform: "uppercase" as const }}>DESTINO</p>
              <p style={{ fontSize: 14, fontWeight: 600, color: colors.textPrimary, margin: "2px 0 0" }}>{dest.name}</p>
            </div>
          </div>
        </div>
      )}

      {fare && (
        <div style={{ background: colors.surface.paperLight, borderRadius: 10, padding: 14, textAlign: "left", border: "1px solid rgba(0,0,0,0.06)" }}>
          <p style={{ fontSize: 10, fontWeight: 600, color: colors.textMuted, margin: "0 0 10px", letterSpacing: "0.05em", textTransform: "uppercase" as const }}>DETALLE</p>
          {[
            ["Distancia", `${fare.distance_km.toFixed(1)} km`],
            ["Tiempo", `${fare.eta_minutes} min`],
          ].map(([label, val]) => (
            <div key={label as string} style={{ display: "flex", justifyContent: "space-between", padding: "3px 0" }}>
              <span style={{ fontSize: 13, color: colors.textSecondary }}>{label}</span>
              <span style={{ fontSize: 13, fontWeight: 500, color: colors.textPrimary }}>{val}</span>
            </div>
          ))}
          <div style={{ height: 1, background: "rgba(0,0,0,0.06)", margin: "8px 0" }} />
          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <span style={{ fontSize: 17, fontWeight: 700, color: colors.textPrimary }}>Total</span>
            <span style={{ fontSize: 17, fontWeight: 700, color: colors.cobalt }}>${fare.total.toFixed(2)}</span>
          </div>
          <p style={{ fontSize: 11, color: colors.textSecondary, textAlign: "right", margin: "4px 0 0" }}>
            vía {paymentMethod === "cash" ? "Efectivo" : "Tarjeta"}
          </p>
        </div>
      )}

      <button type="button" onClick={onNewTrip} aria-label="Iniciar nuevo viaje"
        style={{
          width: "100%", height: 52, background: colors.cobalt, color: "#fff",
          borderRadius: 14, fontSize: 16, fontWeight: 600, fontFamily: "Inter, sans-serif",
          border: "none", cursor: "pointer",
          boxShadow: "0 4px 20px rgba(59,130,246,0.25)",
          transition: "all 0.2s",
          marginTop: 4,
        }}
        onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.01)"; }}
        onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}>
        Nuevo viaje
      </button>
    </div>
  );
}
