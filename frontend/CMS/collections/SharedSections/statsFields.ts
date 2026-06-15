import type { Field } from "payload";

export const statsFields: Field[] = [
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
  {
    name: "items",
    label: "נתונים",
    type: "array",
    minRows: 1,
    labels: {
      singular: "נתון",
      plural: "נתונים",
    },
    fields: [
      {
        name: "value",
        label: "ערך",
        type: "text",
        localized: true,
        required: true,
        admin: {
          description: 'לדוגמה: "+20,000", "4.9", "אלפי"',
        },
      },
      {
        name: "label",
        label: "תווית",
        type: "text",
        localized: true,
        required: true,
        admin: {
          description: 'לדוגמה: "הזמנות", "דירוג לקוחות", "יעדים"',
        },
      },
      {
        name: "icon",
        label: "אייקון",
        type: "upload",
        relationTo: "media",
        admin: {
          description: "אופציונלי: תמונה/אייקון SVG לצד הנתון.",
        },
      },
    ],
  },
];
