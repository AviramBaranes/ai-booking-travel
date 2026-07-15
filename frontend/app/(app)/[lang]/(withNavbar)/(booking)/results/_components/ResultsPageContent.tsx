"use client";

import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { ResultsLoading } from "./ResultsLoading";
import ErrorResultPageContent from "./ErrorPage";
import { BookingStepper } from "../../_components/BookingStepper";
import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { ExpiredSearchGate } from "../../_components/ExpiredSearchGate";
import { CarResults } from "../CarResults";
import { SearchQuery, toSearchRequest } from "../searchQuery";

export function ResultsPageContent({ query }: { query: SearchQuery }) {
  const searchRequest = toSearchRequest(query);
  const { isLoading, error, data } = useAvailableCars(searchRequest, {
    fromCache: false,
  });

  if (isLoading) {
    return <ResultsLoading />;
  }

  if (error || (data && !data.availableVehicles.length)) {
    return <ErrorResultPageContent />;
  }

  return (
    <div className="lg:w-2/3 mx-auto lg:pt-15 pb-6">
      <BookingStepper currentStep="results" />
      <div className="lg:my-4">
        <SearchDataBanner
          pickUpLocationId={query.pickupLocationId}
          dropOffLocationId={query.dropoffLocationId}
          pickUpTime={query.pickupTime}
          dropOffTime={query.dropoffTime}
          pickUpDate={query.pickupDate}
          dropOffDate={query.dropoffDate}
          driverAge={query.driverAge}
          couponCode={query.couponCode}
          searchRequest={searchRequest}
          showButton
        />
      </div>
      <ExpiredSearchGate>
        <CarResults searchRequest={searchRequest} />
      </ExpiredSearchGate>
    </div>
  );
}
