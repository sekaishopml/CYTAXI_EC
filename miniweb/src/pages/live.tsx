import React, { useEffect, useState } from "react";
import { Layout } from "@/components/layout/layout";
import { subscribeToTrip, TrackingUpdate } from "@/services/tracking";

export default function LiveTrackingPage() {
  const [tripId, setTripId] = useState("");
  const [connected, setConnected] = useState(false);
  const [update, setUpdate] = useState<TrackingUpdate | null>(null);
  const [completed, setCompleted] = useState(false);
  const [steps, setSteps] = useState<string[]>(["Waiting for driver..."]);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const tid = params.get("trip_id");
    if (!tid) { setConnected(false); return; }
    setTripId(tid);

    const unsubscribe = subscribeToTrip(tid, (data) => {
      setUpdate(data);
      setConnected(true);
      if (data.type === "trip_started") {
        setSteps(prev => [...prev, "Trip started"]);
      }
      if (data.type === "location_update") {
        setSteps(prev => {
          const last = prev[prev.length - 1];
          const newStep = `Location: ${data.driver?.lat?.toFixed(4)}, ${data.driver?.lng?.toFixed(4)} | ETA: ${Math.ceil((data.eta_seconds || 0) / 60)}min`;
          if (last !== newStep) return [...prev, newStep];
          return prev;
        });
      }
      if (data.type === "trip_completed") {
        setCompleted(true);
        setSteps(prev => [...prev, "Trip completed!"]);
      }
    });

    return unsubscribe;
  }, []);

  const driver = update?.driver;

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">{completed ? "Trip Completed!" : "Live Tracking"}</h1>

        {!connected && !tripId && (
          <div className="card text-center py-8 text-muted-foreground">No trip ID provided. Use ?trip_id= to track a trip.</div>
        )}

        {!connected && tripId && (
          <div className="card text-center py-12 space-y-3">
            <div className="animate-spin mx-auto rounded-full h-10 w-10 border-3 border-primary border-t-transparent" />
            <p className="text-muted-foreground">Connecting...</p>
          </div>
        )}

        {connected && driver && (
          <div className="card space-y-4">
            <div className={`text-center p-3 rounded-lg ${completed ? "bg-accent/10" : "bg-primary/10"}`}>
              <p className="text-2xl font-mono font-bold">
                {driver.lat.toFixed(6)}, {driver.lng.toFixed(6)}
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                {completed ? "Final position" : "Live position"}
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="text-center p-3 bg-muted rounded-lg">
                <p className="text-2xl font-bold">{Math.ceil((update?.eta_seconds || 0) / 60)}</p>
                <p className="text-xs text-muted-foreground">Minutes ETA</p>
              </div>
              <div className="text-center p-3 bg-muted rounded-lg">
                <p className="text-2xl font-bold">{completed ? "✅" : "🚗"}</p>
                <p className="text-xs text-muted-foreground">{completed ? "Arrived" : "Moving"}</p>
              </div>
            </div>

            <div className="bg-muted p-4 rounded-lg space-y-2">
              <p className="font-semibold text-sm">Driver Info</p>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <span className="text-muted-foreground">Name:</span><span>{driver.name}</span>
                <span className="text-muted-foreground">Vehicle:</span><span>{driver.vehicle}</span>
                <span className="text-muted-foreground">Plate:</span><span className="font-mono">{driver.plate}</span>
                <span className="text-muted-foreground">Rating:</span><span>⭐ {driver.rating?.toFixed(1)}</span>
              </div>
            </div>
          </div>
        )}

        {completed && (
          <div className="card text-center bg-accent/5">
            <div className="text-5xl mb-3">🎉</div>
            <h2 className="text-lg font-semibold">Trip Completed</h2>
            <p className="text-sm text-muted-foreground mt-1">Thank you for riding with CYTAXI</p>
          </div>
        )}

        <details className="card">
          <summary className="font-semibold text-sm cursor-pointer">Trip Log ({steps.length} events)</summary>
          <div className="mt-3 space-y-1 max-h-48 overflow-y-auto">
            {steps.map((s, i) => (
              <p key={i} className="text-xs font-mono text-muted-foreground">[{i}] {s}</p>
            ))}
          </div>
        </details>
      </div>
    </Layout>
  );
}
