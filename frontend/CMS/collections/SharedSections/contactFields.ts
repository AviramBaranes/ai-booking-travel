import { Field } from "payload";

export const contactBlockFields: Field[]  = [
  {
    name: "eyebrow",
    label: "תג / כותרת עליונה קטנה",
    type: "text",
    localized: true,
    admin: {
      description: 'טקסט קצר המוצג מעל הכותרת הראשית. לדוגמה: "למה לבחור בנו".',
    },
  },
  {
    name: "title",
    label: "כותרת הבלוק",
    type: "text",
    localized: true,
  },
  {
    name: "subtitle",
    label: "כותרת משנה",
    type: "textarea",
    localized: true,
  },
  {
    name: "contactForm",
    label: "טופס יצירת קשר",
    type: "relationship",
    relationTo: "forms",
    required: true,
  },
  {
    name: "contactInfo",
    label: "פרטי יצירת קשר",
    type: "array",
    fields: [
      {
        name: "icon",
        label: "אייקון",
        type: "relationship",
        relationTo: "media",
        required: true,
      },
      {
        name: "title",
        label: "כותרת",
        type: "text",
        localized: true,
        required: true,
      },
      {
        name: "content",
        label: "ערך",
        type: "richText",
        localized: true,
        required: true,
      },
    ],
  },
];
