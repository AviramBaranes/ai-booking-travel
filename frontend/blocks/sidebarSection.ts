import type { Block } from "payload";

export const sidebarSectionBlock: Block = {
  slug: "sidebarSection",
  interfaceName: "SidebarSectionBlock",
  labels: {
    singular: "אזור עם עוגן",
    plural: "אזורים עם עוגן",
  },
  fields: [
    {
      name: "sections",
      label: "אזורים",
      type: "array",
      minRows: 1,
      labels: {
        singular: "אזור",
        plural: "אזורים",
      },
      fields: [
        {
          name: "title",
          label: "כותרת",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "content",
          label: "תוכן",
          type: "richText",
          localized: true,
          required: true,
        },
      ],
    },
  ],
};
