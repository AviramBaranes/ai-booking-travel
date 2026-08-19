import type { MetadataRoute } from "next";
import { getCachedPayload } from "@/shared/server/cms";
import {
  SUPPORTED_LANGS,
  type SupportedLang,
} from "@/shared/constants/supported_langs";
import { POSTS_PER_PAGE } from "@/shared/constants/blog";
import { absoluteUrl } from "@/shared/seo/site";
import {
  DEFAULT_LANG,
  hasTranslation,
  localePath,
  localizedPaths,
  type LocalizedValue,
} from "@/shared/seo/alternates";

export const revalidate = 3600;

/**
 * Shape returned by a `locale: "all"` query: localized fields come back keyed
 * by locale with no fallback applied, which is what lets us tell a real
 * translation from a Hebrew value bleeding into `/en/`.
 */
type AllLocalesDoc = {
  slug?: LocalizedValue;
  updatedAt?: string;
};

type EntryOptions = {
  lastModified?: string;
  changeFrequency?: MetadataRoute.Sitemap[number]["changeFrequency"];
  priority?: number;
};

/**
 * One entry per translated locale, each carrying the same `languages` map.
 * hreflang must be reciprocal — every URL referenced as an alternate has to
 * appear in the sitemap as its own `<url>` and point back.
 */
function toEntries(
  pathByLang: Partial<Record<SupportedLang, string>>,
  options: EntryOptions = {},
): MetadataRoute.Sitemap {
  const present = SUPPORTED_LANGS.filter((lang) => pathByLang[lang]);
  if (present.length === 0) return [];

  const languages: Record<string, string> | undefined =
    present.length > 1
      ? {
          ...Object.fromEntries(
            present.map((lang) => [lang, absoluteUrl(pathByLang[lang]!)]),
          ),
          "x-default": absoluteUrl(
            pathByLang[DEFAULT_LANG] ?? pathByLang[present[0]]!,
          ),
        }
      : undefined;

  return present.map((lang) => ({
    url: absoluteUrl(pathByLang[lang]!),
    ...options,
    ...(languages ? { alternates: { languages } } : {}),
  }));
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const payload = await getCachedPayload();

  const [homepage, pages, posts] = await Promise.all([
    payload.findGlobal({ slug: "homepage", locale: "all", depth: 0 }),
    payload.find({
      collection: "pages",
      locale: "all",
      depth: 0,
      pagination: false,
      select: { slug: true, updatedAt: true },
    }),
    payload.find({
      collection: "blog-posts",
      locale: "all",
      depth: 0,
      pagination: false,
      select: { slug: true, updatedAt: true },
    }),
  ]);

  const entries: MetadataRoute.Sitemap = [];

  // Homepage — the `homepage` global rendered at /{lang}
  const homeTitle = (homepage as { title?: LocalizedValue }).title;
  const homePaths: Partial<Record<SupportedLang, string>> = {};
  for (const lang of SUPPORTED_LANGS) {
    if (lang === DEFAULT_LANG || hasTranslation(homeTitle, lang)) {
      homePaths[lang] = localePath(lang);
    }
  }
  entries.push(
    ...toEntries(homePaths, {
      lastModified: (homepage as { updatedAt?: string }).updatedAt ?? undefined,
      changeFrequency: "weekly",
      priority: 1,
    }),
  );

  // CMS pages
  for (const doc of pages.docs as AllLocalesDoc[]) {
    const paths = localizedPaths(doc.slug, (lang, slug) =>
      localePath(lang, slug),
    );
    // The homepage is a global rendered at /{lang}; a "home"-slugged page
    // would duplicate it.
    if (Object.values(paths).some((path) => path.endsWith("/home"))) continue;

    entries.push(
      ...toEntries(paths, {
        lastModified: doc.updatedAt,
        changeFrequency: "monthly",
        priority: 0.8,
      }),
    );
  }

  // Blog listing — /{lang}/blog/page/{n}. The listing renders in both locales
  // regardless of per-post translation, so pair them unconditionally.
  const totalPages = Math.max(Math.ceil(posts.docs.length / POSTS_PER_PAGE), 1);
  for (let page = 1; page <= totalPages; page++) {
    const paths = Object.fromEntries(
      SUPPORTED_LANGS.map((lang) => [
        lang,
        localePath(lang, "blog", "page", String(page)),
      ]),
    ) as Record<SupportedLang, string>;

    entries.push(
      ...toEntries(paths, {
        changeFrequency: "weekly",
        priority: page === 1 ? 0.7 : 0.4,
      }),
    );
  }

  // Blog posts
  for (const doc of posts.docs as AllLocalesDoc[]) {
    const paths = localizedPaths(doc.slug, (lang, slug) =>
      localePath(lang, "blog", slug),
    );
    entries.push(
      ...toEntries(paths, {
        lastModified: doc.updatedAt,
        changeFrequency: "monthly",
        priority: 0.6,
      }),
    );
  }

  return entries;
}
