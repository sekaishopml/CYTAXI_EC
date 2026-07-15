"use client";
import { useEffect, useRef, useState } from "react";
import dynamic from "next/dynamic";
import { gsap } from "@/services/gsap";
import { useTripFlow } from "@/hooks/useTripFlow";
import { PickUpStep } from "@/components/states/PickUpStep";
import { FormState } from "@/components/states/FormState";
import { ConfirmState } from "@/components/states/ConfirmState";
import { TrackingState } from "@/components/states/TrackingState";
import { CompletedState } from "@/components/states/CompletedState";

const MapView = dynamic(() => import("@/components/MapView").then(m => ({ default: m.MapView })), { ssr: false });

export default function HomePage() {
  const flow = useTripFlow();
  const prevState = useRef(flow.state);
  const [keyboardOpen, setKeyboardOpen] = useState(false);
  const [keyboardHeight, setKeyboardHeight] = useState(0);

  // Detectar teclado movil
  useEffect(() => {
    const vv = window.visualViewport;
    if (!vv) return;
    const handleResize = () => {
      const diff = window.innerHeight - vv.height;
      if (diff > 100) {
        setKeyboardOpen(true);
        setKeyboardHeight(diff);
      } else {
        setKeyboardOpen(false);
        setKeyboardHeight(0);
      }
    };
    vv.addEventListener("resize", handleResize);
    return () => vv.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const content = flow.contentRef.current;
    if (!content || prevState.current === flow.state) return;

    if (flow.state === "searching" || flow.state === "driver_found" || flow.state === "in_progress") {
      gsap.fromTo(content, { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.5, ease: "back.out(1.2)" });
    } else if (flow.state === "completed") {
      gsap.fromTo(content, { opacity: 0, scale: 0.95 }, { opacity: 1, scale: 1, duration: 0.5, ease: "back.out(1.2)" });
    } else if (content) {
      gsap.fromTo(content, { opacity: 0, y: 15 }, { opacity: 1, y: 0, duration: 0.3 });
    }
    prevState.current = flow.state;
  }, [flow.state]);

  useEffect(() => {
    gsap.fromTo(flow.sheetRef.current, { opacity: 0, y: 40 }, { opacity: 1, y: 0, duration: 0.5, ease: "power3.out", delay: 0.2 });
  }, []);

  const showPin = flow.state === "pickup_select";
  const mapInteractive = flow.state === "pickup_select";
  const showMap = flow.state !== "completed";
  const mapH = keyboardOpen ? "30dvh" : "50dvh";

  return (
    <div style={{ width: "100vw", height: "100dvh", overflow: "hidden", display: "flex", flexDirection: "column", position: "fixed", inset: 0, touchAction: "none" }}>
      {/* Map */}
      {showMap && (
        <div style={{ width: "100%", height: mapH, position: "relative", flexShrink: 0, touchAction: mapInteractive ? "auto" : "none" }}>
          <MapView onCenterChange={flow.handleCenterChange} showPin={showPin} interactive={mapInteractive} />
          <div style={{ position: "absolute", bottom: 0, left: 0, right: 0, height: 32, background: "linear-gradient(transparent, rgba(248,249,251,0.95))" }} />
        </div>
      )}

      {/* Bottom panel */}
      <div ref={flow.sheetRef} style={{ width: "100%", flex: 1, overflow: "hidden", background: "var(--uk-background)", display: "flex", flexDirection: "column", touchAction: "auto", WebkitOverflowScrolling: "touch" }}>
        <div ref={flow.contentRef} style={{ flex: 1, overflowY: "auto", overflowX: "hidden" }}>
          {flow.state === "pickup_select" && <PickUpStep {...flow.pickupStepProps} />}
          {flow.state === "input" && <FormState {...flow.destStepProps} />}
          {flow.state === "confirm" && <ConfirmState {...flow.confirmProps} />}
          {["searching", "driver_found", "in_progress"].includes(flow.state) && <TrackingState {...flow.trackingProps} />}
          {flow.state === "completed" && <CompletedState {...flow.completedProps} />}
        </div>
      </div>

      {/* Bottom nav — visible solo si no hay teclado */}
      {!keyboardOpen && (
        <nav style={{
          position: "fixed", bottom: 0, width: "100%", zIndex: 50,
          background: "rgba(255,255,255,0.9)", backdropFilter: "blur(16px)",
          boxShadow: "0px -4px 20px rgba(0,0,0,0.05)",
          display: "flex", justifyContent: "space-around", alignItems: "center", padding: "8px 0 12px",
        }}>
          {[
            { icon: "home", label: "Inicio", active: ["pickup_select", "input", "confirm"].includes(flow.state) },
            { icon: "analytics", label: "Actividad", active: false },
            { icon: "directions_car", label: "Viajes", active: ["searching", "driver_found", "in_progress", "completed"].includes(flow.state) },
            { icon: "account_balance_wallet", label: "Billetera", active: false },
            { icon: "person", label: "Perfil", active: false },
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
              <span style={{ marginTop: 2 }}>{tab.label}</span>
            </a>
          ))}
        </nav>
      )}
    </div>
  );
}
