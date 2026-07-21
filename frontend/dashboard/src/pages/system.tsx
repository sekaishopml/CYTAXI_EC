import { useState, useEffect } from "react";
import { Layout } from "../components/Layout";

interface EngineStatus {
  name: string; port: number; status: "healthy" | "degraded" | "down"; uptime: string; version: string; region: string;
}

interface MetricCard {
  label: string; value: string; unit: string; trend: "up" | "down" | "stable";
}

export default function SystemPage() {
  const [engines] = useState<EngineStatus[]>([
    { name: "API Gateway", port: 8000, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Trip Engine", port: 8087, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Pricing Engine", port: 8088, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Matching Engine", port: 8089, status: "healthy", uptime: "14d 7h", version: "0.9.0", region: "us-east" },
    { name: "Payment Engine", port: 8091, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Customer Engine", port: 8085, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Driver Engine", port: 8086, status: "degraded", uptime: "12d 3h", version: "1.0.0", region: "us-east" },
    { name: "Geo Engine", port: 8082, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
    { name: "Trust Engine", port: 8092, status: "healthy", uptime: "14d 7h", version: "0.9.0", region: "us-east" },
    { name: "Admin Engine", port: 8094, status: "healthy", uptime: "7d 2h", version: "1.0.0", region: "us-east" },
    { name: "Analytics Engine", port: 8093, status: "healthy", uptime: "14d 7h", version: "0.9.0", region: "us-east" },
    { name: "Notification Engine", port: 8090, status: "healthy", uptime: "14d 7h", version: "1.0.0", region: "us-east" },
  ]);

  const [metrics] = useState<MetricCard[]>([
    { label: "Peticiones/min", value: "1,247", unit: "req/min", trend: "up" },
    { label: "Latencia p50", value: "42", unit: "ms", trend: "stable" },
    { label: "Latencia p99", value: "187", unit: "ms", trend: "up" },
    { label: "Tasa de error", value: "0.12", unit: "%", trend: "down" },
    { label: "Uso de CPU", value: "34", unit: "%", trend: "stable" },
    { label: "Memoria", value: "2.1", unit: "GB", trend: "up" },
  ]);

  const statusColors: Record<string, string> = { healthy: "#00a152", degraded: "#ffc107", down: "#ba1a1a" };
  const trendIcons: Record<string, string> = { up: "↑", down: "↓", stable: "→" };
  const trendColors: Record<string, string> = { up: "#ba1a1a", down: "#00a152", stable: "#8a8a8a" };

  return (
    <Layout activeSection="system">
      <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 24, letterSpacing: "-0.02em" }}>Estado del Sistema</h2>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(6, 1fr)", gap: 12, marginBottom: 24 }}>
        {metrics.map((m) => (
          <div key={m.label} style={{ background: "#fff", borderRadius: 12, padding: "16px", boxShadow: "0 1px 3px rgba(0,0,0,0.04)" }}>
            <p style={{ fontSize: 10, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", margin: "0 0 6px" }}>{m.label}</p>
            <div style={{ display: "flex", alignItems: "baseline", gap: 6 }}>
              <span style={{ fontSize: 24, fontWeight: 700, color: "#121212", letterSpacing: "-0.03em" }}>{m.value}</span>
              <span style={{ fontSize: 11, color: "#8a8a8a" }}>{m.unit}</span>
              <span style={{ fontSize: 14, color: trendColors[m.trend], marginLeft: "auto" }}>{trendIcons[m.trend]}</span>
            </div>
          </div>
        ))}
      </div>

      <div style={{ background: "#fff", borderRadius: 16, boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)", overflow: "hidden" }}>
        <div style={{ padding: "16px 20px", borderBottom: "1px solid #f0f0f0" }}>
          <h3 style={{ fontSize: 14, fontWeight: 600, margin: 0 }}>Motores del sistema</h3>
        </div>
        <div style={{ overflowX: "auto" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ borderBottom: "1px solid #f0f0f0" }}>
                {["Servicio", "Puerto", "Estado", "Uptime", "Versión", "Región"].map((h) => (
                  <th key={h} style={{ padding: "12px 16px", fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", textAlign: "left" }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {engines.map((e) => (
                <tr key={e.port} style={{ borderBottom: "1px solid #f5f5f5" }}>
                  <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 500 }}>{e.name}</td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a", fontFamily: "monospace" }}>{e.port}</td>
                  <td style={{ padding: "12px 16px" }}>
                    <span style={{ display: "inline-flex", alignItems: "center", gap: 6, padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, background: e.status === "healthy" ? "#e8f5e9" : e.status === "degraded" ? "#fff8e1" : "#ffebee", color: e.status === "healthy" ? "#00a152" : e.status === "degraded" ? "#f57f17" : "#c62828" }}>
                      <span style={{ width: 6, height: 6, borderRadius: "50%", background: statusColors[e.status], display: "inline-block" }} />
                      {e.status === "healthy" ? "Saludable" : e.status === "degraded" ? "Degradado" : "Caído"}
                    </span>
                  </td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#3c4a42" }}>{e.uptime}</td>
                  <td style={{ padding: "12px 16px", fontSize: 12, fontFamily: "monospace" }}>{e.version}</td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a" }}>{e.region}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div style={{ marginTop: 24, display: "grid", gridTemplateColumns: "repeat(2, 1fr)", gap: 16 }}>
        <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
          <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>SLA de disponibilidad</h3>
          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <div style={{ fontSize: 36, fontWeight: 700, color: "#00a152", letterSpacing: "-0.03em" }}>99.87%</div>
            <div style={{ fontSize: 12, color: "#8a8a8a" }}>
              99.9% objetivo<br />Últimos 30 días
            </div>
          </div>
        </div>
        <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
          <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>Último backup</h3>
          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <div style={{ fontSize: 14, fontWeight: 600 }}>Hace 2 horas</div>
            <div style={{ fontSize: 12, color: "#8a8a8a" }}>Automático • Diario a las 03:00</div>
          </div>
        </div>
      </div>
    </Layout>
  );
}
