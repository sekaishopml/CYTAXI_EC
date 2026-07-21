"use client";
interface ScooterCardProps {
  onClick: () => void;
}

const SCOOTER_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#b45309" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="5.5" cy="18" r="3.5" fill="#fef3c7" />
    <circle cx="18.5" cy="18" r="3.5" fill="#fef3c7" />
    <path d="M15 7l2 3h3" />
    <path d="M9 7l-3 6h12" />
    <path d="M5.5 14.5L8 7h7l2 3.5" />
  </svg>
);

export function ScooterCard({ onClick }: ScooterCardProps) {
  return (
    <button type="button" onClick={onClick}
      style={{
        background: "linear-gradient(135deg, #ffffff 0%, #f2f7ff 100%)",
        border: "none",
        borderRadius: "0 24px 24px 0",
        padding: 0,
        cursor: "pointer",
        position: "relative",
        overflow: "hidden",
        textAlign: "left",
        fontFamily: "'Space Grotesk', sans-serif",
        width: "100%",
        height: "100%",
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
      <div style={{ position: "absolute", top: -20, right: -20, width: 80, height: 80, borderRadius: "50%", background: "radial-gradient(circle, rgba(255,255,255,0.25) 0%, transparent 70%)", pointerEvents: "none" }} />
      <div style={{ position: "relative", zIndex: 1, padding: 20, display: "flex", flexDirection: "column", height: "100%", boxSizing: "border-box" as const }}>
        <div style={{
          width: 44, height: 44, borderRadius: 14,
          background: "rgba(217,119,6,0.15)",
          display: "flex", alignItems: "center", justifyContent: "center",
          flexShrink: 0, marginBottom: 14,
          boxShadow: "0 4px 12px rgba(0,0,0,0.06)",
        }}>
          {SCOOTER_ICON}
        </div>
        <p style={{ margin: 0, fontSize: 16, fontWeight: 600, color: "#121212", lineHeight: "1.2", letterSpacing: "-0.01em" }}>
          Scooter
        </p>
        <p style={{ margin: "4px 0 0", fontSize: 12, fontWeight: 400, color: "rgba(18,18,18,0.55)", fontFamily: "'Inter', sans-serif", lineHeight: "1.3" }}>
          Viajes rápidos y económicos
        </p>
      </div>
    </button>
  );
}
