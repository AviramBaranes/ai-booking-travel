import { listPriceOffers } from "@/shared/api/price-offers-api";
import { useQuery } from "@tanstack/react-query";
import { PriceOfferFilters } from "./usePriceOfferFilters";

interface UsePriceOffersParams extends PriceOfferFilters {
  Page: number;
}

export function usePriceOffers(params: UsePriceOffersParams) {
  const queryKey = ["priceOffers", params];
  const suspenseResult = useQuery({
    queryKey,
    queryFn: () =>
      listPriceOffers({
        Page: params.Page,
        Name: params.name || undefined,
        Status: params.status ?? undefined,
      }),
  });

  return { ...suspenseResult, queryKey };
}
