import { getAgentPriceOffer } from "@/shared/api/price-offers-api";
import { useSuspenseQuery } from "@tanstack/react-query";

export function usePriceOffer(priceOfferId: number) {
  const queryKey = ["priceOffer", priceOfferId];
  const suspenseResult = useSuspenseQuery({
    queryKey,
    queryFn: () => getAgentPriceOffer(priceOfferId),
  });

  return { ...suspenseResult, queryKey };
}
