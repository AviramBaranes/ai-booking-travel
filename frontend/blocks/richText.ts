import type { Block } from "payload";

/**
 * Rich Text block — free-form text content.
 * Wraps a lexical rich text editor with width and alignment controls.
 */
export const richTextBlock: Block = {
  slug: "richText",
  interfaceName: "RichTextBlock",
  labels: {
    singular: "טקסט עשיר",
    plural: "בלוקי טקסט",
  },
  fields: [
    {
      name: "content",
      label: "תוכן",
      type: "richText",
      localized: true,
      required: true,
    },
    {
      name: "maxWidth",
      label: "רוחב מקסימלי",
      type: "select",
      defaultValue: "lg",
      required: true,
      options: [
        { label: "צר", value: "sm" },
        { label: "בינוני", value: "md" },
        { label: "רחב", value: "lg" },
        { label: "מסך מלא", value: "full" },
      ],
    },
    {
      name: "textAlign",
      label: "יישור טקסט",
      type: "select",
      defaultValue: "start",
      options: [
        { label: "התחלה (ברירת מחדל)", value: "start" },
        { label: "מרכז", value: "center" },
        { label: "סוף", value: "end" },
      ],
    },
  ],
};
