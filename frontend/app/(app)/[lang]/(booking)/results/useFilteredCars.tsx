import { booking } from "@/shared/client";
import { useMemo } from "react";

type Filter = (car: booking.AvailableVehicle) => boolean;

export function useFilteredCars(
  cars: booking.AvailableVehicle[],
  filters: Filter[],
) {
  const filteredCars = useMemo(() => {
    return cars.filter((car) => filters.every((filter) => filter(car)));
  }, [cars, filters]);

  return filteredCars;
}
