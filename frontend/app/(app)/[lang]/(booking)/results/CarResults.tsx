"use client";

import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { booking } from "@/shared/client";
import { useFilteredCars } from "./useFilteredCars";
import { CarGroupsFilter } from "./_components/filters/CarGroupsFilter";
import { FiltersPanel } from "./_components/filters/FiltersPanel";
import { useCheckboxFilters } from "./_hooks/useCheckboxFilters";
import { useTranslations } from "next-intl";
import { useAcrissCodesFilter } from "./_hooks/useAcrissCodesFilter";

interface CarResultsProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarResults({ searchRequest }: CarResultsProps) {
  const t = useTranslations("booking.results");
  const { data } = useAvailableCars(searchRequest);

  const { acrissFilterFn, selectedGroups, setSelectedGroups } =
    useAcrissCodesFilter();

  const {
    selectedFilters,
    toggleOption,
    clearAll,
    filterFunctions,
    hasActiveFilters,
  } = useCheckboxFilters();

  const cars = data?.availableVehicles ?? [];

  const filteredCars = useFilteredCars(cars, [
    acrissFilterFn,
    ...filterFunctions,
  ]);

  return (
    <div>
      <CarGroupsFilter
        title={t("carGroupsFiltersTitle")}
        selectedGroups={selectedGroups}
        setSelectedGroups={setSelectedGroups}
      />

      <div className="mt-10 grid grid-cols-1 gap-8 lg:grid-cols-[312px_1fr]">
        <FiltersPanel
          cars={cars}
          selectedFilters={selectedFilters}
          onToggle={toggleOption}
          onClear={clearAll}
          hasActiveFilters={hasActiveFilters}
        />

        <div>
          {filteredCars.map((vehicle, i) => (
            <div key={i}>
              {vehicle.carDetails.model}------{vehicle.carDetails.acriss}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
