import { useEffect, useRef, useState } from "react";
import { useReservation } from "./useReservation";

export const PAYMENT_SUCCESS_EVENT = "PAYMENT_COMPLETED";

export function usePaymentSuccess({
  reservationId,
  onSuccess,
}: {
  reservationId: number;
  onSuccess?: () => void;
}) {
  const { data: reservation, refetch } = useReservation(reservationId);
  const [shouldPoll, setShouldPoll] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      if (event.data?.type === PAYMENT_SUCCESS_EVENT) {
        setShouldPoll(true);
        refetch();
      }
    };

    window.addEventListener("message", handleMessage);
    return () => {
      window.removeEventListener("message", handleMessage);
    };
  }, [refetch, setShouldPoll]);

  useEffect(() => {
    if (!shouldPoll) return;

    if (reservation?.paymentStatus === "paid") {
      setShouldPoll(false);
      onSuccess?.();
      return;
    }

    intervalRef.current = setInterval(() => {
      refetch();
    }, 2000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [shouldPoll, reservation?.paymentStatus, refetch]);
}
