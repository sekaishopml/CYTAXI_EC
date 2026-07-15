import React, { useState, useEffect, useCallback } from "react";
import { TripCard } from "@/components/ui/trip_card";
import { useTrip } from "@/contexts/trip";
import { startTrip, updateLocation, finishTrip } from "@/services/tracking";

export default function CurrentTripPage() {
  const { current, startTrip: ctxStart, completeTrip } = useTrip();
  const [tripState, setTripState] = useState<"assigned" | "started" | "completed">("assigned");
  const [eta, setEta] = useState(300);
  const [position, setPosition] = useState({ lat: -0.18, lng: -78.47 });
  const [distanceKm, setDistanceKm] = useState(5.5);
  const [events, setEvents] = useState<string[]>([]);

  const addEvent = (msg: string) => setEvents(prev => [...prev, `${new Date().toLocaleTimeString()}: ${msg}`]);

  // Simulate location when trip is started
  useEffect(() => {
    if (tripState !== "started") return;
    const interval = setInterval(async () => {
      const newLat = position.lat + 0.0001;
      const newLng = position.lng + 0.00015;
      const newDist = Math.max(0, distanceKm - 0.1);
      const newEta = Math.max(0, eta - 2);

      setPosition({ lat: newLat, lng: newLng });
      setDistanceKm(newDist);
      setEta(newEta);

      try {
        await updateLocation("trip_demo", "drv_1000", newLat, newLng);
      } catch (e) {}

      addEvent(`📍 ${newLat.toFixed(4)}, ${newLng.toFixed(4)} | ${Math.ceil(newEta / 60)}min | ${newDist.toFixed(1)}km`);

      if (newDist <= 0.1) {
        handleComplete();
      }
    }, 3000);
    return () => clearInterval(interval);
  }, [tripState, position, eta, distanceKm]);

  const handleStart = async () => {
    try {
      setTripState("started");
      ctxStart();
      await startTrip("trip_demo", "drv_1000");
      addEvent("Trip started");
    } catch (e) { addEvent("Error starting trip"); }
  };

  const handleComplete = async () => {
    try {
      setTripState("completed");
      completeTrip();
      await finishTrip("trip_demo");
      addEvent("Trip completed");
    } catch (e) { addEvent("Error completing trip"); }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Current Trip</h1>
      {!current ? (
        <div className="card text-center py-12 text-muted-foreground">No active trip. Accept a request first.</div>
      ) : (
        <div className="space-y-4">
          <div className="card">
            <div className="flex justify-between"><span className="text-muted-foreground">Status</span><span className={`badge ${tripState === "completed" ? "badge-success" : tripState === "started" ? "badge-warning" : "badge-success"}`}>{tripState}</span></div>
            <div className="flex justify-between mt-2"><span className="text-muted-foreground">From</span><span>{current.pickup}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">To</span><span>{current.destination}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">Fare</span><span className="font-semibold">{current.fare}</span></div>
          </div>

          {tripState === "started" && (
            <div className="card space-y-3">
              <div className="grid grid-cols-3 gap-2 text-center">
                <div className="bg-muted rounded-lg p-2"><p className="text-xl font-bold font-mono">{position.lat.toFixed(4)}<br />{position.lng.toFixed(4)}</p><p className="text-xs text-muted-foreground mt-1">Position</p></div>
                <div className="bg-muted rounded-lg p-2"><p className="text-xl font-bold">{Math.ceil(eta / 60)}</p><p className="text-xs text-muted-foreground mt-1">Min ETA</p></div>
                <div className="bg-muted rounded-lg p-2"><p className="text-xl font-bold">{distanceKm.toFixed(1)}</p><p className="text-xs text-muted-foreground mt-1">KM Left</p></div>
              </div>
              <button onClick={handleComplete} className="btn-danger w-full">Finish Trip</button>
            </div>
          )}

          {tripState === "assigned" && (
            <button onClick={handleStart} className="btn-accent w-full text-lg py-4">Start Trip</button>
          )}

          {tripState === "completed" && (
            <div className="card text-center bg-accent/5"><div className="text-4xl mb-2">✅</div><p className="font-semibold">Trip Completed</p><p className="text-sm text-muted-foreground">Fare: {current.fare}</p></div>
          )}

          <details className="card"><summary className="font-semibold text-sm cursor-pointer">Event Log ({events.length})</summary><div className="mt-2 space-y-1 max-h-40 overflow-y-auto">{events.map((e, i) => <p key={i} className="text-xs font-mono text-muted-foreground">{e}</p>)}</div></details>
        </div>
      )}
    </div>
  );
}
