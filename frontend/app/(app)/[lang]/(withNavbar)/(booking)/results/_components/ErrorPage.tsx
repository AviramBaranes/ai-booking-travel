"use client";

import { useTranslations } from "next-intl";
import { SearchForm } from "../../../_components/home/SearchForm/SearchForm";
import { useLastSearch } from "../../../_components/home/SearchForm/lastSearch";
import { SearchQuery } from "../searchQuery";

export default function ErrorResultPageContent({
  query,
}: {
  query?: SearchQuery;
}) {
  const t = useTranslations("booking.results.error");
  const { lastSearch, raw } = useLastSearch();

  // The URL holds only location ids; take the names from the stored search when it is the same one.
  const pickupLocationName =
    lastSearch?.pickupLocation.id === query?.pickupLocationId
      ? lastSearch?.pickupLocation.name
      : undefined;
  const dropoffLocationName =
    lastSearch?.dropoffLocation.id === query?.dropoffLocationId
      ? lastSearch?.dropoffLocation.name
      : undefined;

  return (
    <div className="mt-20">
      <div className="w-10/12 mx-auto">
        {/* Remount once storage is read, since the form only takes its values on mount. */}
        <SearchForm
          key={raw ?? "none"}
          pickUpLocation={
            query && pickupLocationName
              ? { id: query.pickupLocationId, name: pickupLocationName }
              : undefined
          }
          dropOffLocation={
            query && dropoffLocationName
              ? { id: query.dropoffLocationId, name: dropoffLocationName }
              : undefined
          }
          pickUpDate={query?.pickupDate}
          dropOffDate={query?.dropoffDate}
          pickUpTime={query?.pickupTime}
          dropOffTime={query?.dropoffTime}
          driverAge={query?.driverAge}
          couponCode={query?.couponCode}
        />
      </div>
      <div className="p-30 text-center">
        <h4 className="type-h4 text-navy">{t("noResults")}</h4>
      </div>
    </div>
  );
}
