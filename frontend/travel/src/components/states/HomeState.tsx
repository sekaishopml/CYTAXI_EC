"use client";
import { TravelHero } from "@/components/hero/TravelHero";
import { MapPreview } from "@/components/map/MapPreview";
import { RecentTrips } from "@/features/trip/ui/RecentTrips";

interface HomeStateProps {
  onStartTrip: () => void;
}

function getGreeting() {
  const h = new Date().getHours();
  if (h < 12) return "Buenos días";
  if (h < 18) return "Buenas tardes";
  return "Buenas noches";
}

const header = (
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
);

export function HomeState({ onStartTrip }: HomeStateProps) {
  return (
    <TravelHero
      header={header}
      mapPreview={<MapPreview onClick={onStartTrip} shape="p-top" noShadow />}
      onStartTrip={onStartTrip}
      recentTrips={<RecentTrips />}
    />
  );
}
