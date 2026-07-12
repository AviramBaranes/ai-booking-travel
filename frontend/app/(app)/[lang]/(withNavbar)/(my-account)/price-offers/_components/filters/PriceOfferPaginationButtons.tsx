"use client";

import { PaginationButtons } from "../../../_components/PaginationButtons";
import { useTranslations } from "next-intl";
import { usePriceOfferFilters } from "../../_hooks/usePriceOfferFilters";
import { usePriceOffers } from "../../_hooks/usePriceOffers";

export function PriceOfferPaginationButtons() {
  const t = useTranslations("MyAccount.priceOffers");
  const { lang, searchParams, filters, page } = usePriceOfferFilters();
  const { data } = usePriceOffers({
    Page: page,
    ...filters,
  });

  if (!data) return null;

  return (
    <PaginationButtons
      total={data.total}
      page={page}
      searchParams={searchParams}
      basePath={`/${lang}/price-offers`}
      previousLabel={t("pagination.prev")}
      nextLabel={t("pagination.next")}
    />
  );
}
