import zod from "zod";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { applyVoucher } from "@/shared/api/reservations";
import { useReservation } from "../_hooks/useReservation";
import { useMemo, useState } from "react";
import { AppError, isAppError } from "@/shared/api/AppError";
import { Loading } from "@/shared/components/Loading";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";

const applyVoucherSchema = zod.object({
  voucherCode: zod.string().min(1, "requiredField"),
});

type ApplyVoucherFormValues = zod.infer<typeof applyVoucherSchema>;

export function VoucherForm({ reservationId }: { reservationId: number }) {
  const t = useTranslations("MyAccount.reservation");
  const tErrors = useTranslations("ApiErrors");
  const { refetch } = useReservation(reservationId);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ApplyVoucherFormValues>({
    resolver: zodResolver(applyVoucherSchema),
  });

  const { mutate, isPending, error } = useMutation({
    mutationFn: (data: ApplyVoucherFormValues) =>
      applyVoucher(reservationId, data.voucherCode),
    onSuccess: () => refetch(),
  });

  const translatedError = useMemo(() => {
    if (!error) return null;

    if (isAppError(error)) {
      return tErrors(error.code);
    }

    return tErrors("internal_error");
  }, [error]);

  function onSubmit(data: ApplyVoucherFormValues) {
    mutate(data);
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input
        type="text"
        placeholder={t("enterVoucherCode")}
        className="mb-2 border-brand py-6"
        {...register("voucherCode")}
      />
      {errors.voucherCode?.message && (
        <ErrorDisplay>{tErrors(errors.voucherCode.message)}</ErrorDisplay>
      )}
      <Button
        variant="brand"
        type="submit"
        className="type-label w-full py-6 font-bold"
        disabled={isPending}
      >
        {isPending ? <Loading /> : t("apply")}
      </Button>
      {!!translatedError && <ErrorDisplay>{translatedError}</ErrorDisplay>}
    </form>
  );
}
