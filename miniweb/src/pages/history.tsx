import React, { useEffect, useState } from "react";
import { Layout } from "@/components/layout/layout";
import { getPaymentHistory, PaymentRecord } from "@/services/payments";

export default function PaymentHistoryPage() {
  const [payments, setPayments] = useState<PaymentRecord[]>([]);

  useEffect(() => {
    getPaymentHistory().then(d => setPayments(d.payments)).catch(() => {});
  }, []);

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">Payment History</h1>
        {payments.length === 0 && <p className="card text-center py-8 text-muted-foreground">No payments yet</p>}
        <div className="space-y-3">
          {payments.map(p => (
            <div key={p.id} className="card">
              <div className="flex justify-between items-start mb-2">
                <span className="text-sm text-muted-foreground">#{p.id.slice(-6)}</span>
                <span className={`badge text-xs ${p.status === "paid" ? "badge-success" : p.status === "refunded" ? "badge-danger" : "badge-warning"}`}>
                  {p.status}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-1 text-sm">
                <span className="text-muted-foreground">Amount:</span><span className="font-semibold">${p.amount.toFixed(2)}</span>
                <span className="text-muted-foreground">Method:</span><span className="capitalize">{p.method}</span>
                <span className="text-muted-foreground">Trip:</span><span>{p.trip_id}</span>
                <span className="text-muted-foreground">Date:</span><span>{new Date(p.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </Layout>
  );
}
