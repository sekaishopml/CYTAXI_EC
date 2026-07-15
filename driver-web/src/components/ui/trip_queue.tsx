import React from "react";
import { TripCard } from "./trip_card";

interface TripQueueProps {
  trips: Array<{ id: string; pickup: string; destination: string; fare: string; eta: number }>;
  onAccept: (id: string) => void; onReject: (id: string) => void;
}
export function TripQueue({ trips, onAccept, onReject }: TripQueueProps) {
  if (trips.length === 0) {
    return <div className="card text-center py-12 text-muted-foreground">No trip requests. Stay online to receive new trips.</div>;
  }
  return (
    <div className="space-y-3" role="list" aria-label="Trip queue">
      {trips.map(t => (
        <TripCard key={t.id} {...t} onAccept={() => onAccept(t.id)} onReject={() => onReject(t.id)} />
      ))}
    </div>
  );
}
