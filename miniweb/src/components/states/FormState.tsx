"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";
const INPUT = "#e1e2e4";

interface FormStateProps {
  name: string; setName: Dispatch<SetStateAction<string>>;
  destQuery: string; setDestQuery: Dispatch<SetStateAction<string>>;
  destSuggestions: Place[];
  dest: Place | null;
  onSearch: (q: string) => void;
  onSelect: (place: Place) => void;
  onConfirm: () => void;
  loading: boolean;
  pickupAddress: string;
}

export function FormState({ name, setName, destQuery, setDestQuery, destSuggestions, dest, onSearch, onSelect, onConfirm, loading, pickupAddress }: FormStateProps) {
  return (
    <div style={{ padding: "4px 14px 0", display: "flex", flexDirection: "column", gap: 6 }}>

      <p style={{ fontSize: 15, fontWeight: 600, color: T1, margin: 0 }}>
        Buenos días, <span style={{ color: G }}>{name || "Alex"}</span>
      </p>

      {/* Pickup — auto desde pin del mapa */}
      <div style={{
        background: "rgba(255,255,255,0.92)", backdropFilter: "blur(16px)",
        borderRadius: 14, boxShadow: "0px 4px 20px rgba(0,0,0,0.05)",
        border: "1px solid #bbcabf33", padding: 10,
        display: "flex", flexDirection: "column", gap: 8,
        position: "relative", overflow: "hidden",
      }}>
        {/* Pickup (solo lectura) */}
        <div style={{ display: "flex", alignItems: "center", gap: 8, background: `${G}08`, borderRadius: 10, padding: "8px 10px", border: "1px solid #bbcabf33" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: G, fontVariationSettings: "'FILL' 1" }}>my_location</span>
          <div style={{ flex: 1, minWidth: 0 }}>
            <label style={{ fontSize: 10, fontWeight: 600, letterSpacing: "0.05em", color: T2, fontFamily: "Inter" }}>Recoger en</label>
            <p style={{ fontSize: 13, fontWeight: 500, color: T1, margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
              {pickupAddress || "Arrastra el mapa para ajustar"}
            </p>
          </div>
          <span style={{ fontSize: 10, color: G, fontWeight: 600, letterSpacing: "0.03em", fontFamily: "Inter" }}>PIN 📌</span>
        </div>

        {/* Línea vertical */}
        <div style={{ position: "absolute", left: 28, top: 48, height: 24, width: 2, background: `${G}40` }} />

        {/* Destino */}
        <div style={{ display: "flex", alignItems: "center", gap: 8, background: INPUT, borderRadius: 10, padding: "8px 10px", border: "1px solid #006c494D" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: G }}>search</span>
          <div style={{ flex: 1 }}>
            <label style={{ fontSize: 10, fontWeight: 600, letterSpacing: "0.05em", color: G, fontFamily: "Inter" }}>¿A dónde vas?</label>
            <input style={{ background: "transparent", border: "none", padding: 0, fontSize: 14, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
              placeholder="Buscar destino" value={destQuery}
              onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value); }} autoFocus />
          </div>
        </div>

        {/* Sugerencias destino */}
        {destSuggestions.length > 0 && !dest && (
          <div style={{ background: "#fff", borderRadius: 10, border: "1px solid #bbcabf", overflow: "hidden", marginTop: -4 }}>
            {destSuggestions.slice(0, 2).map((p, i) => (
              <div key={i} onClick={() => onSelect(p)}
                style={{ padding: "6px 10px", display: "flex", alignItems: "center", gap: 8, cursor: "pointer", borderBottom: i === 0 ? "1px solid #bbcabf" : "none", fontSize: 13, color: T1 }}>
                <span className="material-symbols-outlined" style={{ fontSize: 16, color: G, flexShrink: 0 }}>trip</span>
                {p.name}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Buscar viaje */}
      <button onClick={onConfirm} disabled={!dest || loading}
        style={{
          width: "100%", height: 42, background: G, color: "#fff",
          borderRadius: 9999, display: "flex", alignItems: "center", justifyContent: "center",
          fontSize: 15, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: "pointer",
          boxShadow: "0px 4px 12px rgba(0,108,73,0.3)",
          opacity: (!dest || loading) ? 0.4 : 1, transition: "opacity 0.2s"
        }}
      >{loading ? "Calculando..." : "Buscar viaje"}</button>
    </div>
  );
}
