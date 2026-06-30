"use client";

import { formatPrice } from "@/shared/utils/formatPrice";
import { CarDetailsSection } from "../../../../_components/CarDetailsSection";
import { CostBreakdownSection } from "../../../../_components/CostBreakdownSection";
import { IncludedSection } from "../../../../_components/IncludedSection";
import { RentalSummary } from "../../../../_components/RentalSummary";
import { usePriceOffer } from "../../_hooks/usePriceOffer";
import { HeaderSection } from "./HeaderSection";
import { useTranslations } from "next-intl";
import { PayAtPickupSection } from "../../../../_components/PayAtPickupSection";
import { useAddonsGallery } from "@/shared/hooks/useAddonsGallery";

export function PriceOfferSummary({ priceOfferId }: { priceOfferId: number }) {
  const t = useTranslations("MyAccount.priceOffer.summary.labels");
  const { data: priceOffer } = usePriceOffer(priceOfferId);
  const { data: addonsGallery } = useAddonsGallery();

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <HeaderSection priceOffer={priceOffer} />
      <RentalSummary
        dropoffDate={priceOffer.dropoffDate}
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
        rentalDays={priceOffer.rentalDays}
        isAutomatic={priceOffer.carDetails.isAutoGear}
      />
      <IncludedSection planInclusions={priceOffer.planInclusions} />

      <PayAtPickupSection
        currency={priceOffer.currencyCode}
        payAtPickup={priceOffer.payAtPickup}
        addonsGallery={addonsGallery}
      />

      <CostBreakdownSection
        showDisclaimer
        priceBefDesc={priceOffer.priceBefDesc}
        discountAmount={0}
        erpPrice={priceOffer.erpPrice}
        totalPrice={priceOffer.totalPrice}
        currencyCode={priceOffer.currencyCode}
      />

      <div className="text-white bg-brand py-3 5 px-5 flex justify-between items-center rounded-xl">
        <span className="type-paragraph">{t("offeredPrice")}</span>
        <h4 className="type-h4">
          {formatPrice(priceOffer.offeredPrice, priceOffer.offeredCurrencyCode)}
        </h4>
      </div>
    </div>
  );
}
