export const fontFamily = {
  sans: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'JetBrains Mono', 'Fira Code', monospace",
  display: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
} as const;

export const fontSize = {
  xs: 10,
  sm: 12,
  base: 13,
  md: 14,
  lg: 16,
  xl: 18,
  xxl: 20,
  h3: 22,
  h2: 26,
  h1: 32,
  hero: 42,
  display: 56,
} as const;

export const fontWeight = {
  regular: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
  extrabold: 800,
} as const;

export const lineHeight = {
  tight: 1.15,
  normal: 1.3,
  relaxed: 1.5,
  loose: 1.7,
} as const;

export const letterSpacing = {
  tight: "-0.02em",
  normal: "-0.01em",
  wide: "0.01em",
  wider: "0.04em",
  widest: "0.08em",
} as const;

export const typography = {
  display: {
    fontSize: fontSize.display,
    fontWeight: fontWeight.extrabold,
    lineHeight: lineHeight.tight,
    letterSpacing: letterSpacing.tight,
  },
  hero: {
    fontSize: fontSize.hero,
    fontWeight: fontWeight.bold,
    lineHeight: lineHeight.tight,
    letterSpacing: letterSpacing.tight,
  },
  h1: {
    fontSize: fontSize.h1,
    fontWeight: fontWeight.bold,
    lineHeight: lineHeight.tight,
    letterSpacing: letterSpacing.tight,
  },
  h2: {
    fontSize: fontSize.h2,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.normal,
    letterSpacing: letterSpacing.normal,
  },
  h3: {
    fontSize: fontSize.h3,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.normal,
    letterSpacing: letterSpacing.normal,
  },
  body: {
    fontSize: fontSize.md,
    fontWeight: fontWeight.medium,
    lineHeight: lineHeight.relaxed,
    letterSpacing: letterSpacing.normal,
  },
  bodySm: {
    fontSize: fontSize.base,
    fontWeight: fontWeight.regular,
    lineHeight: lineHeight.relaxed,
    letterSpacing: letterSpacing.normal,
  },
  label: {
    fontSize: fontSize.xs,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.normal,
    letterSpacing: letterSpacing.wider,
    textTransform: "uppercase" as const,
  },
  price: {
    fontSize: fontSize.xxl,
    fontWeight: fontWeight.bold,
    lineHeight: lineHeight.tight,
    letterSpacing: letterSpacing.tight,
  },
  caption: {
    fontSize: fontSize.sm,
    fontWeight: fontWeight.regular,
    lineHeight: lineHeight.relaxed,
    letterSpacing: letterSpacing.normal,
  },
} as const;
