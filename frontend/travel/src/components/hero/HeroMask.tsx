"use client";

interface HeroMaskProps {
  shape?: "p-top";
}

const P_TOP_CLIP = "polygon(24px 0, calc(100% - 24px) 0, 100% 24px, 100% 180px, 88% 179.9px, 76% 179.4px, 64% 178.4px, 52% 176.9px, 40% 174.8px, 28% 172.1px, 16% 168.7px, 0 163px, 0 24px, 24px 0)";

export function HeroMask({ shape }: HeroMaskProps) {
  if (!shape) return null;
  return (
    <div style={{
      position: "absolute",
      inset: 0,
      clipPath: P_TOP_CLIP,
      background: "transparent",
      pointerEvents: "none",
    }} />
  );
}
