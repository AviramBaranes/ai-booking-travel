import { booking } from "@/shared/client";
import { useMemo, useState } from "react";

export type PlansCountFilter = "2" | "3" | null;
export type AddOnsFilter = "has" | "not" | null;

export function useDevFilters() {
  const isDevelopment = process.env.NODE_ENV === "development";

  const [plansCountFilter, setPlansCountFilter] =
    useState<PlansCountFilter>(null);
  const [addOnsFilter, setAddOnsFilter] = useState<AddOnsFilter>(null);

  function togglePlansCount(value: Exclude<PlansCountFilter, null>) {
    setPlansCountFilter((current) => (current === value ? null : value));
  }

  function toggleAddOns(value: Exclude<AddOnsFilter, null>) {
    setAddOnsFilter((current) => (current === value ? null : value));
  }

  const filterFn = useMemo(() => {
    return (car: booking.AvailableVehicle) => {
      if (!isDevelopment) {
        return true;
      }

      if (
        plansCountFilter !== null &&
        car.plans.length !== Number(plansCountFilter)
      ) {
        return false;
      }

      const addOnsCount = car.addOns?.length ?? 0;

      if (addOnsFilter === "has" && addOnsCount === 0) {
        return false;
      }

      if (addOnsFilter === "not" && addOnsCount > 0) {
        return false;
      }

      return true;
    };
  }, [addOnsFilter, isDevelopment, plansCountFilter]);

  return {
    isDevelopment,
    plansCountFilter,
    addOnsFilter,
    togglePlansCount,
    toggleAddOns,
    filterFn,
  };
}
