import React, { useState, useEffect } from "react";
import { getDriverEarnings } from "@/services/payments";

export default function DashboardPage() {
  return null; // keep existing dashboard, just add earnings widget
}

export function EarningsWidget() {
  const [earnings, setEarnings] = useState({ trips_completed: 0, total_earnings: 0, net_earnings: 0 });

  useEffect(() => {
    getDriverEarnings("drv_1000").then(setEarnings).catch(() => {});
  }, []);

  return (
    <div className="card">
      <h3 className="font-semibold mb-3">Earnings</h3>
      <div className="grid grid-cols-3 gap-2 text-center">
        <div className="bg-muted p-2 rounded">
          <p className="text-lg font-bold">{earnings.trips_completed}</p>
          <p className="text-xs text-muted-foreground">Trips</p>
        </div>
        <div className="bg-muted p-2 rounded">
          <p className="text-lg font-bold">${earnings.total_earnings.toFixed(0)}</p>
          <p className="text-xs text-muted-foreground">Total</p>
        </div>
        <div className="bg-muted p-2 rounded">
          <p className="text-lg font-bold text-accent">${earnings.net_earnings.toFixed(0)}</p>
          <p className="text-xs text-muted-foreground">Net</p>
        </div>
      </div>
    </div>
  );
}
