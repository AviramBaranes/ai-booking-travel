import { BlogPost } from "@/payload-types";
import { SectionHeader } from "../blocks/SectionHeader";
import Image from "next/image";
import { Populated } from "@/shared/types/payload";
import Link from "next/link";
import { Button } from "@/components/ui/button";

interface RelatedPostsProps {
  pillText: string;
  title: string;
  subtitle?: string;
  posts: BlogPost[];
  lang: string;
  showButton?: boolean;
}

type FeaturedImage = Populated<BlogPost["featuredImage"]>;

export function getCardImageUrl(post: BlogPost): string {
  return (
    (post.featuredImage as FeaturedImage).sizes?.blogCard?.url ||
    (post.featuredImage as FeaturedImage).url ||
    ""
  );
}

export function RelatedPosts({
  pillText,
  title,
  subtitle,
  posts,
  lang,
  showButton = true,
}: RelatedPostsProps) {
  if (posts.length === 0) return null;

  return (
    <>
      <div className="my-10">
        <SectionHeader pillText={pillText} title={title} subtitle={subtitle} />
      </div>
      <div className="flex w-full flex-nowrap gap-6 overflow-x-auto pb-4">
        {posts.map((post) => (
          <div
            key={post.id}
            className="flex shrink-0 basis-[85%] items-stretch sm:basis-90 lg:basis-[calc((100%-4.5rem)/4)]"
          >
            {post.featuredImage && (
              <div className="flex h-full flex-col justify-between gap-4 rounded-xl border border-border p-4 shadow-card">
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
                  <span className="type-paragraph font-semibold text-brand-blue hover:underline">
                    קרא עוד &gt;
                  </span>
                </Link>
              </div>
            )}
          </div>
        ))}
      </div>
      {showButton && (
        <Link href={`/${lang}/blog`} className="text-center w-full">
          <Button variant="outline" className="mx-auto mt-6 rounded-md p-5">
            {lang === "he" ? "לכל הכתבות" : "All posts"}
          </Button>
        </Link>
      )}
    </>
  );
}
