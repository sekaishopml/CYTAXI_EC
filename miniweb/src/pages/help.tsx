import React from "react";
import { Layout } from "@/components/layout/layout";

export default function HelpPage() {
  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Help & Support</h1>
        <div className="card space-y-3">
          <h2 className="font-semibold">How to request a trip</h2>
          <p className="text-sm text-muted">Enter your pickup and destination, then tap Request Trip.</p>
        </div>
        <div className="card space-y-3">
          <h2 className="font-semibold">Payment methods</h2>
          <p className="text-sm text-muted">We support card payments, cash, and wallet transfers.</p>
        </div>
        <div className="card space-y-3">
          <h2 className="font-semibold">Contact us</h2>
          <p className="text-sm text-muted">WhatsApp: +593 99 999 9999</p>
        </div>
      </section>
    </Layout>
  );
}
