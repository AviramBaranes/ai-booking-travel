"use client";

import { Button } from "@/components/ui/button";
import { X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import {
  PriceOfferFilterKey,
  usePriceOfferFilters,
} from "../../_hooks/usePriceOfferFilters";

export function ClearFilterRow() {
  const router = useRouter();
  const { lang, searchParams, activeFilters } = usePriceOfferFilters();
  const t = useTranslations("MyAccount.priceOffers");
  const tStatus = useTranslations("MyAccount.priceOffer.summary.status");
  const labelByKey: Record<PriceOfferFilterKey, string> = {
    name: t("namePlaceholder"),
    status: t("statusPlaceholder"),
  };

  const clearFilter = (key: PriceOfferFilterKey) => {
    const nextQuery = new URLSearchParams(searchParams.toString());
    nextQuery.delete(key);

    const queryString = nextQuery.toString();
    const basePath = `/${lang}/price-offers`;
    router.push(queryString ? `${basePath}?${queryString}` : basePath);
  };

  if (activeFilters.length === 0) {
    return null;
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {activeFilters.map((filter) => {
        const value =
          filter.key === "status" ? tStatus(filter.value) : filter.value;

        return (
          <Button
            key={filter.key}
            variant="outline"
            type="button"
            onClick={() => clearFilter(filter.key)}
            className="type-paragraph font-normal flex items-center gap-1"
          >
            <X className="size-4" />

            <span className="font-semibold">
              {labelByKey[filter.key]}
              {": "}
            </span>
            {value}
          </Button>
        );
      })}
    </div>
  );
}