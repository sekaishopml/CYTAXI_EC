const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  allowedDevOrigins: ["64.176.219.221", "localhost", "127.0.0.1"],
  env: {
    NEXT_PUBLIC_GOOGLE_MAPS_KEY: process.env.GOOGLE_MAPS_KEY || "AIzaSyDtY41w4LgAdPi5tXcNr_75Pq7hmFrXgxU",
  },
  // Prevent caching on mobile browsers
  async headers() {
    return [
      { source: "/:path*", headers: [
        { key: "Cache-Control", value: "no-cache, no-store, must-revalidate" },
        { key: "Pragma", value: "no-cache" },
        { key: "Expires", value: "0" },
      ]},
    ];
  },
  async rewrites() {
    return [
      { source: "/api/:path*", destination: "http://localhost:8000/api/:path*" },
    ];
  },
};
module.exports = nextConfig;
