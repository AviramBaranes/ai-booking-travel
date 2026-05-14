import { MapPin, Cable, Cog, Users, KeyRound } from "lucide-react";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type TranslateFn = (key: string, params?: any) => string;

export interface FilterConfig {
  id:
    | "locationType"
    | "transmissionType"
    | "isElectric"
    | "passengersCount"
    | "rentalCompany";
  icon: typeof MapPin;
  filterKey: string;
  titleKey: string;
  getOptionLabel: (value: string, t: TranslateFn) => string;
}

export const FILTERS_LIST: FilterConfig[] = [
  {
    id: "locationType",
    icon: MapPin,
    filterKey: "locationDetails.locationType",
    titleKey: "booking.results.filters.pickupLocation.title",
    getOptionLabel: (value, t) =>
      t(`booking.results.filters.pickupLocation.${value}`),
  },
  {
    id: "transmissionType",
    icon: Cog,
    filterKey: "carDetails.isAutoGear",
    titleKey: "booking.results.filters.transmissionType.title",
    getOptionLabel: (value, t) =>
      t(
        value === "true"
          ? "booking.results.filters.transmissionType.automatic"
          : "booking.results.filters.transmissionType.manual",
      ),
  },
  {
    id: "isElectric",
    icon: Cable,
    filterKey: "carDetails.isElectric",
    titleKey: "booking.results.filters.electricCar.title",
    getOptionLabel: (value, t) =>
      t(
        value === "true"
          ? "booking.results.filters.electricCar.fullElectric"
          : "booking.results.filters.electricCar.regular",
      ),
  },
  {
    id: "passengersCount",
    icon: Users,
    filterKey: "carDetails.seats",
    titleKey: "booking.results.filters.passengersCount.title",
    getOptionLabel: (value, t) =>
      t("booking.results.filters.passengersCount.option", { count: value }),
  },
  {
    id: "rentalCompany",
    icon: KeyRound,
    filterKey: "carDetails.supplierName",
    titleKey: "booking.results.filters.rentalCompany.title",
    getOptionLabel: (value) => value,
  },
];
