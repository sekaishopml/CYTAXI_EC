import { useState } from "react";
import { Layout } from "../components/Layout";

const PLANS = [
  { name: "Gratuito", price: "$0", drivers: 3, commission: "5%", features: ["Rastreo en vivo", "3 conductores", "Pagos básicos"], popular: false },
  { name: "Starter", price: "$99", drivers: 50, commission: "3%", features: ["Rastreo en vivo", "50 conductores", "Analítica básica", "Pagos automatizados"], popular: false },
  { name: "Business", price: "$299", drivers: 200, commission: "2%", features: ["Todo lo de Starter", "200 conductores", "Reportes avanzados", "SLA prioritario", "API pública"], popular: true },
  { name: "Enterprise", price: "$999", drivers: "Ilimitados", commission: "1%", features: ["Todo ilimitado", "Soporte dedicado", "On-premise opcional", "IA avanzada", "SSO + MFA"], popular: false },
];

export default function BillingPage() {
  const [period, setPeriod] = useState<"monthly" | "yearly">("monthly");

  return (
    <Layout activeSection="billing">
      <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 24, letterSpacing: "-0.02em" }}>Facturación y Planes</h2>

      <div style={{ display: "flex", gap: 4, marginBottom: 24, background: "#fff", borderRadius: 12, padding: 4, boxShadow: "0 1px 2px rgba(0,0,0,0.04)", width: "fit-content" }}>
        {(["monthly", "yearly"] as const).map((p) => (
          <button key={p} onClick={() => setPeriod(p)}
            style={{
              padding: "8px 20px", borderRadius: 10, fontSize: 13, fontWeight: 600,
              border: "none", cursor: "pointer", fontFamily: "Inter",
              background: period === p ? "#121212" : "transparent",
              color: period === p ? "#fff" : "#8a8a8a",
            }}>
            {p === "monthly" ? "Mensual" : "Anual (ahorra 15%)"}
          </button>
        ))}
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 16, marginBottom: 32 }}>
        {PLANS.map((plan) => (
          <div key={plan.name} style={{
            background: "#fff", borderRadius: 16, padding: "24px", boxShadow: plan.popular ? "0 0 0 2px #00a152, 0 4px 16px rgba(0,0,0,0.08)" : "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)",
            position: "relative",
          }}>
            {plan.popular && <div style={{ position: "absolute", top: 12, right: 12, background: "#00a152", color: "#fff", fontSize: 10, fontWeight: 700, padding: "2px 8px", borderRadius: 6, textTransform: "uppercase", letterSpacing: "0.04em" }}>Popular</div>}
            <p style={{ fontSize: 14, fontWeight: 600, margin: "0 0 4px" }}>{plan.name}</p>
            <p style={{ fontSize: 28, fontWeight: 700, margin: "8px 0", letterSpacing: "-0.03em" }}>
              {period === "yearly" && plan.price !== "$0" ? `$${Math.round(parseInt(plan.price.replace("$", "")) * 10 * 0.85 / 100)}` : plan.price}
              <span style={{ fontSize: 12, fontWeight: 400, color: "#8a8a8a" }}>/{period === "monthly" ? "mes" : "año"}</span>
            </p>
            <p style={{ fontSize: 12, color: "#8a8a8a", margin: "4px 0 12px" }}>Hasta {plan.drivers} conductores • {plan.commission} comisión</p>
            <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
              {plan.features.map((f) => (
                <div key={f} style={{ display: "flex", alignItems: "center", gap: 8, fontSize: 12 }}>
                  <span style={{ color: "#00a152" }}>✓</span>
                  <span>{f}</span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      <div style={{ background: "#fff", borderRadius: 16, padding: "20px", boxShadow: "0 1px 3px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)" }}>
        <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 16 }}>Resumen de ingresos (este mes)</h3>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 16 }}>
          <div>
            <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 4 }}>Suscripciones</p>
            <p style={{ fontSize: 22, fontWeight: 700 }}>$1,098</p>
          </div>
          <div>
            <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 4 }}>Comisiones</p>
            <p style={{ fontSize: 22, fontWeight: 700 }}>$847</p>
          </div>
          <div>
            <p style={{ fontSize: 11, fontWeight: 600, color: "#8a8a8a", textTransform: "uppercase", letterSpacing: "0.04em", marginBottom: 4 }}>Total</p>
            <p style={{ fontSize: 22, fontWeight: 700, color: "#00a152" }}>$1,945</p>
          </div>
        </div>
      </div>
    </Layout>
  );
}
