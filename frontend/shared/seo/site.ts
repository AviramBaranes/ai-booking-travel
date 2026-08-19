/**
 * Single source of truth for the public origin.
 *
 * Everything derives from one variable, `NEXT_PUBLIC_SITE_URL`, so this works
 * on any host — no platform-specific variables. Note that `NODE_ENV` is not
 * usable here: `next build` sets it to "production" for the stage build too,
 * so it cannot tell the canonical site apart from a staging copy.
 *
 * Set per environment:
 *   local     unset (defaults to http://localhost:3000)
 *   stage     https://<stage-host>
 *   production https://aibookingtravel.com
 */

/** The one origin allowed to be indexed. */
const PRODUCTION_URL = "https://aibookingtravel.com";

const stripTrailingSlash = (url: string) => url.replace(/\/$/, "");

export const SITE_URL = stripTrailingSlash(
  process.env.NEXT_PUBLIC_SITE_URL?.trim() || "http://localhost:3000",
);

/**
 * Only the canonical production origin may be indexed. Any other origin —
 * stage, a preview build, localhost, or a deploy that forgot to set the
 * variable — serves `Disallow: /` and `noindex`.
 *
 * This fails safe in both directions: a misconfigured stage can never claim to
 * be production, and a misconfigured production is merely absent from the
 * index rather than competing with itself.
 */
export const IS_INDEXABLE = SITE_URL === PRODUCTION_URL;

if (
  process.env.NODE_ENV === "production" &&
  !process.env.NEXT_PUBLIC_SITE_URL
) {
  console.warn(
    "[seo] NEXT_PUBLIC_SITE_URL is not set — falling back to localhost. " +
      "Canonical URLs will be wrong and the site will not be indexable.",
  );
}

export const absoluteUrl = (path: string) =>
  `${SITE_URL}/${path.replace(/^\//, "")}`;

/** Brand name used in title templates, OpenGraph and Organization JSON-LD. */
export const SITE_NAME = "AI Booking Travel";
