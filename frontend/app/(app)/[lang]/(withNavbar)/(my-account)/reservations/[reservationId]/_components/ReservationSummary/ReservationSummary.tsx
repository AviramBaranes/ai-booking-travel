"use client";

import { HeaderSection } from "./HeaderSection";
import { CarDetailsSection } from "../../../../_components/CarDetailsSection";
import { IncludedSection } from "../../../../_components/IncludedSection";
import { CostBreakdownSection } from "../../../../_components/CostBreakdownSection";
import { RentalSummary } from "../../../../_components/RentalSummary";
import { useReservation } from "../../_hooks/useReservation";

export function ReservationSummary({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation } = useReservation(reservationId);

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <HeaderSection reservation={reservation} />
      <RentalSummary
        dropoffDate={reservation.dropoffDate}
        dropoffTime={reservation.dropoffTime}
        dropoffLocationName={reservation.dropoffLocationName}
        pickupDate={reservation.pickupDate}
        pickupTime={reservation.pickupTime}
        pickupLocationName={reservation.pickupLocationName}
      />
      <CarDetailsSection
        brand={reservation.carDetails.supplierName}
        model={reservation.carDetails.model}
        carType={reservation.carDetails.carType}
        rentalDays={reservation.rentalDays}
        isAutomatic={reservation.carDetails.isAutoGear}
      />
      <IncludedSection planInclusions={reservation.planInclusions} />
      <CostBreakdownSection
        priceBefDesc={reservation.priceBefDesc}
        discountAmount={reservation.discountAmount}
        erpPrice={reservation.erpPrice}
        totalPrice={reservation.totalPrice}
        currencyCode={reservation.currencyCode}
      />
    </div>
  );
}
