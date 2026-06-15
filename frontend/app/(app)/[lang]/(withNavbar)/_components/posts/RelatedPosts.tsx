import { BlogPost } from "@/payload-types";
import { SectionHeader } from "../blocks/SectionHeader";
import Image from "next/image";
import { Populated } from "@/shared/types/payload";
import Link from "next/link";
import clsx from "clsx";
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
      <div
        className={clsx("flex flex-wrap", {
          "justify-start": posts.length < 4,
          "justify-between": posts.length >= 4,
        })}
      >
        {posts.map((post) => (
          <div key={post.id} className="w-1/4 items-stretch flex">
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
        {showButton && (
          <Link
            href={`/${lang}/blog`}
            className="text-center w-full"
          >
            <Button variant="outline" className="mx-auto mt-6 rounded-md p-5">
              לכל הכתבות
            </Button>
          </Link>
        )}
      </div>
    </>
  );
}
