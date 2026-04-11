import type { GlobalConfig } from "payload";

export const SuppliersGallery: GlobalConfig = {
  slug: "suppliersGallery",
  label: "גלריית ספקים",
  fields: [
    {
      name: "suppliers",
      label: "ספקים",
      type: "array",
      minRows: 1,
      fields: [
        {
          name: "name",
          label: "שם",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "media",
          label: "מדיה",
          type: "upload",
          relationTo: "media",
          required: true,
        },
      ],
    },
  ],
};
