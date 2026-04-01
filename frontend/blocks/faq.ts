import type { Block } from "payload";

/**
 * FAQ block — accordion question-and-answer block organized by categories.
 * From Figma: used on FAQ page with sections "לפני ההשכרה", "בזמן ההשכרה", "אחרי ההשכרה".
 * Works as an inline block (page-specific) and can also be referenced via sharedSectionRef.
 */
export const faqBlock: Block = {
  slug: "faq",
  interfaceName: "FAQBlock",
  labels: {
    singular: "שאלות נפוצות",
    plural: "בלוקי שאלות נפוצות",
  },
  fields: [
    {
      name: "eyebrow",
      label: "תג / כותרת עליונה קטנה",
      type: "text",
      localized: true,
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
      type: "text",
      localized: true,
    },
    {
      name: "categories",
      label: "קטגוריות",
      type: "array",
      minRows: 1,
      labels: {
        singular: "קטגוריה",
        plural: "קטגוריות",
      },
      admin: {
        description:
          'כל קטגוריה מציגה כותרת וקבוצת שאלות-תשובות. לדוגמה: "לפני ההשכרה".',
      },
      fields: [
        {
          name: "heading",
          label: "כותרת קטגוריה",
          type: "text",
          localized: true,
          admin: {
            description:
              "לדוגמה: לפני ההשכרה / בזמן ההשכרה / אחרי ההשכרה. ניתן להשאיר ריק אם אין קטגוריות.",
          },
        },
        {
          name: "items",
          label: "שאלות ותשובות",
          type: "array",
          minRows: 1,
          labels: {
            singular: "שאלה",
            plural: "שאלות",
          },
          fields: [
            {
              name: "question",
              label: "שאלה",
              type: "text",
              localized: true,
              required: true,
            },
            {
              name: "answer",
              label: "תשובה",
              type: "richText",
              localized: true,
              required: true,
            },
          ],
        },
      ],
    },
  ],
};
