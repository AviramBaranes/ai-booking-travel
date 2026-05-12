"use client";

import { CarDetailsSection } from "../../../../_components/CarDetailsSection";
import { CostBreakdownSection } from "../../../../_components/CostBreakdownSection";
import { IncludedSection } from "../../../../_components/IncludedSection";
import { RentalSummary } from "../../../../_components/RentalSummary";
import { usePriceOffer } from "../../_hooks/usePriceOffer";
import { HeaderSection } from "./HeaderSection";

export function PriceOfferSummary({ priceOfferId }: { priceOfferId: number }) {
  const { data: priceOffer } = usePriceOffer(priceOfferId);

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <HeaderSection priceOffer={priceOffer} />
      <RentalSummary
        dropoffDate={priceOffer.returnDate}
        dropoffTime={priceOffer.dropoffTime}
        dropoffLocationName={priceOffer.dropoffLocationName}
        pickupDate={priceOffer.pickupDate}
        pickupTime={priceOffer.pickupTime}
        pickupLocationName={priceOffer.pickupLocationName}
      />
      <CarDetailsSection
        brand={priceOffer.carDetails.supplierName}
        model={priceOffer.carDetails.model}
        carType={priceOffer.carDetails.carType}
        rentalDays={0}
      />
      <IncludedSection planInclusions={priceOffer.planInclusions} />
      <CostBreakdownSection
        priceBefDesc={priceOffer.priceBefDesc}
        discountAmount={0}
        erpPrice={priceOffer.erpPrice}
        totalPrice={priceOffer.totalPrice}
        currencyCode={priceOffer.currencyCode}
      />
    </div>
  );
}
