"use client";
import type { Place, DriverInfo } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

interface DestinationStateProps {
  dest: Place | null;
  driver: DriverInfo | null;
  onComplete: () => void;
}

export function DestinationState({ dest, driver, onComplete }: DestinationStateProps) {
  return (
    <div style={{ padding: "24px 20px", textAlign: "center" }}>
      <div style={{
        width: 72, height: 72, borderRadius: "50%",
        background: colors.cobaltBg,
        display: "flex", alignItems: "center", justifyContent: "center",
        margin: "0 auto 16px",
      }}>
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke={colors.green} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z"/><circle cx="12" cy="10" r="3"/>
        </svg>
      </div>

      <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: "0 0 4px", letterSpacing: "-0.01em" }}>
        Llegando a tu destino
      </p>
      <p style={{ fontSize: 13, color: colors.textMuted, margin: "0 0 16px" }}>
        {dest?.name || "Prepárate para bajar"}
      </p>

      {dest && (
        <div style={{
          background: colors.surface.paperLight, borderRadius: 10,
          boxShadow: shadows.card, border: "1px solid rgba(0,0,0,0.06)",
          padding: "14px 16px", textAlign: "left",
        }}>
          <p style={{ fontSize: 10, fontWeight: 600, color: colors.textMuted, textTransform: "uppercase" as const, letterSpacing: "0.05em", margin: "0 0 4px" }}>
            Destino
          </p>
          <p style={{ fontSize: 15, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>
            {dest.name}
          </p>
          <p style={{ fontSize: 12, color: colors.textSecondary, margin: "2px 0 0" }}>
            {dest.address}
          </p>
        </div>
      )}

      {driver && (
        <div style={{
          marginTop: 12, display: "flex", justifyContent: "center", gap: 12,
        }}>
          <div style={{ textAlign: "center" }}>
            <p style={{ fontSize: 13, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>{driver.name}</p>
            <p style={{ fontSize: 11, color: colors.textMuted, margin: "2px 0 0" }}>Tu conductor</p>
          </div>
        </div>
      )}

      <button type="button" onClick={onComplete} aria-label="Finalizar viaje"
        style={{
          width: "100%", marginTop: 16, height: 52,
          background: colors.green, color: "#fff",
          borderRadius: 14, fontSize: 16, fontWeight: 600,
          border: "none", cursor: "pointer", fontFamily: "Inter, sans-serif",
          boxShadow: "0 4px 20px rgba(59,130,246,0.25)",
          transition: "all 0.2s",
        }}
        onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.01)"; }}
        onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}>
        Finalizar viaje
      </button>
    </div>
  );
}
