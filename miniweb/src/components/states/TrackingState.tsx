"use client";
import { TripState, DriverInfo, Place } from "@/types";

const G = "#006c49";
const T1 = "#191c1e";
const T2 = "#3c4a42";
const SURFACE = "#edeef0";
const ERR = "#ba1a1a";
const ERR_BG = "#ffdad6";

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
}

export function TrackingState({ state, driver, eta, route, paymentMethod, onCancel, onRejectDriver }: TrackingStateProps) {
  const driverPhoto = "https://lh3.googleusercontent.com/aida-public/AB6AXuCkHhJFLUV3YsxXFylEOTUJy2z4lY_LCg9OoNenlSm_K-ZxKIkS9pQ_fdp981WEoFsla2qlTjop9e_QlOvKVTB_5InZjlT-19WQ3Lud2rbaohDgGg0IGYHSEm_leWW44fU7MKi6axbn51drsGLkfBYn3xsO6BrI0CmJAmuUNy9K_R1-OovQ5pbx9r7C4T_i08qo3ZQrfjBmVFg_UHQiBpZp_qO7JdGdJNdkiDwACi1XOPZD5m-ALdow1g";

  if (state === "searching") {
    return (
      <div style={{ padding: "20px", textAlign: "center" }}>
        <div style={{ padding: "16px 0" }}>
          {/* Pulso animado */}
          <div style={{ position: "relative", width: 56, height: 56, margin: "0 auto 20px" }}>
            <div style={{ position: "absolute", inset: 0, borderRadius: "50%", background: `${G}20`, animation: "pulse 2s cubic-bezier(0.4,0,0.6,1) infinite" }} />
            <div style={{ position: "absolute", inset: 4, borderRadius: "50%", background: `${G}15`, animation: "pulse 2s cubic-bezier(0.4,0,0.6,1) 0.3s infinite" }} />
            <div style={{ position: "absolute", inset: 8, borderRadius: "50%", background: `${G}40`, display: "flex", alignItems: "center", justifyContent: "center" }}>
              <span style={{ fontSize: 20 }}>🔍</span>
            </div>
          </div>
          <p style={{ fontSize: 18, fontWeight: 600, marginBottom: 4, color: T1 }}>Buscando conductor</p>
          <p style={{ fontSize: 13, color: "#8a8a8a" }}>Conectando con los mejores conductores cercanos</p>
        </div>

        {/* Shimmer skeleton cards */}
        <div style={{ display: "flex", flexDirection: "column", gap: 10, marginTop: 8 }}>
          {[1,2,3].map(i => (
            <div key={i} style={{ display: "flex", alignItems: "center", gap: 12, padding: "14px", borderRadius: 12, background: "rgba(255,255,255,0.7)", overflow: "hidden", position: "relative" }}>
              <div style={{ width: 44, height: 44, borderRadius: "50%", background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)", backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out infinite" }} />
              <div style={{ flex: 1 }}>
                <div style={{ height: 12, borderRadius: 6, background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)", backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out infinite", width: "60%", marginBottom: 6 }} />
                <div style={{ height: 10, borderRadius: 5, background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)", backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out 0.2s infinite", width: "80%" }} />
              </div>
              <div style={{ width: 40, height: 10, borderRadius: 5, background: "linear-gradient(90deg, #eee 25%, #f5f5f5 50%, #eee 75%)", backgroundSize: "200% 100%", animation: "shimmer 1.5s ease-in-out 0.1s infinite" }} />
            </div>
          ))}
        </div>

        {onCancel && (
          <button onClick={onCancel}
            style={{ marginTop: 16, padding: "12px 28px", background: "transparent", color: "#ba1a1a", border: "1px solid rgba(186,26,26,0.2)", borderRadius: 12, fontSize: 14, fontWeight: 600, fontFamily: "Inter", cursor: "pointer", transition: "all 0.15s" }}>
            Cancelar
          </button>
        )}
      </div>
    );
  }

  if (!driver) return null;

  return (
    <div style={{ padding: "12px 20px 20px" }}>
      {/* ETA */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: 14 }}>
        <div>
          <p style={{ fontSize: 20, fontWeight: 600, color: T1, margin: 0 }}>
            Llega en <span style={{ color: G }}>{eta <= 5 ? "Llegando..." : `${Math.ceil(eta / 60)} min`}</span>
          </p>
          <p style={{ fontSize: 13, color: T2, margin: "2px 0 0", display: "flex", alignItems: "center", gap: 4 }}>
            <span className="material-symbols-outlined" style={{ fontSize: 14 }}>distance</span>
            {route?.distance_km?.toFixed(1) || "0.5"} km de distancia
          </p>
        </div>
        <span style={{ background: `${G}1A`, color: G, padding: "4px 10px", borderRadius: 9999, fontSize: 11, fontWeight: 600, display: "flex", alignItems: "center", gap: 4 }}>
          <span style={{ width: 6, height: 6, borderRadius: "50%", background: G, animation: "pulse 1.5s infinite" }} />
          En vivo
        </span>
      </div>

      {/* Tarjeta conductor */}
      <div style={{ background: "#f8f9fb", borderRadius: 14, padding: 14, display: "flex", alignItems: "center", gap: 14, border: "1px solid #edeef0", marginBottom: 14 }}>
        <div style={{ position: "relative", width: 60, height: 60, flexShrink: 0 }}>
          <img src={driverPhoto} alt="Conductor" style={{ width: "100%", height: "100%", borderRadius: "50%", objectFit: "cover", border: "2px solid #fff" }} />
          <span style={{ position: "absolute", bottom: -3, right: -3, background: "#fff", borderRadius: "50%", padding: 2, boxShadow: "0 1px 4px rgba(0,0,0,0.1)", fontSize: 10, fontWeight: 700 }}>⭐ {driver.rating.toFixed(1)}</span>
        </div>
        <div style={{ flex: 1 }}>
          <p style={{ fontSize: 17, fontWeight: 600, margin: 0 }}>{driver.name}</p>
          <p style={{ fontSize: 13, color: T2, margin: "1px 0" }}>{driver.vehicle}</p>
          <span style={{ display: "inline-block", background: SURFACE, padding: "2px 8px", borderRadius: 4, fontSize: 11, fontWeight: 600, letterSpacing: "0.08em" }}>{driver.plate}</span>
        </div>
      </div>

      {/* Timeline dinámico */}
      <div style={{ marginBottom: 14 }}>
        <TimelineStep label="Viaje solicitado" active={true} completed={true} />
        <TimelineStep label="Conductor asignado" active={true} completed={state !== "driver_found"} current={state === "driver_found"} />
        <TimelineStep label="Conductor en camino" active={state === "in_progress" || state === "completed"} completed={state === "completed"} current={state === "in_progress"} />
        {state === "in_progress" && <TimelineStep label="En viaje" active={true} current={false} />}
        <TimelineStep label="Destino alcanzado" active={state === "completed"} completed={state === "completed"} />
      </div>

      {/* Botones de acción */}
      <div style={{ display: "flex", gap: 10 }}>
        <button style={{ flex: 1, height: 50, background: SURFACE, color: T1, fontSize: 15, fontWeight: 600, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", gap: 6, border: "1px solid #d9dadc", cursor: "pointer", fontFamily: "Inter" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18 }}>call</span>
          Llamar
        </button>
        <button style={{ flex: 1, height: 50, background: SURFACE, color: T1, fontSize: 15, fontWeight: 600, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", gap: 6, border: "1px solid #d9dadc", cursor: "pointer", fontFamily: "Inter" }}>
          <span className="material-symbols-outlined" style={{ fontSize: 18 }}>chat</span>
          Mensaje
        </button>
        <button onClick={onRejectDriver} style={{ width: 50, height: 50, background: ERR_BG, color: ERR, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", border: "none", cursor: "pointer", flexShrink: 0 }}>
          <span className="material-symbols-outlined" style={{ fontSize: 20 }}>close</span>
        </button>
      </div>

      <div style={{ background: SURFACE, borderRadius: 10, padding: "10px 14px", marginTop: 12, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <span style={{ fontSize: 13, color: T2 }}>Pago</span>
        <span style={{ fontSize: 13, fontWeight: 600 }}>{paymentMethod === "cash" ? "💵 Efectivo" : "💳 Tarjeta"}</span>
      </div>
    </div>
  );
}

function TimelineStep({ label, active, completed, current }: { label: string; active: boolean; completed?: boolean; current?: boolean }) {
  return (
    <div style={{ display: "flex", alignItems: "flex-start", gap: 10, opacity: current || completed ? 1 : 0.45, marginBottom: 12 }}>
      <div style={{ width: 10, height: 10, borderRadius: "50%", background: completed ? G : current ? G : T2, marginTop: 4, flexShrink: 0, boxShadow: current ? `0 0 0 4px ${G}20` : "none", animation: current ? "pulse 1.5s infinite" : "none" }} />
      <p style={{ fontSize: 13, fontWeight: 500, color: current ? G : T1, margin: 0 }}>{label}</p>
    </div>
  );
}
