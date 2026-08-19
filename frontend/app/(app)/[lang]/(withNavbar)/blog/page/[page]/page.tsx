import { notFound } from "next/navigation";
import { cache } from "react";
import type { Metadata } from "next";
import type { BlogPost } from "@/payload-types";
import { POSTS_PER_PAGE } from "@/shared/constants/blog";
import { buildAlternates, localePath } from "@/shared/seo/alternates";
import {
  SUPPORTED_LANGS,
  type SupportedLang,
} from "@/shared/constants/supported_langs";
import { getCachedPayload } from "@/shared/server/cms";

import Image from "next/image";
import { Populated } from "@/shared/types/payload";
import { SectionHeader } from "../../../_components/blocks/SectionHeader";
import { RefreshRouteOnSave } from "../../../_components/LivePreview/RefreshRouteOnSave";
import { getCardImageUrl } from "../../../_components/posts/RelatedPosts";
import Link from "next/link";
import { BlogPagination } from "../_components/BlogPagination";

type FeaturedImage = Populated<BlogPost["featuredImage"]>;
export const revalidate = 3600;

const getBlogSettings = cache(async (lang: SupportedLang) => {
  const payload = await getCachedPayload();
  return payload.findGlobal({ slug: "site-settings", locale: lang });
});

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { lang, page } = await params;
  if (!SUPPORTED_LANGS.includes(lang as SupportedLang)) return {};

  const pageNumber = Number(page);
  const settings = await getBlogSettings(lang as SupportedLang);

  // Paginated pages stay indexable and self-canonical — rel prev/next is
  // deprecated, and canonicalling them all to page 1 hides posts from Google.
  const suffix =
    pageNumber > 1
      ? ` — ${lang === "he" ? `עמוד ${pageNumber}` : `Page ${pageNumber}`}`
      : "";

  const pathByLang = Object.fromEntries(
    SUPPORTED_LANGS.map((code) => [
      code,
      localePath(code, "blog", "page", String(pageNumber)),
    ]),
  ) as Record<SupportedLang, string>;

  return {
    title: `${settings?.title ?? "Blog"}${suffix}`,
    description: settings?.subtitle ?? undefined,
    alternates: buildAlternates(lang as SupportedLang, pathByLang),
  };
}

type Props = {
  params: Promise<{
    lang: string;
    page: string;
  }>;
};

export async function generateStaticParams() {
  const payload = await getCachedPayload();

  const params = await Promise.all(
    SUPPORTED_LANGS.map(async (lang) => {
      const result = await payload.find({
        collection: "blog-posts",
        locale: lang,
        draft: false,
        depth: 0,
        limit: POSTS_PER_PAGE,
        page: 1,
        select: {
          slug: true,
        },
      });

      const totalPages = Math.max(result.totalPages ?? 1, 1);

      return Array.from({ length: totalPages }, (_, index) => ({
        lang,
        page: String(index + 1),
      }));
    }),
  );

  return params.flat();
}

async function getBlogPostsPage(lang: SupportedLang, page: number) {
  const payload = await getCachedPayload();

  const [posts, settings] = await Promise.all([
    payload.find({
      collection: "blog-posts",
      locale: lang,
      draft: false,
      depth: 1,
      limit: POSTS_PER_PAGE,
      page,
      sort: "-publishedAt",
      select: {
        title: true,
        slug: true,
        excerpt: true,
        featuredImage: true,
        publishedAt: true,
      },
    }),

    payload.findGlobal({
      slug: "site-settings",
      locale: lang,
    }),
  ]);

  return {
    posts: posts.docs as BlogPost[],
    currentPage: posts.page ?? page,
    totalPages: posts.totalPages ?? 1,
    hasNextPage: posts.hasNextPage,
    hasPrevPage: posts.hasPrevPage,
    nextPage: posts.nextPage,
    prevPage: posts.prevPage,

    featuredImage: settings?.featuredImage,
    pillText: settings?.pillText,
    title: settings?.title,
    subtitle: settings?.subtitle,
  };
}

export default async function BlogPaginatedPage({ params }: Props) {
  const { lang, page } = await params;

  if (!SUPPORTED_LANGS.includes(lang as SupportedLang)) {
    notFound();
  }

  const currentPage = Number(page);

  if (!Number.isInteger(currentPage) || currentPage < 1) {
    notFound();
  }

  const data = await getBlogPostsPage(lang as SupportedLang, currentPage);

  if (currentPage > data.totalPages) {
    notFound();
  }

  const image = data.featuredImage as Populated<BlogPost["featuredImage"]>;
  const { pillText, title, subtitle, posts } = data;

  return (
    <>
      <RefreshRouteOnSave />
      <div>
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

        <div className="my-10 mx-5">
          <SectionHeader
            pillText={pillText}
            title={title}
            subtitle={subtitle}
          />
        </div>

        <div className="flex flex-wrap flex-col lg:flex-row justify-start lg:w-2/3 lg:mx-auto mx-5">
          {posts.map((post) => (
            <div key={post.id} className="lg:w-1/4 items-stretch flex">
              {post.featuredImage && (
                <div className="p-4 shadow-card m-3 rounded-xl border border-border flex flex-col justify-between gap-4">
                  <div className="relative aspect-275/195 w-full overflow-hidden rounded-xl">
                    <Image
                      src={getCardImageUrl(post)}
                      alt={(post.featuredImage as FeaturedImage).alt}
                      fill
                      className="object-cover"
                    />
                  </div>
                  <h6 className="type-h6 text-navy">{post.title}</h6>
                  <p className="type-paragraph text-text-secondary">
                    {post.excerpt}
                  </p>
                  <Link href={`/${lang}/blog/${post.slug}`}>
                    <span className="type-paragraph text-brand-blue font-semibold hover:underline">
                      קרא עוד &gt;
                    </span>
                  </Link>
                </div>
              )}
            </div>
          ))}
        </div>
        <div className="my-5">
          <BlogPagination
            lang={lang}
            currentPage={data.currentPage}
            totalPages={data.totalPages}
          />
        </div>
      </div>
    </>
  );
}
