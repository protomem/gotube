/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  rewrites: async () => {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/:path*", // TODO: Move to .env
      },
    ];
  },
};

export default nextConfig;
