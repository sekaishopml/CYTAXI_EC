export const transitions = {
  fast: "all 0.15s cubic-bezier(0.4, 0, 0.2, 1)",
  smooth: "all 0.25s cubic-bezier(0.4, 0, 0.2, 1)",
  spring: "all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)",
} as const;

export const easing = {
  standard: [0.4, 0, 0.2, 1] as readonly number[],
  spring: [0.34, 1.56, 0.64, 1] as readonly number[],
  decelerate: [0.16, 1, 0.3, 1] as readonly number[],
} as const;
