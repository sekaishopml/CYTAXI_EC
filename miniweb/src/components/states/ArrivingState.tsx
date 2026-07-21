"use client";
import { useState, useEffect } from "react";
import type { DriverInfo, RoutePayload } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

interface ArrivingStateProps {
  driver: DriverInfo | null;
  eta: number;
  route: RoutePayload | null;
  arrived: boolean;
  onArrived: () => void;
  onCancel?: () => void;
}

export function ArrivingState({ driver, eta, route, arrived, onArrived, onCancel }: ArrivingStateProps) {
  const [elapsed, setElapsed] = useState(0);

  useEffect(() => {
    if (arrived) return;
    const interval = setInterval(() => setElapsed((p) => p + 1), 1000);
    return () => clearInterval(interval);
  }, [arrived]);

  const remaining = Math.max(0, eta - elapsed);

  return (
    <div style={{ padding: "24px 20px", textAlign: "center" }}>
      <div style={{ position: "relative", width: 80, height: 80, margin: "0 auto 16px" }}>
        {!arrived && (
          <svg width="80" height="80" viewBox="0 0 80 80" style={{ position: "absolute", inset: 0, transform: "rotate(-90deg)" }}>
            <circle cx="40" cy="40" r="36" fill="none" stroke="rgba(0,0,0,0.06)" strokeWidth="4" />
            <circle cx="40" cy="40" r="36" fill="none" stroke={colors.cobalt} strokeWidth="4"
              strokeDasharray={`${2 * Math.PI * 36}`}
              strokeDashoffset={`${2 * Math.PI * 36 * (1 - Math.min(1, elapsed / Math.max(eta, 1)))}`}
              strokeLinecap="round"
              style={{ transition: "stroke-dashoffset 1s cubic-bezier(0.4, 0, 0.2, 1)" }} />
          </svg>
        )}
        <div style={{
          position: "absolute", inset: 8, borderRadius: "50%",
          background: colors.cobaltBg,
          display: "flex", alignItems: "center", justifyContent: "center",
        }}>
          <span style={{ fontSize: 28 }}>{arrived ? "✅" : "🚗"}</span>
        </div>
      </div>

      <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: "0 0 4px", letterSpacing: "-0.01em" }}>
        {arrived ? "El conductor ha llegado" : "Conductor en camino"}
      </p>
      <p style={{ fontSize: 13, color: colors.textMuted, margin: 0 }}>
        {arrived
          ? "Te está esperando, sal a la calle"
          : `Llega en ${Math.ceil(remaining / 60)} min • ${route?.distance_km?.toFixed(1) || "0.5"} km`}
      </p>

      {driver && (
        <div style={{
          marginTop: 16, background: colors.surface.paperLight, borderRadius: 10,
          boxShadow: shadows.card, border: "1px solid rgba(0,0,0,0.06)",
          padding: "14px 16px", display: "flex", alignItems: "center", gap: 12,
        }}>
          <div style={{
            width: 48, height: 48, borderRadius: "50%",
            background: "linear-gradient(135deg, #dbeafe, #bfdbfe)",
            display: "flex", alignItems: "center", justifyContent: "center",
            flexShrink: 0, fontSize: 20, border: "2px solid #fff", boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
          }}>
            {driver.photo ? (
              <img src={driver.photo} alt="" style={{ width: "100%", height: "100%", borderRadius: "50%", objectFit: "cover" }} />
            ) : "👤"}
          </div>
          <div style={{ flex: 1, textAlign: "left" }}>
            <p style={{ fontSize: 15, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>{driver.name}</p>
            <p style={{ fontSize: 12, color: colors.textMuted, margin: "2px 0 0" }}>{driver.vehicle} • {driver.plate}</p>
          </div>
          <span style={{ fontSize: 13, fontWeight: 600, color: colors.cobalt, display: "flex", alignItems: "center", gap: 3 }}>
            ⭐ {driver.rating.toFixed(1)}
          </span>
        </div>
      )}

      <div style={{ display: "flex", gap: 10, marginTop: 16 }}>
        {!arrived && (
          <button type="button" onClick={onArrived} aria-label="Confirmar llegada del conductor"
            style={{
              flex: 1, height: 52, background: colors.cobalt, color: "#fff",
              borderRadius: 14, fontSize: 15, fontWeight: 600, border: "none",
              cursor: "pointer", fontFamily: "Inter, sans-serif",
              boxShadow: "0 4px 20px rgba(59,130,246,0.25)",
              transition: "all 0.2s",
            }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.01)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}>
            El conductor llegó
          </button>
        )}
        {arrived && (
          <button type="button" onClick={onArrived} aria-label="Iniciar viaje"
            style={{
              flex: 1, height: 52, background: colors.cobalt, color: "#fff",
              borderRadius: 14, fontSize: 15, fontWeight: 600, border: "none",
              cursor: "pointer", fontFamily: "Inter, sans-serif",
              boxShadow: "0 4px 20px rgba(59,130,246,0.25)",
              transition: "all 0.2s",
            }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.01)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}>
            Iniciar viaje
          </button>
        )}
        {onCancel && (
          <button type="button" onClick={onCancel} aria-label="Cancelar viaje"
            style={{
              height: 52, padding: "0 20px", background: "transparent",
              color: colors.danger, border: `1px solid ${colors.dangerBg}`,
              borderRadius: 14, fontSize: 14, fontWeight: 600,
              cursor: "pointer", fontFamily: "Inter, sans-serif",
              transition: "all 0.15s",
            }}
            onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(186,26,26,0.04)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}>
            Cancelar
          </button>
        )}
      </div>
    </div>
  );
}
