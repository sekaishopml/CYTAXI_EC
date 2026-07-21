// CYTAXI Design System — Single source of truth for all visual tokens

export const colors = {
  // Brand
  green: "#006c49",
  greenLight: "#00a152",
  greenBright: "#00e676",
  greenDark: "#004d35",

  // Primary text
  textPrimary: "#191c1e",
  textSecondary: "#3c4a42",
  textMuted: "#8a8a8a",

  // Status
  blue: "#448aff",
  blueLight: "#82b1ff",
  blueDark: "#1565c0",
  red: "#ba1a1a",
  redBg: "#ffdaD6",

  // Surfaces
  surfaceWhite: "#ffffff",
  surfaceGlass: "rgba(255,255,255,0.88)",
  surfaceCard: "rgba(255,255,255,0.7)",
  surfaceBg: "#f9f9f9",
  surfaceDark: "#121212",
  surfaceSkeleton: "#eee",

  // Borders
  borderLight: "rgba(0,0,0,0.04)",
  borderMedium: "rgba(0,0,0,0.08)",
  borderGreen: "rgba(0,108,73,0.2)",

  // Map
  mapDark: "#1c1c1e",
} as const;

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 20,
  xxl: 24,
} as const;

export const radius = {
  sm: 8,
  md: 12,
  lg: 14,
  xl: 16,
  pill: 20,
  full: 9999,
} as const;

export const shadows = {
  card: "0 1px 3px rgba(0,0,0,0.04), 0 6px 20px rgba(0,0,0,0.04)",
  button: "0 2px 8px rgba(0,0,0,0.12)",
  buttonGreen: "0 4px 20px rgba(0,108,73,0.2)",
  float: "0 8px 32px rgba(0,0,0,0.12)",
  glass: "0 8px 40px rgba(0,0,0,0.15), inset 0 1px 0 rgba(255,255,255,0.9)",
} as const;

export const typography = {
  h1: { fontSize: 19, fontWeight: 700, letterSpacing: "-0.02em" } as const,
  h2: { fontSize: 16, fontWeight: 600, letterSpacing: "-0.01em" } as const,
  body: { fontSize: 14, fontWeight: 500, letterSpacing: "-0.01em" } as const,
  bodySm: { fontSize: 13, fontWeight: 400 } as const,
  label: { fontSize: 10, fontWeight: 600, letterSpacing: "0.04em", textTransform: "uppercase" } as const,
  price: { fontSize: 22, fontWeight: 700, letterSpacing: "-0.02em" } as const,
} as const;

export const glass = {
  surface: "rgba(255,255,255,0.88)",
  blur: "blur(24px) saturate(180%)",
  border: "1px solid rgba(0,0,0,0.04)",
} as const;

export const transitions = {
  fast: "all 0.15s cubic-bezier(0.4, 0, 0.2, 1)",
  smooth: "all 0.25s cubic-bezier(0.4, 0, 0.2, 1)",
  spring: "all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)",
} as const;
