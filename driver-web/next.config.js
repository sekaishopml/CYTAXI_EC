/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: [
    "@cytaxi/ui",
    "@cytaxi/api-client",
    "@cytaxi/design-tokens",
    "@cytaxi/fonts",
  ],
  async rewrites() {
    return [{ source: "/api/:path*", destination: "http://localhost:8000/api/:path*" }];
  },
};
module.exports = nextConfig;
