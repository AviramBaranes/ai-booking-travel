import { revalidateOnGlobalChange } from "@/CMS/hooks/revalidate";
import type { GlobalConfig } from "payload";
import { benefitsBlock, faqBlock, createSharedSectionRefBlock } from "../blocks";

export const Homepage: GlobalConfig = {
  slug: "homepage",
  label: "דף הבית",
  admin: {
    description:
      "תוכן דף הבית. הוסיפו בלוקים בסדר הרצוי ושלבו אזורים משותפים לפי הצורך.",
    group: "תוכן",
  },
  hooks: {
    afterChange: [revalidateOnGlobalChange],
  },
  fields: [
    {
      type: "tabs",
      tabs: [
        {
          label: "גיבור",
          fields: [
            {
              name: "featuredImage",
              label: "תמונה ראשית",
              type: "upload",
              relationTo: "media",
            },
            {
              name: "title",
              label: "כותרת ראשית",
              type: "text",
              localized: true,
              required: true,
            },
            {
              name: "subtitle",
              label: "כותרת משנה",
              type: "textarea",
              localized: true,
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
          ],
        },
        {
          label: "פריסה",
          fields: [
            {
              name: "layout",
              label: "בלוקים",
              type: "blocks",
              minRows: 0,
              admin: {
                description:
                  "בנו את דף הבית על ידי הוספת בלוקים. שלבו אזורים משותפים עם תוכן ייחודי לדף הבית.",
              },
              blocks: [createSharedSectionRefBlock("secRef"), benefitsBlock, faqBlock],
            },
          ],
        },
      ],
    },
  ],
};
