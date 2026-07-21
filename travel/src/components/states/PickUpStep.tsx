"use client";
import { colors, radius } from "@cytaxi/design-tokens";

interface PickUpStepProps {
  onConfirm: () => void;
  address: string;
  loading: boolean;
}

export function PickUpStep({ onConfirm, address, loading }: PickUpStepProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "14px 18px 12px", gap: 12 }}>
      <div style={{ textAlign: "left" }}>
        <p style={{ margin: 0, fontSize: 18, fontWeight: 600, fontFamily: "'Space Grotesk', sans-serif", color: colors.textPrimary, letterSpacing: "-0.02em" }}>
          ¿Dónde te recogemos?
        </p>
        <p style={{ fontSize: 12, color: colors.textMuted, margin: "4px 0 0", fontWeight: 400, fontFamily: "'Inter', sans-serif" }}>
          Arrastra el mapa para ajustar el punto exacto
        </p>
      </div>

      <div style={{
        background: colors.surface.paperLight,
        borderRadius: radius.md, border: "1px solid rgba(0,0,0,0.06)",
        padding: "12px 14px", display: "flex", alignItems: "center", gap: 10,
      }}>
        <div style={{
          width: 24, height: 24, borderRadius: "50%", background: colors.cobalt,
          flexShrink: 0, display: "flex", alignItems: "center", justifyContent: "center",
          color: "#fff", fontSize: 11, fontWeight: 700, fontFamily: "'Inter', sans-serif",
        }}>A</div>
        <div style={{ flex: 1, minWidth: 0 }}>
          <p style={{ fontSize: 9, fontWeight: 500, color: colors.textMuted, margin: 0, letterSpacing: "0.06em", textTransform: "uppercase", fontFamily: "'JetBrains Mono', monospace" }}>
            ORIGEN
          </p>
          <p style={{ fontSize: 14, fontWeight: 500, color: colors.textPrimary, margin: "2px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" as const, fontFamily: "'Inter', sans-serif" }}>
            {address || "Ubicación actual"}
          </p>
        </div>
      </div>

      <button type="button" onClick={onConfirm} disabled={loading}
        aria-label="Confirmar ubicación de recogida"
        style={{
          width: "100%", height: 48, background: loading ? "#8e96a0" : colors.cobalt, color: "#fff",
          borderRadius: radius.xs, fontSize: 15, fontWeight: 600, fontFamily: "'Inter', sans-serif",
          border: "none", cursor: loading ? "not-allowed" : "pointer",
          transition: "background 0.15s ease, opacity 0.15s ease",
          opacity: loading ? 0.6 : 1,
        }}>
        {loading ? (
          <span style={{ display: "inline-flex", alignItems: "center", gap: 8 }}>
            <span style={{ width: 16, height: 16, border: "2px solid rgba(255,255,255,0.3)", borderTopColor: "#fff", borderRadius: "50%", animation: "spin 0.7s linear infinite" }} />
            Detectando...
          </span>
        ) : "Confirmar ubicación"}
      </button>
    </div>
  );
}
