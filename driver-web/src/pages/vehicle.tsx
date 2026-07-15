import React from "react";
import { VehicleCard } from "@/components/ui/cards";

const mockVehicles = [
  { plate: "ABC-1234", brand: "Toyota", model: "Corolla", year: 2023, color: "White", type: "Standard", active: true },
  { plate: "XYZ-5678", brand: "Hyundai", model: "Tucson", year: 2022, color: "Silver", type: "XL", active: false },
];

export default function VehiclePage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Vehicles</h1>
      <div className="space-y-3">{mockVehicles.map(v => <VehicleCard key={v.plate} {...v} />)}</div>
      <button className="btn-primary w-full">Add Vehicle</button>
    </div>
  );
}
