"use client";

import { HeaderSection } from "./HeaderSection";
import { CarDetailsSection } from "../../../../_components/CarDetailsSection";
import { IncludedSection } from "../../../../_components/IncludedSection";
import { CostBreakdownSection } from "../../../../_components/CostBreakdownSection";
import { RentalSummary } from "../../../../_components/RentalSummary";
import { useReservation } from "../../_hooks/useReservation";
import { PayAtPickupSection } from "../../../../_components/PayAtPickupSection";
import { useAddonsGallery } from "@/shared/hooks/useAddonsGallery";
import { useParams } from "next/navigation";

export function ReservationSummary({
  reservationId,
}: {
  reservationId: number;
}) {
  const { lang } = useParams();
  const { data: reservation } = useReservation(reservationId);
  const { data: addonsGallery } = useAddonsGallery();

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
      />
      <CarDetailsSection
        brand={reservation.carDetails.supplierName}
        model={reservation.carDetails.model}
        carType={reservation.carDetails.carType}
        rentalDays={reservation.rentalDays}
        isAutomatic={reservation.carDetails.isAutoGear}
      />
      <IncludedSection planInclusions={reservation.planInclusions} />
      <PayAtPickupSection
        currency={reservation.currencyCode}
        fees={reservation.payAtPickup.fees}
        selectedAddons={reservation.payAtPickup.selectedAddons?.map((addon) => {
          const addonGalleryItem = addonsGallery.addons?.find(
            (item) => item.addonId === addon.id.toString(),
          );
          return {
            ...addon,
            name: addonGalleryItem
              ? lang === "he"
                ? addonGalleryItem.hebrewName
                : addonGalleryItem.englishName
              : addon.name,
          };
        })}
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
