import z from "zod";

const ENGLISH_ONLY_REGEX = /^[A-Z\s]+$/;

export function orderFormSchema(t: (key: string) => string) {
  return z.object({
    driverTitle: z.enum(["Mr", "Ms"], { error: t("validation.required") }),
    driverFirstName: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required"))
      .regex(ENGLISH_ONLY_REGEX, t("validation.englishOnly")),
    driverLastName: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required"))
      .regex(ENGLISH_ONLY_REGEX, t("validation.englishOnly")),
    flightNumber: z.string().optional(),
    termsAccepted: z.literal(true, {
      error: t("validation.mustAcceptTerms"),
    }),
  });
}

export type OrderFormValues = z.infer<ReturnType<typeof orderFormSchema>>;
