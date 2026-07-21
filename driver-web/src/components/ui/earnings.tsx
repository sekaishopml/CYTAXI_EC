import React, { useState, useEffect } from "react";
import { ApiClient } from "@cytaxi/api-client";

const client = new ApiClient({
  baseUrl: typeof window !== "undefined"
    ? `${window.location.protocol}//${window.location.host}/api/v1`
    : "http://localhost:8000",
});

interface EarningsData {
  trips_completed?: number;
  total_earnings?: number;
  net_earnings?: number;
  trips?: Array<{ id: string; fare: number; created_at?: string }>;
}

async function getDriverEarnings(driverId: string): Promise<EarningsData> {
  try {
    const res = await client.request<EarningsData>({
      method: "GET",
      path: `/driver/drivers/${driverId}/earnings`,
    });
    return res.data;
  } catch {
    return { trips_completed: 0, total_earnings: 0, net_earnings: 0, trips: [] };
  }
}

export function EarningsWidget({ driverId }: { driverId: string }) {
  const [earnings, setEarnings] = useState<EarningsData>({
    trips_completed: 0,
    total_earnings: 0,
    net_earnings: 0,
    trips: [],
  });

  useEffect(() => {
    if (driverId) getDriverEarnings(driverId).then(setEarnings).catch(() => {});
  }, [driverId]);

  return (
    <div className="bg-white rounded-2xl p-5 border border-gray-100">
      <h3 className="font-semibold mb-3 text-gray-900">Ganancias</h3>
      <div className="grid grid-cols-3 gap-2 text-center">
        <div className="bg-gray-50 p-2 rounded">
          <p className="text-lg font-bold">{earnings.trips_completed ?? 0}</p>
          <p className="text-xs text-gray-500">Viajes</p>
        </div>
        <div className="bg-gray-50 p-2 rounded">
          <p className="text-lg font-bold">${(earnings.total_earnings ?? 0).toFixed(0)}</p>
          <p className="text-xs text-gray-500">Total</p>
        </div>
        <div className="bg-gray-50 p-2 rounded">
          <p className="text-lg font-bold text-brand-green">${(earnings.net_earnings ?? 0).toFixed(0)}</p>
          <p className="text-xs text-gray-500">Neto</p>
        </div>
      </div>
    </div>
  );
}
