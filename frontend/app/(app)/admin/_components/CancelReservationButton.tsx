"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { cancelReservationByAdmin } from "@/shared/api/reservations-api";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";

interface CancelReservationButtonProps {
  reservationId: number;
}

/**
 * CancelReservationButton cancels a reservation from the admin portal. The admin endpoint
 * skips the cancellation-window and credit-card checks that stop agents and customers, and
 * it neither adjusts balance due nor refunds the payment — so the confirmation spells out
 * that the action cannot be undone.
 */
export function CancelReservationButton({
  reservationId,
}: CancelReservationButtonProps) {
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const queryClient = useQueryClient();

  const { mutate, isPending, error } = useMutation({
    mutationFn: () => cancelReservationByAdmin(reservationId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["reservation-detail", reservationId],
      });
      setIsConfirmOpen(false);
    },
  });

  const tError = useTranslatedError(error);

  return (
    <>
      <Button
        variant="destructive"
        className="h-8 px-6 py-0"
        onClick={() => setIsConfirmOpen(true)}
      >
        בטל הזמנה
      </Button>

      {/* A cancellation in flight must finish before the confirmation can be dismissed. */}
      <Dialog
        open={isConfirmOpen}
        onOpenChange={(open) => !isPending && setIsConfirmOpen(open)}
      >
        <DialogContent dir="rtl" showCloseButton={false}>
          <DialogHeader>
            <DialogTitle className="text-start">
              האם אתה בטוח שברצונך לבטל הזמנה זאת?
            </DialogTitle>
            <DialogDescription className="text-start">
              פעולה זאת בלתי הפיכה.
            </DialogDescription>
          </DialogHeader>

          <ErrorDisplay>{tError}</ErrorDisplay>

          <DialogFooter className="sm:justify-start">
            <Button
              variant="destructive"
              className="h-8 px-6 py-0"
              loading={isPending}
              onClick={() => mutate()}
            >
              כן, בטל
            </Button>
            <Button
              variant="outline"
              className="h-8 px-6"
              disabled={isPending}
              onClick={() => setIsConfirmOpen(false)}
            >
              התחרטתי
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
