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
      <div style={{ fontSize: 44, marginBottom: 6 }}>✅</div>
      <p style={{ fontSize: 20, fontWeight: 600, marginBottom: 14 }}>Viaje completado</p>

      {pickup && dest && (
        <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 12, marginBottom: 14, textAlign: "left" }}>
          <div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 4 }}>
            <span className="material-symbols-outlined" style={{ fontSize: 16, color: "#006c49" }}>location_on</span>
            <p style={{ fontSize: 13, margin: 0 }}>{pickup.name}</p>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
            <span className="material-symbols-outlined" style={{ fontSize: 16, color: "#276ef1" }}>trip</span>
            <p style={{ fontSize: 13, margin: 0 }}>{dest.name}</p>
          </div>
        </div>
      )}

      {fare && (
        <div style={{ background: "#f6f6f6", borderRadius: 14, padding: 14, marginBottom: 14, textAlign: "left" }}>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}><span style={{ fontSize: 13, color: "#3c4a42" }}>Distancia</span><span style={{ fontSize: 13 }}>{fare.distance_km.toFixed(1)} km</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}><span style={{ fontSize: 13, color: "#3c4a42" }}>Tiempo</span><span style={{ fontSize: 13 }}>{fare.eta_minutes} min</span></div>
          <div style={{ height: 1, background: "#d9dadc", margin: "6px 0" }} />
          <div style={{ display: "flex", justifyContent: "space-between" }}><span style={{ fontSize: 18, fontWeight: 700 }}>Total</span><span style={{ fontSize: 18, fontWeight: 700, color: "#006c49" }}>${fare.total.toFixed(2)}</span></div>
          <p style={{ fontSize: 11, color: "#3c4a42", textAlign: "right", margin: "2px 0 0" }}>via {paymentMethod === "cash" ? "Efectivo" : "Tarjeta"}</p>
        </div>
      )}

      {!rated ? (
        <div style={{ marginBottom: 14 }}>
          <p style={{ fontSize: 13, color: "#3c4a42", marginBottom: 6 }}>Califica a tu conductor</p>
          <div style={{ display: "flex", justifyContent: "center", gap: 6 }}>
            {[1, 2, 3, 4, 5].map(s => (
              <button key={s} onClick={() => { setRating(s); setRated(true); }}
                style={{ fontSize: 28, background: "none", border: "none", cursor: "pointer", transition: "transform 0.2s" }}
              >{s <= rating ? "⭐" : "☆"}</button>
            ))}
          </div>
        </div>
      ) : (
        <p style={{ fontSize: 14, color: "#16a34a", fontWeight: 600, marginBottom: 14 }}>⭐ {rating} estrellas — ¡Gracias!</p>
      )}

      <button onClick={onNewTrip} style={{ width: "100%", height: 44, background: "#006c49", color: "#fff", borderRadius: 9999, fontSize: 15, fontWeight: 600, fontFamily: "Inter", border: "none", cursor: "pointer", boxShadow: "0px 4px 12px rgba(0,108,73,0.3)" }}>
        Nuevo viaje
      </button>
    </div>
  );
}
