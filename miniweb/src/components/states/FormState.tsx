"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

const STITCH_GREEN = "#006c49";

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
    <div style={{ padding: "12px 20px 24px" }}>
      <div className="sheet-handle" />
      <h2 className="text-headline-mobile" style={{ marginBottom: 16 }}>Where to?</h2>

      {/* Phone */}
      <div style={{ marginBottom: 12 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Phone number</label>
        <div className="input-field" style={{ display: "flex", alignItems: "center", gap: 10, padding: 0, overflow: "hidden" }}>
          <span className="material-symbols-outlined" style={{ marginLeft: 16, fontSize: 20, color: "var(--uk-on-surface-variant)" }}>phone</span>
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} type="tel" placeholder="+593 99 999 9999" value={phone} onChange={e => setPhone(e.target.value)} />
        </div>
      </div>

      {/* Name */}
      <div style={{ marginBottom: 12 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Your name</label>
        <div className="input-field" style={{ display: "flex", alignItems: "center", gap: 10, padding: 0, overflow: "hidden" }}>
          <span className="material-symbols-outlined" style={{ marginLeft: 16, fontSize: 20, color: "var(--uk-on-surface-variant)" }}>person</span>
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} placeholder="Your name" value={name} onChange={e => setName(e.target.value)} />
        </div>
      </div>

      {/* Pickup */}
      <div style={{ marginBottom: 12 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Pickup</label>
        <div className="input-field" style={{ display: "flex", alignItems: "center", gap: 10, padding: 0, overflow: "hidden" }}>
          <div className="dot dot-pickup" style={{ marginLeft: 16, background: STITCH_GREEN }} />
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} placeholder="Where are you?" value={pickupQuery} onChange={e => { setPickupQuery(e.target.value); onSearch(e.target.value, true); }} autoFocus />
        </div>
        {pickupSuggestions.length > 0 && !pickup && (
          <div className="suggestions">
            {pickupSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} className="suggestion-item" onClick={() => onSelect(p, true)}>
                <span className="material-symbols-outlined" style={{ fontSize: 18, color: STITCH_GREEN, flexShrink: 0 }}>location_on</span>
                <div><p className="text-body-md" style={{ fontWeight: 500 }}>{p.name}</p><p className="text-label-sm text-muted">{p.address}</p></div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Destination */}
      <div style={{ marginBottom: 16 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Destination</label>
        <div className="input-field" style={{ display: "flex", alignItems: "center", gap: 10, padding: 0, overflow: "hidden" }}>
          <span className="material-symbols-outlined" style={{ marginLeft: 16, fontSize: 20, color: "#276ef1" }}>trip</span>
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} placeholder="Where to?" value={destQuery} onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value, false); }} />
        </div>
        {destSuggestions.length > 0 && !dest && (
          <div className="suggestions">
            {destSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} className="suggestion-item" onClick={() => onSelect(p, false)}>
                <span className="material-symbols-outlined" style={{ fontSize: 18, color: "#276ef1", flexShrink: 0 }}>location_on</span>
                <div><p className="text-body-md" style={{ fontWeight: 500 }}>{p.name}</p><p className="text-label-sm text-muted">{p.address}</p></div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Vehicle type cards */}
      <div style={{ marginBottom: 16 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 8 }}>Vehicle type</label>
        <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
          {[
            { id: "standard", name: "Standard", icon: "directions_car", eta: "3 min", price: "$5.50" },
            { id: "xl", name: "XL", icon: "airport_shuttle", eta: "5 min", price: "$8.20" },
            { id: "premium", name: "Premium", icon: "electric_car", eta: "4 min", price: "$10.50" },
          ].map(v => (
            <div key={v.id} className="vehicle-card" style={{ padding: "14px 16px", borderColor: "transparent", background: STITCH_GREEN + "08" }}>
              <div style={{ width: 44, height: 44, borderRadius: 12, background: "#f6f6f6", display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0 }}>
                <span className="material-symbols-outlined" style={{ fontSize: 22 }}>{v.icon}</span>
              </div>
              <div style={{ flex: 1 }}><p className="text-body-md" style={{ fontWeight: 600 }}>{v.name}</p><p className="text-label-sm text-muted">{v.eta} away</p></div>
              <p className="text-title" style={{ fontSize: 18 }}>{v.price}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Payment method */}
      <div style={{ marginBottom: 16 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 8 }}>Payment</label>
        <div style={{ display: "flex", gap: 10 }}>
          <button onClick={() => setPaymentMethod("cash")}
            style={{ flex: 1, padding: 12, borderRadius: 12, fontSize: 14, fontWeight: 600, border: "none", cursor: "pointer", transition: "all 0.2s", fontFamily: "Inter",
              background: paymentMethod === "cash" ? STITCH_GREEN : "#f6f6f6",
              color: paymentMethod === "cash" ? "#fff" : "#191c1e" }}
          >💵 Cash</button>
          <button onClick={() => setPaymentMethod("card")}
            style={{ flex: 1, padding: 12, borderRadius: 12, fontSize: 14, fontWeight: 600, border: "none", cursor: "pointer", transition: "all 0.2s", fontFamily: "Inter",
              background: paymentMethod === "card" ? STITCH_GREEN : "#f6f6f6",
              color: paymentMethod === "card" ? "#fff" : "#191c1e" }}
          >💳 Card</button>
        </div>
      </div>

      {/* Confirm */}
      <button onClick={onConfirm} disabled={!pickup || !dest || !phone || loading} style={{ width: "100%", padding: 16, borderRadius: 14, fontSize: 17, fontWeight: 600, border: "none", cursor: "pointer", fontFamily: "Inter", transition: "opacity 0.2s",
        background: STITCH_GREEN, color: "#fff", opacity: (!pickup || !dest || !phone || loading) ? 0.4 : 1 }}>
        {loading ? "Calculating..." : "Request your ride"}
      </button>
    </div>
  );
}
