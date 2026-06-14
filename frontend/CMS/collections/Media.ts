import type { CollectionConfig } from "payload";

export const Media: CollectionConfig = {
  slug: "media",
  labels: {
    plural: "מדיה",
    singular: "מדיה",
  },
  access: {
    read: () => true,
  },
  fields: [
    {
      name: "alt",
      type: "text",
      required: true,
    },
  ],
  upload: {
    imageSizes: [
      {
        name: "blogHero",
        width: 900,
        height: 380,
        position: "centre",
      },
      {
        name: "blogCard",
        width: 520,
        height: 320,
        position: "centre",
      },
      {
        name: "og",
        width: 1200,
        height: 630,
        position: "centre",
      },
    ],
  },
};
