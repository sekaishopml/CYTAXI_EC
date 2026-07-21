export type BillingPeriod = "monthly" | "yearly";
export type SubscriptionStatus = "active" | "trialing" | "past_due" | "canceled" | "incomplete";
export type InvoiceStatus = "draft" | "open" | "paid" | "void" | "uncollectible";

export interface PlanPrice {
  id: string;
  planId: string;
  currency: string;
  amount: number;
  period: BillingPeriod;
  driversIncluded: number;
  pricePerAdditionalDriver: number;
  tripCommissionPercent: number;
}

export interface Plan {
  id: string;
  name: string;
  description: string;
  prices: PlanPrice[];
  features: string[];
  isPublic: boolean;
}

export interface Subscription {
  id: string;
  tenantId: string;
  planId: string;
  status: SubscriptionStatus;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  trialEnd?: string;
  cancelAtPeriodEnd: boolean;
  driverCount: number;
  metadata: Record<string, string>;
}

export interface Invoice {
  id: string;
  subscriptionId: string;
  tenantId: string;
  number: string;
  status: InvoiceStatus;
  currency: string;
  subtotal: number;
  tax: number;
  total: number;
  paidAt?: string;
  dueDate: string;
  lines: InvoiceLine[];
}

export interface InvoiceLine {
  description: string;
  quantity: number;
  unitAmount: number;
  amount: number;
  type: "subscription" | "additional_driver" | "trip_commission" | "adjustment" | "coupon";
}

export interface Coupon {
  id: string;
  code: string;
  type: "percentage" | "fixed_amount";
  value: number;
  maxRedemptions: number;
  redemptionCount: number;
  expiresAt?: string;
  isActive: boolean;
}

export interface TaxRate {
  id: string;
  label: string;
  percentage: number;
  country: string;
  region?: string;
  isActive: boolean;
}

export const DEFAULT_PLANS: Plan[] = [
  {
    id: "plan_free",
    name: "Gratuito",
    description: "Para cooperativas pequeñas que inician",
    isPublic: true,
    features: ["Rastreo en vivo", "3 conductores", "Pagos básicos"],
    prices: [{
      id: "price_free_monthly", planId: "plan_free", currency: "USD",
      amount: 0, period: "monthly", driversIncluded: 3,
      pricePerAdditionalDriver: 5, tripCommissionPercent: 5,
    }],
  },
  {
    id: "plan_starter",
    name: "Starter",
    description: "Para cooperativas en crecimiento",
    isPublic: true,
    features: ["Rastreo en vivo", "50 conductores", "Analítica básica", "Pagos automatizados"],
    prices: [{
      id: "price_starter_monthly", planId: "plan_starter", currency: "USD",
      amount: 9900, period: "monthly", driversIncluded: 50,
      pricePerAdditionalDriver: 3, tripCommissionPercent: 3,
    }, {
      id: "price_starter_yearly", planId: "plan_starter", currency: "USD",
      amount: 99900, period: "yearly", driversIncluded: 50,
      pricePerAdditionalDriver: 2.5, tripCommissionPercent: 2.5,
    }],
  },
  {
    id: "plan_business",
    name: "Business",
    description: "Para cooperativas consolidadas",
    isPublic: true,
    features: ["Todo lo de Starter", "200 conductores", "Reportes avanzados", "SLA prioritario", "API pública"],
    prices: [{
      id: "price_business_monthly", planId: "plan_business", currency: "USD",
      amount: 29900, period: "monthly", driversIncluded: 200,
      pricePerAdditionalDriver: 2, tripCommissionPercent: 2,
    }, {
      id: "price_business_yearly", planId: "plan_business", currency: "USD",
      amount: 299900, period: "yearly", driversIncluded: 200,
      pricePerAdditionalDriver: 1.5, tripCommissionPercent: 1.5,
    }],
  },
  {
    id: "plan_enterprise",
    name: "Enterprise",
    description: "Para grandes operaciones multi-ciudad",
    isPublic: false,
    features: ["Todo ilimitado", "Soporte dedicado", "On-premise opcional", "IA avanzada", "SSO + MFA", "SLA 99.9%"],
    prices: [{
      id: "price_enterprise_monthly", planId: "plan_enterprise", currency: "USD",
      amount: 99900, period: "monthly", driversIncluded: 9999,
      pricePerAdditionalDriver: 0, tripCommissionPercent: 1,
    }],
  },
];

export function calculateSubscriptionAmount(
  plan: Plan,
  period: BillingPeriod,
  driverCount: number,
  coupon?: Coupon,
  taxRate?: TaxRate,
): { subtotal: number; discount: number; tax: number; total: number } {
  const price = plan.prices.find(p => p.period === period);
  if (!price) return { subtotal: 0, discount: 0, tax: 0, total: 0 };

  let subtotal = price.amount;
  const additionalDrivers = Math.max(0, driverCount - price.driversIncluded);
  subtotal += additionalDrivers * price.pricePerAdditionalDriver * 100;

  let discount = 0;
  if (coupon && coupon.isActive) {
    if (coupon.type === "percentage") {
      discount = Math.round(subtotal * coupon.value / 100);
    } else {
      discount = coupon.value;
    }
    discount = Math.min(discount, subtotal);
  }

  const afterDiscount = subtotal - discount;
  const tax = taxRate ? Math.round(afterDiscount * taxRate.percentage / 100) : 0;
  const total = afterDiscount + tax;

  return { subtotal, discount, tax, total };
}
