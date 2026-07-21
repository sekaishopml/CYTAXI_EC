import React from "react";
import { TripCard } from "@/components/ui/trip_card";
import { useTrip } from "@/contexts/trip";

const mockHistory = [
  { id: "h1", pickup: "Airport", destination: "Hotel Zone", fare: "$25.00", eta: 0, status: "completed" },
  { id: "h2", pickup: "Mall", destination: "Downtown", fare: "$8.50", eta: 0, status: "cancelled" },
];

export default function TripHistoryPage() {
  const { history } = useTrip();
  const trips = history.length > 0 ? history : mockHistory;
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Trip History</h1>
      <div className="space-y-3">{trips.map(t => <TripCard key={t.id} {...t} />)}</div>
    </div>
  );
}
