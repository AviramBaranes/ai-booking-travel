"use client";

import { SuppliersGallery } from "@/payload-types";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { useFilteredCars } from "./useFilteredCars";
import { CarCard } from "./_components/carCard/CarCard";
import { CarGroupsFilter } from "./_components/filters/CarGroupsFilter";
import { FiltersPanel } from "./_components/filters/FiltersPanel";
import { useCheckboxFilters } from "./_hooks/useCheckboxFilters";
import { useAcrissCodesFilter } from "./_hooks/useAcrissCodesFilter";
import { DevFilters } from "./_components/filters/DevFilters";
import { useDevFilters } from "./_hooks/useDevFilters";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

interface CarResultsProps {
  supplierGallery: SuppliersGallery;
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarResults({
  searchRequest,
  supplierGallery,
}: CarResultsProps) {
  const t = useTranslations("booking.results");
  const { data } = useAvailableCars(searchRequest);
  const {
    acrissFilterFn,
    selectedGroups,
    setSelectedGroups,
    clearAcrissFilters,
  } = useAcrissCodesFilter();

  const {
    isDevelopment,
    plansCountFilter,
    addOnsFilter,
    togglePlansCount,
    toggleAddOns,
    filterFn: devFilterFn,
  } = useDevFilters();

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
    devFilterFn,
  ]);

  return (
    <div>
      <CarGroupsFilter
        title={t("carGroupsFiltersTitle")}
        selectedGroups={selectedGroups}
        setSelectedGroups={setSelectedGroups}
      />

      {isDevelopment && (
        <DevFilters
          plansCountFilter={plansCountFilter}
          addOnsFilter={addOnsFilter}
          onPlansCountChange={togglePlansCount}
          onAddOnsChange={toggleAddOns}
        />
      )}

      <div className="mt-10 flex gap-6 justify-between">
        <div className="w-1/4">
          <FiltersPanel
            cars={cars}
            selectedFilters={selectedFilters}
            onToggle={toggleOption}
            onClear={clearAll}
            hasActiveFilters={hasActiveFilters}
          />
        </div>

        {filteredCars.length ? (
          <div className="w-3/4 flex flex-col gap-6">
            {filteredCars.map((vehicle) => (
              <CarCard
                key={vehicle.id}
                daysCount={data?.daysCount ?? 0}
                supplierGallery={supplierGallery}
                vehicle={vehicle}
                searchRequest={searchRequest}
              />
            ))}
          </div>
        ) : (
          <div className="p-20 text-center flex flex-col items-center gap-4">
            <h4 className="type-h4 text-navy">{t("error.filterNoResults")}</h4>
            <Button
              variant="outline"
              type="button"
              onClick={() => {
                clearAcrissFilters();
                clearAll();
              }}
              className="type-paragraph font-normal flex items-center gap-1 mt-6"
            >
              <X />
              {t("clearFilters")}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
