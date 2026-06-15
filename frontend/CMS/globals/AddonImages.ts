import type { GlobalConfig } from "payload";

export const AddonImagesGlobal: GlobalConfig = {
  slug: "addonsGallery",
  label: "גלריית תוספים",
  fields: [
    {
      name: "addons",
      labels: {
        singular: "תוסף",
        plural: "תוספים",
      },
      type: "array",
      fields: [
        {
          name: "addonId",
          label: "מזהה תוסף בפלקס",
          type: "text",
          required: true,
        },
        {
          name: "englishName",
          label: "שם התוסף באנגלית",
          type: "text",
          required: true,
        },
        {
          name: "hebrewName",
          label: "שם התוסף בעברית",
          type: "text",
          required: true,
        },
        {
          name: "image",
          label: "תמונה",
          type: "upload",
          relationTo: "media",
          required: true,
        },
      ],
    },
  ],
};
