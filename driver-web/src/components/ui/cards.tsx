import React from "react";

interface VehicleCardProps { plate: string; brand: string; model: string; year: number; color: string; type: string; active: boolean; }

export function VehicleCard({ plate, brand, model, year, color, type, active }: VehicleCardProps) {
  return (
    <div className="card flex items-center gap-4">
      <div className="w-14 h-14 rounded-xl bg-muted flex items-center justify-center text-xl">◆</div>
      <div className="flex-1">
        <p className="font-semibold">{brand} {model} ({year})</p>
        <p className="text-sm text-muted-foreground">{plate} &middot; {color} &middot; {type}</p>
      </div>
      <span className={`badge ${active ? "badge-success" : "badge-danger"}`}>{active ? "Active" : "Inactive"}</span>
    </div>
  );
}

interface DocumentCardProps { name: string; type: string; status: string; expiresAt?: string; }

export function DocumentCard({ name, type, status, expiresAt }: DocumentCardProps) {
  return (
    <div className="card flex items-center gap-4">
      <div className="w-14 h-14 rounded-xl bg-muted flex items-center justify-center text-xl">◈</div>
      <div className="flex-1">
        <p className="font-semibold">{name}</p>
        <p className="text-sm text-muted-foreground">{type}{expiresAt ? ` · Expires: ${expiresAt}` : ""}</p>
      </div>
      <span className={`badge ${status === "verified" ? "badge-success" : status === "rejected" ? "badge-danger" : "badge-warning"}`}>
        {status}
      </span>
    </div>
  );
}
