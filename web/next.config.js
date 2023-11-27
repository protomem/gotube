/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  rewrites: async () => {
    return [
      {
        source: "/api/:path*",
        destination: `http://${process.env.API_ADDR}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
