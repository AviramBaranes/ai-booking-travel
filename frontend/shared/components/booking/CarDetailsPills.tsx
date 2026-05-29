import Image from "next/image";
import type { ElementType } from "react";
import { useTranslations } from "next-intl";
import { Cable, User } from "lucide-react";

import { broker } from "@/shared/client";

type BookingResultsT = ReturnType<typeof useTranslations<"booking.results">>;
type CarDetailPillKey =
  | "seats"
  | "doors"
  | "isAutoGear"
  | "isElectric"
  | "hasAC"
  | "bags";

type PillIcon =
  | { icon: ElementType; imageSrc?: never }
  | { imageSrc: string; icon?: never };

type PillConfig = PillIcon & {
  key: CarDetailPillKey;
  getLabel: (carDetails: broker.CarDetails, t: BookingResultsT) => string | number | null;
};

const CAR_DETAILS_PILLS = [
  {
    key: "seats",
    icon: User,
    getLabel: ({ seats }) => (seats > 0 ? seats : null),
  },
  {
    key: "doors",
    imageSrc: "/assets/icons/Doors.svg",
    getLabel: ({ doors }) => (doors > 0 ? doors : null),
  },
  {
    key: "isAutoGear",
    imageSrc: "/assets/icons/Gear.svg",
    getLabel: ({ isAutoGear }, t) =>
      isAutoGear ? t("carDetails.auto") : t("carDetails.manual"),
  },
  {
    key: "isElectric",
    icon: Cable,
    getLabel: ({ isElectric }, t) =>
      isElectric ? t("carDetails.electric") : null,
  },
  {
    key: "hasAC",
    imageSrc: "/assets/icons/AC.svg",
    getLabel: ({ hasAC }, t) => (hasAC ? t("carDetails.ac") : null),
  },
  {
    key: "bags",
    imageSrc: "/assets/icons/Bags.svg",
    getLabel: ({ bags }) => (bags > 0 ? bags : null),
  },
] satisfies PillConfig[];

export function CarDetailsPills({
  carDetails,
}: {
  carDetails: broker.CarDetails;
}) {
  const t = useTranslations("booking.results");
  return (
    <div className="flex flex-wrap items-center gap-2 mt-4">
      {CAR_DETAILS_PILLS.map((config) => {
          const { key, icon: Icon } = config;
          const label = config.getLabel(carDetails, t);
          if (label === null) return null;

          return (
            <div
              key={key}
              className="flex items-center gap-1 bg-[#E7E9F5] px-4 py-1 rounded-full text-sm font-normal"
            >
              {Icon ? (
                <Icon size={16} className="text-black/80" />
              ) : (
                <Image
                  src={config.imageSrc}
                  alt={`${key} icon`}
                  width={16}
                  height={16}
                  className="w-4 h-4"
                />
              )}
              <span className="text-navy">{label}</span>
            </div>
          );
        }
      )}
    </div>
  );
}
