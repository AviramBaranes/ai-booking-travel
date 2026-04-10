"use client";

import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { booking } from "@/shared/client";
import { useFilteredCars } from "./useFilteredCars";
import { CarGroupsFilter } from "./_components/CarGroupsFilter";
import { useTranslations } from "next-intl";
import { useState } from "react";

interface CarResultsProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarResults({ searchRequest }: CarResultsProps) {
  const t = useTranslations("booking.results");
  const { data } = useAvailableCars(searchRequest);
  const [acrissFilters, setAcrissFilters] = useState<Set<string>>(new Set());

  const filteredCars = useFilteredCars(data?.availableVehicles ?? [], [
    (car) =>
      acrissFilters.size === 0 || acrissFilters.has(car.carDetails.acriss),
  ]);

  return (
    <div>
      <CarGroupsFilter
        title={t("carGroupsFiltersTitle")}
        setAcrissCodes={setAcrissFilters}
      />
      {filteredCars.map((vehicle, i) => (
        <div key={i}>
          {vehicle.carDetails.model}------{vehicle.carDetails.acriss}
        </div>
      ))}
    </div>
  );
}
