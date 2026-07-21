"use client";
import { ReactNode, RefObject, useRef, useState, useLayoutEffect, useCallback, useMemo } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { RideState, stateAnimations } from "@cytaxi/ride-machine";
import { colors } from "@cytaxi/design-tokens";

interface BottomSheetProps {
  state: RideState;
  prevState: RideState | null;
  direction: "forward" | "back";
  isTransitioning: boolean;
  sheetRef: RefObject<HTMLDivElement>;
  showNavbar: boolean;
  keyboardOpen: boolean;
  keyboardH?: number;
  children: ReactNode;
}

const HANDLE_H = 16;
const NAVBAR_H = 72;

const sheetVariants = {
  enter: (dir: string) => ({
    opacity: 0,
    y: dir === "back" ? -12 : 22,
    scale: 0.99,
  }),
  center: { opacity: 1, y: 0, scale: 1 },
  exit: (dir: string) => ({
    opacity: 0,
    y: dir === "back" ? 22 : -12,
    scale: 0.99,
  }),
};

const NOISE_SVG = `url("data:image/svg+xml,%3Csvg viewBox='0 0 512 512' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.06'/%3E%3C/svg%3E")`;

function randDur(base: number) {
  return base + Math.random() * 1.5;
}

function Blob({ index, state }: { index: number; state: RideState }) {
  const isInput = state === "input";
  const isConfirm = state === "confirm";
  const isActive = ["searching", "driver_found", "arriving", "arrived", "in_progress", "destination"].includes(state);
  const isPickup = state === "pickup_select";
  const isHome = state === "travel_home";

  const d = useMemo(() => {
    const configs = [
      { x: -30, y: 20, s: 95 },
      { x: -15, y: 55, s: 90 },
      { x: 0,   y: 15, s: 85 },
      { x: 15,  y: 40, s: 80 },
      { x: 30,  y: 10, s: 85 },
      { x: 40,  y: 50, s: 90 },
      { x: 50,  y: 25, s: 95 },
    ];
    const c = configs[index % 7];
    return {
      x: c.x, y: c.y, s: c.s,
      dur: [6 + index * 0.5, 5 + index * 0.4, 7 + index * 0.4, 6 + index * 0.5],
      drift: [25 + index * 4, -22 - index * 3, 20 + index * 3, -18 - index * 4],
    };
  }, [index]);

  const useMixed = isConfirm || isActive;
  const isBlue = index % 2 === 0;
  const color = useMixed
    ? (isBlue ? "#2563eb" : "#9333ea")
    : isHome
      ? "#3b82f6"
      : (isInput ? "#9333ea" : "#2563eb");
  const color2 = useMixed
    ? (isBlue ? "#60a5fa" : "#c084fc")
    : isHome
      ? "#93c5fd"
      : (isInput ? "#c084fc" : "#60a5fa");
  const opacity = isPickup ? 0.7 : isInput ? 0.75 : useMixed ? 0.8 : isHome ? 0.15 : 0.25;

  return (
    <motion.div
      style={{
        position: "absolute",
        width: `${d.s}%`,
        height: `${d.s}%`,
        borderRadius: "50%",
        filter: "blur(45px)",
        mixBlendMode: "screen",
        willChange: "transform, background",
        pointerEvents: "none",
        top: `${d.y}%`,
        left: `${d.x}%`,
      }}
      animate={{
        x: [0, d.drift[0], d.drift[1], d.drift[2], d.drift[3], 0],
        y: [0, d.drift[3], d.drift[0], d.drift[2], d.drift[1], 0],
        scale: [1, 1.15, 0.9, 1.1, 0.95, 1],
        background: [
          `radial-gradient(circle, ${color}dd, ${color2}44 50%, transparent 70%)`,
          `radial-gradient(circle, ${color2}cc, ${color}55 50%, transparent 70%)`,
          `radial-gradient(circle, ${color}bb, ${color2}33 50%, transparent 70%)`,
          `radial-gradient(circle, ${color2}dd, ${color}44 50%, transparent 70%)`,
          `radial-gradient(circle, ${color}cc, ${color2}55 50%, transparent 70%)`,
          `radial-gradient(circle, ${color}dd, ${color2}44 50%, transparent 70%)`,
        ],
      }}
      transition={{
        x: { duration: d.dur[0], repeat: Infinity, ease: [0.45, 0.05, 0.55, 0.95] },
        y: { duration: d.dur[1], repeat: Infinity, ease: [0.45, 0.05, 0.55, 0.95] },
        scale: { duration: d.dur[2], repeat: Infinity, ease: [0.45, 0.05, 0.55, 0.95] },
        background: { duration: d.dur[3], repeat: Infinity, ease: [0.45, 0.05, 0.55, 0.95] },
      }}
    />
  );
}

export function BottomSheet({
  state, prevState, direction, isTransitioning,
  sheetRef, showNavbar, keyboardOpen, keyboardH = 0, children,
}: BottomSheetProps) {
  const anim = stateAnimations[state];
  const contentRef = useRef<HTMLDivElement>(null);
  const [measuredH, setMeasuredH] = useState(0);

  const measure = useCallback(() => {
    const el = contentRef.current;
    if (!el) return;
    requestAnimationFrame(() => {
      if (!el) return;
      setMeasuredH(el.clientHeight);
    });
  }, []);

  useLayoutEffect(() => {
    measure();
    const el = contentRef.current;
    if (!el) return;
    const ro = new ResizeObserver(measure);
    ro.observe(el);
    return () => ro.disconnect();
  }, [state, children, keyboardOpen, measure]);

  const targetH = measuredH > 0 ? measuredH + HANDLE_H : undefined;
  const archH = targetH ?? NAVBAR_H + 16;

  const isPulsing = state === "confirm" || ["searching", "driver_found", "arriving", "arrived", "in_progress", "destination"].includes(state);

  const [blobs] = useState(() => Array.from({ length: 7 }, (_, i) => i));

  return (
    <>
      <motion.div
        style={{
          position: "fixed",
          left: 0, right: 0, bottom: 0,
          zIndex: 498,
          pointerEvents: "none",
          borderRadius: "64px 64px 0 0",
        }}
        animate={{ height: archH }}
        transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
      >
        {blobs.map(i => (
          <Blob key={i} index={i} state={state} />
        ))}
        <motion.div
          style={{
            position: "absolute",
            inset: 0,
            borderRadius: "64px 64px 0 0",
            background: "linear-gradient(180deg, rgba(255,255,255,0.12) 0%, rgba(255,255,255,0.04) 60%, rgba(0,0,0,0.01) 100%)",
          }}
          animate={isPulsing ? { opacity: [0.7, 1, 0.7] } : { opacity: 0.9 }}
          transition={isPulsing ? { duration: 2.5, repeat: Infinity, ease: "easeInOut" } : { duration: 0.5 }}
        />
        <div
          style={{
            position: "absolute",
            inset: 0,
            borderRadius: "64px 64px 0 0",
            backgroundImage: NOISE_SVG,
            backgroundRepeat: "repeat",
            backgroundSize: "180px 180px",
            opacity: 0.08,
            mixBlendMode: "soft-light",
          }}
        />
      </motion.div>

      <motion.div
        ref={sheetRef}
        role="region"
        aria-label="Panel de información del viaje"
        initial={false}
        animate={{ height: targetH }}
        transition={{ duration: (anim?.duration || 300) / 1000 * 0.5, ease: [0.16, 1, 0.3, 1] }}
        style={{
          width: "auto",
          overflow: "visible",
          background: "rgba(245, 246, 248, 0.75)",
          backdropFilter: "blur(24px) saturate(1.2)",
          WebkitBackdropFilter: "blur(24px) saturate(1.2)",
          borderTopLeftRadius: 64,
          borderTopRightRadius: 64,
          borderBottomLeftRadius: 0,
          borderBottomRightRadius: 0,
          position: "fixed",
          left: 16,
          right: 16,
          bottom: 0,
          zIndex: 500,
        }}
      >
        <div style={{
          width: 32, height: 3, borderRadius: 2,
          background: "rgba(0,0,0,0.14)",
          margin: "9px auto 4px",
          flexShrink: 0,
          pointerEvents: "none",
        }} />

        <AnimatePresence mode="wait" custom={direction}>
          <motion.div
            key={state}
            ref={contentRef}
            custom={direction}
            variants={sheetVariants}
            initial="enter"
            animate="center"
            exit="exit"
            transition={{
              type: "tween",
              duration: ((anim?.duration || 300) / 1000) * 0.6,
              ease: [0.16, 1, 0.3, 1],
            }}
            style={{
              overflowY: "auto",
              WebkitOverflowScrolling: "touch",
              maxHeight: keyboardOpen
                ? `calc(100dvh - ${keyboardH + 20}px)`
                : showNavbar
                  ? `calc(100dvh - ${NAVBAR_H + 12}px)`
                  : "calc(100dvh - 24px)",
            }}
          >
            {children}
          </motion.div>
        </AnimatePresence>
      </motion.div>
    </>
  );
}
