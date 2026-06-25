import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";
import { withPayload } from "@payloadcms/next/withPayload";

const nextConfig: NextConfig = {
  /* config options here */
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**",
      },
    ],
    localPatterns: [
      {
        pathname: "/api/media/file/**",
      },
      {
        // This allows all standard assets from your public folder (like /logo.png)
        pathname: "/**",
      },
    ],
  },
};

const withNextIntl = createNextIntlPlugin();
const config = withPayload(withNextIntl(nextConfig));

delete (
  config.experimental as Record<string, unknown> | undefined
)?.enableServerFastRefresh;

export default config;