/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  typescript: {
    tsconfigPath: './tsconfig.json',
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  },
  outputFileTracingRoot: __dirname,
  
  // Headers for cache control
  headers: async () => [
    {
      source: '/api/:path*',
      headers: [
        { key: 'Cache-Control', value: 'no-store' },
      ],
    },
    {
      source: '/(.*)',
      headers: [
        { key: 'Cache-Control', value: 's-maxage=3600, stale-while-revalidate=86400' },
      ],
    },
  ],
}

module.exports = nextConfig
