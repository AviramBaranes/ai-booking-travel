import type { Block } from "payload";

/**
 * CTA block — simple button for Thank You, 404, confirmation pages etc.
 */
export const ctaBlock: Block = {
  slug: "cta",
  interfaceName: "CTABlock",
  labels: {
    singular: "כפתור קריאה לפעולה",
    plural: "כפתורי קריאה לפעולה",
  },
  fields: [
    {
      name: "label",
      label: "טקסט הכפתור",
      type: "text",
      localized: true,
      required: true,
    },
    {
      name: "page",
      label: "יעד הכפתור",
      type: "relationship",
      relationTo: "pages",
      required: true,
    },
    {
      name: "backgroundColor",
      label: "צבע רקע",
      type: "select",
      defaultValue: "brand",
      options: [
        { label: "מותג (כחול כהה)", value: "brand" },
        { label: "כתום / הדגשה", value: "accent" },
        { label: "לבן", value: "white" },
      ],
    },
    {
      name: "textColor",
      label: "צבע טקסט",
      type: "select",
      defaultValue: "white",
      options: [
        { label: "לבן", value: "white" },
        { label: "כהה", value: "dark" },
      ],
    },
  ],
};
