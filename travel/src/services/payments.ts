const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export interface FeeBreakdown {
  base: number; distance: number; time: number; subtotal: number;
  tax: number; commission_pct: number; total: number; currency: string;
}

export interface PaymentRecord {
  id: string; trip_id: string; customer_id: string; driver_id: string;
  amount: number; method: string; status: string; fare_details: FeeBreakdown; created_at: string;
}

export interface Receipt {
  id: string; payment_id: string; trip_id: string; customer_id: string; driver_id: string;
  amount: number; method: string; fare_details: FeeBreakdown; date: string; status: string;
}

export async function createPayment(tripId: string, customerId: string, driverId: string, distanceKm: number, durationSec: number): Promise<{ payment: PaymentRecord }> {
  const res = await fetch(`${GATEWAY}/api/v1/payments`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId, customer_id: customerId, driver_id: driverId, distance_km: distanceKm, duration_sec: durationSec }),
  });
  return res.json();
}

export async function confirmPayment(paymentId: string, method: string): Promise<{ payment: PaymentRecord; receipt: Receipt }> {
  const res = await fetch(`${GATEWAY}/api/v1/payments/confirm`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ payment_id: paymentId, method }),
  });
  return res.json();
}

export async function getPayment(paymentId: string): Promise<PaymentRecord> {
  const res = await fetch(`${GATEWAY}/api/v1/payments/${paymentId}`);
  return res.json();
}

export async function getReceipt(receiptId: string): Promise<Receipt> {
  const res = await fetch(`${GATEWAY}/api/v1/receipts/${receiptId}`);
  return res.json();
}

export async function getPaymentHistory(): Promise<{ payments: PaymentRecord[] }> {
  const res = await fetch(`${GATEWAY}/api/v1/payments/history`);
  return res.json();
}

export async function getDriverEarnings(driverId: string): Promise<{ trips_completed: number; total_earnings: number; net_earnings: number }> {
  const res = await fetch(`${GATEWAY}/api/v1/payments/driver/${driverId}/earnings`);
  return res.json();
}

export async function refundPayment(paymentId: string, reason: string): Promise<any> {
  const res = await fetch(`${GATEWAY}/api/v1/payments/refund`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ payment_id: paymentId, reason }),
  });
  return res.json();
}
