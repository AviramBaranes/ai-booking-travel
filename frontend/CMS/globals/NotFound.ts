import type { GlobalConfig } from "payload";
import { revalidateOnGlobalChange } from "../hooks/revalidate";

export const NotFoundConfig: GlobalConfig = {
  slug: "not-found",
  label: "404 - לא נמצא",
  hooks: {
    afterChange: [revalidateOnGlobalChange],
  },
  fields: [
    {
      name: "title",
      label: "כותרת",
      type: "text",
      localized: true,
    },
    {
      name: "subtitle",
      label: "תת-כותרת",
      type: "text",
      localized: true,
    },
    {
      name: "buttonText",
      label: "טקסט כפתור",
      type: "text",
      localized: true,
    },
  ],
};
