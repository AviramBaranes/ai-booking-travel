import type { CollectionBeforeChangeHook, CollectionConfig } from "payload";
import { slugify } from "../Pages";
import { revalidateOnCollectionChange } from "@/CMS/hooks/revalidate";

export const BlogCategories: CollectionConfig = {
  slug: "blog-categories",
  labels: {
    singular: "קטגוריה",
    plural: "קטגוריות",
  },
  admin: {
    useAsTitle: "title",
    defaultColumns: ["title", "slug", "updatedAt"],
  },
  defaultSort: "title",
  hooks: {
    beforeChange: [generateSlugFromCategoryTitle()],
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
              label: "שם הקטגוריה",
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
                  "החלק שיופיע בכתובת של עמוד הקטגוריה. לדוגמה: car-rental-guides.",
              },
            },
            {
              name: "description",
              label: "תיאור",
              type: "textarea",
              localized: true,
              admin: {
                description:
                  "תיאור קצר של הקטגוריה. יכול להופיע בעמוד הקטגוריה ולעזור ל־SEO.",
              },
            },
            {
              name: "image",
              label: "תמונה",
              type: "upload",
              relationTo: "media",
              admin: {
                description:
                  "תמונה אופציונלית לעמוד הקטגוריה או לכרטיסיות קטגוריה באתר.",
              },
            },
          ],
        },
        {
          label: "SEO",
          fields: [
            {
              name: "seo",
              label: "SEO",
              type: "group",
              fields: [
                {
                  name: "title",
                  label: "כותרת SEO",
                  type: "text",
                  localized: true,
                  admin: {
                    description:
                      "כותרת ייעודית למנועי חיפוש. אם ריק, ניתן להשתמש בשם הקטגוריה.",
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
                      "תיאור קצר למנועי חיפוש ושיתופים. מומלץ עד 220 תווים.",
                  },
                },
                {
                  name: "image",
                  label: "תמונת SEO",
                  type: "upload",
                  relationTo: "media",
                  admin: {
                    description:
                      "תמונה ייעודית לשיתופים. אם ריק, ניתן להשתמש בתמונת הקטגוריה.",
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

function generateSlugFromCategoryTitle(): CollectionBeforeChangeHook {
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