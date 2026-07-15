"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

const G = "#006c49";  // primary green
const T1 = "#191c1e"; // on-surface
const T2 = "#3c4a42"; // on-surface-variant
const BG = "#f8f9fb"; // surface/background
const INPUT = "#e1e2e4"; // surface-container-highest
const CARD_BG = "#ffffff/90"; // glass card
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

export function FormState({ phone, setPhone, name, setName, pickupQuery, setPickupQuery, destQuery, setDestQuery, pickupSuggestions, destSuggestions, pickup, dest, onSearch, onSelect, onConfirm, loading, paymentMethod, setPaymentMethod }: FormStateProps) {
  return (
    <div style={{ padding: "8px 16px 20px", display: "flex", flexDirection: "column", gap: 16 }}>

      {/* Greeting */}
      <div>
        <h1 style={{ fontSize: 20, fontWeight: 600, lineHeight: 1.4, color: T1 }}>
          Good morning, <span style={{ color: G }}>{name || "Alex"}</span>
        </h1>
      </div>

      {/* Glassmorphic Interaction Card */}
      <div style={{
        background: "rgba(255,255,255,0.9)",
        backdropFilter: "blur(20px)",
        borderRadius: 16,
        boxShadow: "0px 4px 20px rgba(0,0,0,0.05)",
        border: `1px solid ${OUTLINE}4D`,
        padding: 16,
        display: "flex",
        flexDirection: "column",
        gap: 12,
        position: "relative",
        overflow: "hidden"
      }}>

        {/* Pickup field */}
        <div style={{ display: "flex", alignItems: "center", gap: 12, background: INPUT, borderRadius: 12, padding: 12, border: `1px solid ${OUTLINE}4D`, transition: "border-color 0.2s" }}>
          <button style={{ background: "none", border: "none", padding: 0, display: "flex", alignItems: "center", justifyContent: "center", color: G, cursor: "pointer" }}>
            <span className="material-symbols-outlined" style={{ fontSize: 22, fontVariationSettings: "'FILL' 1" }}>my_location</span>
          </button>
          <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
            <label style={{ fontSize: 12, fontWeight: 600, letterSpacing: "0.05em", color: T2, marginBottom: 2, fontFamily: "Inter" }}>Current Location</label>
            <input
              style={{ background: "transparent", border: "none", padding: 0, fontSize: 16, fontWeight: 400, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
              placeholder="Search pickup location"
              type="text"
              value={pickupQuery}
              onChange={e => { setPickupQuery(e.target.value); onSearch(e.target.value, true); }}
              autoFocus
            />
          </div>
        </div>

        {/* Vertical divider */}
        <div style={{ position: "absolute", left: 33, top: 56, bottom: 56, width: 2, background: `${G}33` }} />

        {/* Destination field */}
        <div style={{ display: "flex", alignItems: "center", gap: 12, background: INPUT, borderRadius: 12, padding: 12, border: `1px solid ${G}4D`, boxShadow: "0px 1px 3px rgba(0,0,0,0.05)", transition: "border-color 0.2s" }}>
          <div style={{ display: "flex", alignItems: "center", justifyContent: "center", color: G, paddingLeft: 1 }}>
            <span className="material-symbols-outlined" style={{ fontSize: 22 }}>search</span>
          </div>
          <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
            <label style={{ fontSize: 12, fontWeight: 600, letterSpacing: "0.05em", color: G, marginBottom: 2, fontFamily: "Inter" }}>Where to?</label>
            <input
              style={{ background: "transparent", border: "none", padding: 0, fontSize: 16, fontWeight: 400, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
              placeholder="Search destination"
              type="text"
              value={destQuery}
              onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value, false); }}
            />
          </div>
        </div>

        {/* Phone number */}
        <div style={{ display: "flex", alignItems: "center", gap: 12, background: INPUT, borderRadius: 12, padding: 12, border: `1px solid ${OUTLINE}4D` }}>
          <span className="material-symbols-outlined" style={{ fontSize: 20, color: T2 }}>phone</span>
          <input style={{ background: "transparent", border: "none", padding: 0, fontSize: 16, fontFamily: "Inter", color: T1, width: "100%", outline: "none" }}
            type="tel" placeholder="Phone number" value={phone} onChange={e => setPhone(e.target.value)} />
        </div>

        {/* Suggestions dropdown */}
        {pickupSuggestions.length > 0 && !pickup && (
          <div style={{ background: "#fff", borderRadius: 12, border: `1px solid ${OUTLINE}`, overflow: "hidden", position: "relative", zIndex: 10, marginTop: -8 }}>
            {pickupSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} onClick={() => onSelect(p, true)} style={{ padding: "10px 14px", display: "flex", alignItems: "center", gap: 10, cursor: "pointer", borderBottom: i < 3 ? `1px solid ${OUTLINE}` : "none", background: "transparent", transition: "background 0.15s" }}
                onMouseEnter={e => e.currentTarget.style.background = "#f8f9fb"} onMouseLeave={e => e.currentTarget.style.background = "transparent"}>
                <span className="material-symbols-outlined" style={{ fontSize: 18, color: G, flexShrink: 0 }}>location_on</span>
                <div><p style={{ fontSize: 14, fontWeight: 500, color: T1, margin: 0 }}>{p.name}</p><p style={{ fontSize: 12, color: T2, margin: 0 }}>{p.address}</p></div>
              </div>
            ))}
          </div>
        )}

        {destSuggestions.length > 0 && !dest && (
          <div style={{ background: "#fff", borderRadius: 12, border: `1px solid ${OUTLINE}`, overflow: "hidden", position: "relative", zIndex: 10, marginTop: -8 }}>
            {destSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} onClick={() => onSelect(p, false)} style={{ padding: "10px 14px", display: "flex", alignItems: "center", gap: 10, cursor: "pointer", borderBottom: i < 3 ? `1px solid ${OUTLINE}` : "none", background: "transparent", transition: "background 0.15s" }}
                onMouseEnter={e => e.currentTarget.style.background = "#f8f9fb"} onMouseLeave={e => e.currentTarget.style.background = "transparent"}>
                <span className="material-symbols-outlined" style={{ fontSize: 18, color: G, flexShrink: 0 }}>trip</span>
                <div><p style={{ fontSize: 14, fontWeight: 500, color: T1, margin: 0 }}>{p.name}</p><p style={{ fontSize: 12, color: T2, margin: 0 }}>{p.address}</p></div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Spacer */}
      <div style={{ flex: 1 }} />

      {/* Find a Ride button */}
      <button onClick={onConfirm} disabled={!pickup || !dest || !phone || loading}
        style={{
          width: "100%", height: 56, background: G, color: "#fff",
          borderRadius: 9999, display: "flex", alignItems: "center", justifyContent: "center",
          fontSize: 20, fontWeight: 600, fontFamily: "Inter",
          border: "none", cursor: "pointer",
          boxShadow: "0px 4px 12px rgba(0,108,73,0.3)",
          opacity: (!pickup || !dest || !phone || loading) ? 0.4 : 1,
          transition: "opacity 0.2s, transform 0.1s"
        }}
      >
        {loading ? "Calculating..." : "Find a Ride"}
      </button>
    </div>
  );
}
