import { describe, it, expect, beforeEach } from "vitest";
import {
  enqueueAction, dequeueAction, getPendingActions,
  clearQueue, incrementRetry, isOnline,
} from "../offline-queue";

beforeEach(() => {
  localStorage.clear();
});

describe("offline-queue", () => {
  it("enqueues and retrieves actions", () => {
    enqueueAction("TRIP_REQUEST", { dest: "test" });
    enqueueAction("LOCATION_UPDATE", { lat: -2.18, lng: -79.88 });

    const pending = getPendingActions();
    expect(pending).toHaveLength(2);
    expect(pending[0].type).toBe("TRIP_REQUEST");
    expect(pending[1].type).toBe("LOCATION_UPDATE");
  });

  it("dequeues a specific action by id", () => {
    enqueueAction("TRIP_REQUEST", {});
    const [first] = getPendingActions();
    dequeueAction(first.id);
    expect(getPendingActions()).toHaveLength(0);
  });

  it("increments retry count", () => {
    enqueueAction("RATING_SUBMIT", { score: 5 });
    const [action] = getPendingActions();
    expect(action.retries).toBe(0);
    incrementRetry(action.id);
    expect(getPendingActions()[0].retries).toBe(1);
  });

  it("excludes actions at max retries from pending", () => {
    enqueueAction("PAYMENT_CONFIRM", {});
    const [action] = getPendingActions();
    for (let i = 0; i < 3; i++) incrementRetry(action.id);
    expect(getPendingActions()).toHaveLength(0);
  });

  it("clears entire queue", () => {
    enqueueAction("TRIP_REQUEST", {});
    enqueueAction("STATUS_CHANGE", {});
    clearQueue();
    expect(getPendingActions()).toHaveLength(0);
  });

  it("handles corrupt localStorage gracefully", () => {
    localStorage.setItem("cytaxi_offline_queue", "corrupt");
    expect(getPendingActions()).toEqual([]);
  });

  it("reports online status", () => {
    expect(typeof isOnline()).toBe("boolean");
  });
});
