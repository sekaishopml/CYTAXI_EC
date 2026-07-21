export const shadows = {
  none: "none",
  card: "0 1px 3px rgba(15,20,25,0.08), 0 1px 2px rgba(15,20,25,0.04)",
  cardHover: "0 4px 12px rgba(15,20,25,0.10), 0 2px 4px rgba(15,20,25,0.06)",
  button: "0 2px 8px rgba(15,20,25,0.12)",
  buttonGreen: "0 4px 14px rgba(0,179,107,0.3)",
  float: "0 8px 24px rgba(15,20,25,0.12)",
  floatHover: "0 12px 32px rgba(15,20,25,0.16)",
  glass: "0 8px 32px rgba(15,20,25,0.10), inset 0 1px 0 rgba(255,255,255,0.5)",
  modal: "0 24px 48px rgba(15,20,25,0.20)",
  dropdown: "0 4px 16px rgba(15,20,25,0.12)",
  toast: "0 4px 16px rgba(15,20,25,0.14)",
  sheet: "0 -4px 24px rgba(15,20,25,0.10)",
  limeGlow: "0 0 30px rgba(0,179,107,0.25)",
  limeGlowSubtle: "0 0 20px rgba(0,179,107,0.15)",
} as const;

export const blur = {
  none: 0,
  sm: "blur(4px)",
  md: "blur(12px)",
  lg: "blur(24px)",
  xl: "blur(40px)",
} as const;
