"use client";

import { useTranslations } from "next-intl";
import { usePriceOfferFilters } from "../_hooks/usePriceOfferFilters";
import { usePriceOffers } from "../_hooks/usePriceOffers";

export function PriceOfferResultsCounter() {
  const t = useTranslations("MyAccount.priceOffers");
  const { filters, page } = usePriceOfferFilters();
  const {
    data: {
      total,
      priceOffers: { length },
    },
  } = usePriceOffers({
    Page: page,
    ...filters,
  });

  return (
    <p className="text-xs text-text-secondary">
      {t("showingXResults", {
        count: length,
        total: total,
      })}
    </p>
  );
}