import { price_offer } from "@/shared/client";
import { HeaderSection } from "./HeaderSection";
import { RentalSummary } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/RentalSummary";
import { CarDetailsSection } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/CarDetailsSection";
import { IncludedSection } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/IncludedSection";
import { getTranslations } from "next-intl/server";
import { formatPrice } from "@/shared/utils/formatPrice";
import { PayAtPickupSection } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/PayAtPickupSection";
import Image from "next/image";

export async function ClientOfferSummary({
  offer,
  lang,
}: {
  offer: price_offer.GetPriceOfferResponse;
  lang: string;
}) {
  const t = await getTranslations("MyAccount");

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <HeaderSection offer={offer} lang={lang} />
      <RentalSummary
        dropoffDate={offer.dropoffDate}
        dropoffTime={offer.dropoffTime}
        dropoffLocationName={offer.dropoffLocationName}
        pickupDate={offer.pickupDate}
        pickupTime={offer.pickupTime}
        pickupLocationName={offer.pickupLocationName}
      />
      <CarDetailsSection
        brand={offer.carDetails.supplierName}
        model={offer.carDetails.model}
        carType={offer.carDetails.carType}
        rentalDays={offer.rentalDays}
        isAutomatic={offer.carDetails.isAutoGear}
      />
      <IncludedSection planInclusions={offer.planInclusions} />

      <PayAtPickupSection
        fees={offer.payAtPickup.fees}
        selectedAddons={offer.payAtPickup.selectedAddons ?? []}
        currency={offer.currencyCode}
      />

      {offer.isErpIncluded && (
        <div className="mt-2 rounded-md p-1 flex items-start gap-2 font-bold">
          <Image
            src="/assets/icons/stamp.gif"
            alt="stamp"
            width={24}
            height={24}
            className="w-6 h-6"
          />
          {t("summary.erpIncluded")}
        </div>
      )}
      <div className="text-white bg-brand py-3 5 px-5 flex justify-between items-center rounded-xl">
        <span className="type-paragraph">
          {t("priceOffer.summary.labels.totalToPay")}
        </span>
        <h4 className="type-h4">
          {formatPrice(offer.totalPrice, offer.currencyCode)}
        </h4>
      </div>
      <div className="border-yellow-400 border my-4 rounded-xl bg-yellow-50 text-yellow-800 p-4">
        ⚠️ {t("summary.priceOfferDisclaimer")}
      </div>
    </div>
  );
}
