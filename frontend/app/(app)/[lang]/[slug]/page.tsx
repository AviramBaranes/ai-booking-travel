import { cache } from "react";
import { notFound } from "next/navigation";
import { getPayload } from "payload";
import config from "@payload-config";
import type { Metadata } from "next";
import type { Page } from "@/payload-types";
import type { Populated } from "@/shared/types/payload";
import Image from "next/image";
import { BlocksRenderer } from "../components/blocks/BlocksRenderer";
import { PagesDecorations } from "../components/decorations/PagesDecorations";

type Props = {
  params: Promise<{ lang: string; slug: string }>;
};

const getPage = cache(
  async (slug: string, lang: string): Promise<Page | null> => {
    const payload = await getPayload({ config });
    const result = await payload.find({
      collection: "pages",
      where: { slug: { equals: slug } },
      locale: lang as "he" | "en",
      draft: false,
      limit: 1,
    });

    return (result.docs[0] as Page) ?? null;
  },
);

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
      <BlocksRenderer blocks={page.layout} />
    </main>
  );
}
