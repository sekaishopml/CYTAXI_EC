"use client";
import { AnimatePresence, motion } from "framer-motion";
import { colors } from "@/theme";

interface LocationOverlayProps {
  showPin: boolean;
  state: string;
  centerAddr: string;
  isDragging: boolean;
  onGPS: () => void;
}

export function LocationOverlay({ showPin, state, centerAddr, isDragging, onGPS }: LocationOverlayProps) {
  return (
    <>
      <AnimatePresence>
        {showPin && (
          <>
            <div style={{
              position: "absolute", top: "50%", left: "50%",
              zIndex: 700, pointerEvents: "none",
              transform: "translate(-50%, -100%)",
              filter: "drop-shadow(0 4px 8px rgba(0,0,0,0.3))",
            }}>
              <motion.div
                key={`pin-${state}`}
                initial={{ opacity: 0, scale: 0.3, y: 10 }}
                animate={{
                  opacity: 1,
                  scale: isDragging ? 1.15 : 1,
                  y: isDragging ? -5 : 0,
                }}
                exit={{ opacity: 0, scale: 0.3, y: 10 }}
                transition={{ type: "spring", stiffness: 480, damping: 22, mass: 0.6 }}
              >
                <svg width="32" height="44" viewBox="0 0 34 44" fill="none">
                  <defs>
                    {state === "input" ? (
                      <linearGradient id="pinGradB" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                        <stop offset="0%" stopColor="#93c5fd" /><stop offset="50%" stopColor="#60a5fa" /><stop offset="100%" stopColor="#3b82f6" />
                      </linearGradient>
                    ) : (
                      <linearGradient id="pinGradA" x1="17" y1="0" x2="17" y2="44" gradientUnits="userSpaceOnUse">
                        <stop offset="0%" stopColor="#93c5fd" /><stop offset="50%" stopColor="#3b82f6" /><stop offset="100%" stopColor="#2563eb" />
                      </linearGradient>
                    )}
                  </defs>
                  <path d="M17 0 C26 0 34 8 34 16 C34 26 26 33 20 41 C18.5 42.5 15.5 42.5 14 41 C8 33 0 26 0 16 C0 8 8 0 17 0 Z"
                    fill={state === "input" ? "url(#pinGradB)" : "url(#pinGradA)"} />
                  <circle cx="17" cy="15" r="7" fill="white" opacity="0.93" />
                  <circle cx="17" cy="15" r="3.5" fill={state === "input" ? "#3b82f6" : "#2563eb"} />
                </svg>
              </motion.div>
            </div>

            <motion.div
              key={`label-${state}`}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -6, scale: 0.95 }}
              transition={{ type: "spring", stiffness: 300, damping: 25, delay: 0.1 }}
              style={{
                position: "absolute",
                top: "calc(50% - 82px)",
                left: "50%",
                transform: "translateX(-50%)",
                zIndex: 710,
                background: "rgba(255,255,255,0.92)",
                backdropFilter: "blur(16px) saturate(180%)",
                WebkitBackdropFilter: "blur(16px) saturate(180%)",
                borderRadius: 10,
                padding: "6px 14px",
                boxShadow: "0 2px 12px rgba(0,0,0,0.1), 0 0 0 1px rgba(0,0,0,0.04)",
                fontSize: 12, fontWeight: 500,
                fontFamily: "'Inter', sans-serif",
                color: "#121212",
                whiteSpace: "nowrap",
                maxWidth: 220,
                overflow: "hidden",
                textOverflow: "ellipsis",
                pointerEvents: "none",
                lineHeight: "1.3",
              }}
            >
              {state === "pickup_select" ? (
                <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
                  <span style={{
                    width: 16, height: 16, borderRadius: "50%",
                    background: "#2563eb", color: "#fff", fontSize: 9, fontWeight: 700,
                    display: "inline-flex", alignItems: "center", justifyContent: "center",
                    flexShrink: 0,
                  }}>A</span>
                  {centerAddr}
                </span>
              ) : (
                <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
                  <span style={{
                    width: 16, height: 16, borderRadius: "50%",
                    background: "#3b82f6", color: "#fff", fontSize: 9, fontWeight: 700,
                    display: "inline-flex", alignItems: "center", justifyContent: "center",
                    flexShrink: 0,
                  }}>B</span>
                  {centerAddr}
                </span>
              )}
            </motion.div>
          </>
        )}
      </AnimatePresence>

      <button type="button" onClick={onGPS} aria-label="Mi ubicación"
        style={{
          position: "absolute", right: 12, bottom: 14, zIndex: 690,
          width: 36, height: 36, borderRadius: 10,
          background: "rgba(255,255,255,0.95)", border: "none",
          boxShadow: "0 2px 8px rgba(0,0,0,0.12)", cursor: "pointer",
          display: "flex", alignItems: "center", justifyContent: "center",
          color: colors.cobalt,
        }}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
          <circle cx="12" cy="12" r="3" /><path d="M12 2v4m0 12v4m10-10h-4M6 12H2" />
        </svg>
      </button>
    </>
  );
}
