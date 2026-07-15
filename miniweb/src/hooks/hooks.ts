import { useState, useCallback } from "react";
import { api } from "@/services/api";

export function useTripRequest() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const requestTrip = useCallback(async (origin: string, destination: string) => {
    setLoading(true);
    setError(null);
    try {
      const result = await api.trips.get("new");
      return result;
    } catch (e) {
      setError((e as Error).message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return { requestTrip, loading, error };
}

export function useProfile() {
  const [loading, setLoading] = useState(false);

  const fetchProfile = useCallback(async (customerId: string) => {
    setLoading(true);
    try {
      return await api.customer.profile(customerId);
    } finally {
      setLoading(false);
    }
  }, []);

  return { fetchProfile, loading };
}
