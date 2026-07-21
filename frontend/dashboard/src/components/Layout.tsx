import { useRouter } from "next/router";

const navItems = [
  { label: "Dashboard", icon: "📊", href: "/", section: "overview" },
  { label: "Tenants", icon: "🏢", href: "/tenants", section: "tenants" },
  { label: "Facturación", icon: "💰", href: "/billing", section: "billing" },
  { label: "Sistema", icon: "⚙️", href: "/system", section: "system" },
];

export function Layout({ children, activeSection }: { children: React.ReactNode; activeSection: string }) {
  const router = useRouter();

  return (
    <div style={{ fontFamily: "'Inter', system-ui, sans-serif", minHeight: "100vh", background: "#f5f5f7", color: "#121212" }}>
      <aside style={{ position: "fixed", left: 0, top: 0, bottom: 0, width: 220, background: "#121212", color: "#fff", padding: "24px 16px", display: "flex", flexDirection: "column" }}>
        <div style={{ fontSize: 18, fontWeight: 700, marginBottom: 4, letterSpacing: "-0.03em" }}>CYTAXI</div>
        <div style={{ fontSize: 10, color: "rgba(255,255,255,0.3)", marginBottom: 32, textTransform: "uppercase", letterSpacing: "0.08em" }}>Super Admin</div>
        <nav style={{ display: "flex", flexDirection: "column", gap: 2 }}>
          {navItems.map((item) => (
            <a key={item.href} href={item.href} onClick={(e) => { e.preventDefault(); router.push(item.href); }}
              style={{
                display: "flex", alignItems: "center", gap: 10, padding: "10px 12px", borderRadius: 10,
                fontSize: 13, fontWeight: 500, cursor: "pointer",
                background: activeSection === item.section ? "rgba(255,255,255,0.08)" : "transparent",
                color: activeSection === item.section ? "#fff" : "rgba(255,255,255,0.6)",
                textDecoration: "none", transition: "all 0.15s",
              }}>
              <span>{item.icon}</span>
              <span>{item.label}</span>
            </a>
          ))}
        </nav>
        <div style={{ marginTop: "auto", padding: "12px", borderRadius: 10, background: "rgba(255,255,255,0.04)", fontSize: 11, color: "rgba(255,255,255,0.4)" }}>
          v1.0.0 • Enterprise
        </div>
      </aside>
      <main style={{ marginLeft: 220, padding: "24px 32px" }}>
        {children}
      </main>
    </div>
  );
}
