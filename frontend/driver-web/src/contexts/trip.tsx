import React, { createContext, useContext, useState, useCallback } from "react";

interface Trip { id: string; status: string; pickup: string; destination: string; fare: string; eta: number; }

interface TripCtx {
  queue: Trip[]; current: Trip | null; history: Trip[];
  accept: (id: string) => void; reject: (id: string) => void; startTrip: () => void; completeTrip: () => void;
}
const TripContext = createContext<TripCtx | undefined>(undefined);

export function TripProvider({ children }: { children: React.ReactNode }) {
  const [queue, setQueue] = useState<Trip[]>([]);
  const [current, setCurrent] = useState<Trip | null>(null);
  const [history, setHistory] = useState<Trip[]>([]);

  const accept = useCallback((id: string) => {
    const trip = queue.find(t => t.id === id);
    if (trip) { setQueue(q => q.filter(t => t.id !== id)); setCurrent({ ...trip, status: "accepted" }); }
  }, [queue]);
  const reject = useCallback((id: string) => { setQueue(q => q.filter(t => t.id !== id)); }, []);
  const startTrip = useCallback(() => { if (current) setCurrent({ ...current, status: "started" }); }, [current]);
  const completeTrip = useCallback(() => {
    if (current) { setHistory(h => [{ ...current, status: "completed" }, ...h]); setCurrent(null); }
  }, [current]);

  return <TripContext.Provider value={{ queue, current, history, accept, reject, startTrip, completeTrip }}>{children}</TripContext.Provider>;
}
export function useTrip() { const c = useContext(TripContext); if (!c) throw new Error("useTrip"); return c; }
