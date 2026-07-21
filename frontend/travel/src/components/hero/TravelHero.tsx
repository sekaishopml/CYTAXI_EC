"use client";
import { ReactNode } from "react";
import { ModuleGrid } from "@/components/modules/ModuleGrid";
import { ScooterCard } from "@/components/modules/ScooterCard";
import { DeliveryCard } from "@/components/modules/DeliveryCard";
import { OfficeCard } from "@/components/modules/OfficeCard";

interface TravelHeroProps {
  header: ReactNode;
  mapPreview: ReactNode;
  onStartTrip: () => void;
  onDelivery?: () => void;
  onOffice?: () => void;
  recentTrips: ReactNode;
}

const UNIFIED_GRADIENT = "linear-gradient(135deg, #ffffff 0%, #f2f7ff 100%)";

export function TravelHero({
  header, mapPreview, onStartTrip,
  onDelivery, onOffice, recentTrips,
}: TravelHeroProps) {
  const scooterExtension = (
    <div style={{
      position: "absolute",
      bottom: "100%",
      left: 0, right: 0,
      height: 17,
      background: UNIFIED_GRADIENT,
      pointerEvents: "none",
    }} />
  );

  return (
    <ModuleGrid
      header={header}
      mapPreview={mapPreview}
      scooterExtension={scooterExtension}
      scooterCard={<ScooterCard onClick={onStartTrip} />}
      deliveryCard={<DeliveryCard onClick={onDelivery || (() => {})} />}
      officeCard={<OfficeCard onClick={onOffice || (() => {})} />}
      recentTrips={recentTrips}
    />
  );
}
