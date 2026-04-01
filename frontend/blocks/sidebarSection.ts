import type { Block } from "payload";

/**
 * Sidebar Section block — a titled content section with an anchor ID.
 *
 * Used on pages like About and FAQ where the frontend generates a sticky
 * sidebar navigation from all sidebarSection blocks in the layout.
 * The `title` serves as both the section heading and the sidebar link label.
 */
export const sidebarSectionBlock: Block = {
  slug: "sidebarSection",
  interfaceName: "SidebarSectionBlock",
  labels: {
    singular: "אזור עם עוגן",
    plural: "אזורים עם עוגן",
  },
  fields: [
    {
      name: "anchor",
      label: "מזהה עוגן (URL)",
      type: "text",
      required: true,
      admin: {
        description:
          "מזהה ייחודי שישמש כ-#anchor בכתובת URL ובסרגל הניווט הצדדי. " +
          "אותיות לועזיות קטנות ומקפים בלבד. לדוגמה: about-us, our-vision.",
      },
    },
    {
      name: "title",
      label: "כותרת",
      type: "text",
      localized: true,
      required: true,
      admin: {
        description: "כותרת האזור — תוצג כקישור בסרגל הצד ובתור כותרת הסעיף.",
      },
    },
    {
      name: "content",
      label: "תוכן",
      type: "richText",
      localized: true,
      required: true,
    },
  ],
};
