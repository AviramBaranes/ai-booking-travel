import { availability } from "@/shared/client";
import { useMemo } from "react";
import { CAR_GROUPS_FILTERS } from "../../_components/_constants/carGroupsFilters";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";

export function useAcrissCodesFilter() {
  const selectedGroups = useBookingSessionStore((state) => state.carGroupFilters);
  const clearAcrissFilters = useBookingSessionStore((state) => state.clearCarGroupFilters);

  const acrissCodes = useMemo(() => {
    return new Set(
      Array.from(selectedGroups).flatMap((groupName) => {
        const group = CAR_GROUPS_FILTERS.find((g) => g.name === groupName);
        return group ? group.acrissCodes : [];
      }),
    );
  }, [selectedGroups]);

  const filterFunction = (car: availability.AvailableVehicle) => {
    return acrissCodes.size === 0 || acrissCodes.has(car.carDetails.acriss);
  };

  
  return {
    selectedGroups,
    clearAcrissFilters,
    acrissFilterFn: filterFunction,
  };
}
