import { notFound } from "next/navigation";
import { cache } from "react";
import type { Metadata } from "next";
import type { Page } from "@/payload-types";
import type { Populated } from "@/shared/types/payload";
import Image from "next/image";
import { BlocksRenderer } from "../_components/blocks/BlocksRenderer";
import { PagesDecorations } from "../_components/decorations/PagesDecorations";
import { RefreshRouteOnSave as PayloadLivePreview } from "../_components/LivePreview/RefreshRouteOnSave";
import {
  SUPPORTED_LANGS,
  SupportedLang,
} from "@/shared/constants/supported_langs";
import { getCachedPayload } from "@/shared/server/cms";
import { absoluteUrl, SITE_NAME } from "@/shared/seo/site";
import { ogImages } from "@/shared/seo/metadata";
import { JsonLd } from "@/shared/seo/JsonLd";
import { breadcrumbSchema, faqSchema } from "@/shared/seo/schema";
import {
  buildAlternates,
  hasTranslation,
  localePath,
  localizedPaths,
  type LocalizedValue,
} from "@/shared/seo/alternates";

type Props = {
  params: Promise<{ lang: string; slug: string }>;
};

const getPage = cache(async (slug: string, lang: string): Promise<Page | null> => {
  const payload = await getCachedPayload();
  const result = await payload.find({
    collection: "pages",
    where: { slug: { equals: slug } },
    locale: lang as SupportedLang,
    draft: false,
    limit: 1,
  });

  return (result.docs[0] as Page) ?? null;
});

/**
 * Sibling slugs for hreflang. Queried with `locale: "all"`, which applies no
 * fallback — the only way to tell a real English page from one that would
 * silently serve Hebrew.
 */
const getPageSlugs = cache(async (id: number) => {
  const payload = await getCachedPayload();
  const doc = await payload.findByID({
    collection: "pages",
    id,
    locale: "all",
    depth: 0,
  });
  return (doc as unknown as { slug?: LocalizedValue }).slug;
});

export const revalidate = 3600;
export async function generateStaticParams() {
  const payload = await getCachedPayload();
  // `locale: "all"` applies no fallback. Querying per-locale instead returns
  // the Hebrew slug for untranslated pages, and those /en/<hebrew-slug> params
  // prerender as 404s — the where clause in getPage does not fall back, so the
  // lookup finds nothing.
  const result = await payload.find({
    collection: "pages",
    locale: "all",
    draft: false,
    limit: 1000,
    depth: 0,
    select: {
      slug: true,
    },
  });

  return result.docs.flatMap((page) => {
    const slugs = (page as unknown as { slug?: LocalizedValue }).slug;
    return SUPPORTED_LANGS.filter((lang) => hasTranslation(slugs, lang)).map(
      (lang) => ({ lang, slug: slugs![lang]!.trim() }),
    );
  });
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { lang, slug } = await params;
  const page = await getPage(decodeURIComponent(slug), lang);
  if (!page) return {};

  const slugs = await getPageSlugs(page.id);
  const pathByLang = localizedPaths(slugs, (code, pageSlug) =>
    localePath(code, pageSlug),
  );

  const title = page.meta?.title ?? page.title;
  const description = page.meta?.description ?? page.excerpt ?? undefined;

  return {
    title,
    description,
    // For an untranslated page this canonicals /en/… back to the Hebrew
    // original instead of letting the two compete as duplicates.
    alternates: buildAlternates(lang as SupportedLang, pathByLang),
    openGraph: {
      title,
      description,
      type: "article",
      url: pathByLang[lang as SupportedLang]
        ? absoluteUrl(pathByLang[lang as SupportedLang]!)
        : undefined,
      images: ogImages(page.meta?.image, page.featuredImage),
    },
  };
}

export default async function SlugPage({ params }: Props) {
  const { lang, slug } = await params;
  const page = await getPage(decodeURIComponent(slug), lang);

  if (!page) notFound();

  const image = page.featuredImage as Populated<Page["featuredImage"]>;

  const faq = faqSchema(page.layout);
  const breadcrumbs = breadcrumbSchema([
    { name: SITE_NAME, url: absoluteUrl(localePath(lang as SupportedLang)) },
    {
      name: page.title,
      url: absoluteUrl(localePath(lang as SupportedLang, page.slug)),
    },
  ]);

  return (
    <>
      <PayloadLivePreview />
      <JsonLd data={breadcrumbs} />
      {faq && <JsonLd data={faq} />}
      <div className="relative">
        {page.includeBgDecorations && <PagesDecorations />}
        {image?.url && (
          <Image
            src={image.url}
            alt={image.alt}
            width={image.width ?? 1200}
            height={image.height ?? 630}
            style={{ width: "100%", height: "auto" }}
            priority
          />
        )}
        {page.renderTitle && (
          <div className="lg:w-4/10 lg:mx-auto mx-5 pb-8 mt-10">
            <h3 className="type-h3 mb-4 pb-8 text-navy">{page.title}</h3>
            <hr />
          </div>
        )}
        <BlocksRenderer blocks={page.layout} />
      </div>
    </>
  );
}
