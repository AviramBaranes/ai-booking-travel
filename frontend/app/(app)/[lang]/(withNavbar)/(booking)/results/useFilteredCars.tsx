import { availability } from "@/shared/client";
import { useMemo } from "react";

type Filter = (car: availability.AvailableVehicle) => boolean;

export function useFilteredCars(
  cars: availability.AvailableVehicle[],
  filters: Filter[],
) {
  const filteredCars = useMemo(() => {
    return cars.filter((car) => filters.every((filter) => filter(car)));
  }, [cars, filters]);

  return filteredCars;
}
