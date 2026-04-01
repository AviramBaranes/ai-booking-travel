import type { Block } from "payload";

export const heroBlock: Block = {
  slug: "hero",
  interfaceName: "HeroBlock",
  labels: {
    singular: "בלוק הירו",
    plural: "בלוקי הירו",
  },
  fields: [
    {
      name: "variant",
      label: "סגנון",
      type: "select",
      defaultValue: "centered",
      required: true,
      options: [
        { label: "מרוכז (כותרת + תיאור במרכז)", value: "centered" },
        { label: "תמונה מימין", value: "image-right" },
        { label: "תמונה משמאל", value: "image-left" },
        { label: "רקע מלא", value: "full-bg" },
      ],
    },
    {
      name: "title",
      label: "כותרת",
      type: "text",
      localized: true,
      required: true,
    },
    {
      name: "description",
      label: "תיאור",
      type: "richText",
      localized: true,
    },
    {
      name: "media",
      label: "תמונה / וידאו",
      type: "relationship",
      relationTo: "media",
    },
  ],
};
