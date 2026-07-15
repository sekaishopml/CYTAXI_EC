import React from "react";

interface ProfileCardProps {
  name: string;
  phone: string;
  email?: string;
  tripsCount: number;
}

export function ProfileCard({ name, phone, email, tripsCount }: ProfileCardProps) {
  return (
    <div className="card">
      <div className="flex items-center gap-4 mb-4">
        <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center text-lg font-bold">
          {name.charAt(0)}
        </div>
        <div>
          <h2 className="font-semibold">{name}</h2>
          <p className="text-sm text-muted">{phone}</p>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-2 text-sm">
        <p className="text-muted">Trips</p>
        <p>{tripsCount}</p>
        <p className="text-muted">Email</p>
        <p>{email || "—"}</p>
      </div>
    </div>
  );
}
