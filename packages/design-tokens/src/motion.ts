export const duration = {
  instant: 0,
  hover: 120,
  fade: 180,
  drawer: 250,
  modal: 220,
  toast: 180,
  sheet: 300,
  page: 350,
  hero: 600,
  slow: 800,
} as const;

export const easing = {
  linear: [0, 0, 1, 1] as const,
  easeIn: [0.4, 0, 1, 1] as const,
  easeOut: [0, 0, 0.2, 1] as const,
  easeInOut: [0.4, 0, 0.2, 1] as const,
  cobalt: [0.16, 1, 0.3, 1] as const,
  spring: [0.34, 1.56, 0.64, 1] as const,
  bounce: [0.18, 1.3, 0.4, 1] as const,
  smooth: [0.23, 1, 0.32, 1] as const,
} as const;

export const transition = {
  hover: `all ${duration.hover}ms cubic-bezier(${easing.easeInOut})`,
  fade: `all ${duration.fade}ms cubic-bezier(${easing.easeOut})`,
  drawer: `transform ${duration.drawer}ms cubic-bezier(${easing.easeOut})`,
  modal: `opacity ${duration.modal}ms cubic-bezier(${easing.easeOut})`,
  toast: `all ${duration.toast}ms cubic-bezier(${easing.spring})`,
  sheet: `transform ${duration.sheet}ms cubic-bezier(${easing.easeOut})`,
  spring: `all ${duration.page}ms cubic-bezier(${easing.spring})`,
  fast: `all 0.15s cubic-bezier(${easing.easeOut})`,
  smooth: `all 0.25s cubic-bezier(${easing.smooth})`,
} as const;

export const transitions = transition;

export const zIndex = {
  base: 1,
  dropdown: 100,
  sticky: 200,
  navbar: 300,
  modal: 400,
  sheet: 500,
  toast: 600,
  tooltip: 700,
  overlay: 800,
  max: 9999,
} as const;
