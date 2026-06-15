import type { Field } from "payload";

export const suppliersFields: Field[] = [
  {
    name: "pillText",
    label: "תג עליון",
    type: "text",
    localized: true,
    admin: {
      description: 'לדוגמה: "השותפים שלנו"',
    },
  },
  {
    name: "title",
    label: "כותרת",
    type: "text",
    localized: true,
    required: true,
    admin: {
      description: 'לדוגמה: "חברות ההשכרה המובילות בעולם"',
    },
  },
  {
    name: "subtitle",
    label: "כותרת משנה",
    type: "text",
    localized: true,
  },
  {
    name: "logos",
    label: "לוגואים של חברות השכרה",
    type: "array",
    labels: {
      singular: "לוגו",
      plural: "לוגואים",
    },
    admin: {
      description: "לוגו לכל חברת השכרה. השתמשו ב-SVG או PNG שקוף.",
    },
    fields: [
      {
        name: "logo",
        label: "תמונת לוגו",
        type: "upload",
        relationTo: "media",
        required: true,
      },
    ],
  },
];
