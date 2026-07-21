export type SoundName =
  | "trip_request"
  | "trip_accepted"
  | "driver_arriving"
  | "driver_arrived"
  | "trip_started"
  | "trip_completed"
  | "payment_success"
  | "payment_failed"
  | "error"
  | "message_received"
  | "button_tap";

export type HapticStyle = "light" | "medium" | "heavy" | "selection" | "success" | "warning" | "error";

interface SoundConfig {
  name: SoundName;
  file?: string;
  duration: number;
  haptic?: HapticStyle;
  volume: number;
  description: string;
}

export const SOUNDS: Record<SoundName, SoundConfig> = {
  trip_request: { name: "trip_request", duration: 1500, haptic: "heavy", volume: 0.8, description: "New trip request for driver" },
  trip_accepted: { name: "trip_accepted", duration: 500, haptic: "success", volume: 0.7, description: "Driver accepted the trip" },
  driver_arriving: { name: "driver_arriving", duration: 1000, haptic: "medium", volume: 0.6, description: "Driver is arriving" },
  driver_arrived: { name: "driver_arrived", duration: 800, haptic: "heavy", volume: 0.7, description: "Driver has arrived" },
  trip_started: { name: "trip_started", duration: 600, haptic: "medium", volume: 0.6, description: "Trip has started" },
  trip_completed: { name: "trip_completed", duration: 1200, haptic: "success", volume: 0.8, description: "Trip completed" },
  payment_success: { name: "payment_success", duration: 800, haptic: "success", volume: 0.7, description: "Payment confirmed" },
  payment_failed: { name: "payment_failed", duration: 600, haptic: "error", volume: 0.5, description: "Payment failed" },
  error: { name: "error", duration: 400, haptic: "error", volume: 0.5, description: "Error occurred" },
  message_received: { name: "message_received", duration: 300, haptic: "light", volume: 0.4, description: "New chat message" },
  button_tap: { name: "button_tap", duration: 100, haptic: "selection", volume: 0.2, description: "UI button tap" },
};

export function playSound(name: SoundName, volume?: number): void {
  if (typeof window === "undefined") return;
  const config = SOUNDS[name];
  if (!config) return;

  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.connect(gain);
    gain.connect(ctx.destination);

    const vol = volume ?? config.volume;
    gain.gain.value = vol * 0.3;

    const frequencies: Record<string, number> = {
      trip_request: 440, trip_accepted: 880, driver_arriving: 660,
      driver_arrived: 990, trip_started: 770, trip_completed: 1100,
      payment_success: 880, payment_failed: 220, error: 180,
      message_received: 550, button_tap: 1000,
    };

    osc.frequency.value = frequencies[name] || 440;
    osc.type = name === "error" || name === "payment_failed" ? "sawtooth" : "sine";
    osc.start();
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + config.duration / 1000);
    osc.stop(ctx.currentTime + config.duration / 1000);

    setTimeout(() => ctx.close(), config.duration + 100);
  } catch {}
}

export function triggerHaptic(style: HapticStyle): void {
  if (typeof window === "undefined") return;
  try {
    const nav = (navigator as any);
    if (nav.vibrate) {
      const patterns: Record<HapticStyle, number[]> = {
        light: [5], medium: [10, 30, 10], heavy: [15, 50, 15, 50, 15],
        selection: [5, 5, 5], success: [10, 30, 20], warning: [20, 50, 20], error: [30, 50, 30, 50, 30],
      };
      nav.vibrate(patterns[style]);
    }
  } catch {}
}

export function playSoundWithHaptic(name: SoundName, volume?: number): void {
  const config = SOUNDS[name];
  playSound(name, volume);
  if (config?.haptic) triggerHaptic(config.haptic);
}

export function getSoundConfig(name: SoundName): SoundConfig | undefined {
  return SOUNDS[name];
}
