import { useEffect, useCallback, useRef } from "react";
import type { RideEventName, EventPayloadMap } from "./types";
import { getGlobalBus, EventBus } from "./event-bus";

export function useEvent<E extends RideEventName>(
  event: E,
  handler: (payload: EventPayloadMap[E]) => void,
  deps: unknown[] = [],
  bus?: EventBus
): void {
  const handlerRef = useRef(handler);
  handlerRef.current = handler;

  useEffect(() => {
    const b = bus ?? getGlobalBus();
    const unsub = b.on(event, (payload) => {
      handlerRef.current(payload as EventPayloadMap[E]);
    });
    return unsub;
  }, [event, bus, ...deps]);
}

export function useEmitEvent(bus?: EventBus) {
  const b = bus ?? getGlobalBus();
  return useCallback(
    <E extends RideEventName>(event: E, payload: EventPayloadMap[E]) => {
      b.emit(event, payload);
    },
    [b]
  );
}
