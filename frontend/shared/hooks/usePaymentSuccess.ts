import { useEffect, useRef, useState } from "react";

export const PAYMENT_SUCCESS_EVENT = "PAYMENT_COMPLETED";

export function usePaymentSuccess({
  onSuccess,
  onLoadingStart,
  refetch,
  isSuccess,
  isFailure,
  onFailure,
}: {
  onSuccess?: () => void;
  onFailure?: () => void;
  onLoadingStart?: () => void;
  refetch: () => void;
  isSuccess: boolean;
  isFailure?: boolean;
}) {
  const [shouldPoll, setShouldPoll] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      if (event.data?.type === PAYMENT_SUCCESS_EVENT) {
        setShouldPoll(true);
        onLoadingStart?.();
        refetch();
      }
    };

    window.addEventListener("message", handleMessage);
    return () => {
      window.removeEventListener("message", handleMessage);
    };
  }, [refetch, setShouldPoll, onLoadingStart]);

  useEffect(() => {
    if (!shouldPoll) return;

    if (isSuccess) {
      setShouldPoll(false);
      onSuccess?.();
      return;
    }

    if (isFailure) {
      setShouldPoll(false);
      onFailure?.();
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
  }, [shouldPoll, isSuccess, isFailure, refetch, onSuccess, onFailure]);
}
