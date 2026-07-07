import type { GlobalConfig } from "payload";
import { revalidateOnGlobalChange } from "../hooks/revalidate";

export const SuppliersGallery: GlobalConfig = {
  slug: "suppliersGallery",
  label: "גלריית ספקים",
  access: {
    read: () => true,
  },
  hooks: {
    afterChange: [revalidateOnGlobalChange],
  },
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
