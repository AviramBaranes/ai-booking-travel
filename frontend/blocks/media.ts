import type { Block } from "payload";

/**
 * Media block — standalone image or video with optional caption.
 */
export const mediaBlock: Block = {
  slug: "media",
  interfaceName: "MediaBlock",
  labels: {
    singular: "מדיה",
    plural: "בלוקי מדיה",
  },
  fields: [
    {
      name: "media",
      label: "קובץ מדיה",
      type: "relationship",
      relationTo: "media",
      required: true,
    },
    {
      name: "caption",
      label: "כיתוב",
      type: "text",
      localized: true,
    },
    {
      name: "aspectRatio",
      label: "יחס גובה-רוחב",
      type: "select",
      defaultValue: "16:9",
      required: true,
      options: [
        { label: "16:9", value: "16:9" },
        { label: "4:3", value: "4:3" },
        { label: "1:1", value: "1:1" },
        { label: "אוטומטי", value: "auto" },
      ],
    },
    {
      name: "fullBleed",
      label: "רוחב מלא (edge-to-edge)",
      type: "checkbox",
      defaultValue: false,
    },
  ],
};
