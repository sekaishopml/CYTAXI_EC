import React, { useState } from "react";
import { Layout } from "@/components/layout/layout";
import { executeFullJourney, JourneyResponse } from "@/services/journey";

export default function HomePage() {
  const [phone, setPhone] = useState("");
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<JourneyResponse | null>(null);
  const [step, setStep] = useState<"input" | "confirm" | "result">("input");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setStep("confirm");
  };

  const handleConfirm = async () => {
    setLoading(true);
    const res = await executeFullJourney({
      phone,
      passenger_name: "Customer",
      origin_address: origin,
      origin_lat: -0.180653,
      origin_lng: -78.467838,
      dest_address: destination,
      dest_lat: -0.229850,
      dest_lng: -78.524950,
    });
    setResult(res);
    setLoading(false);
    setStep("result");
  };

  if (step === "result" && result) {
    return (
      <Layout>
        <div className="space-y-6">
          <h1 className="text-2xl font-bold">{result.success ? "Trip Requested!" : "Something went wrong"}</h1>
          {result.success && (
            <div className="card space-y-4">
              <div className="flex items-center justify-center text-5xl">✅</div>
              <p className="text-center text-muted-foreground">Your trip request was created successfully. A driver will be assigned shortly.</p>
              {result.fare_estimate && (
                <div className="bg-muted p-4 rounded-lg space-y-2">
                  <h3 className="font-semibold">Fare Estimate</h3>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <span className="text-muted-foreground">Base:</span><span>${result.fare_estimate.base?.toFixed(2)}</span>
                    <span className="text-muted-foreground">Distance:</span><span>${result.fare_estimate.distance?.toFixed(2)}</span>
                    <span className="text-muted-foreground font-semibold">Total:</span><span className="font-semibold">${result.fare_estimate.total?.toFixed(2)}</span>
                  </div>
                </div>
              )}
              <div className="space-y-1 text-sm">
                <h4 className="font-semibold">Journey Steps</h4>
                {result.timeline.map((t, i) => (
                  <div key={i} className="flex justify-between">
                    <span>{t.step}</span>
                    <span className={t.status === "completed" ? "text-accent" : "text-danger"}>{t.status} ({t.duration})</span>
                  </div>
                ))}
              </div>
              <button onClick={() => { setStep("input"); setResult(null); }} className="btn-primary w-full">New Trip</button>
            </div>
          )}
          {!result.success && (
            <div className="card text-center">
              <p className="text-danger mb-4">{result.error}</p>
              <button onClick={() => setStep("input")} className="btn-primary">Try Again</button>
            </div>
          )}
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Where to?</h1>
        {step === "input" ? (
          <form onSubmit={handleSubmit} className="card space-y-4">
            <input className="input" type="tel" placeholder="Your phone number" value={phone} onChange={e => setPhone(e.target.value)} required />
            <input className="input" placeholder="Pickup location" value={origin} onChange={e => setOrigin(e.target.value)} required />
            <input className="input" placeholder="Destination" value={destination} onChange={e => setDestination(e.target.value)} required />
            <button type="submit" className="btn-primary w-full">Continue</button>
          </form>
        ) : (
          <div className="card space-y-4">
            <h2 className="font-semibold">Confirm your trip</h2>
            <div className="space-y-2 text-sm">
              <p><span className="text-muted-foreground">From:</span> {origin}</p>
              <p><span className="text-muted-foreground">To:</span> {destination}</p>
              <p><span className="text-muted-foreground">Phone:</span> {phone}</p>
            </div>
            <p className="text-xs text-muted-foreground">Estimated fare will be calculated upon confirmation</p>
            <div className="flex gap-3">
              <button onClick={() => setStep("input")} className="btn-secondary flex-1">Back</button>
              <button onClick={handleConfirm} disabled={loading} className="btn-primary flex-1">
                {loading ? "Processing..." : "Confirm Trip"}
              </button>
            </div>
          </div>
        )}
      </section>
    </Layout>
  );
}
