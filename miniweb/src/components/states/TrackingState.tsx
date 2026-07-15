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
}

export function TrackingState({ state, driver, eta, route, paymentMethod }: TrackingStateProps) {
  const driverPhoto = "https://lh3.googleusercontent.com/aida-public/AB6AXuCkHhJFLUV3YsxXFylEOTUJy2z4lY_LCg9OoNenlSm_K-ZxKIkS9pQ_fdp981WEoFsla2qlTjop9e_QlOvKVTB_5InZjlT-19WQ3Lud2rbaohDgGg0IGYHSEm_leWW44fU7MKi6axbn51drsGLkfBYn3xsO6BrI0CmJAmuUNy9K_R1-OovQ5pbx9r7C4T_i08qo3ZQrfjBmVFg_UHQiBpZp_qO7JdGdJNdkiDwACi1XOPZD5m-ALdow1g";

  if (state === "searching") {
    return (
      <div style={{ padding: "16px 20px", textAlign: "center" }}>
        <div style={{ padding: "24px 0" }}>
          <div style={{ width: 44, height: 44, border: "4px solid #e7e8ea", borderTopColor: G, borderRadius: "50%", animation: "spin 0.8s linear infinite", margin: "0 auto 16px" }} />
          <p style={{ fontSize: 18, fontWeight: 600, marginBottom: 6 }}>Buscando conductor...</p>
          <p style={{ fontSize: 13, color: T2 }}>Conectando con vehículos cercanos</p>
        </div>
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
            Llega en <span style={{ color: G }}>{Math.ceil(eta / 60)} min</span>
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

      {/* Timeline */}
      <div style={{ marginBottom: 14 }}>
        <TimelineStep icon="check" label="Viaje solicitado" sublabel="14:15" active={true} completed={true} />
        <TimelineStep icon="person" label="Conductor asignado" sublabel="14:16" active={true} completed={true} />
        <TimelineStep icon="dot" label="Conductor llegando" sublabel="Acercándose al origen" active={true} current={true} />
        <TimelineStep icon="location_on" label="Recogida" sublabel="Esperando llegada" active={true} completed={false} />
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
        <button style={{ width: 50, height: 50, background: ERR_BG, color: ERR, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", border: "none", cursor: "pointer", flexShrink: 0 }}>
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

function TimelineStep({ icon, label, sublabel, active, completed, current }: { icon: string; label: string; sublabel: string; active: boolean; completed?: boolean; current?: boolean }) {
  return (
    <div style={{ display: "flex", alignItems: "flex-start", gap: 10, position: "relative", opacity: current || completed ? 1 : 0.5, marginBottom: 14 }}>
      <div style={{ position: "relative", zIndex: 10, width: 22, height: 22, borderRadius: "50%", background: current ? "#fff" : G, border: current ? `2px solid ${G}` : "none", display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0 }}>
        {current ? (
          <div style={{ width: 7, height: 7, borderRadius: "50%", background: G, animation: "pulse 1.5s infinite" }} />
        ) : (
          <span className="material-symbols-outlined" style={{ fontSize: 12, color: "#fff" }}>{icon}</span>
        )}
      </div>
      <div>
        <p style={{ fontSize: 13, fontWeight: 600, color: current ? G : T1, margin: 0 }}>{label}</p>
        <p style={{ fontSize: 11, color: T2, margin: 0 }}>{sublabel}</p>
      </div>
    </div>
  );
}
