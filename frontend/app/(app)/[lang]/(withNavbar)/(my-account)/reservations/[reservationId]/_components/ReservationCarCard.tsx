"use client";

import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { useTranslations } from "next-intl";
import { VoucherForm } from "./VoucherForm";
import { useReservation } from "../_hooks/useReservation";
import { ReservationPaymentDialog } from "./ReservationPaymentDialog";
import { Suspense, useEffect, useState } from "react";
import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";

export function ReservationCarCard({
  reservationId,
}: {
  reservationId: number;
}) {
  const {
    data: reservation,
    refetch,
    isLoading,
  } = useReservation(reservationId);
  const t = useTranslations("MyAccount.reservation");
  const [isInitialized, setIsInitialized] = useState(false);
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);

  useEffect(() => {
    if (isInitialized) return;

    if (!isLoading && reservation) {
      setIsInitialized(true);
      setShowPaymentDialog(
        reservation.reservationStatus === "booked" &&
          reservation.paymentStatus === "unpaid",
      );
    }
  }, [isLoading, reservation, isInitialized, setIsInitialized]);

  if (isLoading || !reservation) {
    return <SelectedCarCardSkeleton />;
  }

  return (
    <div className="sticky top-24">
      <SelectedCarCardWrapper>
        <SelectedCarHeader carDetails={reservation.carDetails} />
        <FreeCancellationBadge
          pickupDate={reservation.pickupDate}
          pickupTime={reservation.pickupTime}
          text={t("freeCancellation")}
        />
        {reservation.reservationStatus === "booked" && (
          <VoucherForm reservationId={reservationId} refetch={refetch} />
        )}
        {showPaymentDialog && (
          <ReservationPaymentDialog
            reservationId={reservationId}
            setShow={setShowPaymentDialog}
          />
        )}
      </SelectedCarCardWrapper>
    </div>
  );
}
