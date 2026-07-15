import React from "react";

export default function HelpPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Help & Support</h1>
      <div className="card">
        <h2 className="font-semibold">How to receive trips</h2>
        <p className="text-sm text-muted-foreground mt-1">Keep your status online. Trip requests appear in the queue automatically.</p>
      </div>
      <div className="card">
        <h2 className="font-semibold">Earnings</h2>
        <p className="text-sm text-muted-foreground mt-1">Earnings are calculated after each completed trip. Payouts are processed weekly.</p>
      </div>
      <div className="card">
        <h2 className="font-semibold">Contact Support</h2>
        <p className="text-sm text-muted-foreground mt-1">WhatsApp: +593 99 999 9999 &middot; Email: support@cytaxi.app</p>
      </div>
    </div>
  );
}
