import { useState } from "react";
import zod from "zod";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { applyVoucher } from "@/shared/api/reservations";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useReservation } from "../_hooks/useReservation";

const applyVoucherSchema = zod.object({
  voucherCode: zod.string().min(1, "requiredField"),
});

type ApplyVoucherFormValues = zod.infer<typeof applyVoucherSchema>;

export function VoucherForm({
  reservationId,
  refetch,
}: {
  reservationId: number;
  refetch?: () => void;
}) {
  const t = useTranslations("MyAccount.reservation.voucher");
  const tErrors = useTranslations("ApiErrors");
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

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
    onSuccess: () => {
      setSuccessMessage(t("successMessage"));
      setTimeout(() => {
        refetch?.();
        setSuccessMessage(null);
      }, 2000);
    },
  });

  const translatedError = useTranslatedError(error);

  function onSubmit(data: ApplyVoucherFormValues) {
    mutate(data);
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="print:hidden">
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
        loading={isPending}
      >
        {t("apply")}
      </Button>
      {!!translatedError && <ErrorDisplay>{translatedError}</ErrorDisplay>}
      {successMessage && <SuccessBadge text={successMessage} />}
    </form>
  );
}
