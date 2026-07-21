import { describe, it, expect, beforeEach } from "vitest";
import { saveSession, loadSession, clearSession, isSessionValid } from "../state-recovery";

beforeEach(() => {
  localStorage.clear();
});

describe("state-recovery", () => {
  it("saves and loads a valid session", () => {
    saveSession({ state: "input", pickupAddress: "Av. 9 de Octubre" });
    const loaded = loadSession();
    expect(loaded).not.toBeNull();
    expect(loaded!.state).toBe("input");
    expect(loaded!.pickupAddress).toBe("Av. 9 de Octubre");
    expect(loaded!.v).toBe(5);
  });

  it("returns null when no session exists", () => {
    expect(loadSession()).toBeNull();
  });

  it("clears existing session", () => {
    saveSession({ state: "confirm" });
    clearSession();
    expect(loadSession()).toBeNull();
  });

  it("rejects expired sessions (older than 24h)", () => {
    const expired = {
      v: 5,
      state: "input",
      savedAt: Date.now() - 86400001,
    };
    localStorage.setItem("cytaxi_journey", JSON.stringify(expired));
    expect(loadSession()).toBeNull();
  });

  it("rejects wrong version", () => {
    localStorage.setItem(
      "cytaxi_journey",
      JSON.stringify({ v: 4, state: "input", savedAt: Date.now() }),
    );
    expect(loadSession()).toBeNull();
  });

  it("rejects malformed JSON", () => {
    localStorage.setItem("cytaxi_journey", "not-json");
    expect(loadSession()).toBeNull();
  });

  it("isSessionValid returns true for pre-trip states", () => {
    const data = { v: 5, state: "input", savedAt: Date.now() };
    expect(isSessionValid(data as any)).toBe(true);
  });

  it("isSessionValid returns false for null/undefined", () => {
    expect(isSessionValid(null)).toBe(false);
    expect(isSessionValid(undefined as any)).toBe(false);
  });
});
