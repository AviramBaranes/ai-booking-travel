import type { MetadataRoute } from "next";
import { IS_INDEXABLE, SITE_URL } from "@/shared/seo/site";

export default function robots(): MetadataRoute.Robots {
  // Stage and previews must never be indexed alongside production.
  if (!IS_INDEXABLE) {
    return { rules: { userAgent: "*", disallow: "/" } };
  }

  return {
    rules: [
      {
        userAgent: "*",
        disallow: [
          // Back office and CMS
          "/cms",
          "/api/",
          "/admin",
          "/accounting",
          // Booking funnel — query-param driven, no standalone SEO value
          "/*/results",
          "/*/plans",
          "/*/order",
          // Customer surfaces. /*/offers/ is intentionally absent: WhatsApp
          // respects robots.txt, so blocking it would kill the link preview.
          // It uses noindex instead.
          "/*/reservations",
          "/*/price-offers",
          "/*/profile",
          "/*/password-reset",
        ],
      },
    ],
    sitemap: `${SITE_URL}/sitemap.xml`,
    host: SITE_URL,
  };
}
