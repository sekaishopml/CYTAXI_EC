import React from "react";
import { AvailabilityToggle } from "@/components/ui/availability_toggle";
import { TripQueue } from "@/components/ui/trip_queue";
import { useAvailability } from "@/contexts/availability";
import { useTrip } from "@/contexts/trip";
import { useAuth } from "@/contexts/auth";

const mockQueue = [{ id: "t1", pickup: "Downtown", destination: "Mall", fare: "$12.50", eta: 5 }];

export default function DashboardPage() {
  const { driver } = useAuth();
  const { available } = useAvailability();
  const { queue, current, history, accept, reject } = useTrip();
  const trips = queue.length > 0 ? queue : mockQueue;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <AvailabilityToggle />
      {current && (
        <section>
          <h2 className="text-lg font-semibold mb-3">Current Trip</h2>
          <div className="card bg-accent/5 border-accent"><TripCard {...current} /></div>
        </section>
      )}
      {!current && available && (
        <section>
          <h2 className="text-lg font-semibold mb-3">Trip Queue ({trips.length})</h2>
          <TripQueue trips={trips} onAccept={accept} onReject={reject} />
        </section>
      )}
      <section className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div className="card text-center"><p className="text-2xl font-bold">{history.length}</p><p className="text-xs text-muted-foreground">Trips Today</p></div>
        <div className="card text-center"><p className="text-2xl font-bold">{driver?.rating || "—"}</p><p className="text-xs text-muted-foreground">Rating</p></div>
        <div className="card text-center"><p className="text-2xl font-bold">{queue.length}</p><p className="text-xs text-muted-foreground">In Queue</p></div>
        <div className="card text-center"><p className="text-2xl font-bold">$0</p><p className="text-xs text-muted-foreground">Today&apos;s Earnings</p></div>
      </section>
    </div>
  );
}
