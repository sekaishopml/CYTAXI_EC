"use client";
import { colors, radius } from "@cytaxi/design-tokens";

interface RecentTrip {
  id: string;
  from: string;
  to: string;
  date: string;
  price: string;
  status: string;
}

const MOCK_TRIPS: RecentTrip[] = [
  { id: "PO123RT", from: "Av. 9 de Octubre 100", to: "Mall del Sol", date: "Hoy, 14:30", price: "$8.50", status: "Completado" },
  { id: "RO213KS", from: "Parque Histórico", to: "Plaza de la Administración", date: "Ayer, 10:15", price: "$6.20", status: "Completado" },
];

export function RecentTrips() {
  return (
    <div>
      <div style={{
        display: "flex", justifyContent: "space-between", alignItems: "center",
        padding: "16px 0 12px",
      }}>
        <p style={{
          margin: 0, fontSize: 16, fontWeight: 600,
          fontFamily: "'Space Grotesk', sans-serif",
          color: "#121212", letterSpacing: "-0.01em",
        }}>
          Viajes recientes
        </p>
        <button type="button" onClick={() => {}} aria-label="Ver historial completo"
          style={{
            background: "transparent", border: "none", cursor: "pointer",
            fontSize: 12, fontWeight: 500, color: colors.cobalt,
            fontFamily: "'Inter', sans-serif", padding: "6px 12px",
            borderRadius: 10,
            transition: "background 0.15s",
          }}
          onMouseEnter={e => { e.currentTarget.style.background = "rgba(59,130,246,0.08)"; }}
          onMouseLeave={e => { e.currentTarget.style.background = "transparent"; }}
        >
          Ver historial →
        </button>
      </div>

      <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        {MOCK_TRIPS.map(trip => (
          <div
            key={trip.id}
            style={{
              background: "#ffffff",
              borderRadius: 20,
              padding: "16px 16px 14px",
              cursor: "pointer",
              transition: "transform 0.2s ease, box-shadow 0.2s ease",
              boxShadow: "0 2px 12px rgba(0,0,0,0.06)",
            }}
            onMouseEnter={e => {
              e.currentTarget.style.transform = "translateY(-1px)";
              e.currentTarget.style.boxShadow = "0 6px 24px rgba(0,0,0,0.1)";
            }}
            onMouseLeave={e => {
              e.currentTarget.style.transform = "translateY(0)";
              e.currentTarget.style.boxShadow = "0 2px 12px rgba(0,0,0,0.06)";
            }}
          >
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
              <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
                <span style={{
                  fontSize: 10, fontWeight: 600, padding: "2px 8px",
                  borderRadius: 6, background: "rgba(59,130,246,0.08)",
                  color: colors.cobalt,
                  fontFamily: "'JetBrains Mono', monospace", letterSpacing: "0.04em",
                }}>
                  {trip.id}
                </span>
                <span style={{
                  fontSize: 10, color: "rgba(0,0,0,0.3)",
                  fontFamily: "'Inter', sans-serif",
                }}>
                  {trip.date}
                </span>
              </div>
              <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
                <span style={{
                  fontSize: 10, color: "rgba(0,0,0,0.3)",
                  fontFamily: "'Inter', sans-serif",
                }}>
                  {trip.status}
                </span>
                <span style={{
                  fontSize: 15, fontWeight: 700, color: "#121212",
                  fontFamily: "'Space Grotesk', sans-serif",
                }}>
                  {trip.price}
                </span>
              </div>
            </div>

            <div style={{ display: "flex", gap: 12, alignItems: "flex-start" }}>
              <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 2, paddingTop: 4 }}>
                <div style={{
                  width: 10, height: 10, borderRadius: "50%",
                  background: "linear-gradient(135deg, #2563eb, #3b82f6)",
                  flexShrink: 0,
                  boxShadow: "0 2px 4px rgba(37,99,235,0.3)",
                }} />
                <div style={{ width: 2, height: 24, background: "rgba(0,0,0,0.08)", borderRadius: 1, flexShrink: 0 }} />
                <div style={{ width: 10, height: 10, borderRadius: "50%", background: colors.cobaltLight, flexShrink: 0 }} />
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{
                  margin: 0, fontSize: 13, fontWeight: 500,
                  color: "#121212", fontFamily: "'Inter', sans-serif",
                }}>
                  {trip.from}
                </p>
                <p style={{
                  margin: "8px 0 0", fontSize: 12, fontWeight: 400,
                  color: "rgba(0,0,0,0.45)", fontFamily: "'Inter', sans-serif",
                }}>
                  {trip.to}
                </p>
              </div>
            </div>

            <button
              type="button"
              onClick={e => { e.stopPropagation(); }}
              aria-label="Pedir de nuevo"
              style={{
                marginTop: 12, width: "100%", height: 38,
                background: "linear-gradient(135deg, #f8fafc, #f1f5f9)",
                border: "1px solid rgba(0,0,0,0.06)",
                borderRadius: 12, fontSize: 12, fontWeight: 600,
                fontFamily: "'Inter', sans-serif", color: "#121212",
                cursor: "pointer", transition: "all 0.2s ease",
              }}
              onMouseEnter={e => {
                e.currentTarget.style.background = "linear-gradient(135deg, #f1f5f9, #e2e8f0)";
              }}
              onMouseLeave={e => {
                e.currentTarget.style.background = "linear-gradient(135deg, #f8fafc, #f1f5f9)";
              }}
            >
              Pedir de nuevo
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
