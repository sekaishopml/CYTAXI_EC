import React from "react";

interface BookingFormProps {
  onSubmit: (data: { origin: string; destination: string; type: string }) => void;
  loading?: boolean;
}

export function BookingForm({ onSubmit, loading }: BookingFormProps) {
  const [origin, setOrigin] = React.useState("");
  const [destination, setDestination] = React.useState("");
  const [type, setType] = React.useState("standard");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ origin, destination, type });
  };

  return (
    <form onSubmit={handleSubmit} className="card space-y-4" role="form">
      <h2 className="text-lg font-semibold">Request a Trip</h2>
      <input
        className="input"
        type="text"
        placeholder="Pickup location"
        value={origin}
        onChange={(e) => setOrigin(e.target.value)}
        required
        aria-label="Pickup location"
      />
      <input
        className="input"
        type="text"
        placeholder="Destination"
        value={destination}
        onChange={(e) => setDestination(e.target.value)}
        required
        aria-label="Destination"
      />
      <select
        className="input"
        value={type}
        onChange={(e) => setType(e.target.value)}
        aria-label="Vehicle type"
      >
        <option value="standard">Standard</option>
        <option value="xl">XL</option>
        <option value="premium">Premium</option>
      </select>
      <button type="submit" className="btn-primary w-full" disabled={loading}>
        {loading ? "Requesting..." : "Request Trip"}
      </button>
    </form>
  );
}
