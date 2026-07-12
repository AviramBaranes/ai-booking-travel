import { getAgentPriceOffer } from "@/shared/api/price-offers-api";
import { useQuery } from "@tanstack/react-query";

export function usePriceOffer(priceOfferId: number) {
  const queryKey = ["priceOffer", priceOfferId];
  const suspenseResult = useQuery({
    queryKey,
    queryFn: () => getAgentPriceOffer(priceOfferId),
  });

  return { ...suspenseResult, queryKey };
}
