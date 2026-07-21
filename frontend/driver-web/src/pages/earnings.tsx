import React from "react";
import { useAuth } from "@/contexts/auth";
import { EarningsWidget } from "@/components/ui/earnings";

export default function EarningsPage() {
  const { driver } = useAuth();
  const driverId = driver?.id || "drv";

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Ganancias</h1>
        <p className="text-sm text-gray-500">Tu actividad y pagos como conductor</p>
      </div>

      <EarningsWidget driverId={driverId} />

      <div className="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 className="font-semibold mb-3 text-gray-900">Pagos recientes</h3>
        <p className="text-sm text-gray-500">
          Los pagos se liquidan semanalmente. Pronto verás aquí el desglose por viaje y los comprobantes.
        </p>
      </div>
    </div>
  );
}
