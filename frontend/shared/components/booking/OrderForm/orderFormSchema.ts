import z from "zod";

const ENGLISH_ONLY_REGEX = /^[A-Z\s]+$/;
const PHONE_REGEX = /^05\d{8}$/;

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
