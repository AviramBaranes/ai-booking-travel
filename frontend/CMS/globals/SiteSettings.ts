import type { GlobalConfig } from "payload";

export const SiteSettings: GlobalConfig = {
  slug: "site-settings",
  label: "הגדרות האתר",
  fields: [
    {
      type: "collapsible",
      label: "הגדרות פוסט",
      fields: [
        {
          name: "rpPillText",
          label: "תצוגת פוסטים דומים - תג עליון",
          type: "text",
          localized: true,
        },
        {
          name: "rpTitle",
          label: "תצוגת פוסטים דומים - כותרת",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "rpSubtitle",
          label: "תצוגת פוסטים דומים - כותרת משנה",
          type: "text",
          localized: true,
        },
      ],
    },
    {
      type: "collapsible",
      label: "הגדרות עמוד בלוג",
      fields: [
        {
          name: "featuredImage",
          label: "תמונה ראשית",
          type: "upload",
          relationTo: "media",
        },
        {
          name: "pillText",
          label: "תג עליון",
          type: "text",
          localized: true,
        },
        {
          name: "title",
          label: "כותרת",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "subtitle",
          label: "כותרת משנה",
          type: "text",
          localized: true,
        },
      ],
    },
  ],
};
