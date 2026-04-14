import { useQuery } from "@tanstack/react-query";
import { booking } from "../client";
import { searchAvailableCars } from "../api/booking-api";

export const bookingKeys = {
  availability: (params: booking.SearchAvailabilityRequest) =>
    ["booking", "availability", params] as const,
};

interface UseAvailableCarsOptions {
  fromCache?: boolean;
}

export function useAvailableCars(
  params: booking.SearchAvailabilityRequest,
  opts?: UseAvailableCarsOptions,
) {
  return useQuery({
    queryKey: bookingKeys.availability(params),
    queryFn: () => searchAvailableCars(params),
    staleTime: 15 * 60 * 1000,
    enabled: opts?.fromCache !== true,
  });
}
