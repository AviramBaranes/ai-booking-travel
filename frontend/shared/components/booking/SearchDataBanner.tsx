"use client";

import { useState } from "react";
import {
  SearchForm,
  SearchFormFields,
} from "@/app/(app)/[lang]/_components/home/SearchForm/SearchForm";
import { booking } from "@/shared/client";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { SearchDataBannerDisplay } from "./SearchDataBannerDisplay";
import { SearchDataFormWrapper } from "./SearchDataFormWrapper";

interface SearchDataBannerProps extends Omit<
  SearchFormFields,
  "pickUpLocation" | "dropOffLocation"
> {
  pickUpLocationId: number;
  dropOffLocationId: number;
  searchRequest: booking.SearchAvailabilityRequest;
  showButton?: boolean;
  fromCache?: boolean;
}

export function SearchDataBanner({
  pickUpLocationId,
  dropOffLocationId,
  pickUpTime,
  dropOffTime,
  pickUpDate,
  dropOffDate,
  driverAge,
  couponCode,
  showButton,
  searchRequest,
  fromCache,
}: SearchDataBannerProps) {
  const { data } = useAvailableCars(searchRequest, { fromCache });
  const pickUpLocationName = data?.pickupLocationName ?? "";
  const dropOffLocationName = data?.dropoffLocationName ?? "";
  const [showForm, setShowForm] = useState(false);

  if (showForm) {
    return (
      <SearchDataFormWrapper onClose={() => setShowForm(false)}>
        <SearchForm
          className="w-full"
          pickUpLocation={{
            id: pickUpLocationId,
            name: pickUpLocationName,
          }}
          dropOffLocation={{
            id: dropOffLocationId,
            name: dropOffLocationName,
          }}
          pickUpDate={pickUpDate}
          dropOffDate={dropOffDate}
          pickUpTime={pickUpTime}
          dropOffTime={dropOffTime}
          couponCode={couponCode}
          driverAge={driverAge}
        />
      </SearchDataFormWrapper>
    );
  }

  return (
    <SearchDataBannerDisplay
      pickUpLocationName={pickUpLocationName}
      dropOffLocationName={dropOffLocationName || pickUpLocationName}
      pickUpDate={pickUpDate}
      dropOffDate={dropOffDate}
      pickUpTime={pickUpTime}
      dropOffTime={dropOffTime}
      driverAge={driverAge}
      showButton={showButton}
      onEditClick={() => setShowForm(true)}
    />
  );
}
