import { notFound } from "next/navigation";
import type { Metadata } from "next";
import type { BlogCategory, BlogPost, Form, Media } from "@/payload-types";
import type { Populated } from "@/shared/types/payload";
import Image from "next/image";
import { BlocksRenderer } from "../../_components/blocks/BlocksRenderer";
import { PagesDecorations } from "../../_components/decorations/PagesDecorations";
import { RefreshRouteOnSave as PayloadLivePreview } from "../../_components/LivePreview/RefreshRouteOnSave";
import {
  SUPPORTED_LANGS,
  SupportedLang,
} from "@/shared/constants/supported_langs";
import { getCachedPayload } from "@/shared/server/cms";
import { PayloadFormRenderer } from "@/shared/components/forms/FormRenderer";

type Props = {
  params: Promise<{ lang: string; slug: string }>;
};

const getPost = async (
  slug: string,
  lang: string,
): Promise<BlogPost | null> => {
  const payload = await getCachedPayload();
  const result = await payload.find({
    collection: "blog-posts",
    where: { slug: { equals: slug } },
    locale: lang as SupportedLang,
    draft: false,
    limit: 1,
  });

  return (result.docs[0] as BlogPost) ?? null;
};

export const revalidate = 3600;
export async function generateStaticParams() {
  const payload = await getCachedPayload();
  const params = await Promise.all(
    SUPPORTED_LANGS.map(async (lang) => {
      const result = await payload.find({
        collection: "blog-posts",
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
  const post = await getPost(decodeURIComponent(slug), lang);
  if (!post) return {};
  return {
    title: post.seo?.title ?? post.title,
    description: post.seo?.description ?? post.excerpt ?? undefined,
  };
}

export default async function SlugPage({ params }: Props) {
  const { lang, slug } = await params;
  const post = await getPost(decodeURIComponent(slug), lang);

  if (!post) notFound();

  const image = post.featuredImage as Populated<BlogPost["featuredImage"]>;

  return (
    <>
      <PayloadLivePreview />
      <div className="relative overflow-x-clip">
        <div className="bg-navy py-20 px-72 flex flex-col gap-5.5">
          <p className="type-paragraph text-white/55">
            ראשי / {(post.category as BlogCategory).title}
          </p>
          <h3 className="type-h3 text-white w-2/3">{post.title}</h3>
          <p className="type-paragraph text-white font-semibold">
            {post.excerpt}
          </p>
        </div>
        <PagesDecorations />

        <div className="w-2/3 mx-auto">
          <div className="shadow-card bg-white p-10 flex justify-between gap-24 rounded-b-xl">
            <div className="w-7/10">
              {image?.url && (
                <div className="overflow-hidden rounded-2xl shadow-[0_16px_40px_rgba(15,23,42,0.18)]">
                  <Image
                    src={image.url}
                    alt={image.alt}
                    width={image.width ?? 780}
                    height={image.height ?? 280}
                    className="h-auto w-full object-cover"
                    priority
                  />
                </div>
              )}
              <BlocksRenderer blocks={post.layout} />
              <h5 className="type-h5 text-navy my-5">תגיות:</h5>
              <div className="flex items-start gap-4   flex-wrap">
                {post.tags?.map((tag) => (
                  <span
                    key={tag.id}
                    className="type-paragraph text-navy bg-brand-blue/5 rounded-full px-4 py-1.5"
                  >
                    {tag.tag}
                  </span>
                ))}
              </div>
            </div>
            <div className="w-3/10">
              <div className="sticky top-24 flex flex-col gap-6">
                <>
                  {post.form && (
                    <div className="rounded-xl border border-border-light pb-8">
                      <h6 className="type-h6 bg-navy text-white font-semibold mb-4 rounded-t-xl p-6">
                        {(post.form as Form).title}
                      </h6>
                      <div className="px-4">
                        <PayloadFormRenderer form={post.form as Form} />
                      </div>
                    </div>
                  )}
                  {post.banner && post.banner.image && (
                    <a
                      href={post.banner.link ?? "#"}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="block"
                    >
                      <Image
                        src={(post.banner.image as Media)?.url ?? ""}
                        alt={(post.banner.image as Media)?.alt}
                        width={(post.banner.image as Media)?.width ?? 300}
                        height={(post.banner.image as Media)?.height ?? 250}
                        className="h-auto w-full object-cover rounded-xl"
                        priority
                      />
                    </a>
                  )}
                </>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
