"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod/dist/zod.js";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { applyVoucherByAdmin } from "@/shared/api/reservations-api";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";

const adminVoucherSchema = z.object({
  voucher: z.string().trim().min(1, "יש להזין מספר שובר"),
});
type AdminVoucherFormData = z.infer<typeof adminVoucherSchema>;

interface AdminVoucherFormProps {
  reservationId: number;
}

/**
 * AdminVoucherForm vouchers a booked agent reservation on the agent's behalf. The backend
 * checks the agent's credit, adds the price to their balance due and emails them the voucher,
 * just as when the agent vouchers it.
 */
export function AdminVoucherForm({ reservationId }: AdminVoucherFormProps) {
  const queryClient = useQueryClient();

  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<AdminVoucherFormData>({
    resolver: zodResolver(adminVoucherSchema),
  });

  const { mutate, isPending, error } = useMutation({
    mutationFn: (values: AdminVoucherFormData) =>
      applyVoucherByAdmin(reservationId, values.voucher),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["reservation-detail", reservationId],
      });
    },
  });

  const tError = useTranslatedError(error);

  return (
    <form
      onSubmit={handleSubmit((values) => mutate(values))}
      className="flex flex-wrap items-center gap-3"
    >
      <Input
        type="text"
        placeholder="מספר שובר"
        className="h-8 w-52"
        isError={!!errors.voucher}
        {...register("voucher")}
      />
      <Button
        type="submit"
        variant="brand"
        className="h-8 px-6"
        loading={isPending}
      >
        כרטס הזמנה
      </Button>

      <div className="w-full">
        <ErrorDisplay>{errors.voucher?.message ?? tError}</ErrorDisplay>
      </div>
    </form>
  );
}
