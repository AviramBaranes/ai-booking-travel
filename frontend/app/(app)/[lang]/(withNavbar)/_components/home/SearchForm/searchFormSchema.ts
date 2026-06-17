import z from "zod";

export function searchSchema(t: (key: string) => string) {
  const baseSchema = z.object({
    isDropoffDifferentLoc: z.boolean().optional(),
    pickupLocation: z.int({ error: t("validation.required") }),
    dropoffLocation: z.int({ error: t("validation.required") }).optional(),
    pickupDate: z.date({ error: t("validation.required") }),
    pickupTime: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required")),
    dropoffDate: z.date({ error: t("validation.required") }),
    dropoffTime: z
      .string({ error: t("validation.required") })
      .min(1, t("validation.required")),
    driverAge: z.number().min(18),
    couponCode: z.string().optional(),
  });

  const crossFieldSchema = z
    .object({
      isDropoffDifferentLoc: z.boolean().optional(),
      dropoffLocation: z.any().optional(),
      pickupDate: z.any().optional(),
      dropoffDate: z.any().optional(),
      pickupTime: z.any().optional(),
      dropoffTime: z.any().optional(),
    })
    .superRefine((data, ctx) => {
      if (data.isDropoffDifferentLoc && !data.dropoffLocation) {
        ctx.addIssue({
          code: "custom",
          path: ["dropoffLocation"],
          message: t("validation.required"),
        });
      }

      const now = new Date();
      now.setHours(0, 0, 0, 0);

      if (data.pickupDate instanceof Date && data.pickupDate < now) {
        ctx.addIssue({
          code: "custom",
          path: ["pickupDate"],
          message: t("validation.dateInPast"),
        });
      }

      if (data.dropoffDate instanceof Date && data.dropoffDate < now) {
        ctx.addIssue({
          code: "custom",
          path: ["dropoffDate"],
          message: t("validation.dateInPast"),
        });
      }

      if (data.pickupDate instanceof Date && data.dropoffDate instanceof Date) {
        if (data.dropoffDate < data.pickupDate) {
          ctx.addIssue({
            code: "custom",
            path: ["dropoffDate"],
            message: t("validation.dropoffBeforePickup"),
          });
        }

        if (
          data.dropoffDate.toDateString() === data.pickupDate.toDateString() &&
          data.pickupTime &&
          data.dropoffTime &&
          data.dropoffTime <= data.pickupTime
        ) {
          ctx.addIssue({
            code: "custom",
            path: ["dropoffTime"],
            message: t("validation.dropoffTimeBeforePickup"),
          });
        }
      }

      const [hours, minutes] = data.pickupTime.split(":").map(Number);

      const pickupDateTime = new Date(data.pickupDate);
      pickupDateTime.setHours(hours, minutes, 0, 0);

      if (pickupDateTime < new Date()) {
        ctx.addIssue({
          code: "custom",
          path: ["pickupTime"],
          message: t("validation.dateInPast"),
        });
      }
    });

  return baseSchema.and(crossFieldSchema);
}

export type SearchFormValues = z.infer<ReturnType<typeof searchSchema>>;
