"use client";
import { useState } from "react";
import { FareBreakdown, DriverInfo, Place } from "@/types";

interface CompletedStateProps {
  fare: FareBreakdown | null;
  driver: DriverInfo | null;
  pickup: Place | null;
  dest: Place | null;
  onNewTrip: () => void;
  paymentMethod: "cash" | "card";
}

export function CompletedState({ fare, driver, pickup, dest, onNewTrip, paymentMethod }: CompletedStateProps) {
  const [rating, setRating] = useState(0);
  const [rated, setRated] = useState(false);

  return (
    <div style={{ padding: "12px 20px 24px", textAlign: "center" }}>
      <div className="sheet-handle" />
      <div style={{ fontSize: 56, marginBottom: 8 }}>✅</div>
      <h2 className="text-headline-mobile" style={{ marginBottom: 16 }}>Trip Completed</h2>

      {/* Route */}
      {pickup && dest && (
        <div style={{ background: "var(--uk-input-bg)", borderRadius: 16, padding: 14, marginBottom: 16, textAlign: "left" }}>
          <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 6 }}>
            <div className="dot dot-pickup" /><p className="text-body-md">{pickup.name}</p>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <div className="dot dot-dest" /><p className="text-body-md">{dest.name}</p>
          </div>
        </div>
      )}

      {/* Fare */}
      {fare && (
        <div style={{ background: "var(--uk-input-bg)", borderRadius: 16, padding: 16, marginBottom: 16, textAlign: "left" }}>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Distance</span><span>{fare.distance_km.toFixed(1)} km</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 6 }} className="text-body-md"><span className="text-muted">Time</span><span>{fare.eta_minutes} min</span></div>
          <div className="divider" />
          <div style={{ display: "flex", justifyContent: "space-between" }} className="text-headline-mobile"><span>Total</span><span>${fare.total.toFixed(2)}</span></div>
          <div className="text-label-sm text-muted" style={{ textAlign: "right", marginTop: 4 }}>via {paymentMethod === "cash" ? "Cash" : "Card"}</div>
        </div>
      )}

      {/* Rating */}
      {!rated ? (
        <div style={{ marginBottom: 16 }}>
          <p className="text-body-md text-muted" style={{ marginBottom: 8 }}>Rate your driver</p>
          <div style={{ display: "flex", justifyContent: "center", gap: 8 }}>
            {[1, 2, 3, 4, 5].map(s => (
              <button key={s} onClick={() => { setRating(s); setRated(true); }}
                style={{ fontSize: 36, background: "none", border: "none", cursor: "pointer", transition: "transform 0.2s" }}
                onMouseEnter={e => (e.currentTarget.style.transform = "scale(1.2)")}
                onMouseLeave={e => (e.currentTarget.style.transform = "scale(1)")}
              >{s <= rating ? "⭐" : "☆"}</button>
            ))}
          </div>
        </div>
      ) : (
        <p className="text-body-md" style={{ color: "#16a34a", fontWeight: 600, marginBottom: 16 }}>⭐ {rating} stars — Thank you!</p>
      )}

      <button onClick={onNewTrip} className="btn btn-primary">New Trip</button>
    </div>
  );
}
