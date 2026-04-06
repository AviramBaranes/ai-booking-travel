import type { GlobalConfig } from "payload";
import {
  richTextBlock,
  faqBlock,
  sharedSectionRefBlock,
  sidebarSectionBlock,
} from "../blocks";

export const Homepage: GlobalConfig = {
  slug: "homepage",
  label: "דף הבית",
  admin: {
    description:
      "תוכן דף הבית. הוסיפו בלוקים בסדר הרצוי ושלבו אזורים משותפים לפי הצורך.",
    group: "תוכן",
  },
  versions: {
    drafts: true,
  },
  fields: [
    {
      type: "tabs",
      tabs: [
        {
          label: "פריסה",
          fields: [
            {
              name: "layout",
              label: "בלוקים",
              type: "blocks",
              minRows: 1,
              admin: {
                description:
                  "בנו את דף הבית על ידי הוספת בלוקים. שלבו אזורים משותפים (ניוזלטר, חברות השכרה, סטטיסטיקות) עם תוכן ייחודי לדף הבית.",
              },
              blocks: [
                richTextBlock,
                faqBlock,
                sharedSectionRefBlock,
                sidebarSectionBlock,
              ],
            },
          ],
        },
      ],
    },
  ],
};
