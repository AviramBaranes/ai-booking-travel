import { booking } from "@/shared/client";
import { useMemo, useState } from "react";
import {
  FILTERS_LIST,
  type FilterConfig,
} from "../../_components/_constants/filtersList";
import { getNestedValue, toFilterValue } from "./useFiltersOptions";

export type SelectedFilters = Map<FilterConfig["id"], Set<string>>;

type CarFilter = (car: booking.AvailableVehicle) => boolean;

export function useCheckboxFilters() {
  const [selectedFilters, setSelectedFilters] = useState<SelectedFilters>(
    new Map(),
  );

  function toggleOption(filterId: FilterConfig["id"], value: string) {
    setSelectedFilters((prev) => {
      const next = new Map(prev);
      const nextSet = new Set(next.get(filterId) ?? []);

      if (nextSet.has(value)) {
        nextSet.delete(value);
      } else {
        nextSet.add(value);
      }

      if (nextSet.size === 0) {
        next.delete(filterId);
      } else {
        next.set(filterId, nextSet);
      }

      return next;
    });
  }

  function clearAll() {
    setSelectedFilters(new Map());
  }

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

        return (car: booking.AvailableVehicle) => {
          const value = getNestedValue(car, filter.filterKey);
          const normalized = toFilterValue(value);
          return normalized !== null && values.has(normalized);
        };
      })
      .filter((filter): filter is CarFilter => filter !== null);
  }, [selectedFilters]);

  const hasActiveFilters = selectedFilters.size > 0;

  return {
    selectedFilters,
    toggleOption,
    clearAll,
    filterFunctions,
    hasActiveFilters,
  };
}
