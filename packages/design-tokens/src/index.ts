export { colors } from "./colors";
export { spacing, radius } from "./spacing";
export { fontFamily, fontSize, fontWeight, lineHeight, letterSpacing, typography } from "./typography";
export { shadows, blur } from "./shadows";
export { breakpoints, mediaQuery } from "./breakpoints";
export { duration, easing, transition, transitions, zIndex } from "./motion";

export const glass = {
  surface: "rgba(255,255,255,0.88)",
  blur: "blur(24px) saturate(180%)",
  border: "1px solid rgba(0,0,0,0.04)",
} as const;
