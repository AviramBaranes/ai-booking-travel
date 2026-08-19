import type { Homepage } from "@/payload-types";
import { Metadata } from "next";
import { cache } from "react";
import { notFound } from "next/navigation";
import { Populated } from "@/shared/types/payload";
import { Hero } from "./_components/home/Hero";
import { BlocksRenderer } from "./_components/blocks/BlocksRenderer";
import { HomepageDecorations } from "./_components/decorations/HomepageDecorations";
import {
  SUPPORTED_LANGS,
  type SupportedLang,
} from "@/shared/constants/supported_langs";
import { getCachedPayload } from "@/shared/server/cms";
import { absoluteUrl } from "@/shared/seo/site";
import { ogImages } from "@/shared/seo/metadata";
import { JsonLd } from "@/shared/seo/JsonLd";
import {
  faqSchema,
  organizationSchema,
  websiteSchema,
} from "@/shared/seo/schema";
import {
  buildAlternates,
  DEFAULT_LANG,
  hasTranslation,
  localePath,
  type LocalizedValue,
} from "@/shared/seo/alternates";

type Props = {
  params: Promise<{ lang: string }>;
};

export const revalidate = 3600;
export async function generateStaticParams() {
  const params = SUPPORTED_LANGS.map((locale) => ({ lang: locale }));
  return params;
}

const getHomepage = cache(async (lang: string): Promise<Homepage | null> => {
  const payload = await getCachedPayload();
  const result = await payload.findGlobal({
    slug: "homepage",
    locale: lang as "he" | "en",
    draft: false,
  });

  return result;
});

/** Social profiles for the Organization schema's `sameAs`. */
const getSocialLinks = cache(async (lang: SupportedLang): Promise<string[]> => {
  const payload = await getCachedPayload();
  const footer = await payload.findGlobal({ slug: "footer", locale: lang });
  return (footer?.socialsLinks ?? [])
    .map((social) => social.link)
    .filter((link): link is string => Boolean(link?.startsWith("http")));
});

/**
 * `locale: "all"` applies no fallback, so a missing `title` here means the
 * locale genuinely has no translation rather than silently serving Hebrew.
 */
const getHomepageLocales = cache(async () => {
  const payload = await getCachedPayload();
  return payload.findGlobal({ slug: "homepage", locale: "all", depth: 0 });
});

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { lang } = await params;
  const homepage = await getHomepage(lang);
  if (!homepage) return {};

  const localized = (await getHomepageLocales()) as { title?: LocalizedValue };
  const pathByLang: Partial<Record<SupportedLang, string>> = {};
  for (const code of SUPPORTED_LANGS) {
    if (code === DEFAULT_LANG || hasTranslation(localized.title, code)) {
      pathByLang[code] = localePath(code);
    }
  }

  const title = homepage.meta?.title ?? homepage.title;
  const description =
    homepage.meta?.description ?? homepage.excerpt ?? undefined;

  return {
    title,
    description,
    alternates: buildAlternates(lang as SupportedLang, pathByLang),
    openGraph: {
      title,
      description,
      url: absoluteUrl(localePath(lang as SupportedLang)),
      images: ogImages(homepage.meta?.image, homepage.featuredImage),
    },
  };
}

export default async function Homepage({ params }: Props) {
  const { lang } = await params;
  const homepage = await getHomepage(lang);

  if (!homepage) notFound();

  const image = homepage.featuredImage as Populated<Homepage["featuredImage"]>;

  const socials = await getSocialLinks(lang as SupportedLang);
  const faq = faqSchema(homepage.layout);

  return (
    <main className="relative overflow-hidden">
      <JsonLd data={organizationSchema(lang as SupportedLang, socials)} />
      <JsonLd
        data={websiteSchema(
          lang as SupportedLang,
          absoluteUrl(localePath(lang as SupportedLang)),
        )}
      />
      {faq && <JsonLd data={faq} />}
      <HomepageDecorations />
      <Hero
        lang={lang}
        image={image}
        title={homepage.title}
        subtitle={homepage.subtitle ?? ""}
      />
      <BlocksRenderer
        blocks={homepage.layout}
        faqClassName="w-7/10 max-w-full"
      />
    </main>
  );
}
