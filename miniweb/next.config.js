const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  allowedDevOrigins: ["64.176.219.221"],
  async rewrites() {
    return [
      { source: "/api/:path*", destination: "http://64.176.219.221/api/:path*" },
    ];
  },
};
module.exports = nextConfig;
