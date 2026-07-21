"use client";
import { colors, radius } from "@cytaxi/design-tokens";
import { MapPreview } from "@/components/MapPreview";
import { RecentTrips } from "@/components/RecentTrips";

interface TravelHomeProps {
  onStartTrip: () => void;
}

function getGreeting() {
  const h = new Date().getHours();
  if (h < 12) return "Buenos días";
  if (h < 18) return "Buenas tardes";
  return "Buenas noches";
}

const UNIFIED_GRADIENT = "linear-gradient(135deg, #ffffff 0%, #f2f7ff 100%)";
const ENVIOS_GRADIENT = "linear-gradient(135deg, #f0fdf4 0%, #e6f7ec 100%)";

function Card({
  icon,
  title,
  subtitle,
  gradient,
  iconBg,
  shape,
  fill,
  onClick,
}: {
  icon: React.ReactNode;
  title: string;
  subtitle: string;
  gradient: string;
  iconBg: string;
  shape?: "p-top" | "p-stem" | "b-stem" | "b-bottom" | "envios";
  fill?: boolean;
  onClick: () => void;
}) {
  const rad =
    shape === "p-top" ? "24px 24px 0 0" :
    shape === "p-stem" ? "0 24px 24px 0" :
    shape === "b-bottom" ? "0 0 24px 24px" :
    shape === "envios" ? "24px" :
    "24px";

  return (
    <button
      type="button"
      onClick={onClick}
      style={{
        background: gradient,
        border: "none",
        borderRadius: rad,
        padding: 0,
        cursor: "pointer",
        position: "relative",
        overflow: "hidden",
        textAlign: "left",
        fontFamily: "'Space Grotesk', sans-serif",
        width: fill ? "100%" : undefined,
        height: fill ? "100%" : undefined,
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
      <div style={{ position: "relative", zIndex: 1, padding: 20, display: "flex", flexDirection: "column", height: fill ? "100%" : undefined, boxSizing: "border-box" as const }}>
        <div style={{
          width: 44, height: 44, borderRadius: 14,
          background: iconBg,
          display: "flex", alignItems: "center", justifyContent: "center",
          flexShrink: 0, marginBottom: 14,
          boxShadow: "0 4px 12px rgba(0,0,0,0.06)",
        }}>
          {icon}
        </div>
        <p style={{
          margin: 0, fontSize: 16, fontWeight: 600,
          color: "#121212", lineHeight: "1.2", letterSpacing: "-0.01em",
        }}>
          {title}
        </p>
        <p style={{
          margin: "4px 0 0", fontSize: 12, fontWeight: 400,
          color: "rgba(18,18,18,0.55)", fontFamily: "'Inter', sans-serif",
          lineHeight: "1.3",
        }}>
          {subtitle}
        </p>
      </div>
    </button>
  );
}

const PIN_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="#2563eb" stroke="#fff" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z" />
    <circle cx="12" cy="10" r="3" fill="#fff" />
  </svg>
);

const SCOOTER_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#b45309" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="5.5" cy="18" r="3.5" fill="#fef3c7" />
    <circle cx="18.5" cy="18" r="3.5" fill="#fef3c7" />
    <path d="M15 7l2 3h3" />
    <path d="M9 7l-3 6h12" />
    <path d="M5.5 14.5L8 7h7l2 3.5" />
  </svg>
);

const PACKAGE_ICON = (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#047857" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d="M16.5 9.4l-9-5.19M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 002 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
    <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
    <line x1="12" y1="22.08" x2="12" y2="12" />
  </svg>
);

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

export function TravelHome({ onStartTrip }: TravelHomeProps) {
  return (
    <div style={{
      height: "100dvh",
      overflowY: "auto",
      WebkitOverflowScrolling: "touch",
      background: "linear-gradient(180deg, #f0f4f8 0%, #e8edf2 100%)",
    }}>
      <div style={{ padding: "20px 20px 0" }}>
        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 16 }}>
          <div>
            <p style={{
              margin: 0, fontSize: 20, fontWeight: 600,
              fontFamily: "'Space Grotesk', sans-serif",
              color: "#121212", letterSpacing: "-0.02em",
            }}>
              {getGreeting()}
            </p>
            <p style={{
              margin: "2px 0 0", fontSize: 13, fontWeight: 400,
              fontFamily: "'Inter', sans-serif",
              color: "rgba(18,18,18,0.5)",
            }}>
              ¿A dónde vamos hoy?
            </p>
          </div>
          <div style={{
            width: 44, height: 44, borderRadius: "50%",
            background: "linear-gradient(135deg, #3b82f6, #60a5fa)",
            display: "flex", alignItems: "center", justifyContent: "center",
            color: "#fff", fontSize: 18, fontWeight: 600,
            fontFamily: "'Inter', sans-serif",
            flexShrink: 0,
            boxShadow: "0 4px 12px rgba(59,130,246,0.3)",
          }}>
            N
          </div>
        </div>
      </div>

      <div style={{
        display: "flex",
        flexDirection: "column",
        gap: 0,
        marginTop: 12,
        padding: "0 20px",
      }}>
        <div style={{
          borderRadius: "24px",
          overflow: "hidden",
          boxShadow: "0 8px 32px rgba(0,0,0,0.12)",
        }}>
          <MapPreview onClick={onStartTrip} shape="p-top" noShadow />
          <div style={{
            display: "flex",
            gap: 12,
          }}>
            <div style={{
              flex: 1,
              maxWidth: "55%",
              minHeight: 150,
              display: "flex",
              position: "relative",
            }}>
              <div style={{
                position: "absolute",
                bottom: "100%",
                left: 0,
                right: 0,
                height: 17,
                background: UNIFIED_GRADIENT,
                pointerEvents: "none",
              }} />
              <Card
                icon={SCOOTER_ICON}
                title="Scooter"
                subtitle="Viajes rápidos y económicos"
                gradient={UNIFIED_GRADIENT}
                iconBg="rgba(217,119,6,0.15)"
                shape="p-stem"
                fill
                onClick={() => {}}
              />
            </div>
            <div style={{ marginLeft: "auto", display: "flex", alignItems: "center", padding: "14px 0" }}>
              <Card
                icon={PACKAGE_ICON}
                title="Envíos"
                subtitle="Paquetes y mensajería"
                gradient={ENVIOS_GRADIENT}
                iconBg="rgba(4,120,87,0.12)"
                shape="envios"
                onClick={() => {}}
              />
            </div>
          </div>
        </div>
        <div style={{ marginTop: 12 }}>
          <Card
            icon={BUILDING_ICON}
            title="Oficina"
            subtitle="Viajes corporativos para tu equipo"
            gradient={UNIFIED_GRADIENT}
            iconBg="rgba(2,132,199,0.12)"
            shape="b-bottom"
            onClick={() => {}}
          />
        </div>
      </div>

      <div style={{ padding: "0 20px 24px" }}>
        <RecentTrips />
      </div>
    </div>
  );
}
