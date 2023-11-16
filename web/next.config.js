/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Dest from env
  rewrites: async () => {
    return [
      {
        source: "/api/:path*",
        destination: process.env.API_ADDR + "/api/:path*",
      },
    ];
  },
};

module.exports = nextConfig;
