"use client";

import { HeaderSection } from "./HeaderSection";
import { CarDetailsSection } from "../../../../_components/CarDetailsSection";
import { IncludedSection } from "../../../../_components/IncludedSection";
import { CostBreakdownSection } from "../../../../_components/CostBreakdownSection";
import { RentalSummary } from "../../../../_components/RentalSummary";
import { useReservation } from "../../_hooks/useReservation";
import { PayAtPickupSection } from "../../../../_components/PayAtPickupSection";
import { useAddonsGallery } from "@/shared/hooks/useAddonsGallery";
import { SummarySkeleton } from "@/shared/components/booking/SummarySkeleton";

export function ReservationSummary({
  reservationId,
}: {
  reservationId: number;
}) {
  const { data: reservation, isLoading } = useReservation(reservationId);
  const { data: addonsGallery } = useAddonsGallery();

  if (isLoading || !reservation) {
    return <SummarySkeleton />;
  }

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
        flightNumber={reservation.flightNumber}
        supplierInfo={{
          termsAndConditions: reservation.supplierTerms,
          pickupDetails: reservation.pickupDetails,
          dropoffDetails: reservation.dropoffDetails,
        }}
      />
      <CarDetailsSection
        brand={reservation.carDetails.supplierName}
        model={reservation.carDetails.model}
        carType={reservation.carDetails.carType}
        rentalDays={reservation.rentalDays}
        isAutomatic={reservation.carDetails.isAutoGear}
        excess={reservation.excess}
        excessCurrency={reservation.excessCurrency}
      />
      <IncludedSection planInclusions={reservation.planInclusions} />
      <PayAtPickupSection
        currency={reservation.currencyCode}
        payAtPickup={reservation.payAtPickup}
        addonsGallery={addonsGallery}
      />
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
