"use client";

import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { useTranslations } from "next-intl";
import { VoucherForm } from "./VoucherForm";
import { useReservation } from "../_hooks/useReservation";

export function ReservationCarCard({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation, refetch } = useReservation(reservationId);
  const t = useTranslations("MyAccount.reservation");

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
      </SelectedCarCardWrapper>
    </div>
  );
}
