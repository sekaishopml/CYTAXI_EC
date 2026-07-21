"use client";
interface OfficeCardProps {
  onClick: () => void;
}

const BUILDING_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#0369a1" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <rect x="4" y="2" width="16" height="20" rx="2" ry="2" fill="#f0f9ff" />
    <line x1="9" y1="6" x2="9" y2="6.01" strokeWidth="2.5" />
    <line x1="15" y1="6" x2="15" y2="6.01" strokeWidth="2.5" />
    <line x1="9" y1="10" x2="9" y2="10.01" strokeWidth="2.5" />
    <line x1="15" y1="10" x2="15" y2="10.01" strokeWidth="2.5" />
    <line x1="9" y1="14" x2="9" y2="14.01" strokeWidth="2.5" />
    <line x1="15" y1="14" x2="15" y2="14.01" strokeWidth="2.5" />
    <path d="M9 18h6v4H9z" />
  </svg>
);

export function OfficeCard({ onClick }: OfficeCardProps) {
  return (
    <button type="button" onClick={onClick}
      style={{
        background: "linear-gradient(135deg, #ffffff 0%, #f2f7ff 100%)",
        border: "none",
        borderRadius: "0 0 24px 24px",
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
          background: "rgba(2,132,199,0.12)",
          display: "flex", alignItems: "center", justifyContent: "center",
          flexShrink: 0, marginBottom: 14,
          boxShadow: "0 4px 12px rgba(0,0,0,0.06)",
        }}>
          {BUILDING_ICON}
        </div>
        <p style={{ margin: 0, fontSize: 16, fontWeight: 600, color: "#121212", lineHeight: "1.2", letterSpacing: "-0.01em" }}>
          Oficina
        </p>
        <p style={{ margin: "4px 0 0", fontSize: 12, fontWeight: 400, color: "rgba(18,18,18,0.55)", fontFamily: "'Inter', sans-serif", lineHeight: "1.3" }}>
          Viajes corporativos para tu equipo
        </p>
      </div>
    </button>
  );
}
