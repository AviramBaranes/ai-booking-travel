import type { CollectionConfig, Field } from "payload";

const newsletterFields: Field[] = [
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

const suppliersFields: Field[] = [
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
        type: "relationship",
        relationTo: "media",
        required: true,
      },
    ],
  },
];

/**
 * Stats / Trust Numbers section fields.
 *
 * From Figma (About page): Row of 4 key metrics with icons.
 * Observed examples: "+20,000 לרכישות", "4.9 דירוג לקוח", "4.9 דירוג פרש", "אלפי הזמנות".
 * Values may contain symbols and prefixes like "+" so they're stored as localized text.
 */
const statsFields: Field[] = [
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
        name: "caption",
        label: "תיאור נוסף",
        type: "text",
        localized: true,
        admin: {
          description: "אופציונלי: שורת הסבר קצרה מתחת לתווית.",
        },
      },
      {
        name: "icon",
        label: "אייקון",
        type: "relationship",
        relationTo: "media",
        admin: {
          description: "אופציונלי: תמונה/אייקון SVG לצד הנתון.",
        },
      },
    ],
  },
];

// ─── Collection ────────────────────────────────────────────────────────────────

export const SharedSections: CollectionConfig = {
  slug: "sharedSections",
  labels: {
    singular: "אזור משותף",
    plural: "אזורים משותפים",
  },
  admin: {
    useAsTitle: "internalTitle",
    defaultColumns: ["internalTitle", "type", "_status", "updatedAt"],
    description:
      "אזורים שניתן לשלב בכל עמוד או בדף הבית. עריכה כאן תשתקף בכל מקום שהאזור מוטמע.",
  },
  defaultSort: "internalTitle",
  versions: {
    drafts: true,
  },
  fields: [
    // ── Identification ─────────────────────────────────────────────────────
    {
      name: "internalTitle",
      label: "שם פנימי",
      type: "text",
      required: true,
      admin: {
        description: "שם לזיהוי בממשק הניהול בלבד. לא מוצג בפרונטאנד.",
      },
    },
    {
      name: "type",
      label: "סוג אזור",
      type: "select",
      required: true,
      admin: {
        description: "בחרו את סוג האזור. הבחירה קובעת אילו שדות יוצגו.",
      },
      options: [
        { label: "ניוזלטר", value: "newsletter" },
        { label: "חברות השכרה", value: "suppliers" },
        { label: "סטטיסטיקות", value: "stats" },
      ],
    },

    // ── Newsletter ─────────────────────────────────────────────────────────
    {
      name: "newsletter",
      label: "ניוזלטר",
      type: "group",
      admin: {
        condition: (data) => data?.type === "newsletter",
      },
      fields: newsletterFields,
    },

    // ── Suppliers ─────────────────────────────────────────────────────────
    {
      name: "suppliers",
      label: "חברות השכרה",
      type: "group",
      admin: {
        condition: (data) => data?.type === "suppliers",
      },
      fields: suppliersFields,
    },

    // ── Stats ──────────────────────────────────────────────────────────────
    {
      name: "stats",
      label: "סטטיסטיקות",
      type: "group",
      admin: {
        condition: (data) => data?.type === "stats",
      },
      fields: statsFields,
    },
  ],
};
