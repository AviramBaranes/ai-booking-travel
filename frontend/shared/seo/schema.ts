import "server-only";
import { convertLexicalToPlaintext } from "@payloadcms/richtext-lexical/plaintext";
import type { BlogPost, FAQBlock, Media } from "@/payload-types";
import type { SupportedLang } from "@/shared/constants/supported_langs";
import { CONTACT_PHONE } from "@/shared/constants/contact";
import { absoluteUrl, SITE_NAME, SITE_URL } from "./site";

type Schema = Record<string, unknown>;

const localeTag = (lang: SupportedLang) => (lang === "he" ? "he-IL" : "en-US");

/** Absolute URL for a Payload upload field, or undefined at depth 0. */
function mediaUrl(image: number | Media | null | undefined): string | undefined {
  if (!image || typeof image !== "object" || !image.url) return undefined;
  return image.url.startsWith("http") ? image.url : absoluteUrl(image.url);
}

export function organizationSchema(
  lang: SupportedLang,
  sameAs: string[] = [],
): Schema {
  return {
    "@context": "https://schema.org",
    "@type": "Organization",
    "@id": `${SITE_URL}/#organization`,
    name: SITE_NAME,
    url: SITE_URL,
    logo: absoluteUrl("/logo.png"),
    inLanguage: localeTag(lang),
    ...(sameAs.length ? { sameAs } : {}),
    contactPoint: [
      {
        "@type": "ContactPoint",
        telephone: CONTACT_PHONE,
        contactType: "customer service",
        availableLanguage: ["he", "en"],
      },
    ],
  };
}

export function websiteSchema(lang: SupportedLang, homeUrl: string): Schema {
  return {
    "@context": "https://schema.org",
    "@type": "WebSite",
    "@id": `${SITE_URL}/#website`,
    name: SITE_NAME,
    url: homeUrl,
    inLanguage: localeTag(lang),
    publisher: { "@id": `${SITE_URL}/#organization` },
  };
}

export function blogPostingSchema(
  post: BlogPost,
  url: string,
  lang: SupportedLang,
): Schema {
  return {
    "@context": "https://schema.org",
    "@type": "BlogPosting",
    headline: post.seo?.title ?? post.title,
    description: post.seo?.description ?? post.excerpt ?? undefined,
    image: mediaUrl(post.seo?.image) ?? mediaUrl(post.featuredImage),
    // `publishedAt` is populated by a hook gated on `_status`, which never
    // fires because drafts are not enabled — fall back to createdAt.
    datePublished: post.publishedAt ?? post.createdAt,
    dateModified: post.updatedAt,
    inLanguage: localeTag(lang),
    mainEntityOfPage: { "@type": "WebPage", "@id": url },
    author: { "@id": `${SITE_URL}/#organization` },
    publisher: { "@id": `${SITE_URL}/#organization` },
  };
}

export function breadcrumbSchema(
  items: { name: string; url: string }[],
): Schema {
  return {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: items.map((item, index) => ({
      "@type": "ListItem",
      position: index + 1,
      name: item.name,
      item: item.url,
    })),
  };
}

/** Any Payload blocks array — pages, posts and the homepage each have a different union. */
type LayoutBlocks = readonly { blockType: string }[] | null | undefined;

/** Avoids importing `lexical` directly; it is only a transitive dependency. */
type LexicalData = Parameters<typeof convertLexicalToPlaintext>[0]["data"];

/**
 * FAQPage schema from the `faq` blocks on a page.
 *
 * Google requires the markup to match content visible on the page, so this
 * only reads blocks that BlocksRenderer actually renders. Returns undefined
 * when there is no FAQ, so callers can skip the script entirely.
 */
export function faqSchema(...layouts: LayoutBlocks[]): Schema | undefined {
  const faqBlocks = layouts
    .flatMap((layout) => layout ?? [])
    .filter((block): block is FAQBlock & { blockType: "faq" } =>
      "blockType" in block ? block.blockType === "faq" : false,
    );

  const entities = faqBlocks
    .flatMap((block) => block.categories ?? [])
    .flatMap((category) => category.items ?? [])
    .map((item) => {
      const answer = convertLexicalToPlaintext({
        data: item.answer as unknown as LexicalData,
      }).trim();
      if (!item.question || !answer) return null;
      return {
        "@type": "Question",
        name: item.question,
        acceptedAnswer: { "@type": "Answer", text: answer },
      };
    })
    .filter((entity): entity is NonNullable<typeof entity> => entity !== null);

  if (entities.length === 0) return undefined;

  return {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: entities,
  };
}
