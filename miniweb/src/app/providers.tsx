import React from "react";
import { AuthProvider } from "@/contexts/auth_context";
import { TripProvider } from "@/contexts/trip_context";
import { QueryClient, QueryClientProvider } from "react-query";

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <TripProvider>{children}</TripProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}
