import type { Block } from "payload";

/**
 * Benefits block — marketing "why us" section with a grid of benefit cards.
 * Each card has an image/icon, title, and subtitle.
 * Examples: vehicle variety, brands, 24/7 support.
 */
export const relatedPostsBlock: Block = {
  slug: "related-posts",
  interfaceName: "RelatedPostsBlock",
  labels: {
    singular: "פוסטים קשורים",
    plural: "בלוקי פוסטים קשורים",
  },
  fields: [
    {
      name: "eyebrow",
      label: "תג / כותרת עליונה קטנה",
      type: "text",
      localized: true,
      admin: {
        description:
          'טקסט קצר המוצג מעל הכותרת הראשית. לדוגמה: "למה לבחור בנו".',
      },
    },
    {
      name: "title",
      label: "כותרת הבלוק",
      type: "text",
      localized: true,
    },
    {
      name: "subtitle",
      label: "כותרת משנה",
      type: "textarea",
      localized: true,
    },
    {
      name: "relatedPosts",
      label: "בחירת פוסטים",
      type: "relationship",
      relationTo: "blog-posts",
      hasMany: true,
      maxRows: 4,
    },
  ],
};
