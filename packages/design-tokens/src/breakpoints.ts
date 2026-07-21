export const breakpoints = {
  mobile: 0,
  tablet: 768,
  laptop: 1024,
  desktop: 1280,
  wide: 1440,
} as const;

export const mediaQuery = {
  mobile: `@media (min-width: ${breakpoints.mobile}px)`,
  tablet: `@media (min-width: ${breakpoints.tablet}px)`,
  laptop: `@media (min-width: ${breakpoints.laptop}px)`,
  desktop: `@media (min-width: ${breakpoints.desktop}px)`,
  wide: `@media (min-width: ${breakpoints.wide}px)`,
  mobileOnly: `@media (max-width: ${breakpoints.tablet - 1}px)`,
  tabletOnly: `@media (min-width: ${breakpoints.tablet}px) and (max-width: ${breakpoints.laptop - 1}px)`,
  hover: `@media (hover: hover)`,
  dark: `@media (prefers-color-scheme: dark)`,
  reduceMotion: `@media (prefers-reduced-motion: reduce)`,
} as const;
