import { createContext, useContext } from "react";

export type TenantPlan = "free" | "starter" | "business" | "enterprise";

export interface TenantBranding {
  primaryColor: string;
  secondaryColor: string;
  logoUrl: string;
  faviconUrl: string;
  appName: string;
  termsUrl: string;
  privacyUrl: string;
}

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  plan: TenantPlan;
  isActive: boolean;
  maxDrivers: number;
  maxVehicles: number;
  locale: string;
  timezone: string;
  domain: string;
  branding: TenantBranding;
  features: string[];
  createdAt: number;
  updatedAt: number;
}

export const defaultTenant: Tenant = {
  id: "default",
  name: "CYTAXI",
  slug: "cytaxi",
  plan: "enterprise",
  isActive: true,
  maxDrivers: 1000,
  maxVehicles: 1200,
  locale: "es",
  timezone: "America/Guayaquil",
  domain: "cytaxi.app",
  branding: {
    primaryColor: "#00a152",
    secondaryColor: "#121212",
    logoUrl: "/logos/cytaxi.png",
    faviconUrl: "/favicon.ico",
    appName: "CYTAXI",
    termsUrl: "/terms",
    privacyUrl: "/privacy",
  },
  features: ["all"],
  createdAt: 1700000000,
  updatedAt: 1700000000,
};

export const TenantContext = createContext<Tenant>(defaultTenant);

export function useTenant(): Tenant {
  return useContext(TenantContext);
}

export function useFeature(feature: string): boolean {
  const tenant = useTenant();
  return tenant.features.includes("all") || tenant.features.includes(feature);
}

export function useCanManageTenants(): boolean {
  const tenant = useTenant();
  return tenant.plan === "enterprise" || tenant.plan === "business";
}

export const planLabels: Record<TenantPlan, string> = {
  free: "Gratuito",
  starter: "Starter",
  business: "Business",
  enterprise: "Enterprise",
};

export const planLimits: Record<TenantPlan, { drivers: number; vehicles: number; features: string[] }> = {
  free: { drivers: 10, vehicles: 10, features: ["ride_hailing"] },
  starter: { drivers: 50, vehicles: 60, features: ["ride_hailing", "analytics"] },
  business: { drivers: 200, vehicles: 250, features: ["ride_hailing", "analytics", "scheduling", "reports"] },
  enterprise: { drivers: 9999, vehicles: 9999, features: ["all"] },
};
