import React, { useEffect, useState, useCallback } from "react";
import { useAvailability } from "@/contexts/availability";
import { getDriverRequests, acceptRequest, rejectRequest, DriverRequest } from "@/services/assignment";

export default function TripsPage() {
  const { available } = useAvailability();
  const [requests, setRequests] = useState<DriverRequest[]>([]);
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<string>("loading");
  const [error, setError] = useState<string | null>(null);

  const fetchRequests = useCallback(async () => {
    if (!available) { setRequests([]); return; }
    try {
      setLoading(true);
      const data = await getDriverRequests();
      setRequests(data.requests || []);
      setError(null);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setLoading(false);
    }
  }, [available]);

  useEffect(() => {
    fetchRequests();
    const interval = setInterval(fetchRequests, 5000);
    return () => clearInterval(interval);
  }, [fetchRequests]);

  const handleAccept = async (requestId: string) => {
    try {
      setStatus("accepting");
      const result = await acceptRequest(requestId);
      if (result.status === "accepted") {
        alert(`Accepted! Driver: ${result.driver?.name}, Vehicle: ${result.driver?.vehicle}, Plate: ${result.driver?.plate}`);
        fetchRequests();
      }
    } catch (e) {
      setError((e as Error).message);
    } finally { setStatus("idle"); }
  };

  const handleReject = async (requestId: string) => {
    try {
      await rejectRequest(requestId, "Not available");
      fetchRequests();
    } catch (e) {
      setError((e as Error).message);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Trip Requests</h1>
      {!available && <div className="card text-center py-8 text-muted-foreground">You are offline. Go online to receive trips.</div>}
      {error && <div className="card bg-danger/10 text-danger" role="alert"><p>{error}</p></div>}
      {loading && requests.length === 0 && <div className="card text-center py-8 text-muted-foreground">Loading requests...</div>}
      {available && requests.length === 0 && !loading && (
        <div className="card text-center py-12 text-muted-foreground">No trip requests. Waiting for new trips.</div>
      )}
      <div className="space-y-3" role="list">
        {requests.map(req => (
          <div key={req.id} className="card" role="listitem">
            <div className="flex justify-between items-start mb-3">
              <div>
                <p className="font-semibold">{req.pickup} → {req.destination}</p>
                <p className="text-sm text-muted-foreground">{req.fare} &middot; {Math.ceil(req.eta_seconds / 60)} min ETA</p>
              </div>
              <span className="text-xs text-muted-foreground">Expires in {Math.max(0, Math.ceil((new Date(req.expires_at).getTime() - Date.now()) / 1000))}s</span>
            </div>
            <div className="flex gap-2">
              <button onClick={() => handleAccept(req.id)} disabled={status === "accepting"} className="btn-accent flex-1 text-sm py-2">
                Accept
              </button>
              <button onClick={() => handleReject(req.id)} className="btn-danger flex-1 text-sm py-2">Decline</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
