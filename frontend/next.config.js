// next.config.js
/** @type {import('next').NextConfig} */
module.exports = {
  images: {
    // allow Next/Image to fetch from your Go backend
    remotePatterns: [
      {
        protocol: "http",
        hostname: "localhost",
        port: "8080",
        pathname: "/uploads/**",
      },
    ],
  },
};
