"use client";
import { ReactNode } from "react";

interface ModuleGridProps {
  header: ReactNode;
  mapPreview: ReactNode;
  scooterExtension?: ReactNode;
  scooterCard: ReactNode;
  deliveryCard: ReactNode;
  officeCard: ReactNode;
  recentTrips: ReactNode;
}

export function ModuleGrid({
  header, mapPreview, scooterExtension,
  scooterCard, deliveryCard, officeCard, recentTrips,
}: ModuleGridProps) {
  return (
    <div style={{
      height: "100dvh",
      overflowY: "auto",
      WebkitOverflowScrolling: "touch",
      background: "linear-gradient(180deg, #f0f4f8 0%, #e8edf2 100%)",
    }}>
      <div style={{ padding: "20px 20px 0" }}>
        {header}
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
          {mapPreview}
          <div style={{ display: "flex", gap: 12 }}>
            <div style={{
              flex: 1, maxWidth: "55%", minHeight: 150,
              display: "flex", position: "relative",
            }}>
              {scooterExtension}
              {scooterCard}
            </div>
            <div style={{ marginLeft: "auto", display: "flex", alignItems: "center", padding: "14px 0" }}>
              {deliveryCard}
            </div>
          </div>
        </div>
        <div style={{ marginTop: 12 }}>
          {officeCard}
        </div>
      </div>

      <div style={{ padding: "0 20px 24px" }}>
        {recentTrips}
      </div>
    </div>
  );
}
