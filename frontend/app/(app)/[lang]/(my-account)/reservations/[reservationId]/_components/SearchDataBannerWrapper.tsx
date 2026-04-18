"use client";

import { useReservation } from "../_hooks/useReservation";
import { SearchDataBannerDisplay } from "@/shared/components/booking/SearchDataBannerDisplay";
import { SearchDataBannerPrint } from "@/shared/components/booking/SearchDataBannerPrint";

export function SearchDataBannerWrapper({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation } = useReservation(reservationId);

  const bannerProps = {
    driverAge: reservation.driverAge,
    pickUpDate: new Date(reservation.pickupDate),
    dropOffDate: new Date(reservation.returnDate),
    pickUpLocationName: reservation.pickupLocationName,
    dropOffLocationName: reservation.dropoffLocationName,
    pickUpTime: reservation.pickupTime,
    dropOffTime: reservation.dropoffTime,
  };

  return (
    <>
      <div className="print:hidden">
        <SearchDataBannerDisplay {...bannerProps} />
      </div>
      <SearchDataBannerPrint {...bannerProps} />
    </>
  );
}
