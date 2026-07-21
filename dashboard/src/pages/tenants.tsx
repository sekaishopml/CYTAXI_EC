import { useState } from "react";
import { Layout } from "../components/Layout";

interface TenantItem {
  id: string; name: string; slug: string; plan: string; domain: string;
  drivers: number; maxDrivers: number; status: string; locale: string; created: string;
}

const MOCK_TENANTS: TenantItem[] = [
  { id: "tenant_cytaxi", name: "CYTAXI Cooperativa", slug: "cytaxi", plan: "Enterprise", domain: "cytaxi.app", drivers: 14, maxDrivers: 1000, status: "active", locale: "es", created: "2024-01-01" },
  { id: "tenant_demo", name: "Demo Cooperativa", slug: "demo", plan: "Starter", domain: "demo.cytaxi.app", drivers: 5, maxDrivers: 50, status: "active", locale: "es", created: "2024-06-01" },
  { id: "tenant_norte", name: "Taxi Norte", slug: "taxi-norte", plan: "Free", domain: "norte.cytaxi.app", drivers: 0, maxDrivers: 10, status: "inactive", locale: "es", created: "2025-01-15" },
];

export default function TenantsPage() {
  const [tenants] = useState<TenantItem[]>(MOCK_TENANTS);

  return (
    <Layout activeSection="tenants">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 24 }}>
        <div>
          <h2 style={{ fontSize: 22, fontWeight: 700, letterSpacing: "-0.02em", margin: 0 }}>Tenants</h2>
          <p style={{ fontSize: 13, color: "#8a8a8a", margin: "4px 0 0" }}>Gestiona las cooperativas y empresas registradas</p>
        </div>
        <button style={{ padding: "10px 20px", borderRadius: 10, border: "none", background: "#00a152", color: "#fff", fontWeight: 600, fontSize: 13, cursor: "pointer", fontFamily: "Inter" }}>
          + Nuevo tenant
        </button>
      </div>

      <div style={{ background: "#fff", borderRadius: 16, boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)", overflow: "hidden" }}>
        <div style={{ overflowX: "auto" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ borderBottom: "1px solid #f0f0f0" }}>
                {["Cooperativa", "Slug", "Plan", "Dominio", "Conductores", "Estado", "Registro"].map((h) => (
                  <th key={h} style={{ padding: "14px 16px", fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", textAlign: "left" }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr key={t.id} style={{ borderBottom: "1px solid #f5f5f5" }}>
                  <td style={{ padding: "12px 16px", fontSize: 13, fontWeight: 600 }}>{t.name}</td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a", fontFamily: "monospace" }}>{t.slug}</td>
                  <td style={{ padding: "12px 16px" }}>
                    <span style={{ display: "inline-block", padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, 
                      background: t.plan === "Enterprise" ? "#e8f5e9" : t.plan === "Starter" ? "#e3f2fd" : "#f5f5f5",
                      color: t.plan === "Enterprise" ? "#00a152" : t.plan === "Starter" ? "#1565c0" : "#8a8a8a" }}>
                      {t.plan}
                    </span>
                  </td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#3c4a42" }}>{t.domain}</td>
                  <td style={{ padding: "12px 16px", fontSize: 13 }}>{t.drivers}/{t.maxDrivers}</td>
                  <td style={{ padding: "12px 16px" }}>
                    <span style={{ display: "inline-flex", alignItems: "center", gap: 6, padding: "2px 10px", borderRadius: 999, fontSize: 11, fontWeight: 600, 
                      background: t.status === "active" ? "#e8f5e9" : "#f5f5f5", color: t.status === "active" ? "#00a152" : "#8a8a8a" }}>
                      <span style={{ width: 6, height: 6, borderRadius: "50%", background: t.status === "active" ? "#00a152" : "#8a8a8a", display: "inline-block" }} />
                      {t.status === "active" ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td style={{ padding: "12px 16px", fontSize: 12, color: "#8a8a8a" }}>{t.created}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div style={{ marginTop: 24, display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 16 }}>
        <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
          <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 8 }}>Total tenants</p>
          <p style={{ fontSize: 28, fontWeight: 700, color: "#121212", letterSpacing: "-0.03em" }}>3</p>
        </div>
        <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
          <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 8 }}>Activos</p>
          <p style={{ fontSize: 28, fontWeight: 700, color: "#00a152", letterSpacing: "-0.03em" }}>2</p>
        </div>
        <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
          <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 8 }}>Ingreso mensual</p>
          <p style={{ fontSize: 28, fontWeight: 700, color: "#121212", letterSpacing: "-0.03em" }}>$2,159</p>
        </div>
      </div>
    </Layout>
  );
}
