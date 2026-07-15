import React from "react";
import { DocumentCard } from "@/components/ui/cards";

const mockDocs = [
  { name: "Driver License", type: "driver_license", status: "verified" },
  { name: "Vehicle Registration", type: "vehicle_registration", status: "verified" },
  { name: "Insurance", type: "insurance", status: "pending", expiresAt: "2026-06-15" },
  { name: "Background Check", type: "background_check", status: "in_review" },
];

export default function DocumentsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Documents</h1>
      <div className="space-y-3">{mockDocs.map(d => <DocumentCard key={d.type} {...d} />)}</div>
      <button className="btn-primary w-full">Upload Document</button>
    </div>
  );
}
