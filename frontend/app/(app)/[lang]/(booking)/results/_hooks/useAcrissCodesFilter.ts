import { booking } from "@/shared/client";
import { useMemo, useState } from "react";
import { CAR_GROUPS_FILTERS } from "../../_components/_constants/carGroupsFilters";

export function useAcrissCodesFilter() {
  const [selectedGroups, setSelectedGroups] = useState<Set<string>>(new Set());

  const acrissCodes = useMemo(() => {
    return new Set(
      Array.from(selectedGroups).flatMap((groupName) => {
        const group = CAR_GROUPS_FILTERS.find((g) => g.name === groupName);
        return group ? group.acrissCodes : [];
      }),
    );
  }, [selectedGroups]);

  const filterFunction = (car: booking.AvailableVehicle) => {
    return acrissCodes.size === 0 || acrissCodes.has(car.carDetails.acriss);
  };

  const clearAcrissFilters = () => {
    setSelectedGroups(new Set());
  };

  return {
    selectedGroups,
    setSelectedGroups,
    clearAcrissFilters,
    acrissFilterFn: filterFunction,
  };
}
