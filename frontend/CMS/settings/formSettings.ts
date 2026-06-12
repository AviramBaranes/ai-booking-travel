import type { Block, Field } from "payload";

const fieldTranslations: Record<string, string> = {
  fields: "שדות",
  title: "כותרת",
  name: "שם שדה",
  label: "תווית",
  defaultValue: "ערך ברירת מחדל",
  width: "רוחב שדה באחוזים",
  required: "שדה חובה",

  submitButtonLabel: "טקסט כפתור השליחה",
  confirmationType: "סוג אישור",
  confirmationMessage: "הודעת אישור",
  redirect: "הפניה",
  emails: "אימיילים",
  email: "אימייל",

  options: "אפשרויות",
  value: "ערך",
};

const blockLabels: Record<string, { singular: string; plural: string }> = {
  text: {
    singular: "טקסט",
    plural: "שדות טקסט",
  },
  email: {
    singular: "אימייל",
    plural: "שדות אימייל",
  },
  select: {
    singular: "בחירה",
    plural: "שדות בחירה",
  },
  textarea: {
    singular: "טקסט ארוך",
    plural: "שדות טקסט ארוך",
  },
  checkbox: {
    singular: "תיבת סימון",
    plural: "תיבות סימון",
  },
};

const optionTranslations: Record<string, string> = {
  message: "הודעה",
  redirect: "הפניה",
  email: "אימייל",
  emails: "אימיילים",
};

const translateOptions = (options: any[] | undefined) => {
  if (!Array.isArray(options)) return options;

  return options.map((option) => {
    if (typeof option === "string") {
      return optionTranslations[option] ?? option;
    }

    return {
      ...option,
      label:
        typeof option.value === "string" && optionTranslations[option.value]
          ? optionTranslations[option.value]
          : option.label,
    };
  });
};

const translateFieldDeep = (field: any): any => {
  const translated = {
    ...field,
    label:
      field.name && fieldTranslations[field.name]
        ? fieldTranslations[field.name]
        : field.label,
    options: translateOptions(field.options),
    admin: {
      ...field.admin,
      description: undefined,
    },
  };

  if (Array.isArray(translated.fields)) {
    translated.fields = translated.fields.map(translateFieldDeep);
  }

  if (Array.isArray(translated.blocks)) {
    translated.blocks = translated.blocks.map((block: any) => ({
      ...block,
      labels: blockLabels[block.slug] ?? block.labels,
      fields: Array.isArray(block.fields)
        ? block.fields.map(translateFieldDeep)
        : block.fields,
    }));
  }

  return translated;
};

const phoneField: Block = {
  slug: "phone",
  labels: {
    singular: "טלפון",
    plural: "שדות טלפון",
  },
  fields: [
    {
      name: "name",
      type: "text",
      label: "שם שדה",
      required: true,
    },
    {
      name: "label",
      type: "text",
      label: "תווית",
      localized: true,
    },
    {
      name: "defaultValue",
      type: "text",
      label: "ערך ברירת מחדל",
      localized: true,
    },
    {
      name: "width",
      type: "number",
      label: "רוחב שדה באחוזים",
    },
    {
      name: "required",
      type: "checkbox",
      label: "שדה חובה",
    },
  ],
};

export const formSettings = {
  fields: {
    text: true,
    email: true,
    select: true,
    textarea: true,
    checkbox: true,

    country: false,
    message: false,
    number: false,
    payment: false,
    radio: false,
    state: false,
    phone: phoneField,
  },

  formOverrides: {
    labels: {
      singular: "טופס",
      plural: "טפסים",
    },
    admin: {
      group: "טפסים",
      useAsTitle: "title",
    },
    fields: ({ defaultFields }: { defaultFields: Field[] }) => {
      return defaultFields.map((field: any) => {
        const translated = translateFieldDeep(field);

        if (field.name === "fields") {
          return {
            ...translated,
            label: "שדות",
            labels: {
              singular: "שדה",
              plural: "שדות",
            },
            admin: {
              ...translated.admin,
              initCollapsed: true,
            },
          };
        }

        if (field.name === "emails") {
          return {
            ...translated,
            label: "אימיילים",
            labels: {
              singular: "אימייל",
              plural: "אימיילים",
            },
          };
        }

        return translated;
      });
    },
  },

  formSubmissionOverrides: {
    labels: {
      singular: "שליחת טופס",
      plural: "שליחות טפסים",
    },
    admin: {
      group: "טפסים",
    },
  },
};
