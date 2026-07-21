"use client";
import { getTimelineSteps, RideState } from "@cytaxi/ride-machine";
import { colors } from "@cytaxi/design-tokens";
import { motion } from "framer-motion";

interface TripTimelineProps {
  state: RideState;
}

export function TripTimeline({ state }: TripTimelineProps) {
  const steps = getTimelineSteps(state);

  if (state === "pickup_select" || state === "input" || state === "confirm") return null;

  return (
    <div
      role="progressbar"
      aria-label="Progreso del viaje"
      style={{
        display: "flex",
        alignItems: "center",
        gap: 0,
        padding: "10px 16px",
        overflow: "hidden",
        background: colors.surface.paper,
        borderBottom: "1px solid rgba(0,0,0,0.04)",
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        zIndex: 600,
      }}
    >
      {steps.map((step, i) => (
        <div
          key={step.id}
          style={{
            display: "flex",
            alignItems: "center",
            flex: 1,
            minWidth: 0,
          }}
        >
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: i * 0.04, type: "spring", stiffness: 350, damping: 20 }}
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              gap: 3,
              flex: 1,
            }}
          >
            <div
              style={{
                width: 26,
                height: 26,
                borderRadius: "50%",
                background:
                  step.status === "completed" || step.status === "active"
                    ? colors.cobalt
                    : "rgba(0,0,0,0.04)",
                border: step.status === "pending"
                  ? "1.5px solid rgba(0,0,0,0.1)"
                  : "none",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: 11,
                fontWeight: 700,
                color: step.status === "completed" || step.status === "active"
                  ? "#fff"
                  : colors.textMuted,
                transition: "all 0.3s cubic-bezier(0.4,0,0.2,1)",
                position: "relative",
                boxShadow: step.status === "active"
                  ? `0 0 0 4px ${colors.cobaltBg}`
                  : "none",
              }}
            >
              {step.status === "completed" ? (
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M20 6L9 17l-5-5" />
                </svg>
              ) : step.status === "active" ? (
                <div style={{ width: 7, height: 7, borderRadius: "50%", background: "#fff", animation: "dotPulse 1.5s infinite" }} />
              ) : (
                <div style={{ width: 5, height: 5, borderRadius: "50%", background: "rgba(0,0,0,0.15)" }} />
              )}
            </div>
            <span
              style={{
                fontSize: 8,
                fontWeight: step.status === "active" ? 700 : 500,
                color:
                  step.status === "completed"
                    ? colors.textMuted
                    : step.status === "active"
                    ? colors.cobalt
                    : colors.textMuted,
                textAlign: "center",
                lineHeight: 1.1,
                maxWidth: 56,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
                transition: "color 0.3s",
              }}
            >
              {step.label}
            </span>
          </motion.div>
          {i < steps.length - 1 && (
            <motion.div
              initial={{ background: "rgba(0,0,0,0.06)" }}
              animate={{
                background: step.status === "completed"
                  ? colors.cobalt
                  : "rgba(0,0,0,0.06)",
              }}
              transition={{ duration: 0.5, ease: [0.4, 0, 0.2, 1] }}
              style={{
                flex: 1,
                height: 2,
                margin: "0 2px",
                marginBottom: 18,
                borderRadius: 1,
              }}
            />
          )}
        </div>
      ))}
    </div>
  );
}
