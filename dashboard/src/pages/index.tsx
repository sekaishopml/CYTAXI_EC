import { useState, useEffect } from "react";
import { Layout } from "../components/Layout";

interface Stat { label: string; value: string; icon: string; color: string; }
interface TripRow { id: string; passenger: string; driver: string; origin: string; dest: string; status: string; fare: string; time: string; }
interface Incident { id: string; type: string; desc: string; time: string; severity: "low" | "medium" | "high"; }
interface DriverStatus { id: string; name: string; status: "online" | "busy" | "offline"; trips: number; rating: number; }
interface TenantStat { name: string; slug: string; plan: string; trips: number; drivers: number; revenue: string; status: string; }

const MOCK_TENANTS: TenantStat[] = [
  { name: "CYTAXI Cooperativa", slug: "cytaxi", plan: "Enterprise", trips: 127, drivers: 14, revenue: "$1,847", status: "active" },
  { name: "Demo Cooperativa", slug: "demo", plan: "Starter", trips: 23, drivers: 5, revenue: "$312", status: "active" },
  { name: "Taxi Norte", slug: "taxi-norte", plan: "Free", trips: 0, drivers: 0, revenue: "$0", status: "inactive" },
];

const MOCK_TRIPS: TripRow[] = [
  { id: "TRP-001", passenger: "María G.", driver: "Carlos M.", origin: "Centro", dest: "Mall del Sol", status: "En curso", fare: "$8.50", time: "2 min" },
  { id: "TRP-002", passenger: "José R.", driver: "Ana P.", origin: "Norte", dest: "Aeropuerto", status: "Buscando", fare: "$12.00", time: "5 min" },
  { id: "TRP-003", passenger: "Lucía F.", driver: "Pedro S.", origin: "Urdesa", dest: "Samborondón", status: "Completado", fare: "$6.75", time: "18 min" },
  { id: "TRP-004", passenger: "David L.", driver: "—", origin: "Sur", dest: "Centro", status: "Pendiente", fare: "$5.50", time: "1 min" },
];

const MOCK_DRIVERS: DriverStatus[] = [
  { id: "DRV-001", name: "Carlos M.", status: "online", trips: 12, rating: 4.8 },
  { id: "DRV-002", name: "Ana P.", status: "busy", trips: 8, rating: 4.9 },
  { id: "DRV-003", name: "Pedro S.", status: "online", trips: 15, rating: 4.7 },
  { id: "DRV-004", name: "Luis T.", status: "offline", trips: 3, rating: 4.2 },
  { id: "DRV-005", name: "Sofía M.", status: "busy", trips: 10, rating: 4.5 },
];

const MOCK_INCIDENTS: Incident[] = [
  { id: "INC-001", type: "Pago fallido", desc: "Transacción rechazada TRP-002", time: "Hace 2 min", severity: "medium" },
  { id: "INC-002", type: "Queja", desc: "Conductor reportado por cliente", time: "Hace 15 min", severity: "high" },
  { id: "INC-003", type: "Cancelación", desc: "Viaje cancelado por timeout", time: "Hace 30 min", severity: "low" },
];

const statusColors: Record<string, string> = {
  "En curso": "bg-blue-100 text-blue-700",
  "Buscando": "bg-yellow-100 text-yellow-700",
  "Completado": "bg-green-100 text-green-700",
  "Pendiente": "bg-gray-100 text-gray-600",
};

const driverStatusConfig: Record<string, { color: string; dot: string }> = {
  online: { color: "bg-green-100 text-green-700", dot: "bg-green-500" },
  busy: { color: "bg-yellow-100 text-yellow-700", dot: "bg-yellow-500" },
  offline: { color: "bg-gray-100 text-gray-500", dot: "bg-gray-400" },
};

export default function DashboardPage() {
  const [stats, setStats] = useState<Stat[]>([
    { label: "Viajes hoy", value: "—", icon: "🚗", color: "text-brand-green" },
    { label: "Conductores activos", value: "—", icon: "👤", color: "text-blue-500" },
    { label: "Ingresos hoy", value: "—", icon: "💰", color: "text-yellow-500" },
    { label: "Tenants activos", value: "—", icon: "🏢", color: "text-purple-500" },
  ]);
  const [services, setServices] = useState<{ name: string; port: number; status: string }[]>([]);
  const [selectedTab, setSelectedTab] = useState<"trips" | "drivers" | "incidents" | "tenants">("tenants");

  useEffect(() => {
    setStats([
      { label: "Viajes hoy", value: "150", icon: "🚗", color: "text-brand-green" },
      { label: "Conductores activos", value: "19", icon: "👤", color: "text-blue-500" },
      { label: "Ingresos hoy", value: "$2,159", icon: "💰", color: "text-yellow-500" },
      { label: "Tenants activos", value: "2/3", icon: "🏢", color: "text-purple-500" },
    ]);
    setServices([
      { name: "API Gateway", port: 8000, status: "healthy" },
      { name: "Trip Engine", port: 8087, status: "healthy" },
      { name: "Pricing Engine", port: 8088, status: "healthy" },
      { name: "Matching Engine", port: 8089, status: "healthy" },
      { name: "Geo Engine", port: 8082, status: "healthy" },
      { name: "Admin Engine", port: 8094, status: "healthy" },
    ]);
  }, []);

  return (
    <Layout activeSection="overview">
      <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 24, letterSpacing: "-0.02em" }}>
        Centro de Operaciones Multi-Tenant
      </h2>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 16, marginBottom: 24 }}>
        {stats.map((s) => (
          <div key={s.label} style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
              <div>
                <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 6 }}>{s.label}</p>
                <p style={{ fontSize: 28, fontWeight: 700, color: "#121212", letterSpacing: "-0.03em" }}>{s.value}</p>
              </div>
              <span style={{ fontSize: 24 }}>{s.icon}</span>
            </div>
          </div>
        ))}
      </div>

      <div style={{ display: "flex", gap: 4, marginBottom: 16, background: "#fff", borderRadius: 12, padding: 4, boxShadow: "0 1px 2px rgba(0,0,0,0.04)", width: "fit-content" }}>
        {(["tenants", "trips", "drivers", "incidents"] as const).map((tab) => (
          <button key={tab} onClick={() => setSelectedTab(tab)}
            style={{
              padding: "8px 20px", borderRadius: 10, fontSize: 13, fontWeight: 600,
              border: "none", cursor: "pointer", fontFamily: "Inter",
              background: selectedTab === tab ? "#121212" : "transparent",
              color: selectedTab === tab ? "#fff" : "#8a8a8a",
              transition: "all 0.15s",
            }}>
            {tab === "tenants" && "🏢 Tenants"}
            {tab === "trips" && "🚗 Viajes"}
            {tab === "drivers" && "👤 Conductores"}
            {tab === "incidents" && "🔔 Alertas"}
          </button>
        ))}
      </div>

      <div style={{ background: "#fff", borderRadius: 16, boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)", overflow: "hidden" }}>
        {selectedTab === "tenants" && (
          <div style={{ overflowX: "auto" }}>
            <table style={{ width: "100%", borderCollapse: "collapse" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid #f0f0f0" }}>
                  {["Cooperativa", "Slug", "Plan", "Viajes", "Conductores", "Ingresos", "Estado"].map((h) => (
                    <th key={h} style={{ padding: "14px 16px", fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", textAlign: "left" }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {MOCK_TENANTS.map((t) => (
                  <tr key={t.slug} style={{ borderBottom: "1px solid #f5f5f5" }}>
                    <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>{t.name}</td>
                    <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a", fontFamily: "monospace" }}>{t.slug}</td>
                    <td style={{ padding: "12px 16px" }}>
                      <span style={{ display: "inline-block", padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, background: t.plan === "Enterprise" ? "#e8f5e9" : t.plan === "Starter" ? "#e3f2fd" : "#f5f5f5", color: t.plan === "Enterprise" ? "#00a152" : t.plan === "Starter" ? "#1565c0" : "#8a8a8a" }}>
                        {t.plan}
                      </span>
                    </td>
                    <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>{t.trips}</td>
                    <td style={{ padding: "12px 16px", fontSize: 13 }}>{t.drivers}</td>
                    <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>{t.revenue}</td>
                    <td style={{ padding: "12px 16px" }}>
                      <span style={{ display: "inline-flex", alignItems: "center", gap: 6, padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, background: t.status === "active" ? "#e8f5e9" : "#f5f5f5", color: t.status === "active" ? "#00a152" : "#8a8a8a" }}>
                        <span style={{ width: 6, height: 6, borderRadius: "50%", background: t.status === "active" ? "#00a152" : "#8a8a8a", display: "inline-block" }} />
                        {t.status === "active" ? "Activo" : "Inactivo"}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {selectedTab === "trips" && (
          <div style={{ overflowX: "auto" }}>
            <table style={{ width: "100%", borderCollapse: "collapse" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid #f0f0f0" }}>
                  {["ID", "Pasajero", "Conductor", "Origen", "Destino", "Estado", "Tarifa", "Tiempo"].map((h) => (
                    <th key={h} style={{ padding: "14px 16px", fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", textAlign: "left" }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {MOCK_TRIPS.map((t) => (
                  <tr key={t.id} style={{ borderBottom: "1px solid #f5f5f5" }}>
                    <td style={{ padding: "12px 16px", fontSize: 12, fontWeight: 600, color: "#121212", fontFamily: "monospace" }}>{t.id}</td>
                    <td style={{ padding: "12px 16px", fontSize: 13 }}>{t.passenger}</td>
                    <td style={{ padding: "12px 16px", fontSize: 13 }}>{t.driver}</td>
                    <td style={{ padding: "12px 16px", fontSize: 12, color: "#3c4a42" }}>{t.origin}</td>
                    <td style={{ padding: "12px 16px", fontSize: 12, color: "#3c4a42" }}>{t.dest}</td>
                    <td style={{ padding: "12px 16px" }}>
                      <span style={{ display: "inline-block", padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, ...(statusColors[t.status] ? { background: statusColors[t.status].split(" ")[0], color: statusColors[t.status].split(" ")[1] } : {}) }}>
                        {t.status}
                      </span>
                    </td>
                    <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>{t.fare}</td>
                    <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a" }}>{t.time}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {selectedTab === "drivers" && (
          <div style={{ overflowX: "auto" }}>
            <table style={{ width: "100%", borderCollapse: "collapse" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid #f0f0f0" }}>
                  {["ID", "Nombre", "Estado", "Viajes", "Calificación"].map((h) => (
                    <th key={h} style={{ padding: "14px 16px", fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", textAlign: "left" }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {MOCK_DRIVERS.map((d) => {
                  const sc = driverStatusConfig[d.status];
                  return (
                    <tr key={d.id} style={{ borderBottom: "1px solid #f5f5f5" }}>
                      <td style={{ padding: "12px 16px", fontSize: 12, fontWeight: 600, fontFamily: "monospace" }}>{d.id}</td>
                      <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 500 }}>{d.name}</td>
                      <td style={{ padding: "12px 16px" }}>
                        <span style={{ display: "inline-flex", alignItems: "center", gap: 6, padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, background: sc.color.split(" ")[0], color: sc.color.split(" ")[1] }}>
                          <span style={{ width: 6, height: 6, borderRadius: "50%", background: sc.dot, display: "inline-block" }} />
                          {d.status === "online" ? "En línea" : d.status === "busy" ? "Ocupado" : "Desconectado"}
                        </span>
                      </td>
                      <td style={{ padding: "12px 16px", fontSize: 13 }}>{d.trips}</td>
                      <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>⭐ {d.rating.toFixed(1)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}

        {selectedTab === "incidents" && (
          <div style={{ padding: "16px" }}>
            {MOCK_INCIDENTS.length === 0 ? (
              <p style={{ textAlign: "center", padding: "40px 0", color: "#8a8a8a", fontSize: 14 }}>No hay incidentes activos</p>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                {MOCK_INCIDENTS.map((inc) => {
                  const sevColors: Record<string, string> = { low: "#8a8a8a", medium: "#ffc107", high: "#ba1a1a" };
                  return (
                    <div key={inc.id} style={{ display: "flex", alignItems: "center", gap: 12, padding: "12px 16px", borderRadius: 12, background: "#f9f9f9" }}>
                      <div style={{ width: 8, height: 8, borderRadius: "50%", background: sevColors[inc.severity], flexShrink: 0 }} />
                      <div style={{ flex: 1 }}>
                        <p style={{ fontSize: 13, fontWeight: 600, margin: 0 }}>{inc.type}</p>
                        <p style={{ fontSize: 12, color: "#8a8a8a", margin: "2px 0 0" }}>{inc.desc}</p>
                      </div>
                      <span style={{ fontSize: 11, color: "#8a8a8a", flexShrink: 0 }}>{inc.time}</span>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </div>

      <div style={{ marginTop: 24, background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
        <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 16, color: "#121212" }}>Estado de servicios</h3>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          {services.map((s) => (
            <div key={s.name} style={{ display: "flex", alignItems: "center", gap: 8, padding: "8px 14px", borderRadius: 8, background: "#f5f5f7", fontSize: 12 }}>
              <span style={{ width: 6, height: 6, borderRadius: "50%", background: s.status === "healthy" ? "#00a152" : "#ba1a1a", display: "inline-block" }} />
              <span style={{ fontWeight: 500 }}>{s.name}</span>
              <span style={{ color: "#8a8a8a" }}>:{s.port}</span>
            </div>
          ))}
        </div>
      </div>
    </Layout>
  );
}
