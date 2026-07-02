"use client";

import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { availability } from "@/shared/client";
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
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { useEffect } from "react";
import { FiltersSheet } from "./_components/filters/FiltersSheet";

interface CarResultsProps {
  searchRequest: availability.SearchAvailabilityParams;
}

export function CarResults({ searchRequest }: CarResultsProps) {
  const t = useTranslations("booking.results");
  const { data } = useAvailableCars(searchRequest);
  const { acrissFilterFn, clearAcrissFilters } = useAcrissCodesFilter();

  const {
    isDevelopment,
    plansCountFilter,
    addOnsFilter,
    togglePlansCount,
    toggleAddOns,
    filterFn: devFilterFn,
  } = useDevFilters();

  const { clearAll, filterFunctions, hasActiveFilters } = useCheckboxFilters();

  const cars = data?.availableVehicles ?? [];

  const filteredCars = useFilteredCars(cars, [
    acrissFilterFn,
    ...filterFunctions,
    devFilterFn,
  ]);

  const { clearSession } = useBookingSessionStore();
  useEffect(() => {
    clearSession();
  }, [clearSession]);

  return (
    <div>
      <FiltersSheet cars={cars} hasActiveFilters={hasActiveFilters} />
      <div className="lg:block hidden">
        <CarGroupsFilter title={t("carGroupsFiltersTitle")} />
      </div>

      {isDevelopment && (
        <DevFilters
          plansCountFilter={plansCountFilter}
          addOnsFilter={addOnsFilter}
          onPlansCountChange={togglePlansCount}
          onAddOnsChange={toggleAddOns}
        />
      )}

      <div className="mt-10 flex gap-6 justify-between">
        <div className="w-1/4 hidden lg:flex">
          <FiltersPanel cars={cars} hasActiveFilters={hasActiveFilters} />
        </div>

        {filteredCars.length ? (
          <div className="lg:w-3/4 mx-5 lg:mx-0 flex flex-col gap-6">
            {filteredCars.map((vehicle) => (
              <CarCard
                key={vehicle.id}
                daysCount={data?.daysCount ?? 0}
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
