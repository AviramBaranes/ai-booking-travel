import type { GlobalConfig } from "payload";
import {
  heroBlock,
  richTextBlock,
  mediaBlock,
  ctaBlock,
  faqBlock,
  sharedSectionRefBlock,
  sidebarSectionBlock,
} from "../blocks";

/**
 * Homepage global.
 *
 * A single-document global (one per locale) that drives the homepage layout.
 * Uses the same block palette as the Pages collection so blocks and shared
 * sections are consistent across the site.
 *
 * Editors compose the homepage by mixing inline blocks (hero, richText, cta, …)
 * and references to shared sections (newsletter, suppliers, stats) in any order.
 */
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
        // ── Tab 1: Layout ────────────────────────────────────────────────────
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

        // ── Tab 2: SEO added automatically by plugin-seo ─────────────────────
      ],
    },
  ],
};
