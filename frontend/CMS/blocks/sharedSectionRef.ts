import type { Block } from "payload";

/**
 * Shared Section Reference block.
 *
 * Allows editors to place any sharedSection (newsletter, suppliers, stats)
 * at any position in a page or homepage layout.
 * Optional per-placement overrides let editors tweak spacing/theme without
 * editing the shared section itself.
 */
export const sharedSectionRefBlock: Block = {
  slug: "sharedSectionRef",
  dbName: "secRef",
  interfaceName: "SharedSectionRefBlock",
  labels: {
    singular: "אזור משותף",
    plural: "אזורים משותפים",
  },
  fields: [
    {
      name: "section",
      label: "אזור משותף",
      type: "relationship",
      relationTo: "sharedSections",
      required: true,
      admin: {
        description:
          "בחרו אזור קיים מהספרייה המשותפת (ניוזלטר, חברות השכרה, סטטיסטיקות וכד׳).",
      },
    },
    {
      name: "overrides",
      label: "עקיפות מקומיות",
      type: "group",
      admin: {
        description: "שינויים חזותיים לעמוד הזה בלבד, מבלי לפגוע באזור המשותף.",
      },
      fields: [
        {
          name: "spacingTop",
          label: "ריווח עליון",
          type: "select",
          defaultValue: "default",
          options: [
            { label: "ברירת מחדל", value: "default" },
            { label: "ללא", value: "none" },
            { label: "קטן", value: "sm" },
            { label: "בינוני", value: "md" },
            { label: "גדול", value: "lg" },
          ],
        },
        {
          name: "spacingBottom",
          label: "ריווח תחתון",
          type: "select",
          defaultValue: "default",
          options: [
            { label: "ברירת מחדל", value: "default" },
            { label: "ללא", value: "none" },
            { label: "קטן", value: "sm" },
            { label: "בינוני", value: "md" },
            { label: "גדול", value: "lg" },
          ],
        },
      ],
    },
  ],
};
