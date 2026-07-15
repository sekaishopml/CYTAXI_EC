import React, { useEffect, useState } from "react";
import { Layout } from "@/components/layout/layout";
import { startMatching } from "@/services/assignment";

interface DriverInfo {
  driver_id: string; name: string; vehicle: string; plate: string;
  distance_meters: number; eta_seconds: number; score: number; rating: number;
}

export default function TripStatusPage() {
  const [phase, setPhase] = useState<"idle" | "searching" | "found" | "error">("idle");
  const [driver, setDriver] = useState<DriverInfo | null>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const tripId = new URLSearchParams(window.location.search).get("trip_id");
    if (tripId && phase === "idle") {
      setPhase("searching");
      startMatching(tripId, -0.18, -78.47)
        .then(data => {
          if (data.candidates && data.candidates.length > 0) {
            setDriver(data.candidates[0]);
            setPhase("found");
          }
        })
        .catch(e => { setMessage(e.message); setPhase("error"); });
    }
  }, []);

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">Trip Status</h1>
        {phase === "searching" && (
          <div className="card text-center py-12 space-y-3">
            <div className="animate-spin mx-auto rounded-full h-10 w-10 border-3 border-primary border-t-transparent" />
            <p className="text-muted-foreground">Searching for available drivers...</p>
          </div>
        )}
        {phase === "found" && driver && (
          <div className="card space-y-4">
            <div className="flex items-center justify-center text-5xl">🚗</div>
            <h2 className="text-lg font-semibold text-center text-accent">Driver Found!</h2>
            <div className="bg-muted p-4 rounded-lg space-y-2">
              <div className="grid grid-cols-2 gap-2 text-sm">
                <span className="text-muted-foreground">Driver:</span><span className="font-medium">{driver.name}</span>
                <span className="text-muted-foreground">Vehicle:</span><span>{driver.vehicle}</span>
                <span className="text-muted-foreground">Plate:</span><span className="font-mono">{driver.plate}</span>
                <span className="text-muted-foreground">ETA:</span><span>{Math.ceil(driver.eta_seconds / 60)} min</span>
                <span className="text-muted-foreground">Rating:</span><span>⭐ {driver.rating.toFixed(1)}</span>
              </div>
            </div>
            <p className="text-xs text-muted-foreground text-center">Your driver is on the way!</p>
          </div>
        )}
        {phase === "error" && (
          <div className="card text-center"><p className="text-danger">{message}</p></div>
        )}
      </div>
    </Layout>
  );
}
