"use client";
import { TripState, DriverInfo, Place } from "@/types";
import { colors, radius, shadows } from "@cytaxi/design-tokens";

interface TrackingStateProps {
  state: TripState;
  driver: DriverInfo | null;
  eta: number;
  route: { distance_km: number; eta_minutes: number } | null;
  paymentMethod: "cash" | "card";
  pickup: Place | null;
  dest: Place | null;
  onCancel?: () => void;
  onRejectDriver?: () => void;
  noDrivers?: boolean;
  onRetry?: () => void;
}

export function TrackingState({ state, driver, eta, route, paymentMethod, onCancel, onRejectDriver, noDrivers, onRetry }: TrackingStateProps) {
  if (state === "searching") {
    if (noDrivers) {
      return (
        <div style={{ padding: "24px 20px", textAlign: "center" }}>
          <div style={{
            width: 72, height: 72, borderRadius: "50%",
            background: "rgba(186,26,26,0.06)",
            display: "flex", alignItems: "center", justifyContent: "center",
            margin: "0 auto 16px",
          }}>
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke={colors.danger} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10"/><path d="M12 8v4m0 4h.01"/>
            </svg>
          </div>
          <p style={{ fontSize: 18, fontWeight: 700, marginBottom: 4, color: colors.textPrimary, letterSpacing: "-0.01em" }}>
            No encontramos conductores
          </p>
          <p style={{ fontSize: 13, color: colors.textMuted, margin: "0 0 20px" }}>
            No hay conductores disponibles cerca. Intenta de nuevo en unos minutos.
          </p>
          <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
            {onRetry && (
              <button type="button" onClick={onRetry} aria-label="Reintentar búsqueda"
                style={{
                  width: "100%", height: 52, background: colors.green, color: "#fff",
                  borderRadius: 14, fontSize: 15, fontWeight: 600, border: "none",
                  cursor: "pointer", fontFamily: "Inter, sans-serif",
                  boxShadow: `0 4px 20px ${colors.cobaltBg}`,
                  transition: "all 0.2s",
                }}>
                Reintentar
              </button>
            )}
            {onCancel && (
              <button type="button" onClick={onCancel} aria-label="Volver"
                style={{
                  width: "100%", height: 44, background: "transparent",
                  color: colors.textPrimary, border: "1px solid rgba(0,0,0,0.1)",
                  borderRadius: 14, fontSize: 14, fontWeight: 600,
                  cursor: "pointer", fontFamily: "Inter, sans-serif",
                  transition: "all 0.15s",
                }}>
                Volver
              </button>
            )}
          </div>
        </div>
      );
    }
    return (
      <div style={{ padding: "24px 20px", textAlign: "center" }}>
        {/* Animated radar */}
        <div style={{ position: "relative", width: 80, height: 80, margin: "0 auto 20px" }}>
          <div style={{ position: "absolute", inset: 0, borderRadius: "50%", border: `2px solid ${colors.green}`, opacity: 0.15, animation: "radarPulse 2s ease-out infinite" }} />
          <div style={{ position: "absolute", inset: 8, borderRadius: "50%", border: `2px solid ${colors.green}`, opacity: 0.2, animation: "radarPulse 2s ease-out 0.5s infinite" }} />
          <div style={{ position: "absolute", inset: 16, borderRadius: "50%", background: colors.cobaltBg, display: "flex", alignItems: "center", justifyContent: "center" }}>
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke={colors.green} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
            </svg>
          </div>
        </div>
        <p style={{ fontSize: 18, fontWeight: 700, marginBottom: 4, color: colors.textPrimary, letterSpacing: "-0.01em" }}>Buscando conductor</p>
        <p style={{ fontSize: 13, color: colors.textMuted, margin: 0 }}>Conectando con los mejores conductores cercanos</p>

        {/* Shimmer skeleton */}
        <div style={{ display: "flex", flexDirection: "column", gap: 10, marginTop: 20 }}>
          {[1,2,3].map(i => (
            <div key={i} style={{
              display: "flex", alignItems: "center", gap: 12, padding: "14px 16px",
              borderRadius: 10, background: colors.surface.paperLight, border: "1px solid rgba(0,0,0,0.06)",
              overflow: "hidden", position: "relative",
            }}>
              <div style={{
                width: 44, height: 44, borderRadius: "50%",
                background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)",
                backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out infinite",
                animationDelay: `${i * 0.1}s`,
              }} />
              <div style={{ flex: 1 }}>
                <div style={{
                  height: 12, borderRadius: 6, width: `${50 + i * 10}%`, marginBottom: 8,
                  background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)",
                  backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out infinite",
                  animationDelay: `${i * 0.15}s`,
                }} />
                <div style={{
                  height: 10, borderRadius: 5, width: `${70 + i * 5}%`,
                  background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)",
                  backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out infinite",
                  animationDelay: `${i * 0.2}s`,
                }} />
              </div>
            </div>
          ))}
        </div>

        {onCancel && (
          <button type="button" onClick={onCancel} aria-label="Cancelar búsqueda"
            style={{
              marginTop: 16, padding: "12px 28px", background: "transparent",
              color: colors.danger, border: `1px solid ${colors.dangerBg}`,
              borderRadius: 12, fontSize: 14, fontWeight: 600, fontFamily: "Inter, sans-serif",
              cursor: "pointer", transition: "all 0.15s",
            }}
            onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(186,26,26,0.04)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}>
            Cancelar
          </button>
        )}
      </div>
    );
  }

  if (!driver) return null;

  return (
    <div style={{ display: "flex", flexDirection: "column", padding: "16px 20px", gap: 14 }}>
      {/* ETA header */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
        <div>
          <p style={{ fontSize: 20, fontWeight: 700, color: colors.textPrimary, margin: 0, letterSpacing: "-0.02em" }}>
            {eta <= 300 ? "Llegando..." : `${Math.ceil(eta / 60)} min`}
          </p>
          <p style={{ fontSize: 13, color: colors.textMuted, margin: "4px 0 0" }}>
            {route?.distance_km?.toFixed(1) || "0.5"} km de distancia
          </p>
        </div>
        <span style={{
          background: colors.cobaltBg, color: colors.cobalt,
          padding: "5px 12px", borderRadius: 999, fontSize: 11, fontWeight: 600,
          display: "flex", alignItems: "center", gap: 5,
        }}>
          <span style={{ width: 6, height: 6, borderRadius: "50%", background: colors.cobalt, animation: "dotPulse 1.5s infinite" }} />
          En vivo
        </span>
      </div>

      {/* Driver card */}
      <div style={{
        background: colors.surface.paperLight, borderRadius: 10,
        boxShadow: shadows.card, border: "1px solid rgba(0,0,0,0.06)",
        padding: 14, display: "flex", alignItems: "center", gap: 12,
      }}>
        <div style={{ position: "relative", width: 52, height: 52, flexShrink: 0 }}>
          <div style={{
            width: "100%", height: "100%", borderRadius: "50%",
            background: "linear-gradient(135deg, #dbeafe, #bfdbfe)",
            display: "flex", alignItems: "center", justifyContent: "center",
            fontSize: 22, border: "2px solid #fff", boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
          }}>
            🚗
          </div>
          <span style={{
            position: "absolute", bottom: -2, right: -2, background: "#fff",
            borderRadius: 999, padding: "2px 6px", boxShadow: "0 1px 4px rgba(0,0,0,0.1)",
            fontSize: 10, fontWeight: 700, display: "flex", alignItems: "center", gap: 2,
          }}>
            ⭐ {driver.rating.toFixed(1)}
          </span>
        </div>
        <div style={{ flex: 1 }}>
          <p style={{ fontSize: 15, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>{driver.name}</p>
          <p style={{ fontSize: 13, color: colors.textMuted, margin: "2px 0 4px" }}>{driver.vehicle}</p>
          <span style={{
            display: "inline-block", background: "rgba(0,0,0,0.04)", padding: "2px 8px",
            borderRadius: 6, fontSize: 11, fontWeight: 600, letterSpacing: "0.06em",
            color: colors.textPrimary,
          }}>{driver.plate}</span>
        </div>
      </div>

      {/* Timeline */}
      <div style={{ display: "flex", flexDirection: "column", gap: 0, padding: "0 4px" }}>
        <TimelineStep label="Viaje solicitado" status="completed" />
        <TimelineStep label="Conductor asignado" status={state === "driver_found" ? "active" : "completed"} />
        <TimelineStep label="Conductor en camino" status={["arriving", "arrived"].includes(state) ? "active" : ["in_progress", "destination", "payment", "rating", "completed"].includes(state) ? "completed" : "pending"} />
        <TimelineStep label="En viaje" status={state === "in_progress" ? "active" : ["destination", "payment", "rating", "completed"].includes(state) ? "completed" : "pending"} />
        <TimelineStep label="Destino" status={["destination", "payment", "rating", "completed"].includes(state) ? "active" : "pending"} last />
      </div>

      {/* Action buttons */}
      <div style={{ display: "flex", gap: 8 }}>
        <button type="button" aria-label="Llamar al conductor"
          style={{
            flex: 1, height: 44, background: "rgba(255,255,255,0.7)", color: colors.textPrimary,
            fontSize: 13, fontWeight: 600, borderRadius: 12,
            display: "flex", alignItems: "center", justifyContent: "center", gap: 6,
            border: "1px solid rgba(0,0,0,0.08)", cursor: "pointer", fontFamily: "Inter, sans-serif",
            transition: "all 0.15s",
          }}
          onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(0,0,0,0.03)"; }}
          onMouseLeave={(e) => { e.currentTarget.style.background = "rgba(255,255,255,0.7)"; }}>
          📞 Llamar
        </button>
        <button type="button" aria-label="Enviar mensaje al conductor"
          style={{
            flex: 1, height: 44, background: "rgba(255,255,255,0.7)", color: colors.textPrimary,
            fontSize: 13, fontWeight: 600, borderRadius: 12,
            display: "flex", alignItems: "center", justifyContent: "center", gap: 6,
            border: "1px solid rgba(0,0,0,0.08)", cursor: "pointer", fontFamily: "Inter, sans-serif",
            transition: "all 0.15s",
          }}
          onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(0,0,0,0.03)"; }}
          onMouseLeave={(e) => { e.currentTarget.style.background = "rgba(255,255,255,0.7)"; }}>
          💬 Mensaje
        </button>
        {onRejectDriver && (
          <button type="button" onClick={onRejectDriver} aria-label="Rechazar conductor"
            style={{
              width: 44, height: 44, background: colors.dangerBg, color: colors.danger,
              borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center",
              border: "none", cursor: "pointer", flexShrink: 0, fontSize: 16, fontWeight: 600,
              transition: "all 0.15s",
            }}
            onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(186,26,26,0.15)"; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = colors.dangerBg; }}>
            ✕
          </button>
        )}
      </div>

      {/* Payment badge */}
      <div style={{
        background: colors.surface.paperLight, borderRadius: 10,
        padding: "10px 16px", display: "flex", justifyContent: "space-between", alignItems: "center",
        border: "1px solid rgba(0,0,0,0.06)",
      }}>
        <span style={{ fontSize: 13, color: colors.textMuted }}>Pago</span>
        <span style={{ fontSize: 13, fontWeight: 600, color: colors.textPrimary }}>{paymentMethod === "cash" ? "💵 Efectivo" : "💳 Tarjeta"}</span>
      </div>
    </div>
  );
}

function TimelineStep({ label, status, last = false }: { label: string; status: "completed" | "active" | "pending"; last?: boolean }) {
  return (
    <div style={{ display: "flex", alignItems: "flex-start", gap: 12 }}>
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        <div style={{
          width: 20, height: 20, borderRadius: "50%",
          background: status === "completed" ? colors.cobalt : status === "active" ? colors.cobalt : "rgba(0,0,0,0.06)",
          border: status === "pending" ? "1.5px solid rgba(0,0,0,0.1)" : "none",
          display: "flex", alignItems: "center", justifyContent: "center",
          boxShadow: status === "active" ? `0 0 0 4px ${colors.cobaltBg}` : "none",
          transition: "all 0.3s",
        }}>
          {status === "completed" && (
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="3" strokeLinecap="round"><path d="M20 6L9 17l-5-5"/></svg>
          )}
          {status === "active" && (
            <div style={{ width: 8, height: 8, borderRadius: "50%", background: "#fff", animation: "dotPulse 1.5s infinite" }} />
          )}
          {status === "pending" && (
            <div style={{ width: 6, height: 6, borderRadius: "50%", background: "rgba(0,0,0,0.15)" }} />
          )}
        </div>
        {!last && (
          <div style={{
            width: 1.5, height: 24,
            background: status === "completed" ? colors.cobalt : "rgba(0,0,0,0.06)",
            transition: "background 0.3s",
          }} />
        )}
      </div>
      <p style={{
        fontSize: 13, fontWeight: status === "active" ? 600 : 500,
        color: status === "active" ? colors.cobalt : status === "completed" ? colors.textPrimary : colors.textMuted,
        margin: "2px 0 0", transition: "color 0.3s",
      }}>{label}</p>
    </div>
  );
}
