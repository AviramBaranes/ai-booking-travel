import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { useState } from "react";
import { X } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { cancelReservation } from "@/shared/api/reservations";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useReservation } from "../../_hooks/useReservation";

export function DeleteReservationButton({
  reservationId,
}: {
  reservationId: number;
}) {
  const t = useTranslations("MyAccount.reservation.summary.cancel");
  const [open, setOpen] = useState(false);
  const { refetch } = useReservation(reservationId);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const {
    mutate: deleteReservation,
    isPending,
    error,
  } = useMutation({
    mutationFn: () => cancelReservation(reservationId),
    onSuccess: () => {
      refetch();
      setSuccessMessage(t("successMessage"));
      setTimeout(() => {
        setSuccessMessage(null);
        setOpen(false);
      }, 2000);
    },
  });

  const translatedError = useTranslatedError(error);

  return (
    <>
      <Button
        variant="ghost"
        className="py-6 px-6 text-border-muted font-semibold flex gap-4 print:hidden"
        onClick={() => setOpen(true)}
      >
        <Image
          src="/assets/icons/trash.svg"
          alt={t("button")}
          width={24}
          height={24}
          className="w-6 h-6"
        />
        {t("button")}
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          className="min-w-1/4 max-w-md p-6 flex flex-col gap-4 bg-white border-border-light/50 rounded-2xl shadow-modal"
          showCloseButton={false}
        >
          <DialogTitle className="type-h5 text-navy flex items-center justify-between p-6 pb-15 border-b border-muted/50">
            <div className="flex items-center gap-4">
              <Image
                src="/assets/icons/trash.svg"
                alt={t("button")}
                width={24}
                height={24}
                className="w-6 h-6"
              />
              {t("title")}
            </div>
            <X
              className="w-6 h-6 cursor-pointer"
              onClick={() => setOpen(false)}
            />
          </DialogTitle>
          <h5 className="type-h5 text-navy">{t("subTitle")}</h5>
          <p className="type-paragraph text-text-secondary">{t("message")}</p>

          <Button
            onClick={() => deleteReservation()}
            variant="destructive"
            loading={isPending}
          >
            {t("button")}
          </Button>
          <div className="p-0 text-center">
            {!!translatedError && (
              <ErrorDisplay>{translatedError}</ErrorDisplay>
            )}
            {successMessage && <SuccessBadge text={successMessage} />}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
