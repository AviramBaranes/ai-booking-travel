import type { CollectionConfig } from "payload";
import { revalidateOnCollectionChange } from "../../hooks/revalidate";
import { newsletterFields } from "./newsletterFields";
import { suppliersFields } from "./suppliersFields";
import { statsFields } from "./statsFields";
import { contactBlockFields } from "./contactFields";

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
  hooks: {
    afterChange: [revalidateOnCollectionChange],
  },
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
        { label: "צור קשר", value: "contact" },
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
    {
      name:"contact",
      label:"צור קשר",
      type:"group",
      admin: {
        condition: (data) => data?.type === "contact",
      },
      fields: contactBlockFields
    }
  ],
};
