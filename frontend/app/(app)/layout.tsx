import type { Metadata } from "next";
import { SITE_URL } from "@/shared/seo/site";

/**
 * Set here rather than only on the `[lang]` layout so the admin, accounting
 * and not-found routes also resolve relative image URLs instead of warning at
 * build time and falling back to localhost.
 */
export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return children;
}
