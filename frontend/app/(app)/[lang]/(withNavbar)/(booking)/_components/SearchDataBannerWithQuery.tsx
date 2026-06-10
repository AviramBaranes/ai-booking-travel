"use client";

import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { useSearchRequest } from "../_hooks/useSearchRequest";

export function SearchDataBannerWithQuery() {
  const { searchRequest, query } = useSearchRequest();

  return (
    <SearchDataBanner
      searchRequest={searchRequest}
      pickUpLocationId={query.pickupLocationId}
      dropOffLocationId={query.dropoffLocationId}
      pickUpTime={query.pickupTime}
      dropOffTime={query.dropoffTime}
      pickUpDate={query.pickupDate}
      dropOffDate={query.dropoffDate}
      driverAge={query.driverAge}
      couponCode={query.couponCode}
      showButton
      fromCache
    />
  );
}
