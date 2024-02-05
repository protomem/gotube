/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  rewrites: async () => {
    return [
      {
        source: "/api/:path*",
        destination: `${process.env.API_ADDR}/:path*`,
      },
    ];
  },
};

export default nextConfig;
