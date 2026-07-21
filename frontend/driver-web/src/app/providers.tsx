import React from "react";
import { AuthProvider } from "@/contexts/auth";
import { TripProvider } from "@/contexts/trip";
import { AvailabilityProvider } from "@/contexts/availability";
import { QueryClient, QueryClientProvider } from "react-query";

const qc = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={qc}>
      <AuthProvider><TripProvider><AvailabilityProvider>{children}</AvailabilityProvider></TripProvider></AuthProvider>
    </QueryClientProvider>
  );
}
