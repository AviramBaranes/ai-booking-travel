import z from "zod";

const ENGLISH_ONLY_REGEX = /^[A-Z\s]+$/;
const PHONE_REGEX = /^05\d{8}$/;
// The driver may not be Israeli, so an international number in E.164 (+ and country code) is accepted too, as the backend's `phone` validator does.
const DRIVER_PHONE_REGEX = /^(05\d{8}|\+[1-9]\d{6,14})$/;

export function orderFormSchema(t: (key: string) => string, requiresCustomerDetails = false) {
  return z.object({
    driverTitle: z.enum(["Mr", "Mrs", "Ms", "Miss", "Dr"], { error: t("validation.required") }),
    driverFirstName: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required"))
      .regex(ENGLISH_ONLY_REGEX, t("validation.englishOnly")),
    driverLastName: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required"))
      .regex(ENGLISH_ONLY_REGEX, t("validation.englishOnly")),
    // Required when a customer books (the same case that requires the customer's details), optional for agents.
    driverPhone: requiresCustomerDetails
      ? z
          .string({ error: t("validation.required") })
          .min(1, t("validation.required"))
          .regex(DRIVER_PHONE_REGEX, t("validation.invalidDriverPhone"))
      : z
          .string()
          .optional()
          .refine((v) => !v || DRIVER_PHONE_REGEX.test(v), t("validation.invalidDriverPhone")),
    flightNumber: z.string().optional(),
    customerFirstName: requiresCustomerDetails
      ? z.string().min(1, t("validation.required"))
      : z.string().optional(),
    customerLastName: requiresCustomerDetails
      ? z.string().min(1, t("validation.required"))
      : z.string().optional(),
    customerPhone: requiresCustomerDetails
      ? z.string().min(1, t("validation.required")).regex(PHONE_REGEX, t("validation.invalidPhone"))
      : z.string().optional(),
    customerEmail: requiresCustomerDetails
      ? z.string().min(1, t("validation.required")).pipe(z.email(t("validation.invalidEmail")))
      : z.string().optional(),
    termsAccepted: z.literal(true, {
      error: t("validation.mustAcceptTerms"),
    }),
  });
}

export type OrderFormValues = z.infer<ReturnType<typeof orderFormSchema>>
