"use client";
import { ReactNode } from "react";

interface ModuleCardProps {
  icon: ReactNode;
  title: string;
  subtitle: string;
  bgColor: string;
  accentColor: string;
  onClick: () => void;
}

export function ModuleCard({ icon, title, subtitle, bgColor, accentColor, onClick }: ModuleCardProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      style={{
        background: bgColor,
        border: accentColor ? `1px solid ${accentColor}22` : "1px solid rgba(0,0,0,0.06)",
        borderRadius: 20,
        padding: "16px 14px",
        cursor: "pointer",
        display: "flex",
        flexDirection: "column",
        gap: 10,
        textAlign: "left",
        fontFamily: "'Space Grotesk', sans-serif",
        transition: "transform 0.15s ease, box-shadow 0.15s ease",
        minHeight: 100,
        position: "relative",
        overflow: "hidden",
      }}
      onMouseEnter={e => { e.currentTarget.style.transform = "scale(1.02)"; e.currentTarget.style.boxShadow = "0 4px 20px rgba(0,0,0,0.08)"; }}
      onMouseLeave={e => { e.currentTarget.style.transform = "scale(1)"; e.currentTarget.style.boxShadow = "none"; }}
    >
      <div style={{
        width: 36, height: 36, borderRadius: 12,
        background: accentColor ? `${accentColor}18` : "rgba(0,0,0,0.04)",
        display: "flex", alignItems: "center", justifyContent: "center",
        flexShrink: 0,
      }}>
        {icon}
      </div>
      <div>
        <p style={{ margin: 0, fontSize: 15, fontWeight: 600, color: "#121212", lineHeight: "1.2" }}>
          {title}
        </p>
        <p style={{ margin: "2px 0 0", fontSize: 12, fontWeight: 400, color: "rgba(0,0,0,0.45)", fontFamily: "'Inter', sans-serif" }}>
          {subtitle}
        </p>
      </div>
    </button>
  );
}
