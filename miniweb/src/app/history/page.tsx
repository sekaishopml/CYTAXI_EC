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

interface TripItem {
  id: string;
  status: string;
  origin: string;
  destination: string;
  fare?: string;
  driver?: string;
  created_at?: string;
}

export default function HistoryPage() {
  const [trips, setTrips] = useState<TripItem[]>([]);
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
        const res = await client.request<{ trips?: TripItem[] }>({
          method: "GET",
          path: `/trip/customers/${customerId}/trips`,
        });
        if (!cancelled) setTrips(res.data.trips || []);
      } catch (e: any) {
        if (!cancelled) setError(e?.message || "No se pudieron cargar tus viajes");
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
        <h1 style={{ fontSize: 20, fontWeight: 700, margin: 0, color: colors.textPrimary }}>Mis viajes</h1>
        <div style={{ width: 22 }} />
      </header>

      {loading && <p style={{ color: colors.textMuted, textAlign: "center", marginTop: 40 }}>Cargando…</p>}
      {error && <p style={{ color: colors.danger, textAlign: "center", marginTop: 40 }}>{error}</p>}

      {!loading && !error && (
        trips.length === 0 ? (
          <p style={{ color: colors.textMuted, fontSize: 14, textAlign: "center", marginTop: 40 }}>
            Aún no tienes viajes. <Link href="/" style={{ color: colors.green }}>Solicita uno</Link>.
          </p>
        ) : (
          <ul style={{ listStyle: "none", padding: 0, margin: 0, display: "flex", flexDirection: "column", gap: 10 }}>
            {trips.map((t) => (
              <li key={t.id} className="card">
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: 6 }}>
                  <span style={{ fontWeight: 600, fontSize: 13, color: statusColor(t.status) }}>{t.status.toUpperCase()}</span>
                  <span style={{ fontSize: 13, color: colors.textMuted }}>{t.fare || "—"}</span>
                </div>
                <p style={{ fontSize: 13, margin: "0 0 2px", color: colors.textPrimary }}><span style={{ color: colors.textMuted }}>Desde:</span> {t.origin}</p>
                <p style={{ fontSize: 13, margin: 0, color: colors.textPrimary }}><span style={{ color: colors.textMuted }}>Hasta:</span> {t.destination}</p>
                {t.created_at && (
                  <p style={{ fontSize: 11, margin: "6px 0 0", color: colors.textMuted }}>{new Date(t.created_at).toLocaleString()}</p>
                )}
              </li>
            ))}
          </ul>
        )
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
