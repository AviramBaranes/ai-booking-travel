import { useQuery } from "@tanstack/react-query";
import { availability } from "../client";
import { searchAvailableCars } from "../api/booking-api";

export const bookingKeys = {
  availability: (params: availability.SearchAvailabilityParams) =>
    ["booking", "availability", params] as const,
};

interface UseAvailableCarsOptions {
  fromCache?: boolean;
}

export function useAvailableCars(
  params: availability.SearchAvailabilityParams,
  opts?: UseAvailableCarsOptions,
) {
  return useQuery({
    queryKey: bookingKeys.availability(params),
    queryFn: () => searchAvailableCars(params),
    staleTime: Infinity,
    enabled: opts?.fromCache !== true,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
