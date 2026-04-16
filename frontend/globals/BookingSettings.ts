import type { GlobalConfig } from "payload";

export const BookingSettings: GlobalConfig = {
  slug: "booking-settings",
  label: "הגדרות הזמנה",
  fields: [
    // ── ERP Section ──
    {
      type: "collapsible",
      label: "ERP",
      fields: [
        {
          name: "erpTitle",
          label: "כותרת",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "erpContent",
          label: "תוכן",
          type: "textarea",
          localized: true,
          required: true,
        },
        {
          name: "erpPopupTitle",
          label: "כותרת חלון קופץ",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "erpPopupContent",
          label: "תוכן חלון קופץ",
          type: "textarea",
          localized: true,
          required: true,
        },
      ],
    },

    // ── Fees Content Section ──
    {
      type: "collapsible",
      label: "תוכן עמלות",
      fields: [
        {
          name: "youngDriverTitle",
          label: "כותרת נהג צעיר",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "youngDriverContent",
          label: "תוכן נהג צעיר",
          type: "textarea",
          localized: true,
          required: true,
        },
        {
          name: "dropoffChargeTitle",
          label: "כותרת עמלת החזר",
          type: "text",
          localized: true,
          required: true,
        },
        {
          name: "dropoffChargeContent",
          label: "תוכן עמלת החזר",
          type: "textarea",
          localized: true,
          required: true,
        },
      ],
    },

    // ── Order Terms Link ──
    {
      name: "orderTermsLink",
      label: "קישור תנאי הזמנה",
      type: "relationship",
      relationTo: "pages",
      required: true,
    },
  ],
};
