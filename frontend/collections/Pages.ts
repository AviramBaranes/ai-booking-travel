import type { CollectionConfig } from "payload";
import {
  heroBlock,
  richTextBlock,
  mediaBlock,
  ctaBlock,
  faqBlock,
  sharedSectionRefBlock,
  sidebarSectionBlock,
} from "../blocks";

export const Pages: CollectionConfig = {
  slug: "pages",
  labels: {
    singular: "עמוד",
    plural: "עמודים",
  },
  admin: {
    useAsTitle: "title",
    defaultColumns: ["title", "slug", "template", "_status", "updatedAt"],
  },
  defaultSort: "-updatedAt",
  versions: {
    drafts: true,
  },
  fields: [
    {
      type: "tabs",
      tabs: [
        {
          label: "תוכן",
          fields: [
            {
              name: "title",
              label: "כותרת",
              type: "text",
              localized: true,
              required: true,
            },
            {
              name: "slug",
              label: "Slug (כתובת URL)",
              type: "text",
              localized: true,
              required: true,
              index: true,
              unique: true,
              admin: {
                description:
                  'הגדירו slug שונה לכל שפה. עמוד הבית ישמר כ-"home" ב-Payload ויומף ל-"/" ב-Next.js.',
              },
            },
            {
              name: "excerpt",
              label: "תקציר",
              type: "textarea",
              localized: true,
              maxLength: 220,
              admin: {
                description: "תיאור קצר לכרטיסיות ותוצאות חיפוש, עד 220 תווים.",
              },
            },
            {
              name: "featuredImage",
              label: "תמונה ראשית",
              type: "upload",
              relationTo: "media",
            },

            {
              name: "layout",
              label: "פריסת בלוקים",
              type: "blocks",
              minRows: 0,
              blocks: [
                heroBlock,
                richTextBlock,
                mediaBlock,
                ctaBlock,
                faqBlock,
                sharedSectionRefBlock,
                sidebarSectionBlock,
              ],
            },
          ],
        },

        {
          label: "הגדרות",
          fields: [
            {
              name: "template",
              label: "תבנית עמוד",
              type: "select",
              defaultValue: "default",
              required: true,
              options: [
                { label: "ברירת מחדל", value: "default" },
                { label: "דף נחיתה", value: "landing" },
                { label: "אודות", value: "about" },
                { label: "שאלות נפוצות", value: "faq" },
                { label: "משפטי / תנאי שימוש", value: "legal" },
                { label: "עזרה", value: "help" },
                { label: "תודה / אישור", value: "thank-you" },
                { label: "עמוד לא נמצא (404)", value: "not-found" },
              ],
            },
            {
              name: "relatedPages",
              label: "עמודים קשורים",
              type: "relationship",
              relationTo: "pages",
              hasMany: true,
            },
            {
              name: "publishedAt",
              label: "תאריך פרסום",
              type: "date",
              admin: {
                readOnly: true,
                date: {
                  pickerAppearance: "dayAndTime",
                },
              },
              hooks: {
                beforeChange: [
                  ({ siblingData, value }) => {
                    if (!value && siblingData?._status === "published") {
                      return new Date();
                    }
                    return value;
                  },
                ],
              },
            },
          ],
        },
      ],
    },
  ],
};
