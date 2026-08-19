import type { Metadata } from "next";
import type { Media } from "@/payload-types";

/** A Payload upload field: a raw id at depth 0, a populated doc at depth ≥ 1. */
type MaybeMedia = number | Media | null | undefined;

/**
 * First populated candidate wins — e.g. `ogImages(page.meta?.image,
 * page.featuredImage)` prefers the editor's dedicated share image and falls
 * back to the hero.
 *
 * Media URLs may be relative (`/api/media/file/…`); `metadataBase` on the root
 * layout resolves them to absolute for share previews.
 */
export function ogImages(
  ...candidates: MaybeMedia[]
): NonNullable<Metadata["openGraph"]>["images"] | undefined {
  for (const candidate of candidates) {
    if (candidate && typeof candidate === "object" && candidate.url) {
      return [
        {
          url: candidate.url,
          width: candidate.width ?? undefined,
          height: candidate.height ?? undefined,
          alt: candidate.alt || undefined,
        },
      ];
    }
  }
  return undefined;
}
