import React from "react";
import { Layout } from "@/components/layout/layout";
import { TripCard } from "@/components/ui/trip_card";
import { useTrip } from "@/contexts/trip_context";

const mockTrips = [
  { id: "t1", status: "completed", origin: "Downtown", destination: "Airport", fare: "$12.50" },
  { id: "t2", status: "cancelled", origin: "Mall", destination: "Home", fare: "$8.00" },
];

export default function TripHistoryPage() {
  const { history } = useTrip();
  const trips = history.length > 0 ? history : mockTrips;

  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Trip History</h1>
        <div className="space-y-3">
          {trips.map((trip) => (
            <TripCard key={trip.id} {...trip} />
          ))}
        </div>
        {trips.length === 0 && (
          <p className="text-muted text-center py-8">No trips yet</p>
        )}
      </section>
    </Layout>
  );
}
