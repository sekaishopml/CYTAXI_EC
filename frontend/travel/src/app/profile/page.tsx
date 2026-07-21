"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { ApiClient } from "@cytaxi/api-client";
import { colors } from "@cytaxi/design-tokens";

const client = new ApiClient({
  baseUrl: typeof window !== "undefined"
    ? `${window.location.protocol}//${window.location.host}/api/v1`
    : "http://localhost:8000",
});

interface CustomerProfile {
  id: string;
  name: string;
  phone: string;
  email?: string;
  trips_count?: number;
}

interface TripItem {
  id: string;
  status: string;
  origin: string;
  destination: string;
  fare?: string;
  driver?: string;
  created_at?: string;
}

export default function ProfilePage() {
  const [profile, setProfile] = useState<CustomerProfile | null>(null);
  const [history, setHistory] = useState<TripItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const phone = typeof window !== "undefined"
    ? localStorage.getItem("cytaxi_phone") || "0000000000"
    : "0000000000";
  const customerId = `cust_${phone}`;

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [pRes, hRes] = await Promise.all([
          client.request<CustomerProfile>({ method: "GET", path: `/customer/customers/${customerId}/profile` }),
          client.request<{ trips?: TripItem[] }>({ method: "GET", path: `/trip/customers/${customerId}/trips` }),
        ]);
        if (cancelled) return;
        setProfile(pRes.data);
        setHistory(hRes.data.trips || []);
      } catch (e: any) {
        if (!cancelled) setError(e?.message || "No se pudieron cargar tus datos");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [customerId]);

  return (
    <main style={{ minHeight: "100dvh", background: colors.surface.bg, padding: "24px 20px 40px" }}>
      <header style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 20 }}>
        <Link href="/" aria-label="Volver" style={{ fontSize: 22, color: colors.textPrimary, textDecoration: "none" }}>←</Link>
        <h1 style={{ fontSize: 20, fontWeight: 700, margin: 0, color: colors.textPrimary }}>Mi perfil</h1>
        <div style={{ width: 22 }} />
      </header>

      {loading && <p style={{ color: colors.textMuted, textAlign: "center", marginTop: 40 }}>Cargando…</p>}
      {error && <p style={{ color: colors.danger, textAlign: "center", marginTop: 40 }}>{error}</p>}

      {!loading && !error && (
        <>
          <div className="card" style={{ display: "flex", alignItems: "center", gap: 14, marginBottom: 8 }}>
            <div style={{
              width: 52, height: 52, borderRadius: "50%",
              background: `linear-gradient(135deg, ${colors.green}, ${colors.greenDark})`,
              display: "flex", alignItems: "center", justifyContent: "center",
              color: "#fff", fontSize: 20, fontWeight: 700, flexShrink: 0,
            }}>{profile?.name?.charAt(0) || "U"}</div>
            <div style={{ flex: 1 }}>
              <p style={{ fontSize: 16, fontWeight: 600, color: colors.textPrimary, margin: 0 }}>{profile?.name || "Usuario"}</p>
              <p style={{ fontSize: 13, color: colors.textMuted, margin: "2px 0 0" }}>{profile?.phone || phone}</p>
              {profile?.email && <p style={{ fontSize: 12, color: colors.textMuted, margin: 0 }}>{profile.email}</p>}
            </div>
            <div style={{ textAlign: "center" }}>
              <p style={{ fontSize: 22, fontWeight: 700, color: colors.green, margin: 0 }}>{profile?.trips_count ?? history.length}</p>
              <p style={{ fontSize: 10, color: colors.textMuted, margin: 0, letterSpacing: "0.05em", textTransform: "uppercase" }}>Viajes</p>
            </div>
          </div>

          <h2 style={{ fontSize: 16, fontWeight: 600, margin: "24px 0 12px", color: colors.textPrimary }}>Viajes recientes</h2>
          {history.length === 0 ? (
            <p style={{ color: colors.textMuted, fontSize: 14 }}>Aún no tienes viajes.</p>
          ) : (
            <ul style={{ listStyle: "none", padding: 0, margin: 0, display: "flex", flexDirection: "column", gap: 10 }}>
              {history.slice(0, 10).map((t) => (
                <li key={t.id} className="card">
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: 6 }}>
                    <span style={{ fontWeight: 600, fontSize: 13, color: statusColor(t.status) }}>{t.status.toUpperCase()}</span>
                    <span style={{ fontSize: 13, color: colors.textMuted }}>{t.fare || "—"}</span>
                  </div>
                  <p style={{ fontSize: 13, margin: "0 0 2px", color: colors.textPrimary }}><span style={{ color: colors.textMuted }}>Desde:</span> {t.origin}</p>
                  <p style={{ fontSize: 13, margin: 0, color: colors.textPrimary }}><span style={{ color: colors.textMuted }}>Hasta:</span> {t.destination}</p>
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </main>
  );
}

function statusColor(s: string): string {
  if (s === "completed") return colors.green;
  if (s === "cancelled") return colors.danger;
  if (s === "started" || s === "in_progress") return colors.blue;
  return colors.textMuted;
}
