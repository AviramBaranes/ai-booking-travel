/**
 * Payload relationship fields are typed as `number | T` where `number` is the
 * raw ID (depth: 0). When fetching with depth ≥ 1 the relation is always
 * populated, so use this utility to strip the `number` variant and work with
 * the actual document type directly.
 *
 * @example
 * type PopulatedPage = {
 *   featuredImage: Populated<Page["featuredImage"]>; // Media | null | undefined
 * }
 */
export type Populated<T> = Exclude<T, number>;

/**
 * Narrows a SharedSection to a specific `type` variant and makes the
 * corresponding data field required. Eliminates the need for runtime guards
 * like `if (section.type !== "newsletter" || !section.newsletter)`.
 *
 * @example
 * function NewsletterSection({ section }: { section: TypedSection<"newsletter"> }) {
 *   section.newsletter.title // ✅ no guard needed
 * }
 */
export type TypedSection<
  T extends import("@/payload-types").SharedSection["type"],
> = Omit<import("@/payload-types").SharedSection, T> & {
  type: T;
} & Required<Pick<import("@/payload-types").SharedSection, T>>;
