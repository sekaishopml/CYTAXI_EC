"use client";
import { Dispatch, SetStateAction } from "react";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";
const INPUT = "#e1e2e4";

interface PickUpStepProps {
  onConfirm: () => void;
  address: string;
  loading: boolean;
}

export function PickUpStep({ onConfirm, address, loading }: PickUpStepProps) {
  return (
    <div style={{ padding: "8px 14px 0", display: "flex", flexDirection: "column", gap: 8 }}>

      {/* Título */}
      <div style={{ textAlign: "center" }}>
        <p style={{ fontSize: 17, fontWeight: 700, color: T1, margin: 0 }}>Fija el punto de partida</p>
        <p style={{ fontSize: 13, color: T2, margin: "2px 0 0" }}>Arrastra el mapa para mover el marcador</p>
      </div>

      {/* Dirección detectada */}
      <div style={{
        background: INPUT, borderRadius: 12, padding: "10px 12px",
        display: "flex", alignItems: "center", gap: 8,
        border: "1px solid #bbcabf33",
      }}>
        <div style={{ width: 12, height: 12, borderRadius: "50%", background: G, flexShrink: 0 }} />
        <p style={{ fontSize: 16, fontWeight: 500, color: T1, flex: 1, margin: 0 }}>
          {address || "Detectando..."}
        </p>
        <span className="material-symbols-outlined" style={{ fontSize: 20, color: T2, flexShrink: 0 }}>search</span>
      </div>

      {/* Botón confirmar */}
      <button onClick={onConfirm} disabled={loading}
        style={{
          width: "100%", height: 48, background: T1, color: "#fff",
          borderRadius: 12, fontSize: 16, fontWeight: 600, fontFamily: "Inter",
          border: "none", cursor: "pointer",
          boxShadow: "0px 4px 12px rgba(0,0,0,0.15)",
          opacity: loading ? 0.4 : 1, transition: "opacity 0.2s"
        }}
      >{loading ? "..." : "Confirmar punto de partida"}</button>
    </div>
  );
}
