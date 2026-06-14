import type { Field } from "payload";

export const newsletterFields: Field[] = [
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
  {
    name: "benefits",
    label: "יתרונות / נקודות מפתח",
    type: "array",
    labels: {
      singular: "יתרון",
      plural: "יתרונות",
    },
    admin: {
      description: "רשימת bullets שתוצג לצד הטופס.",
    },
    fields: [
      {
        name: "text",
        label: "טקסט",
        type: "text",
        localized: true,
        required: true,
      },
    ],
  },
  {
    name: "formTitle",
    label: "כותרת הטופס",
    type: "text",
    localized: true,
  },
  {
    name: "formSubTitle",
    label: "כותרת משנה של הטופס",
    type: "text",
    localized: true,
  },
  {
    name: "emailPlaceholder",
    label: "טקסט placeholder לאימייל",
    type: "text",
    localized: true,
  },
  {
    name: "submitButtonLabel",
    label: "טקסט כפתור שליחה",
    type: "text",
    localized: true,
    required: true,
  },
  {
    name: "consentLabel",
    label: "טקסט הסכמה (לצד הצ׳קבוקס)",
    type: "text",
    localized: true,
  },
  // Privacy link parts — broken into three so editors can link the policy page
  {
    name: "privacyTextBeforeLink",
    label: "טקסט לפני קישור הפרטיות",
    type: "text",
    localized: true,
  },
  {
    name: "privacyLinkLabel",
    label: "טקסט קישור מדיניות הפרטיות",
    type: "text",
    localized: true,
  },
  {
    name: "privacyPage",
    label: "עמוד מדיניות פרטיות",
    type: "relationship",
    relationTo: "pages",
  },
];
