import { booking } from "@/shared/client";
import { HeaderSection } from "./HeaderSection";
import { RentalSummary } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/RentalSummary";
import { CarDetailsSection } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/CarDetailsSection";
import { IncludedSection } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/IncludedSection";
import { getTranslations } from "next-intl/server";
import { formatPrice } from "@/shared/utils/formatPrice";

export async function ClientOfferSummary({
  offer,
  lang,
}: {
  offer: booking.GetPriceOfferResponse;
  lang: string;
}) {
  const t = await getTranslations("MyAccount.priceOffer.summary.labels");

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
      />
      <IncludedSection planInclusions={offer.planInclusions} />
      <div className="text-white mt-6 bg-brand py-3 5 px-5 flex justify-between items-center rounded-xl">
        <span className="type-paragraph">{t("totalToPay")}</span>
        <h4 className="type-h4">
          {formatPrice(offer.totalPrice, offer.currencyCode)}
        </h4>
      </div>
    </div>
  );
}
