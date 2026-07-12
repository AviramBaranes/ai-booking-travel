import { useState } from "react";
import Confetti from "react-confetti";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { getOrderPaymentIframe } from "@/shared/api/bill-api";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useTranslations } from "next-intl";
import { usePaymentSuccess } from "../../../../../../../../shared/hooks/usePaymentSuccess";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
  PopoverTitle,
  PopoverHeader,
} from "@/components/ui/popover";
import { ButtonGroup } from "@/components/ui/button-group";
import { useDirection } from "@/shared/hooks/useDirection";
import { useReservation } from "../_hooks/useReservation";

export function ReservationPaymentDialog({
  reservationId,
  setShow,
}: {
  reservationId: number;
  setShow: (show: boolean) => void;
}) {
  const t = useTranslations("MyAccount.reservation");
  const { data: reservation, refetch } = useReservation(reservationId);
  const [iframeUrl, setIframeUrl] = useState<string | null>(null);
  const [run, setRun] = useState(false);
  const [recycle, setRecycle] = useState(true);
  const { mutate, isPending, error } = useMutation({
    mutationFn: async (isIls: boolean) =>
      getOrderPaymentIframe(reservationId, isIls),
    onSuccess: (data) => {
      setIframeUrl(data.url);
    },
  });
  const tErr = useTranslatedError(error);
  const isPaymentDisabled = !!iframeUrl || run;

  usePaymentSuccess({
    onSuccess: () => {
      setIframeUrl(null);
      setRun(true);
      setTimeout(() => setRecycle(false), 3000);
      setTimeout(() => setShow(false), 5000);
    },
    refetch,
    isSuccess: reservation?.paymentStatus === "paid",
  });

  return (
    <>
      {/* <Popover> */}
      {/* <PopoverTrigger asChild> */}
      {reservation?.reservationStatus === "booked" &&
        reservation.paymentStatus === "unpaid" && (
          <Button
            disabled={isPaymentDisabled}
            loading={isPending}
            variant="payment"
            className="min-w-34"
            onClick={() => mutate(true)}
          >
            {t("payNow")}
          </Button>
        )}
      {/* </PopoverTrigger>
        <PopoverContent align="start" className="w-80 p-3">
          <PopoverHeader className="mb-1">
            <PopoverTitle className="type-paragraph font-semibold">
              באיזה מטבע תרצה לשלם?
            </PopoverTitle>
          </PopoverHeader>
          <ButtonGroup className="w-full" dir={dir === "rtl" ? "ltr" : "rtl"}>
            <Button
              variant="secondary"
              className="flex-1"
              onClick={() => mutate(false)}
            >
              מטבע העסקה
            </Button>
            <Button
              variant="secondary"
              className="flex-1"
              onClick={() => mutate(true)}
            >
              שקל
            </Button>
          </ButtonGroup>
        </PopoverContent>
      </Popover> */}
      <ErrorDisplay>{tErr}</ErrorDisplay>
      <Dialog
        open={!!iframeUrl}
        onOpenChange={() => {
          setIframeUrl(null);
        }}
      >
        <DialogContent className="w-full max-w-3xl!">
          <DialogTitle>
            <p className="type-paragraph mt-4">
              {t("paymentIframeTitle", { bookingId: reservationId })}
            </p>
          </DialogTitle>
          <iframe src={iframeUrl ?? undefined} className="w-full h-160" />
        </DialogContent>
      </Dialog>
      {run && (
        <div className="fixed inset-0 z-9999 pointer-events-none">
          <Confetti
            run={run}
            numberOfPieces={600}
            gravity={0.3}
            width={window.innerWidth}
            height={window.innerHeight}
            recycle={recycle}
          />
        </div>
      )}
    </>
  );
}
