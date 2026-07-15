"use client";
import { TripState, DriverInfo, TrackingUpdate, Place } from "@/types";

interface TrackingStateProps {
  state: TripState;
  driver: DriverInfo | null;
  eta: number;
  route: { distance_km: number; eta_minutes: number } | null;
  tracking: TrackingUpdate | null;
  onStart: () => void;
  paymentMethod: "cash" | "card";
  pickup: Place | null;
  dest: Place | null;
}

export function TrackingState({ state, driver, eta, route, tracking, onStart, paymentMethod }: TrackingStateProps) {
  const driverPhoto = `data:image/svg+xml,${encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="#f6f6f6"/><circle cx="50" cy="35" r="15" fill="#d1d5db"/><ellipse cx="50" cy="80" rx="25" ry="20" fill="#d1d5db"/><circle cx="37" cy="30" r="2" fill="#6b7280"/><circle cx="63" cy="30" r="2" fill="#6b7280"/><path d="M43 40 Q50 48 57 40" stroke="#6b7280" stroke-width="2" fill="none"/><rect x="15" y="90" width="70" height="10" rx="5" fill="#e0e0e0"/></svg>`
  )}`;

  if (state === "searching") {
    return (
      <div style={{ padding: "12px 20px 40px", textAlign: "center" }}>
        <div className="sheet-handle" />
        <div style={{ padding: "32px 0" }}>
          <div className="spinner" style={{ margin: "0 auto 20px" }} />
          <h3 className="text-headline-mobile" style={{ marginBottom: 8 }}>Searching for drivers</h3>
          <p className="text-body-md text-muted">Finding the best available driver near you</p>
        </div>
      </div>
    );
  }

  if (!driver) return null;

  return (
    <div style={{ padding: "12px 20px 24px" }}>
      <div className="sheet-handle" />

      {/* Status */}
      <div style={{ textAlign: "center", marginBottom: 16 }}>
        <span className={`pill ${state === "driver_found" ? "pill-active" : "pill-info"}`}>
          {state === "driver_found" ? "Driver on the way" : "Trip in progress"}
        </span>
      </div>

      {/* Driver info card */}
      <div style={{ display: "flex", alignItems: "center", gap: 14, background: "var(--uk-input-bg)", borderRadius: 16, padding: 16, marginBottom: 16 }}>
        <div className="driver-avatar" style={{ width: 56, height: 56 }}>
          <img src={driverPhoto} alt="Driver" style={{ width: "100%", height: "100%", borderRadius: "50%" }} />
        </div>
        <div style={{ flex: 1 }}>
          <p className="text-body-lg" style={{ fontWeight: 600 }}>{driver.name}</p>
          <p className="text-label-sm text-muted">⭐ {driver.rating.toFixed(1)} · {driver.vehicle}</p>
          <p className="text-label-sm text-muted" style={{ fontFamily: "monospace", letterSpacing: 1 }}>{driver.plate}</p>
        </div>
      </div>

      {/* ETA & Distance */}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 16 }}>
        <div style={{ background: "var(--uk-input-bg)", borderRadius: 14, padding: 14, textAlign: "center" }}>
          <p style={{ fontSize: 28, fontWeight: 700, color: "#276ef1" }}>{Math.ceil(eta / 60)} min</p>
          <p className="text-label-sm text-muted">Arrival</p>
        </div>
        <div style={{ background: "var(--uk-input-bg)", borderRadius: 14, padding: 14, textAlign: "center" }}>
          <p style={{ fontSize: 28, fontWeight: 700 }}>{route?.distance_km?.toFixed(1) || "—"} km</p>
          <p className="text-label-sm text-muted">Distance</p>
        </div>
      </div>

      {/* Payment */}
      <div style={{ background: "var(--uk-input-bg)", borderRadius: 14, padding: "12px 16px", marginBottom: 16, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <span className="text-body-md text-muted">Payment</span>
        <span className="text-label-md">{paymentMethod === "cash" ? "💵 Cash" : "💳 Card"}</span>
      </div>

      {/* Actions */}
      {state === "driver_found" && (
        <button onClick={onStart} className="btn btn-primary">Start Trip</button>
      )}
      <button className="btn btn-tertiary" style={{ marginTop: 8, color: "var(--uk-error)", width: "100%", padding: 14, background: "transparent", border: "none", fontSize: 15, fontWeight: 600, cursor: "pointer", borderRadius: 14 }}>
        Cancel Trip
      </button>
    </div>
  );
}
