"use client";
import { useState, useEffect } from "react";
import type { FareBreakdown } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

interface PaymentStateProps {
  fare: FareBreakdown | null;
  method: "cash" | "card";
  onDone: () => void;
}

export function PaymentState({ fare, method, onDone }: PaymentStateProps) {
  const [processing, setProcessing] = useState(true);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => setProgress(p => Math.min(p + 2, 100)), 40);
    const timer = setTimeout(() => { setProcessing(false); clearInterval(interval); }, 2000);
    return () => { clearInterval(interval); clearTimeout(timer); };
  }, []);

  if (processing) {
    return (
      <div style={{ padding: "40px 20px", textAlign: "center" }}>
        <div style={{ position: "relative", width: 72, height: 72, margin: "0 auto 20px" }}>
          <svg width="72" height="72" viewBox="0 0 72 72" style={{ position: "absolute", inset: 0, transform: "rotate(-90deg)" }}>
            <circle cx="36" cy="36" r="32" fill="none" stroke="rgba(0,0,0,0.06)" strokeWidth="4" />
            <circle cx="36" cy="36" r="32" fill="none" stroke={colors.cobalt} strokeWidth="4"
              strokeDasharray={`${2 * Math.PI * 32}`}
              strokeDashoffset={`${2 * Math.PI * 32 * (1 - progress / 100)}`}
              strokeLinecap="round"
              style={{ transition: "stroke-dashoffset 0.2s cubic-bezier(0.4, 0, 0.2, 1)" }} />
          </svg>
          <div style={{
            position: "absolute", inset: 0, display: "flex",
            alignItems: "center", justifyContent: "center",
          }}>
            <span style={{ fontSize: 24 }}>💳</span>
          </div>
        </div>
        <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: "0 0 4px", letterSpacing: "-0.01em" }}>
          Procesando pago
        </p>
        <p style={{ fontSize: 13, color: colors.textMuted, margin: 0 }}>
          {method === "card" ? "Transacción segura vía tarjeta" : "Preparando efectivo"}
        </p>
        {fare && (
          <p style={{ fontSize: 32, fontWeight: 700, color: colors.cobalt, margin: "16px 0 0", letterSpacing: "-0.02em" }}>
            ${fare.total.toFixed(2)}
          </p>
        )}
      </div>
    );
  }

  return (
    <div style={{ padding: "24px 20px", textAlign: "center" }}>
      <div style={{
        width: 72, height: 72, borderRadius: "50%",
        background: colors.cobaltBg,
        display: "flex", alignItems: "center", justifyContent: "center",
        margin: "0 auto 16px",
      }}>
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke={colors.cobalt} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M20 6L9 17l-5-5"/>
        </svg>
      </div>

      <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: "0 0 4px", letterSpacing: "-0.01em" }}>
        Pago confirmado
      </p>

      {fare && (
        <div style={{
          marginTop: 12, padding: "14px 16px",
          background: colors.surface.paperLight, borderRadius: 10,
          border: "1px solid rgba(0,0,0,0.06)",
        }}>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
            <span style={{ fontSize: 13, color: colors.textSecondary }}>Total</span>
            <span style={{ fontSize: 22, fontWeight: 700, color: colors.cobalt }}>${fare.total.toFixed(2)}</span>
          </div>
          <p style={{ fontSize: 11, color: colors.textMuted, margin: "4px 0 0" }}>
            {method === "cash" ? "💵 Efectivo" : "💳 Tarjeta"} • {fare.distance_km.toFixed(1)} km
          </p>
        </div>
      )}

      <button type="button" onClick={onDone} aria-label="Calificar viaje"
        style={{
          width: "100%", marginTop: 16, height: 52,
          background: colors.cobalt, color: "#fff",
          borderRadius: 14, fontSize: 16, fontWeight: 600,
          border: "none", cursor: "pointer", fontFamily: "Inter, sans-serif",
          boxShadow: "0 4px 20px rgba(59,130,246,0.25)",
          transition: "all 0.2s",
        }}
        onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.01)"; }}
        onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}>
        Calificar viaje
      </button>
    </div>
  );
}
