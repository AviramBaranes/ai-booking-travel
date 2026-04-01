import type { GlobalConfig } from "payload";

export const Footer: GlobalConfig = {
  slug: "footer",
  label: "פוטר (תחתית)",
  fields: [
    // ── Tab 1: First Floor ──
    {
      type: "tabs",
      tabs: [
        {
          label: "קומה ראשונה",
          fields: [
            {
              name: "firstFloorLinks",
              label: "קישורים",
              type: "array",
              minRows: 1,
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
        // ── Tab 2: Second Floor ──
        {
          label: "קומה שנייה",
          fields: [
            {
              name: "socialsTitle",
              label: "כותרת רשתות חברתיות",
              type: "text",
              localized: true,
            },
            {
              name: "socialsLinks",
              label: "קישורי רשתות חברתיות",
              type: "array",
              fields: [
                {
                  name: "label",
                  label: "טקסט",
                  type: "text",
                  localized: true,
                  required: true,
                },
                {
                  name: "link",
                  label: "קישור",
                  type: "text",
                  required: true,
                },
              ],
            },
            {
              name: "linkGroups",
              label: "קבוצות קישורים",
              type: "array",
              fields: [
                {
                  name: "title",
                  label: "כותרת",
                  type: "text",
                  localized: true,
                  required: true,
                },
                {
                  name: "links",
                  label: "קישורים",
                  type: "array",
                  minRows: 1,
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
        },
        // ── Tab 3: Third Floor ──
        {
          label: "קומה שלישית",
          fields: [
            {
              name: "thirdFloorLinks",
              label: "קישורים",
              type: "array",
              minRows: 1,
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
            {
              name: "rights",
              label: "טקסט זכויות יוצרים",
              type: "text",
              localized: true,
              defaultValue: "© AI Booking Travel כל הזכויות שמורות",
            },
          ],
        },
      ],
    },
  ],
};
