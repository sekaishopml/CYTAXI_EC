"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";

interface FormStateProps {
  destQuery: string; setDestQuery: Dispatch<SetStateAction<string>>;
  destSuggestions: Place[];
  dest: Place | null;
  onSearch: (q: string) => void;
  onSelect: (place: Place) => void;
  onConfirm: () => void;
  onBack: () => void;
  onClearDest: () => void;
  loading: boolean;
  pickupAddress: string;
}

export function FormState({ destQuery, setDestQuery, destSuggestions, dest, onSearch, onSelect, onConfirm, onBack, onClearDest, loading, pickupAddress }: FormStateProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "8px 20px 14px", gap: 12 }}>
      {/* Header */}
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        <button onClick={onBack} style={{ background: "none", border: "none", cursor: "pointer", padding: 6, borderRadius: 8, color: T1, display: "flex", alignItems: "center", justifyContent: "center", transition: "background 0.15s" }}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M19 12H5m7-7-7 7 7 7"/></svg>
        </button>
        <p style={{ fontSize: 16, fontWeight: 600, color: T1, margin: 0, letterSpacing: "-0.01em" }}>
          Buenos días, <span style={{ color: G }}>Alex</span>
        </p>
      </div>

      {/* Glass card */}
      <div style={{
        background: "rgba(255,255,255,0.88)", backdropFilter: "blur(24px) saturate(180%)",
        borderRadius: 16, boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 8px 24px rgba(0,0,0,0.04)",
        border: "1px solid rgba(0,0,0,0.04)", overflow: "hidden",
      }}>
        {/* Pickup row */}
        <div onClick={onBack}
          style={{ display: "flex", alignItems: "center", gap: 12, padding: "14px 16px", cursor: "pointer", transition: "background 0.15s" }}>
          <div style={{ width: 10, height: 10, borderRadius: "50%", background: G, flexShrink: 0, boxShadow: `0 0 0 3px ${G}12` }} />
          <div style={{ flex: 1, minWidth: 0 }}>
            <p style={{ fontSize: 10, fontWeight: 600, color: "#9a9a9a", margin: 0, letterSpacing: "0.04em", textTransform: "uppercase" }}>
              Recoger en
            </p>
            <p style={{ fontSize: 14, fontWeight: 500, color: T1, margin: "2px 0 0", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
              {pickupAddress}
            </p>
          </div>
        </div>

        {/* Divider */}
        <div style={{ height: 1, background: "rgba(0,0,0,0.06)", margin: "0 16px" }} />

        {/* Destination row */}
        <div style={{ display: "flex", alignItems: "center", gap: 12, padding: "14px 16px", background: dest ? `${G}06` : "transparent", transition: "background 0.2s" }}>
          <div style={{ width: 10, height: 10, borderRadius: "50%", background: "#448aff", flexShrink: 0, boxShadow: "0 0 0 3px rgba(68,138,255,0.15)" }} />
          <div style={{ flex: 1, minWidth: 0 }}>
            <p style={{ fontSize: 10, fontWeight: 600, color: dest ? G : "#9a9a9a", margin: 0, letterSpacing: "0.04em", textTransform: "uppercase", transition: "color 0.2s" }}>
              ¿A dónde vas?
            </p>
            <input
              style={{ background: "transparent", border: "none", padding: 0, fontSize: 14, fontFamily: "Inter", color: T1, width: "100%", outline: "none", fontWeight: 500, margin: "2px 0 0", letterSpacing: "-0.01em" }}
              placeholder="Buscar destino"
              value={destQuery}
              onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value); }}
              autoFocus
            />
          </div>
          {destQuery && (
            <button onClick={onClearDest} style={{ background: "rgba(0,0,0,0.05)", border: "none", cursor: "pointer", padding: 4, borderRadius: "50%", display: "flex", transition: "background 0.15s" }}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke={T2} strokeWidth="2" strokeLinecap="round"><path d="M18 6 6 18M6 6l12 12"/></svg>
            </button>
          )}
        </div>

        {/* Suggestions */}
        {destSuggestions.length > 0 && !dest && (
          <div style={{ borderTop: "1px solid rgba(0,0,0,0.05)" }}>
            {destSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} onClick={() => onSelect(p)}
                style={{ padding: "12px 16px", display: "flex", alignItems: "center", gap: 12, cursor: "pointer", borderBottom: i < destSuggestions.length - 1 ? "1px solid rgba(0,0,0,0.04)" : "none", transition: "background 0.1s" }}>
                <span style={{ fontSize: 16, flexShrink: 0, opacity: 0.5 }}>📍</span>
                <div style={{ minWidth: 0 }}>
                  <p style={{ margin: 0, fontSize: 14, fontWeight: 500, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", color: T1 }}>{p.name}</p>
                  <p style={{ margin: "2px 0 0", fontSize: 12, color: "#8a8a8a", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{p.address}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Action button */}
      <button onClick={onConfirm} disabled={loading}
        style={{ width: "100%", height: 52, background: dest ? "#121212" : G, color: "#fff", borderRadius: 14, fontSize: 16, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: dest ? "pointer" : "pointer", boxShadow: dest ? "0 2px 8px rgba(0,0,0,0.12)" : "0 4px 20px rgba(0,108,73,0.2)", transition: "all 0.25s cubic-bezier(0.4, 0, 0.2, 1)", letterSpacing: "-0.01em" }}>
        {loading ? "Calculando..." : dest ? "Buscar viaje" : "Selecciona tu destino"}
      </button>
    </div>
  );
}
