import React from "react";
import { useAvailability } from "@/contexts/availability";

export function AvailabilityToggle() {
  const { available, toggle } = useAvailability();
  return (
    <div className="card flex items-center justify-between">
      <div>
        <p className="font-semibold">{available ? "You're online" : "You're offline"}</p>
        <p className="text-sm text-muted-foreground">{available ? "Receiving trip requests" : "Tap to go online and receive trips"}</p>
      </div>
      <button onClick={toggle} className={`${available ? "badge-success" : "badge-danger"} badge text-sm px-4 py-2 cursor-pointer`}>
        {available ? "Online" : "Offline"}
      </button>
    </div>
  );
}
