"use client";
const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";

interface PickUpStepProps {
  onConfirm: () => void;
  address: string;
  loading: boolean;
}

export function PickUpStep({ onConfirm, address, loading }: PickUpStepProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "16px 20px 14px", gap: 14 }}>
      <div style={{ textAlign: "center" }}>
        <p style={{ fontSize: 19, fontWeight: 700, color: T1, margin: 0, letterSpacing: "-0.02em" }}>
          ¿Dónde te recogemos?
        </p>
        <p style={{ fontSize: 13, color: "#8a8a8a", margin: "4px 0 0", fontWeight: 400 }}>
          Arrastra el mapa para ajustar el punto exacto
        </p>
      </div>

      <div style={{
        background: "rgba(255,255,255,0.8)", backdropFilter: "blur(20px)", borderRadius: 14,
        boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 6px 20px rgba(0,0,0,0.04)",
        border: "1px solid rgba(0,0,0,0.05)", padding: "14px 16px",
        display: "flex", alignItems: "center", gap: 12,
      }}>
        <div style={{
          width: 12, height: 12, borderRadius: "50%", background: G,
          flexShrink: 0, boxShadow: `0 0 0 4px ${G}12, 0 0 0 6px ${G}08`,
        }} />
        <div style={{ flex: 1, minWidth: 0 }}>
          <p style={{ fontSize: 11, fontWeight: 600, color: "#9a9a9a", margin: 0, letterSpacing: "0.04em", textTransform: "uppercase" }}>
            Punto de recogida
          </p>
          <p style={{ fontSize: 15, fontWeight: 500, color: T1, margin: "3px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
            {address || "Ubicación actual"}
          </p>
        </div>
      </div>

      <button onClick={onConfirm}
        style={{
          width: "100%", height: 52, background: "#121212", color: "#fff", borderRadius: 14,
          fontSize: 16, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: "pointer",
          boxShadow: "0 2px 8px rgba(0,0,0,0.12)", transition: "all 0.2s cubic-bezier(0.4, 0, 0.2, 1)",
          letterSpacing: "-0.01em",
        }}>
        Confirmar ubicación
      </button>
    </div>
  );
}
