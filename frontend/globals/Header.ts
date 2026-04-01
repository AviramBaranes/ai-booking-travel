import type { GlobalConfig } from "payload";

export const Header: GlobalConfig = {
  slug: "header",
  label: "הדר (תפריט עליון)",
  fields: [
    {
      name: "links",
      label: "קישורים",
      type: "array",
      minRows: 1,
      fields: [
        {
          name: "type",
          label: "סוג קישור",
          type: "select",
          defaultValue: "link",
          required: true,
          options: [
            { label: "קישור רגיל", value: "link" },
            { label: "מגה-תפריט", value: "mega" },
          ],
        },
        // ── Regular link fields ──
        {
          name: "label",
          label: "טקסט",
          type: "text",
          localized: true,
          required: true,
          admin: {
            condition: (_, siblingData) => siblingData?.type === "link",
          },
        },
        {
          name: "page",
          label: "עמוד",
          type: "relationship",
          relationTo: "pages",
          required: true,
          admin: {
            condition: (_, siblingData) => siblingData?.type === "link",
          },
        },
        // ── Mega link fields ──
        {
          name: "megaLabel",
          label: "כותרת מגה-תפריט",
          type: "text",
          localized: true,
          required: true,
          admin: {
            condition: (_, siblingData) => siblingData?.type === "mega",
          },
        },
        {
          name: "megaLinks",
          label: "קישורים",
          type: "array",
          minRows: 1,
          admin: {
            condition: (_, siblingData) => siblingData?.type === "mega",
          },
          fields: [
            {
              name: "label",
              label: "טקסט",
              type: "text",
              localized: true,
              required: true,
            },
            {
              name: "page",
              label: "עמוד",
              type: "relationship",
              relationTo: "pages",
              required: true,
            },
          ],
        },
      ],
    },
  ],
};
