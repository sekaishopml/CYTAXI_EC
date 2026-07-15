"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";
const INPUT = "#e1e2e4";
const OUTLINE = "#bbcabf";

interface FormStateProps {
  phone: string; setPhone: Dispatch<SetStateAction<string>>;
  name: string; setName: Dispatch<SetStateAction<string>>;
  pickupQuery: string; setPickupQuery: Dispatch<SetStateAction<string>>;
  destQuery: string; setDestQuery: Dispatch<SetStateAction<string>>;
  pickupSuggestions: Place[]; destSuggestions: Place[];
  pickup: Place | null; dest: Place | null;
  onSearch: (q: string, isPickup: boolean) => void;
  onSelect: (place: Place, isPickup: boolean) => void;
  onConfirm: () => void;
  loading: boolean;
  paymentMethod: "cash" | "card";
  setPaymentMethod: Dispatch<SetStateAction<"cash" | "card">>;
}

export function FormState({ phone, setPhone, name, setName, pickupQuery, setPickupQuery, destQuery, setDestQuery, pickupSuggestions, destSuggestions, pickup, dest, onSearch, onSelect, onConfirm, loading }: FormStateProps) {
  return (
    <div style={{ padding: "4px 14px 0", display: "flex", flexDirection: "column", gap: 8 }}>

      <p style={{ fontSize: 15, fontWeight: 600, color: T1, margin: 0 }}>
        Buenos días, <span style={{ color: G }}>{name || "Alex"}</span>
      </p>

      <div style={{
        background: "rgba(255,255,255,0.92)",
        backdropFilter: "blur(16px)",
        borderRadius: 14,
        boxShadow: "0px 4px 20px rgba(0,0,0,0.05)",
        border: "1px solid #bbcabf33",
        padding: 10,
        display: "flex", flexDirection: "column", gap: 8,
        position: "relative", overflow: "hidden"
      }}>
        {/* Origen */}
        <div style={{ display: "flex", alignItems: "center", gap: 8, background: INPUT, borderRadius: 10, padding: "8px 10px", border: "1px solid #bbcabf33" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: G, fontVariationSettings: "'FILL' 1" }}>my_location</span>
          <div style={{ flex: 1 }}>
            <label style={{ fontSize: 10, fontWeight: 600, letterSpacing: "0.05em", color: T2, fontFamily: "Inter" }}>Origen</label>
            <input style={{ background: "transparent", border: "none", padding: 0, fontSize: 14, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
              placeholder="Buscar dirección de salida" value={pickupQuery} onChange={e => { setPickupQuery(e.target.value); onSearch(e.target.value, true); }} autoFocus />
          </div>
        </div>

        <div style={{ position: "absolute", left: 28, top: 44, height: 24, width: 2, background: `${G}40` }} />

        {/* Destino */}
        <div style={{ display: "flex", alignItems: "center", gap: 8, background: INPUT, borderRadius: 10, padding: "8px 10px", border: "1px solid ${G}4D" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18, color: G }}>search</span>
          <div style={{ flex: 1 }}>
            <label style={{ fontSize: 10, fontWeight: 600, letterSpacing: "0.05em", color: G, fontFamily: "Inter" }}>Destino</label>
            <input style={{ background: "transparent", border: "none", padding: 0, fontSize: 14, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
              placeholder="Buscar destino" value={destQuery} onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value, false); }} />
          </div>
        </div>

        {/* Teléfono */}
        <div style={{ display: "flex", alignItems: "center", gap: 8, background: INPUT, borderRadius: 10, padding: "8px 10px", border: "1px solid #bbcabf33" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 16, color: T2 }}>phone</span>
          <input style={{ background: "transparent", border: "none", padding: 0, fontSize: 14, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
            type="tel" placeholder="Teléfono" value={phone} onChange={e => setPhone(e.target.value)} />
        </div>

        {pickupSuggestions.length > 0 && !pickup && (
          <div style={{ background: "#fff", borderRadius: 10, border: "1px solid #bbcabf", overflow: "hidden", marginTop: -4 }}>
            {pickupSuggestions.slice(0, 2).map((p, i) => (
              <div key={i} onClick={() => onSelect(p, true)}
                style={{ padding: "6px 10px", display: "flex", alignItems: "center", gap: 8, cursor: "pointer", borderBottom: i === 0 ? "1px solid #bbcabf" : "none", fontSize: 13, color: T1 }}>
                <span className="material-symbols-outlined" style={{ fontSize: 16, color: G, flexShrink: 0 }}>location_on</span>
                {p.name}
              </div>
            ))}
          </div>
        )}

        {destSuggestions.length > 0 && !dest && (
          <div style={{ background: "#fff", borderRadius: 10, border: "1px solid #bbcabf", overflow: "hidden", marginTop: -4 }}>
            {destSuggestions.slice(0, 2).map((p, i) => (
              <div key={i} onClick={() => onSelect(p, false)}
                style={{ padding: "6px 10px", display: "flex", alignItems: "center", gap: 8, cursor: "pointer", borderBottom: i === 0 ? "1px solid #bbcabf" : "none", fontSize: 13, color: T1 }}>
                <span className="material-symbols-outlined" style={{ fontSize: 16, color: G, flexShrink: 0 }}>trip</span>
                {p.name}
              </div>
            ))}
          </div>
        )}
      </div>

      <button onClick={onConfirm} disabled={!pickup || !dest || !phone || loading}
        style={{
          width: "100%", height: 42, background: G, color: "#fff",
          borderRadius: 9999, display: "flex", alignItems: "center", justifyContent: "center",
          fontSize: 15, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: "pointer",
          boxShadow: "0px 4px 12px rgba(0,108,73,0.3)",
          opacity: (!pickup || !dest || !phone || loading) ? 0.4 : 1,
          transition: "opacity 0.2s"
        }}
      >{loading ? "Calculando..." : "Buscar viaje"}</button>
    </div>
  );
}
