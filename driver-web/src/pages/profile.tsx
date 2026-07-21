import React, { useState, useEffect } from "react";
import { useAuth } from "@/contexts/auth";
import { Avatar, Card } from "@cytaxi/ui";
import { ApiClient } from "@cytaxi/api-client";

const client = new ApiClient({
  baseUrl: typeof window !== "undefined"
    ? `${window.location.protocol}//${window.location.host}/api/v1`
    : "http://localhost:8000",
});

interface DriverProfile {
  id: string;
  name: string;
  phone: string;
  email?: string;
  rating?: number;
  status?: string;
  vehicle?: { plate: string; brand: string; model: string };
}

export default function ProfilePage() {
  const { driver } = useAuth();
  const [profile, setProfile] = useState<DriverProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const id = driver?.id;
    if (!id) { setLoading(false); return; }
    let cancelled = false;
    (async () => {
      try {
        const res = await client.request<DriverProfile>({
          method: "GET",
          path: `/driver/drivers/${id}`,
        });
        if (!cancelled) setProfile(res.data);
      } catch {
        if (!cancelled) setProfile({ id, name: driver?.name || "Conductor", phone: driver?.phone || "" });
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [driver?.id, driver?.name, driver?.phone]);

  if (loading) return <p className="text-sm text-gray-500">Cargando…</p>;
  if (!profile) return <p className="text-sm text-gray-500">No se pudieron cargar tus datos.</p>;

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-900">Mi perfil</h1>

      <Card>
        <div className="flex items-center gap-4">
          <Avatar name={profile.name} size={56} />
          <div>
            <p className="font-bold text-gray-900">{profile.name}</p>
            <p className="text-sm text-gray-500">{profile.phone}</p>
            {profile.rating !== undefined && (
              <p className="text-sm text-gray-500">★ {profile.rating.toFixed(1)}</p>
            )}
          </div>
        </div>
      </Card>

      <Card>
        <h3 className="font-semibold mb-2 text-gray-900">Vehículo</h3>
        {profile.vehicle ? (
          <p className="text-sm text-gray-700">
            {profile.vehicle.brand} {profile.vehicle.model} · {profile.vehicle.plate}
          </p>
        ) : (
          <p className="text-sm text-gray-500">Sin vehículo asociado.</p>
        )}
      </Card>

      <Card>
        <h3 className="font-semibold mb-2 text-gray-900">Estado</h3>
        <p className="text-sm text-gray-700 capitalize">{profile.status || "online"}</p>
      </Card>
    </div>
  );
}
