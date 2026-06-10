import { notFound } from "next/navigation";
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

type Props = {
  params: Promise<{ lang: string; slug: string }>;
};

const getPage = async (slug: string, lang: string): Promise<Page | null> => {
  const payload = await getCachedPayload();
  const result = await payload.find({
    collection: "pages",
    where: { slug: { equals: slug } },
    locale: lang as SupportedLang,
    draft: false,
    limit: 1,
  });

  return (result.docs[0] as Page) ?? null;
};

export const revalidate = 3600;
export async function generateStaticParams() {
  const payload = await getCachedPayload();
  const params = await Promise.all(
    SUPPORTED_LANGS.map(async (lang) => {
      const result = await payload.find({
        collection: "pages",
        locale: lang,
        draft: false,
        limit: 1000,
        select: {
          slug: true,
        },
      });

      return result.docs.map((page) => ({
        lang,
        slug: page.slug,
      }));
    }),
  );

  return params.flat();
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { lang, slug } = await params;
  const page = await getPage(decodeURIComponent(slug), lang);
  if (!page) return {};
  return {
    title: page.meta?.title ?? page.title,
    description: page.meta?.description ?? page.excerpt ?? undefined,
  };
}

export default async function SlugPage({ params }: Props) {
  const { lang, slug } = await params;
  const page = await getPage(decodeURIComponent(slug), lang);

  if (!page) notFound();

  const image = page.featuredImage as Populated<Page["featuredImage"]>;

  return (
    <>
      <PayloadLivePreview />
      <main className="relative">
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
          <div className="w-4/10 mx-auto pb-8">
            <h3 className="type-h3 mb-4 pb-8 text-navy">{page.title}</h3>
            <hr />
          </div>
        )}
        <BlocksRenderer blocks={page.layout} />
      </main>
    </>
  );
}
