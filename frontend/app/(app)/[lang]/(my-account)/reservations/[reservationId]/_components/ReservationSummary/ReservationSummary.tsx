"use client";

import { useReservation } from "../../_hooks/useReservation";
import { HeaderSection } from "./HeaderSection";
import { CarDetailsSection } from "./CarDetailsSection";
import { IncludedSection } from "./IncludedSection";
import { CostBreakdownSection } from "./CostBreakdownSection";
import { RentalSummary } from "./RentalSummary";

export function ReservationSummary({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation } = useReservation(reservationId);

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <HeaderSection reservation={reservation} />
      <RentalSummary reservation={reservation} />
      <CarDetailsSection reservation={reservation} />
      <IncludedSection planInclusions={reservation.planInclusions} />
      <CostBreakdownSection reservation={reservation} />
    </div>
  );
}
