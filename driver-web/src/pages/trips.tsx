import React from "react";
import { AvailabilityToggle } from "@/components/ui/availability_toggle";
import { TripQueue } from "@/components/ui/trip_queue";
import { useTrip } from "@/contexts/trip";

export default function TripsPage() {
  const { queue, accept, reject } = useTrip();
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Trip Queue</h1>
      <AvailabilityToggle />
      {queue.length === 0 ? <div className="card text-center py-12 text-muted-foreground">Queue is empty. New trips will appear here.</div>
        : <TripQueue trips={queue} onAccept={accept} onReject={reject} />}
    </div>
  );
}
