const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  allowedDevOrigins: ["64.176.219.221", "localhost", "127.0.0.1"],
  async rewrites() {
    return [
      { source: "/api/:path*", destination: "http://localhost:8000/api/:path*" },
    ];
  },
};
module.exports = nextConfig;
