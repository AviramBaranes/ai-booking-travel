"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod/dist/zod.js";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { createPenalty } from "@/shared/api/reservations-api";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";

const PENALTY_TYPES = [
  { value: "no_show", label: "אי הגעה" },
  { value: "cancellation", label: "ביטול מאוחר" },
];

const createPenaltySchema = z.object({
  // The input registers as a number, so a blank field arrives as NaN rather than a string.
  amount: z
    .number({ error: "יש להזין עלות" })
    .gt(0, "יש להזין עלות גדולה מאפס"),
  type: z.enum(["no_show", "cancellation"]),
});
type CreatePenaltyFormData = z.infer<typeof createPenaltySchema>;

interface CreatePenaltyFormProps {
  reservationId: number;
  currencyCode: string;
}

/**
 * CreatePenaltyForm records the fee the supplier charged on a canceled reservation. The fee is
 * charged in the reservation's own currency, so only the amount and the type are asked for.
 */
export function CreatePenaltyForm({
  reservationId,
  currencyCode,
}: CreatePenaltyFormProps) {
  const queryClient = useQueryClient();

  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<CreatePenaltyFormData>({
    resolver: zodResolver(createPenaltySchema),
    defaultValues: { type: "no_show" },
  });

  const { mutate, isPending, error, isSuccess, data } = useMutation({
    mutationFn: (values: CreatePenaltyFormData) =>
      createPenalty({
        reservationId,
        type: values.type,
        amount: values.amount,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["reservation-detail", reservationId],
      });
    },
  });

  const tError = useTranslatedError(error);

  if (isSuccess && data) {
    return (
      <p className="text-sm text-slate-600">
        נרשם חיוב {PENALTY_TYPES.find((t) => t.value === data.type)?.label} על
        סך {data.amount} {data.currencyCode}.
      </p>
    );
  }

  return (
    <form
      onSubmit={handleSubmit((values) => mutate(values))}
      className="flex flex-wrap items-end gap-3"
    >
      <div className="flex flex-col gap-1.5">
        <Label className="text-sm font-bold text-slate-700">
          עלות ({currencyCode})
        </Label>
        <Input
          type="number"
          step="0.01"
          min="0"
          placeholder="0.00"
          className="w-40"
          isError={!!errors.amount}
          {...register("amount", { valueAsNumber: true })}
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <Label className="text-sm font-bold text-slate-700">סוג</Label>
        <select
          className="h-8 w-48 rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm outline-none focus-visible:border-ring"
          {...register("type")}
        >
          {PENALTY_TYPES.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>

      <Button
        type="submit"
        variant="brand"
        className="h-8 px-6"
        loading={isPending}
      >
        צור חיוב
      </Button>

      <div className="w-full">
        <ErrorDisplay>{errors.amount?.message ?? tError}</ErrorDisplay>
      </div>
    </form>
  );
}
