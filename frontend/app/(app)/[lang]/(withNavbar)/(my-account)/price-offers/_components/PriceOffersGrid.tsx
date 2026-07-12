"use client";

import { useTranslations } from "next-intl";
import { usePriceOfferFilters } from "../_hooks/usePriceOfferFilters";
import { usePriceOffers } from "../_hooks/usePriceOffers";
import { PriceOfferCard } from "./PriceOfferCard";
import { AccountGridSkeleton } from "../../_components/AccountGridSkeleton";

export function PriceOffersGrid() {
  const t = useTranslations("MyAccount.priceOffers");
  const { filters, page } = usePriceOfferFilters();
  const { isLoading, data } = usePriceOffers({
    Page: page,
    ...filters,
  });

  if (isLoading || !data) {
    return <AccountGridSkeleton />;
  }

  const { total, priceOffers } = data;

  if (total === 0 || priceOffers.length === 0) {
    return (
      <div className="p-10 pb-60 mx-auto text-center">
        <h4 className="type-h4 text-navy">{t("noPriceOffersFound")}</h4>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-4 gap-6">
      {priceOffers.map((priceOffer) => (
        <PriceOfferCard priceOffer={priceOffer} key={priceOffer.id} />
      ))}
    </div>
  );
}
