"use client";
interface DeliveryCardProps {
  onClick: () => void;
}

const PACKAGE_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#047857" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d="M16.5 9.4l-9-5.19M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 002 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
    <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
    <line x1="12" y1="22.08" x2="12" y2="12" />
  </svg>
);

export function DeliveryCard({ onClick }: DeliveryCardProps) {
  return (
    <button type="button" onClick={onClick}
      style={{
        background: "linear-gradient(135deg, #f0fdf4 0%, #e6f7ec 100%)",
        border: "none",
        borderRadius: 24,
        padding: 0,
        cursor: "pointer",
        position: "relative",
        overflow: "hidden",
        textAlign: "left",
        fontFamily: "'Space Grotesk', sans-serif",
        width: "100%",
        transition: "transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.3s ease",
        boxShadow: "0 4px 16px rgba(0,0,0,0.05)",
      }}
      onMouseEnter={e => {
        e.currentTarget.style.transform = "translateY(-3px) scale(1.01)";
        e.currentTarget.style.boxShadow = "0 12px 32px rgba(0,0,0,0.1)";
      }}
      onMouseLeave={e => {
        e.currentTarget.style.transform = "translateY(0) scale(1)";
        e.currentTarget.style.boxShadow = "0 4px 16px rgba(0,0,0,0.05)";
      }}
    >
      <div style={{ position: "absolute", bottom: 0, left: 0, right: 0, height: "40%", background: "radial-gradient(ellipse at 50% 100%, rgba(255,255,255,0.4) 0%, transparent 70%)", pointerEvents: "none" }} />
      <div style={{ position: "relative", zIndex: 1, padding: 20, display: "flex", flexDirection: "column", boxSizing: "border-box" as const }}>
        <div style={{
          width: 44, height: 44, borderRadius: 14,
          background: "rgba(4,120,87,0.12)",
          display: "flex", alignItems: "center", justifyContent: "center",
          flexShrink: 0, marginBottom: 14,
          boxShadow: "0 4px 12px rgba(0,0,0,0.06)",
        }}>
          {PACKAGE_ICON}
        </div>
        <p style={{ margin: 0, fontSize: 16, fontWeight: 600, color: "#121212", lineHeight: "1.2", letterSpacing: "-0.01em" }}>
          Envíos
        </p>
        <p style={{ margin: "4px 0 0", fontSize: 12, fontWeight: 400, color: "rgba(18,18,18,0.55)", fontFamily: "'Inter', sans-serif", lineHeight: "1.3" }}>
          Paquetes y mensajería
        </p>
      </div>
    </button>
  );
}
