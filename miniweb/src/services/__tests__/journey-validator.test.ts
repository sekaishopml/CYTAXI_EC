import { describe, it, expect } from "vitest";
import {
  validateJourneyState,
  validatePickupData,
  validateDestinationData,
  validateRouteData,
  validateFareData,
  validateDriverData,
  validateFullJourney,
} from "../journey-validator";

describe("journey-validator", () => {
  describe("validateJourneyState", () => {
    it("validates a known state", () => {
      const r = validateJourneyState("input", null);
      expect(r.valid).toBe(true);
      expect(r.errors).toHaveLength(0);
    });

    it("rejects unknown state", () => {
      const r = validateJourneyState("unknown_state" as any, null);
      expect(r.valid).toBe(false);
      expect(r.errors).toContain("Estado desconocido: unknown_state");
    });

    it("warns on unexpected back-jump", () => {
      const r = validateJourneyState("input", "destination");
      expect(r.warnings).toContain("Salto inesperado: destination → input");
    });
  });

  describe("validatePickupData", () => {
    it("accepts valid pickup data", () => {
      const r = validatePickupData({ address: "Av. 9 de Octubre", coords: { lat: -2.18, lng: -79.88 } });
      expect(r.valid).toBe(true);
    });

    it("rejects empty address", () => {
      const r = validatePickupData({ address: "", coords: { lat: -2.18, lng: -79.88 } });
      expect(r.valid).toBe(false);
      expect(r.errors).toContain("Dirección de recogida requerida");
    });

    it("rejects null coords", () => {
      const r = validatePickupData({ address: "Test", coords: null });
      expect(r.valid).toBe(false);
    });

    it("warns on warning prefix", () => {
      const r = validatePickupData({ address: "⚠️ Sin acceso", coords: { lat: -2.18, lng: -79.88 } });
      expect(r.warnings.length).toBeGreaterThan(0);
    });
  });

  describe("validateDestinationData", () => {
    it("rejects null dest", () => {
      const r = validateDestinationData({ dest: null, destQuery: "" });
      expect(r.valid).toBe(false);
    });

    it("accepts valid place", () => {
      const r = validateDestinationData({
        dest: { name: "Test", address: "123 St", lat: -2.18, lng: -79.88 },
        destQuery: "Test",
      });
      expect(r.valid).toBe(true);
    });
  });

  describe("validateRouteData", () => {
    it("rejects null route", () => {
      const r = validateRouteData({ route: null, pickupCoords: null, dest: null });
      expect(r.valid).toBe(false);
    });

    it("rejects zero distance route", () => {
      const r = validateRouteData({
        route: { distance_km: 0, distance_meters: 0, duration_seconds: 0, eta_minutes: 0, polyline: "", pickup: { lat: 0, lng: 0 }, dest: { lat: 0, lng: 0 } },
        pickupCoords: { lat: -2.18, lng: -79.88 },
        dest: { name: "T", address: "A", lat: -2.19, lng: -79.87 },
      });
      expect(r.valid).toBe(false);
    });
  });

  describe("validateFareData", () => {
    it("rejects null fare", () => {
      const r = validateFareData({ fare: null, route: null });
      expect(r.valid).toBe(false);
    });

    it("rejects zero total", () => {
      const r = validateFareData({
        fare: { base: 0, distance: 0, time: 0, subtotal: 0, total: 0, currency: "USD", distance_km: 5, eta_minutes: 10, pricing_model: "standard" },
        route: { distance_km: 5, distance_meters: 5000, duration_seconds: 600, eta_minutes: 10, polyline: "", pickup: { lat: 0, lng: 0 }, dest: { lat: 0, lng: 0 } },
      });
      expect(r.valid).toBe(false);
    });
  });

  describe("validateDriverData", () => {
    it("warns when driver missing in active state", () => {
      const r = validateDriverData({ driver: null, state: "in_progress" });
      expect(r.errors).toContain("Conductor requerido en estado in_progress");
    });

    it("accepts valid driver", () => {
      const r = validateDriverData({
        driver: { id: "drv_1", name: "Carlos", vehicle: "Toyota", plate: "ABC-123", rating: 4.8, photo: "", lat: -2.18, lng: -79.88, eta_seconds: 300, tier: "gold", trust_score: 90 },
        state: "driver_found",
      });
      expect(r.valid).toBe(true);
    });

    it("warns on missing driver name", () => {
      const r = validateDriverData({
        driver: { id: "drv_1", name: "", vehicle: "", plate: "", rating: 5, photo: "", lat: 0, lng: 0, eta_seconds: 0, tier: undefined, trust_score: 0 },
        state: "driver_found",
      });
      expect(r.warnings).toContain("Conductor sin nombre");
    });
  });

  describe("validateFullJourney", () => {
    it("aggregates all validators", () => {
      const r = validateFullJourney({
        state: "pickup_select",
        prevState: null,
        pickupAddress: "",
        pickupCoords: null,
        dest: null,
        destQuery: "",
        route: null,
        fare: null,
        driver: null,
        tripId: "",
      });
      expect(r.valid).toBe(false);
      expect(r.errors.length).toBeGreaterThan(1);
    });
  });
});
