import type { CollectionBeforeChangeHook, CollectionConfig } from "payload";
import {
  richTextBlock,
  faqBlock,
  createSharedSectionRefBlock,
  sidebarSectionBlock,
} from "../../blocks";
import { revalidateOnCollectionChange } from "../../hooks/revalidate";
import { slugify } from "../Pages";

export const BlogPosts: CollectionConfig = {
  slug: "blog-posts",
  labels: {
    singular: "פוסט",
    plural: "פוסטים",
  },
  admin: {
    useAsTitle: "title",
    defaultColumns: ["title", "slug", "category", "_status", "publishedAt"],
  },
  defaultSort: "-publishedAt",
  versions: {
    drafts: true,
  },
  hooks: {
    beforeChange: [generateSlugFromBlogPostTitle()],
    afterChange: [revalidateOnCollectionChange],
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
                  "החלק שיופיע בכתובת ה־URL של הפוסט. לדוגמה: car-rental-insurance.",
              },
            },
            {
              name: "excerpt",
              label: "תקציר",
              type: "textarea",
              localized: true,
              maxLength: 220,
              admin: {
                description:
                  "תיאור קצר של הפוסט. משמש לכרטיסיות בלוג, תוצאות חיפוש, שיתופים ו־SEO. מומלץ עד 220 תווים.",
              },
            },
            {
              name: "featuredImage",
              label: "תמונה ראשית",
              type: "upload",
              relationTo: "media",
              required: true,
              admin: {
                description:
                  "התמונה הראשית של הפוסט. תופיע בדרך כלל בראש הפוסט, בכרטיסיות בלוג ובשיתופים.",
              },
            },
            {
              name: "layout",
              label: "פריסת בלוקים",
              type: "blocks",
              minRows: 1,
              blocks: [
                richTextBlock,
                faqBlock,
                createSharedSectionRefBlock("secRef_B"),// name needs to be unique and short.
                sidebarSectionBlock,
              ],
              admin: {
                description:
                  "תוכן הפוסט עצמו. ניתן להרכיב את הפוסט מבלוקים כמו טקסט עשיר, שאלות נפוצות, אזורים משותפים ועוד.",
              },
            },
          ],
        },
        {
          label: "שיוך וסיווג",
          fields: [
            {
              name: "category",
              label: "קטגוריה",
              type: "relationship",
              relationTo: "blog-categories",
              required: true,
              admin: {
                description:
                  "הקטגוריה הראשית שאליה הפוסט שייך. משמשת לניווט, סינון, עמודי קטגוריה ופוסטים קשורים.",
              },
            },
            {
              name: "tags",
              label: "תגיות",
              type: "array",
              labels: {
                singular: "תגית",
                plural: "תגיות",
              },
              admin: {
                description:
                  "תגיות הן נושאים נקודתיים שהפוסט עוסק בהם, למשל: ביטוח, שדה תעופה, פיקדון, נהג צעיר. בשלב הזה הן טקסט חופשי ולא אוסף נפרד.",
              },
              fields: [
                {
                  name: "tag",
                  label: "שם התגית",
                  type: "text",
                  required: true,
                },
              ],
            },
            {
              name: "relatedPosts",
              label: "פוסטים קשורים",
              type: "relationship",
              relationTo: "blog-posts",
              hasMany: true,
              admin: {
                description:
                  "בחירה ידנית של פוסטים קשורים. אם השדה ריק, אפשר להציג אוטומטית פוסטים מאותה קטגוריה.",
              },
            },
          ],
        },
        {
          label: "הגדרות",
          fields: [
            {
              name: "form",
              label: "טופס",
              type: "relationship",
              relationTo: "forms",
              required: false,
              hasMany: false,
              admin: {
                description:
                  "טופס אופציונלי שיוצג בעמוד הפוסט, למשל טופס יצירת קשר, ליד או בקשת הצעת מחיר.",
              },
            },
            {
              name: "banner",
              label: "באנר",
              type: "group",
              admin: {
                description:
                  "באנר אופציונלי שיוצג בעמוד הפוסט. כולל תמונה וקישור שאליו המשתמש יועבר בלחיצה.",
              },
              fields: [
                {
                  name: "image",
                  label: "תמונת באנר",
                  type: "upload",
                  relationTo: "media",
                  required: false,
                  admin: {
                    description:
                      "תמונה שתוצג כבאנר בעמוד הפוסט. מומלץ להשתמש בתמונה רחבה שמתאימה לתצוגה בדסקטופ ובמובייל.",
                  },
                },
                {
                  name: "link",
                  label: "קישור",
                  type: "text",
                  required: false,
                  admin: {
                    description:
                      "כתובת שאליה הבאנר יפנה בלחיצה. אפשר להזין קישור פנימי כמו /he/contact או קישור מלא.",
                  },
                },
              ],
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
                description:
                  "נקבע אוטומטית בפעם הראשונה שהפוסט מפורסם. משמש למיון הפוסטים באתר.",
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
            {
              name: "seo",
              label: "SEO",
              type: "group",
              admin: {
                description:
                  "הגדרות ייעודיות למנועי חיפוש ושיתופים. אם השדות ריקים, אפשר להשתמש בכותרת, בתקציר ובתמונה הראשית של הפוסט.",
              },
              fields: [
                {
                  name: "title",
                  label: "כותרת SEO",
                  type: "text",
                  localized: true,
                  admin: {
                    description:
                      "כותרת ייעודית למנועי חיפוש. אם ריק, ניתן להשתמש בכותרת הפוסט.",
                  },
                },
                {
                  name: "description",
                  label: "תיאור SEO",
                  type: "textarea",
                  localized: true,
                  maxLength: 220,
                  admin: {
                    description:
                      "תיאור קצר שיופיע במנועי חיפוש ובשיתופים. מומלץ עד 220 תווים.",
                  },
                },
                {
                  name: "image",
                  label: "תמונת SEO",
                  type: "upload",
                  relationTo: "media",
                  admin: {
                    description:
                      "תמונה ייעודית לשיתופים. אם ריק, ניתן להשתמש בתמונה הראשית של הפוסט.",
                  },
                },
              ],
            },
          ],
        },
      ],
    },
  ],
};

function generateSlugFromBlogPostTitle(): CollectionBeforeChangeHook {
  return ({ data }) => {
    if (!("title" in data) || !data.title) {
      return data;
    }

    return {
      ...data,
      slug: slugify(data.title),
    };
  };
}
