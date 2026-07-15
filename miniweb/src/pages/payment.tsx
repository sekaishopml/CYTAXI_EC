import React, { useState, useEffect } from "react";
import { Layout } from "@/components/layout/layout";
import { createPayment, confirmPayment, FeeBreakdown, Receipt, PaymentRecord } from "@/services/payments";

const methods = [
  { id: "card", label: "Credit/Debit Card", icon: "💳" },
  { id: "cash", label: "Cash", icon: "💵" },
  { id: "wallet", label: "CYTAXI Wallet", icon: "👛" },
  { id: "transfer", label: "Bank Transfer", icon: "🏦" },
];

export default function PaymentPage() {
  const [step, setStep] = useState<"loading" | "review" | "method" | "confirming" | "done">("review");
  const [payment, setPayment] = useState<PaymentRecord | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [selectedMethod, setSelectedMethod] = useState("card");
  const [error, setError] = useState("");

  useEffect(() => {
    // Auto-create payment on load
    createPayment("trip_demo", "cust_00", "drv_1000", 5.5, 900)
      .then(data => setPayment(data.payment))
      .catch(e => setError(e.message));
  }, []);

  const handleConfirm = async () => {
    if (!payment) return;
    setStep("confirming");
    try {
      const data = await confirmPayment(payment.id, selectedMethod);
      setPayment(data.payment);
      setReceipt(data.receipt);
      setStep("done");
    } catch (e) {
      setError((e as Error).message);
      setStep("review");
    }
  };

  const fare = payment?.fare_details;

  if (step === "done" && receipt) {
    return (
      <Layout>
        <div className="space-y-6">
          <h1 className="text-2xl font-bold">Payment Confirmed</h1>
          <div className="card space-y-4">
            <div className="text-center text-5xl">✅</div>
            <h2 className="text-lg font-semibold text-center">Paid ${fare?.total.toFixed(2)}</h2>
            <p className="text-sm text-muted-foreground text-center">via {receipt.method}</p>

            <div className="bg-muted p-4 rounded-lg space-y-2">
              <h3 className="font-semibold text-sm">Receipt #{receipt.id.slice(-6)}</h3>
              <div className="grid grid-cols-2 gap-1 text-sm">
                <span className="text-muted-foreground">Base:</span><span>${fare?.base.toFixed(2)}</span>
                <span className="text-muted-foreground">Distance:</span><span>${fare?.distance.toFixed(2)}</span>
                <span className="text-muted-foreground">Time:</span><span>${fare?.time.toFixed(2)}</span>
                <span className="text-muted-foreground">Subtotal:</span><span>${fare?.subtotal.toFixed(2)}</span>
                <span className="text-muted-foreground">Tax (12%):</span><span>${fare?.tax.toFixed(2)}</span>
                <span className="text-muted-foreground font-semibold">Total:</span><span className="font-bold">${fare?.total.toFixed(2)}</span>
                <span className="text-muted-foreground">Method:</span><span className="capitalize">{receipt.method}</span>
                <span className="text-muted-foreground">Date:</span><span>{new Date(receipt.date).toLocaleString()}</span>
              </div>
            </div>

            <a href="/history" className="btn-primary w-full block text-center">View History</a>
            <a href="/" className="btn-secondary w-full block text-center">New Trip</a>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">{step === "confirming" ? "Processing..." : "Payment"}</h1>

        {error && <div className="card bg-danger/10 text-danger" role="alert">{error}</div>}

        {fare && (
          <div className="card space-y-4">
            <div className="bg-muted p-4 rounded-lg space-y-2">
              <div className="grid grid-cols-2 gap-1 text-sm">
                <span className="text-muted-foreground">Base:</span><span>${fare.base.toFixed(2)}</span>
                <span className="text-muted-foreground">Distance:</span><span>${fare.distance.toFixed(2)}</span>
                <span className="text-muted-foreground">Time:</span><span>${fare.time.toFixed(2)}</span>
                <span className="text-muted-foreground">Subtotal:</span><span>${fare.subtotal.toFixed(2)}</span>
                <span className="text-muted-foreground">Tax:</span><span>${fare.tax.toFixed(2)}</span>
              </div>
              <div className="border-t border-border pt-2 flex justify-between font-bold">
                <span>Total</span><span>${fare.total.toFixed(2)} {fare.currency}</span>
              </div>
            </div>

            <div>
              <h3 className="text-sm font-semibold mb-2">Payment Method</h3>
              <div className="space-y-2">
                {methods.map(m => (
                  <button key={m.id} onClick={() => setSelectedMethod(m.id)}
                    className={`w-full p-3 rounded-lg border text-left flex items-center gap-3 transition ${selectedMethod === m.id ? "border-primary bg-primary/5" : "border-border"}`}>
                    <span className="text-xl">{m.icon}</span>
                    <span className="text-sm">{m.label}</span>
                    {selectedMethod === m.id && <span className="ml-auto text-accent">✓</span>}
                  </button>
                ))}
              </div>
            </div>

            <button onClick={handleConfirm} disabled={step === "confirming"} className="btn-primary w-full">
              {step === "confirming" ? "Processing..." : `Pay $${fare.total.toFixed(2)}`}
            </button>
          </div>
        )}
      </div>
    </Layout>
  );
}
