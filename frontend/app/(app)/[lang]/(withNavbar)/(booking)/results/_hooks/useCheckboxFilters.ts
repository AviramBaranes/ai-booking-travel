import { availability } from "@/shared/client";
import { useMemo, useState } from "react";
import {
  FILTERS_LIST,
  type FilterConfig,
} from "../../_components/_constants/filtersList";
import { getNestedValue, toFilterValue } from "./useFiltersOptions";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";

export type SelectedFilters = Map<FilterConfig["id"], Set<string>>;

type CarFilter = (car: availability.AvailableVehicle) => boolean;

export function useCheckboxFilters() {
  const selectedFilters = useBookingSessionStore((state) => state.checkboxFilters);
  const clearAllCheckboxFilters = useBookingSessionStore((state) => state.clearAllCheckboxFilters);

  const filterFunctions = useMemo<CarFilter[]>(() => {
    const activeEntries = Array.from(selectedFilters.entries()).filter(
      ([, values]) => values.size > 0,
    );

    return activeEntries
      .map(([filterId, values]) => {
        const filter = FILTERS_LIST.find((item) => item.id === filterId);
        if (!filter) {
          return null;
        }

        return (car: availability.AvailableVehicle) => {
          const value = getNestedValue(car, filter.filterKey);
          const normalized = toFilterValue(value);
          return normalized !== null && values.has(normalized);
        };
      })
      .filter((filter): filter is CarFilter => filter !== null);
  }, [selectedFilters]);

  const hasActiveFilters = selectedFilters.size > 0;

  return {
    clearAll: clearAllCheckboxFilters,
    filterFunctions,
    hasActiveFilters,
  };
}
