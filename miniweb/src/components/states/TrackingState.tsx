"use client";
import { TripState, DriverInfo, Place } from "@/types";

interface TrackingStateProps {
  state: TripState;
  driver: DriverInfo | null;
  eta: number;
  route: { distance_km: number; eta_minutes: number } | null;
  paymentMethod: "cash" | "card";
  pickup: Place | null;
  dest: Place | null;
}

const STITCH_GREEN = "#006c49";
const STITCH_TEXT_PRIMARY = "#191c1e";
const STITCH_TEXT_VARIANT = "#3c4a42";
const STITCH_SURFACE_BRIGHT = "#f8f9fb";
const STITCH_SURFACE = "#edeef0";
const STITCH_ERROR = "#ba1a1a";

export function TrackingState({ state, driver, eta, route, paymentMethod, pickup, dest }: TrackingStateProps) {
  const driverPhoto = "https://lh3.googleusercontent.com/aida-public/AB6AXuCkHhJFLUV3YsxXFylEOTUJy2z4lY_LCg9OoNenlSm_K-ZxKIkS9pQ_fdp981WEoFsla2qlTjop9e_QlOvKVTB_5InZjlT-19WQ3Lud2rbaohDgGg0IGYHSEm_leWW44fU7MKi6axbn51drsGLkfBYn3xsO6BrI0CmJAmuUNy9K_R1-OovQ5pbx9r7C4T_i08qo3ZQrfjBmVFg_UHQiBpZp_qO7JdGdJNdkiDwACi1XOPZD5m-ALdow1g";

  if (state === "searching") {
    return (
      <div style={{ padding: "12px 20px 40px", textAlign: "center" }}>
        <div className="sheet-handle" />
        <div style={{ padding: "32px 0" }}>
          <div style={{ width: 48, height: 48, border: "4px solid #e7e8ea", borderTopColor: STITCH_GREEN, borderRadius: "50%", animation: "spin 0.8s linear infinite", margin: "0 auto 20px" }} />
          <h3 style={{ fontSize: 20, fontWeight: 600, marginBottom: 8 }}>Finding your best driver...</h3>
          <p style={{ fontSize: 14, color: STITCH_TEXT_VARIANT }}>Connecting to nearby vehicles</p>
          <button
            onClick={() => { window.location.reload(); }}
            style={{ marginTop: 16, padding: "10px 24px", background: STITCH_SURFACE, color: STITCH_TEXT_PRIMARY, borderRadius: 9999, fontSize: 12, fontWeight: 600, letterSpacing: "0.05em", border: "none", cursor: "pointer" }}
          >SIMULATE MATCH</button>
        </div>
      </div>
    );
  }

  if (!driver) return null;

  return (
    <div style={{ padding: "12px 20px 24px" }}>
      <div className="sheet-handle" />

      {/* ETA Header */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: 16 }}>
        <div>
          <h2 style={{ fontSize: 24, fontWeight: 600, color: STITCH_TEXT_PRIMARY, marginBottom: 4, letterSpacing: "-0.01em" }}>
            Arriving in <span style={{ color: STITCH_GREEN }}>{Math.ceil(eta / 60)} mins</span>
          </h2>
          <p style={{ fontSize: 14, color: STITCH_TEXT_VARIANT, display: "flex", alignItems: "center", gap: 4 }}>
            <span className="material-symbols-outlined" style={{ fontSize: 16 }}>distance</span>
            {route?.distance_km?.toFixed(1) || "0.5"} km away
          </p>
        </div>
        <div style={{ background: `${STITCH_GREEN}1A`, color: STITCH_GREEN, padding: "6px 12px", borderRadius: 9999, fontSize: 12, fontWeight: 600, display: "flex", alignItems: "center", gap: 6 }}>
          <div style={{ width: 8, height: 8, borderRadius: "50%", background: STITCH_GREEN, animation: "pulse 1.5s infinite" }} />
          Live
        </div>
      </div>

      {/* Driver Card */}
      <div style={{ background: STITCH_SURFACE_BRIGHT, borderRadius: 16, padding: 16, display: "flex", alignItems: "center", gap: 16, border: "1px solid #edeef0", marginBottom: 16 }}>
        <div style={{ position: "relative", width: 64, height: 64, flexShrink: 0 }}>
          <img src={driverPhoto} alt="Driver" style={{ width: "100%", height: "100%", borderRadius: "50%", objectFit: "cover", border: "2px solid #fff" }} />
          <div style={{ position: "absolute", bottom: -4, right: -4, background: "#fff", borderRadius: "50%", padding: 2, boxShadow: "0 1px 4px rgba(0,0,0,0.1)", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 11, fontWeight: 700, color: STITCH_TEXT_PRIMARY }}>
            ⭐ {driver.rating.toFixed(1)}
          </div>
        </div>
        <div style={{ flex: 1 }}>
          <h3 style={{ fontSize: 20, fontWeight: 600, color: STITCH_TEXT_PRIMARY }}>{driver.name}</h3>
          <p style={{ fontSize: 14, color: STITCH_TEXT_VARIANT }}>{driver.vehicle}</p>
          <div style={{ marginTop: 4, display: "inline-block", background: STITCH_SURFACE, padding: "2px 8px", borderRadius: 4, fontSize: 12, fontWeight: 600, letterSpacing: "0.1em", color: STITCH_TEXT_PRIMARY }}>
            {driver.plate}
          </div>
        </div>
      </div>

      {/* Vertical Timeline */}
      <div style={{ padding: "0 8px", marginBottom: 16 }}>
        <div style={{ position: "relative", display: "flex", flexDirection: "column", gap: 20 }}>
          <TimelineStep icon="check" label="Ride Requested" sublabel="14:15" active={true} completed={true} />
          <TimelineStep icon="person" label="Driver Accepted" sublabel="14:16" active={true} completed={true} />
          <TimelineStep icon="dot" label="Driver Arriving" sublabel="Approaching location" active={true} current={true} />
          <TimelineStep icon="location_on" label="Pickup" sublabel="Waiting for arrival" active={true} completed={false} />
        </div>
      </div>

      {/* Action Buttons */}
      <div style={{ display: "flex", gap: 12 }}>
        <button style={{ flex: 1, height: 56, background: STITCH_SURFACE, color: STITCH_TEXT_PRIMARY, fontSize: 16, fontWeight: 600, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", gap: 8, border: "1px solid #d9dadc", cursor: "pointer", boxShadow: "0 1px 3px rgba(0,0,0,0.05)" }}>
          <span className="material-symbols-outlined">call</span>
          Call
        </button>
        <button style={{ flex: 1, height: 56, background: STITCH_SURFACE, color: STITCH_TEXT_PRIMARY, fontSize: 16, fontWeight: 600, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", gap: 8, border: "1px solid #d9dadc", cursor: "pointer", boxShadow: "0 1px 3px rgba(0,0,0,0.05)" }}>
          <span className="material-symbols-outlined">chat</span>
          Message
        </button>
        <button style={{ width: 56, height: 56, background: "#ffdad6", color: STITCH_ERROR, borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", border: "none", cursor: "pointer", boxShadow: "0 1px 3px rgba(0,0,0,0.05)", flexShrink: 0 }}>
          <span className="material-symbols-outlined">close</span>
        </button>
      </div>

      {/* Payment method display */}
      <div style={{ background: STITCH_SURFACE, borderRadius: 12, padding: "12px 16px", marginTop: 16, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <span style={{ fontSize: 14, color: STITCH_TEXT_VARIANT }}>Payment</span>
        <span style={{ fontSize: 14, fontWeight: 600 }}>{paymentMethod === "cash" ? "💵 Cash" : "💳 Card"}</span>
      </div>
    </div>
  );
}

function TimelineStep({ icon, label, sublabel, active, completed, current }: { icon: string; label: string; sublabel: string; active: boolean; completed?: boolean; current?: boolean }) {
  return (
    <div style={{ display: "flex", alignItems: "flex-start", gap: 12, position: "relative", opacity: current || completed ? 1 : 0.5 }}>
      <div style={{ position: "relative", zIndex: 10, width: 24, height: 24, borderRadius: "50%", background: current ? "#fff" : STITCH_GREEN, border: current ? `2px solid ${STITCH_GREEN}` : "none", display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0, boxShadow: "0 1px 3px rgba(0,0,0,0.1)" }}>
        {current ? (
          <div style={{ width: 8, height: 8, borderRadius: "50%", background: STITCH_GREEN, animation: "pulse 1.5s infinite" }} />
        ) : (
          <span className="material-symbols-outlined" style={{ fontSize: 14, color: "#fff" }}>{icon}</span>
        )}
      </div>
      <div style={{ paddingTop: 2 }}>
        <p style={{ fontSize: 14, fontWeight: 600, color: current ? STITCH_GREEN : STITCH_TEXT_PRIMARY }}>{label}</p>
        <p style={{ fontSize: 12, color: STITCH_TEXT_VARIANT }}>{sublabel}</p>
      </div>
    </div>
  );
}
