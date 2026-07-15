import React from "react";

interface TripCardProps {
  id: string;
  status: string;
  origin: string;
  destination: string;
  fare: string;
  driver?: string;
  onClick?: () => void;
}

export function TripCard({ id, status, origin, destination, fare, driver, onClick }: TripCardProps) {
  const statusColors: Record<string, string> = {
    completed: "text-green-600",
    cancelled: "text-red-600",
    started: "text-blue-600",
    searching: "text-yellow-600",
  };

  return (
    <button
      onClick={onClick}
      className="card w-full text-left hover:shadow-md transition-shadow"
      aria-label={`Trip ${id} from ${origin} to ${destination}`}
    >
      <div className="flex justify-between items-start mb-2">
        <span className={`font-semibold text-sm ${statusColors[status] || "text-muted"}`}>
          {status.toUpperCase()}
        </span>
        <span className="text-muted text-sm">{fare}</span>
      </div>
      <p className="text-sm mb-1"><span className="text-muted">From:</span> {origin}</p>
      <p className="text-sm"><span className="text-muted">To:</span> {destination}</p>
      {driver && <p className="text-sm mt-2 text-muted">Driver: {driver}</p>}
    </button>
  );
}
