import { booking } from "@/shared/client";
import { useMemo } from "react";
import { FILTERS_LIST } from "../../_components/_constants/filtersList";

export type FilterOption = {
  id: (typeof FILTERS_LIST)[number]["id"];
  titleKey: string;
  icon: (typeof FILTERS_LIST)[number]["icon"];
  getOptionLabel: (typeof FILTERS_LIST)[number]["getOptionLabel"];
  filterKey: string;
  options: string[];
};

export function getNestedValue(source: unknown, path: string): unknown {
  return path.split(".").reduce<unknown>((current, key) => {
    if (current && typeof current === "object") {
      return (current as Record<string, unknown>)[key];
    }
    return undefined;
  }, source);
}

export function toFilterValue(value: unknown): string | null {
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }

  if (typeof value === "number") {
    return value.toString();
  }

  if (typeof value === "string") {
    return value;
  }

  return null;
}

export function useFilterOptions(cars: booking.AvailableVehicle[]) {
  const filtersOptions = useMemo<FilterOption[]>(() => {
    return FILTERS_LIST.map((filter) => {
      const valuesSet = new Set<string>();

      for (const car of cars) {
        const value = getNestedValue(car, filter.filterKey);
        const filterValue = toFilterValue(value);
        if (filterValue !== null) {
          valuesSet.add(filterValue);
        }
      }

      return {
        id: filter.id,
        titleKey: filter.titleKey,
        icon: filter.icon,
        getOptionLabel: filter.getOptionLabel,
        filterKey: filter.filterKey,
        options: Array.from(valuesSet),
      };
    });
  }, [cars]);

  return filtersOptions;
}
