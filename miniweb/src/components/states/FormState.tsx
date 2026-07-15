"use client";
import { Dispatch, SetStateAction } from "react";
import { Place } from "@/types";

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
        <input className="input-field" type="tel" placeholder="+593 99 999 9999" value={phone} onChange={e => setPhone(e.target.value)} />
      </div>

      {/* Name */}
      <div style={{ marginBottom: 12 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Your name</label>
        <input className="input-field" placeholder="Passenger name" value={name} onChange={e => setName(e.target.value)} />
      </div>

      {/* Pickup */}
      <div style={{ marginBottom: 12 }}>
        <label className="text-label-sm text-muted" style={{ display: "block", marginBottom: 6 }}>Pickup</label>
        <div className="input-field" style={{ display: "flex", alignItems: "center", gap: 10, padding: 0, overflow: "hidden" }}>
          <div className="dot dot-pickup" style={{ marginLeft: 16 }} />
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} placeholder="Where are you?" value={pickupQuery} onChange={e => { setPickupQuery(e.target.value); onSearch(e.target.value, true); }} autoFocus />
        </div>
        {pickupSuggestions.length > 0 && !pickup && (
          <div className="suggestions">
            {pickupSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} className="suggestion-item" onClick={() => onSelect(p, true)}>
                <div className="dot dot-pickup" style={{ width: 8, height: 8 }} />
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
          <div className="dot dot-dest" style={{ marginLeft: 16 }} />
          <input style={{ flex: 1, border: "none", background: "transparent", outline: "none", fontSize: 16, fontFamily: "Inter", padding: "16px 0" }} placeholder="Where to?" value={destQuery} onChange={e => { setDestQuery(e.target.value); onSearch(e.target.value, false); }} />
        </div>
        {destSuggestions.length > 0 && !dest && (
          <div className="suggestions">
            {destSuggestions.slice(0, 4).map((p, i) => (
              <div key={i} className="suggestion-item" onClick={() => onSelect(p, false)}>
                <div className="dot dot-dest" style={{ width: 8, height: 8 }} />
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
            { id: "standard", name: "Standard", icon: "🚗", eta: "3 min", price: "$5.50" },
            { id: "xl", name: "XL", icon: "🚐", eta: "5 min", price: "$8.20" },
            { id: "premium", name: "Premium", icon: "🚙", eta: "4 min", price: "$10.50" },
          ].map(v => (
            <div key={v.id} className="vehicle-card" style={{ padding: "14px 16px" }}>
              <div className="vehicle-icon" style={{ width: 44, height: 44, fontSize: 20 }}>{v.icon}</div>
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
            className={`btn ${paymentMethod === "cash" ? "btn-primary" : "btn-secondary"}`}
            style={{ flex: 1, fontSize: 14, padding: 12 }}
          >💵 Cash</button>
          <button onClick={() => setPaymentMethod("card")}
            className={`btn ${paymentMethod === "card" ? "btn-primary" : "btn-secondary"}`}
            style={{ flex: 1, fontSize: 14, padding: 12 }}
          >💳 Card</button>
        </div>
      </div>

      {/* Confirm */}
      <button onClick={onConfirm} disabled={!pickup || !dest || !phone || loading} className="btn btn-primary" style={{ fontSize: 17, padding: 16, borderRadius: 14 }}>
        {loading ? "Calculating..." : "Confirm"}
      </button>
    </div>
  );
}
