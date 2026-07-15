import React from "react";
import { TripCard } from "@/components/ui/trip_card";
import { useTrip } from "@/contexts/trip";

export default function CurrentTripPage() {
  const { current, startTrip, completeTrip } = useTrip();
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Current Trip</h1>
      {current ? (
        <div className="space-y-3">
          <TripCard {...current} />
          <div className="flex gap-3">
            {current.status === "accepted" && <button onClick={startTrip} className="btn-accent flex-1">Start Trip</button>}
            {current.status === "started" && <button onClick={completeTrip} className="btn-primary flex-1">Complete Trip</button>}
          </div>
        </div>
      ) : <div className="card text-center py-12 text-muted-foreground">No active trip. Accept a trip request to begin.</div>}
    </div>
  );
}
