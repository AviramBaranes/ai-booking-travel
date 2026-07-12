"use client";

import { useTranslations } from "next-intl";
import { usePriceOfferFilters } from "../_hooks/usePriceOfferFilters";
import { usePriceOffers } from "../_hooks/usePriceOffers";

export function PriceOfferResultsCounter() {
  const t = useTranslations("MyAccount.priceOffers");
  const { filters, page } = usePriceOfferFilters();
  const { isLoading, data } = usePriceOffers({
    Page: page,
    ...filters,
  });

  if (isLoading || !data) {
    return (
      <p className="text-xs text-text-secondary">
        {t("showingXResults", {
          count: "X",
          total: "X",
        })}
      </p>
    );
  }
  
  const {
    total,
    priceOffers: { length },
  } = data;

  return (
    <p className="text-xs text-text-secondary">
      {t("showingXResults", {
        count: length,
        total: total,
      })}
    </p>
  );
}
