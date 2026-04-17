"use client";

import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { useReservation } from "../_hooks/useReservation";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { useTranslations } from "next-intl";
import { VoucherForm } from "./VoucherForm";

export function ReservationCarCard({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation } = useReservation(reservationId);
  const t = useTranslations("MyAccount.reservation");

  return (
    <SelectedCarCardWrapper>
      <SelectedCarHeader carDetails={reservation.carDetails} />
      <FreeCancellationBadge
        pickupDate={reservation.pickupDate}
        pickupTime={reservation.pickupTime}
        text={t("freeCancellation")}
      />
      {!reservation.voucher && <VoucherForm reservationId={reservationId} />}
    </SelectedCarCardWrapper>
  );
}
