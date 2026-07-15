import React, { createContext, useContext, useState, useCallback } from "react";

interface Trip {
  id: string;
  status: string;
  origin: string;
  destination: string;
  fare: string;
}

interface TripContextType {
  activeTrip: Trip | null;
  history: Trip[];
  setActiveTrip: (trip: Trip | null) => void;
  setHistory: (trips: Trip[]) => void;
  refreshTrips: () => void;
}

const TripContext = createContext<TripContextType | undefined>(undefined);

export function TripProvider({ children }: { children: React.ReactNode }) {
  const [activeTrip, setActiveTrip] = useState<Trip | null>(null);
  const [history, setHistory] = useState<Trip[]>([]);

  const refreshTrips = useCallback(() => {
    // fetch from API Gateway
  }, []);

  return (
    <TripContext.Provider value={{ activeTrip, history, setActiveTrip, setHistory, refreshTrips }}>
      {children}
    </TripContext.Provider>
  );
}

export function useTrip() {
  const ctx = useContext(TripContext);
  if (!ctx) throw new Error("useTrip must be used within TripProvider");
  return ctx;
}
