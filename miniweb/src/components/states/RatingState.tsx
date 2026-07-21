"use client";
import { useState } from "react";
import type { DriverInfo } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

interface RatingStateProps {
  driver: DriverInfo | null;
  onDone: (score: number) => void;
}

export function RatingState({ driver, onDone }: RatingStateProps) {
  const [score, setScore] = useState(0);
  const [hovered, setHovered] = useState(0);
  const [comment, setComment] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = () => {
    if (score === 0) return;
    setSubmitted(true);
    setTimeout(() => onDone(score), 600);
  };

  if (submitted) {
    return (
      <div style={{ padding: "40px 20px", textAlign: "center" }}>
        <div style={{
          width: 72, height: 72, borderRadius: "50%",
          background: colors.cobaltBg,
          display: "flex", alignItems: "center", justifyContent: "center",
          margin: "0 auto 16px",
        }}>
          <span style={{ fontSize: 32 }}>🎉</span>
        </div>
        <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: 0, letterSpacing: "-0.01em" }}>
          ¡Gracias por tu calificación!
        </p>
        <p style={{ fontSize: 13, color: colors.textMuted, marginTop: 4 }}>
          Tu opinión nos ayuda a mejorar
        </p>
      </div>
    );
  }

  return (
    <div style={{ padding: "24px 20px", textAlign: "center" }}>
      <p style={{ fontSize: 18, fontWeight: 700, color: colors.textPrimary, margin: "0 0 16px", letterSpacing: "-0.01em" }}>
        Califica tu viaje
      </p>

      {driver && (
        <div style={{
          display: "flex", alignItems: "center", justifyContent: "center", gap: 12,
          marginBottom: 20,
        }}>
          <div style={{
            width: 48, height: 48, borderRadius: "50%",
            background: "linear-gradient(135deg, #dbeafe, #bfdbfe)",
            display: "flex", alignItems: "center", justifyContent: "center",
            fontSize: 20, border: "2px solid #fff", boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
          }}>
            {driver.photo ? <img src={driver.photo} alt="" style={{ width: "100%", height: "100%", borderRadius: "50%", objectFit: "cover" }} /> : "👤"}
          </div>
          <p style={{ fontSize: 15, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>
            {driver.name}
          </p>
        </div>
      )}

      <div style={{ display: "flex", justifyContent: "center", gap: 8, marginBottom: 20 }}>
        {[1, 2, 3, 4, 5].map((s) => {
          const active = (hovered || score) >= s;
          return (
            <button type="button" key={s} onClick={() => setScore(s)}
              onMouseEnter={() => setHovered(s)}
              onMouseLeave={() => setHovered(0)}
              aria-label={`${s} estrella${s > 1 ? "s" : ""}`}
              style={{
                fontSize: 36, background: "none", border: "none", cursor: "pointer",
                padding: "0 2px", transition: "all 0.2s cubic-bezier(0.34,1.56,0.64,1)",
                transform: active ? "scale(1.15)" : "scale(1)",
                filter: active ? "none" : "grayscale(1) opacity(0.35)",
              }}>
              {active ? "⭐" : "☆"}
            </button>
          );
        })}
      </div>

      {score > 0 && (
        <input aria-label="Comentario sobre el viaje" value={comment}
          onChange={(e) => setComment(e.target.value)}
          placeholder="Cuéntanos cómo fue tu experiencia..."
          style={{
            width: "100%", padding: "12px 14px", borderRadius: 10,
            border: "1px solid rgba(0,0,0,0.08)",
            fontSize: 13, fontFamily: "Inter, sans-serif", outline: "none",
            background: colors.surface.paperLight,
            marginBottom: 16, boxSizing: "border-box" as const,
          }} />
      )}

      <button type="button" onClick={handleSubmit} disabled={score === 0}
        aria-label={score > 0 ? `Enviar ${score} estrellas` : "Selecciona una calificación"}
        style={{
          width: "100%", height: 52,
          background: score > 0 ? colors.cobalt : "rgba(0,0,0,0.06)",
          color: score > 0 ? "#fff" : colors.textMuted,
          borderRadius: 14, fontSize: 16, fontWeight: 600,
          border: "none", cursor: score > 0 ? "pointer" : "default",
          fontFamily: "Inter, sans-serif",
          boxShadow: score > 0 ? "0 4px 20px rgba(59,130,246,0.25)" : "none",
          transition: "all 0.25s",
          opacity: score > 0 ? 1 : 0.6,
        }}>
        {score > 0 ? `Enviar ${score} estrellas` : "Selecciona una calificación"}
      </button>
    </div>
  );
}
