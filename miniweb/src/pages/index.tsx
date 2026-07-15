import React from "react";
import { Layout } from "@/components/layout/layout";
import { BookingForm } from "@/components/ui/booking_form";
import { TripCard } from "@/components/ui/trip_card";
import { useTrip } from "@/contexts/trip_context";

export default function HomePage() {
  const { activeTrip, history } = useTrip();

  const handleBooking = (data: { origin: string; destination: string; type: string }) => {
    console.log("Booking requested:", data);
  };

  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Where to?</h1>
        <BookingForm onSubmit={handleBooking} />
        {activeTrip && (
          <section>
            <h2 className="text-lg font-semibold mb-3">Active Trip</h2>
            <TripCard {...activeTrip} />
          </section>
        )}
        {history.length > 0 && (
          <section>
            <h2 className="text-lg font-semibold mb-3">Recent Trips</h2>
            <div className="space-y-3">
              {history.map((trip) => (
                <TripCard key={trip.id} {...trip} />
              ))}
            </div>
          </section>
        )}
      </section>
    </Layout>
  );
}
