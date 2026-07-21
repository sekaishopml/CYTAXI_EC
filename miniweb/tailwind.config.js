/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", "-apple-system", "BlinkMacSystemFont", "Segoe UI", "Roboto", "sans-serif"],
        mono: ["JetBrains Mono", "Fira Code", "monospace"],
      },
      colors: {
        brand: {
          green: "#006c49",
          "green-light": "#00a152",
          "green-bright": "#00e676",
          "green-dark": "#004d35",
        },
        surface: {
          glass: "rgba(255,255,255,0.88)",
          card: "rgba(255,255,255,0.7)",
          bg: "#f9f9f9",
          skeleton: "#eee",
        },
      },
      borderRadius: {
        pill: "20px",
      },
      boxShadow: {
        card: "0 1px 3px rgba(0,0,0,0.04), 0 6px 20px rgba(0,0,0,0.04)",
        button: "0 2px 8px rgba(0,0,0,0.12)",
        float: "0 8px 32px rgba(0,0,0,0.12)",
      },
      animation: {
        "fade-in": "fadeIn 0.35s ease-out",
        "slide-up": "slideUp 0.35s ease-out",
        "scale-in": "scaleIn 0.4s ease-out",
        pulse: "pulse 2s cubic-bezier(0.4,0,0.6,1) infinite",
        shimmer: "shimmer 1.5s ease-in-out infinite",
        spin: "spin 0.7s linear infinite",
      },
      keyframes: {
        fadeIn: { "0%": { opacity: 0 }, "100%": { opacity: 1 } },
        slideUp: { "0%": { opacity: 0, transform: "translateY(24px)" }, "100%": { opacity: 1, transform: "translateY(0)" } },
        scaleIn: { "0%": { opacity: 0, transform: "scale(0.92)" }, "100%": { opacity: 1, transform: "scale(1)" } },
        pulse: { "0%,100%": { opacity: 1, transform: "scale(1)" }, "50%": { opacity: 0.6, transform: "scale(1.05)" } },
        shimmer: { "0%": { backgroundPosition: "200% 0" }, "100%": { backgroundPosition: "-200% 0" } },
        spin: { to: { transform: "rotate(360deg)" } },
      },
    },
  },
  plugins: [],
};
