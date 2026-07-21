import React from "react";
import { useAuth } from "@/contexts/auth";
import { useAvailability } from "@/contexts/availability";

export function Header() {
  const { driver } = useAuth();
  const { available, toggle } = useAvailability();
  return (
    <header className="sticky top-0 z-30 bg-surface/95 backdrop-blur border-b border-border h-14 flex items-center justify-between px-4">
      <span className="lg:hidden text-lg font-bold">CYTAXI Driver</span>
      <div className="flex items-center gap-4">
        <button onClick={toggle} className={`badge text-sm px-3 py-1 cursor-pointer ${available ? "badge-success" : "badge-danger"}`}>
          {available ? "Online" : "Offline"}
        </button>
        {driver && <span className="text-sm font-medium hidden sm:inline">{driver.name}</span>}
      </div>
    </header>
  );
}
