import type { CollectionConfig } from "payload";

export const AddonImages: CollectionConfig = {
  slug: "addon-images",
  admin: {
    useAsTitle: "hebrewName",
  },
  labels: {
    plural: "תוספים",
    singular: "תוסף",
  },
  fields: [
    {
      name: "addonId",
      label: "מזהה תוסף בפלקס",
      type: "text",
      required: true,
      unique: true,
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
};
