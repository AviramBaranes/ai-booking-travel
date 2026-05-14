import config from "@payload-config";
import type { Homepage } from "@/payload-types";
import { Metadata } from "next";
import { getPayload } from "payload";
import { cache } from "react";
import { notFound } from "next/navigation";
import { Populated } from "@/shared/types/payload";
import { Hero } from "./_components/home/Hero";
import { BlocksRenderer } from "./_components/blocks/BlocksRenderer";
import { HomepageDecorations } from "./_components/decorations/HomepageDecorations";

type Props = {
  params: Promise<{ lang: string }>;
};

const getHomepage = cache(async (lang: string): Promise<Homepage | null> => {
  const payload = await getPayload({ config });
  const result = await payload.findGlobal({
    slug: "homepage",
    locale: lang as "he" | "en",
    draft: false,
  });

  return result;
});

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { lang } = await params;
  const homepage = await getHomepage(lang);
  if (!homepage) return {};
  return {
    title: homepage.meta?.title ?? homepage.title,
    description: homepage.meta?.description ?? homepage.excerpt ?? undefined,
  };
}

export default async function Homepage({ params }: Props) {
  const { lang } = await params;
  const homepage = await getHomepage(lang);

  if (!homepage) notFound();

  const image = homepage.featuredImage as Populated<Homepage["featuredImage"]>;

  return (
    <main className="relative overflow-hidden">
      <HomepageDecorations />
      <Hero
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
