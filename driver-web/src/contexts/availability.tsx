import React, { createContext, useContext, useState } from "react";

interface AvailCtx { available: boolean; toggle: () => void; }
const AvailContext = createContext<AvailCtx | undefined>(undefined);

export function AvailabilityProvider({ children }: { children: React.ReactNode }) {
  const [available, setAvailable] = useState(true);
  return <AvailContext.Provider value={{ available, toggle: () => setAvailable(a => !a) }}>{children}</AvailContext.Provider>;
}
export function useAvailability() { const c = useContext(AvailContext); if (!c) throw new Error("useAvailability"); return c; }
