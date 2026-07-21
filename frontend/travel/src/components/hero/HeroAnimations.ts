export const HERO_ANIMATIONS = {
  map: {
    initial: { opacity: 0, scale: 0.95 },
    animate: { opacity: 1, scale: 1 },
    transition: { duration: 0.6, ease: [0.16, 1, 0.3, 1] },
  },
  modules: {
    initial: { opacity: 0, y: 20 },
    animate: { opacity: 1, y: 0 },
    transition: { duration: 0.5, delay: 0.2, ease: [0.16, 1, 0.3, 1] },
  },
  card: {
    hover: { scale: 1.02, y: -3 },
    tap: { scale: 0.98 },
  },
} as const;
