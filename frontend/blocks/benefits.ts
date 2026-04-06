import type { Block } from "payload";

/**
 * Benefits block — marketing "why us" section with a grid of benefit cards.
 * Each card has an image/icon, title, and subtitle.
 * Examples: vehicle variety, brands, 24/7 support.
 */
export const benefitsBlock: Block = {
  slug: "benefits",
  interfaceName: "BenefitsBlock",
  labels: {
    singular: "יתרונות",
    plural: "בלוקי יתרונות",
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
      name: "items",
      label: "יתרונות",
      type: "array",
      minRows: 1,
      labels: {
        singular: "יתרון",
        plural: "יתרונות",
      },
      fields: [
        {
          name: "image",
          label: "תמונה / אייקון",
          type: "upload",
          relationTo: "media",
        },
        {
          name: "title",
          label: "כותרת",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "subtitle",
          label: "כותרת משנה",
          type: "text",
          localized: true,
        },
      ],
    },
  ],
};
