import React from "react";

interface TripCardProps {
  id: string; pickup: string; destination: string; fare: string; eta: number; status?: string;
  onAccept?: () => void; onReject?: () => void;
}
export function TripCard({ id, pickup, destination, fare, eta, onAccept, onReject }: TripCardProps) {
  return (
    <div className="card" role="listitem">
      <div className="flex justify-between items-start mb-3">
        <div>
          <p className="font-semibold">{pickup} → {destination}</p>
          <p className="text-sm text-muted-foreground">{fare} &middot; {eta} min ETA</p>
        </div>
        <span className="text-xs text-muted-foreground">#{id.slice(0, 8)}</span>
      </div>
      {onAccept && (
        <div className="flex gap-2">
          <button onClick={onAccept} className="btn-accent flex-1 text-sm py-2">Accept</button>
          <button onClick={onReject} className="btn-danger flex-1 text-sm py-2">Decline</button>
        </div>
      )}
    </div>
  );
}
