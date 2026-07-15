"use client";
import { useEffect, useRef } from "react";
import dynamic from "next/dynamic";
import { gsap } from "@/services/gsap";
import { useTripFlow } from "@/hooks/useTripFlow";
import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

const MapView = dynamic(() => import("@/components/MapView").then(m => ({ default: m.MapView })), { ssr: false });

export default function HomePage() {
  const flow = useTripFlow();
  const prevState = useRef(flow.state);

  useEffect(() => {
    const sheet = flow.sheetRef.current;
    const content = flow.contentRef.current;
    if (!sheet || prevState.current === flow.state) return;

    if (flow.state === "searching" || flow.state === "driver_found" || flow.state === "in_progress") {
      gsap.fromTo(content, { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.5, ease: "back.out(1.2)" });
    } else if (flow.state === "completed") {
      gsap.fromTo(content, { opacity: 0, scale: 0.95 }, { opacity: 1, scale: 1, duration: 0.5, ease: "back.out(1.2)" });
    } else {
      if (content) gsap.fromTo(content, { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.3 });
    }
    prevState.current = flow.state;
  }, [flow.state]);

  useEffect(() => {
    gsap.fromTo(flow.sheetRef.current, { opacity: 0, y: 40 }, { opacity: 1, y: 0, duration: 0.6, ease: "power3.out", delay: 0.2 });
  }, []);

  const renderSheet = () => {
    switch (flow.state) {
      case "idle": return <FormState {...flow.formProps} />;
      case "confirm": return <ConfirmState {...flow.confirmProps} />;
      case "searching": case "driver_found": case "in_progress":
        return <TrackingState {...flow.trackingProps} />;
      case "completed": return <CompletedState {...flow.completedProps} />;
      default: return null;
    }
  };

  const showMap = flow.state !== "completed";

  return (
    <div style={{ width: "100vw", height: "100dvh", overflow: "hidden", display: "flex", flexDirection: "column", position: "relative" }}>
      {/* Top: Map (faded background, 50dvh) */}
      {showMap && (
        <div style={{ width: "100%", height: "50dvh", position: "relative", flexShrink: 0 }}>
          <div style={{ position: "absolute", inset: 0, opacity: 0.6 }}>
            <MapView
              pickupCoords={flow.pickup ? { lat: flow.pickup.lat, lng: flow.pickup.lng } : null}
              destCoords={flow.dest ? { lat: flow.dest.lat, lng: flow.dest.lng } : null}
              driverCoords={flow.driver?.lat && flow.driver?.lng ? { lat: flow.driver.lat, lng: flow.driver.lng } : null}
              polyline={flow.route?.polyline || null}
              fitBounds={flow.state === "confirm"}
            />
          </div>
          {/* Gradient overlay at bottom of map */}
          <div style={{ position: "absolute", bottom: 0, left: 0, right: 0, height: 40, background: "linear-gradient(transparent, rgba(248,249,251,0.9))" }} />
        </div>
      )}

      {(flow.state === "driver_found" || flow.state === "in_progress") && flow.eta > 0 && (
        <div style={{ position: "fixed", top: 16, left: 16, zIndex: 8, background: "rgba(0,0,0,0.88)", color: "#fff", padding: "10px 16px", borderRadius: 14, fontSize: 15, fontWeight: 700, fontFamily: "Inter", backdropFilter: "blur(12px)" }}>
          {Math.ceil(flow.eta / 60)} min away
        </div>
      )}

      {/* Bottom panel — sin scroll */}
      <div ref={flow.sheetRef} style={{
        width: "100%", flex: 1,
        overflow: "hidden",
        background: "var(--uk-background)",
        display: "flex", flexDirection: "column",
      }}>
        <div ref={flow.contentRef} style={{ flex: 1, overflow: "hidden" }}>
          {renderSheet()}
        </div>
      </div>

      {/* Bottom nav bar */}
      <nav style={{
        position: "fixed", bottom: 0, width: "100%", zIndex: 50,
        background: "rgba(255,255,255,0.9)", backdropFilter: "blur(16px)",
        boxShadow: "0px -4px 20px rgba(0,0,0,0.05)",
        display: "flex", justifyContent: "space-around", alignItems: "center",
        padding: "8px 0 12px",
      }}>
        {[
          { icon: "home", label: "Home", active: flow.state === "idle" || flow.state === "confirm" },
          { icon: "analytics", label: "Activity", active: false },
          { icon: "directions_car", label: "Trips", active: ["searching", "driver_found", "in_progress", "completed"].includes(flow.state) },
          { icon: "account_balance_wallet", label: "Wallet", active: false },
          { icon: "person", label: "Profile", active: false },
        ].map(tab => (
          <a key={tab.label} href="#" onClick={(e) => e.preventDefault()}
            style={{
              display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center",
              padding: "6px 14px", borderRadius: 12, fontSize: 10, fontWeight: 600, fontFamily: "Inter",
              transition: "all 0.2s",
              ...(tab.active ? { background: "#006c491A", color: "#006c49" } : { color: "#3c4a42" }),
            }}
          >
            <span className="material-symbols-outlined" style={{ fontSize: 22 }}>{tab.icon}</span>
            {tab.label}
          </a>
        ))}
      </nav>
    </div>
  );
}
