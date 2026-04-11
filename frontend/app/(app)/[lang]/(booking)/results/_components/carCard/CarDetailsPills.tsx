import Image from "next/image";
import { useTranslations } from "next-intl";
import { Cable, User } from "lucide-react";

import { booking } from "@/shared/client";

const CAR_DETAILS_PILLS = [
  {
    key: "seats",
    icon: User,
  },
  {
    key: "doors",
    image: "/assets/icons/Doors.svg",
  },
  {
    key: "isAutoGear",
    image: "/assets/icons/Gear.svg",
    translationKey: "auto",
  },
  {
    key: "isElectric",
    icon: Cable,
    translationKey: "electric",
  },
  {
    key: "hasAC",
    image: "/assets/icons/AC.svg",
    translationKey: "ac",
  },
  {
    key: "bags",
    image: "/assets/icons/Bags.svg",
  },
];

export function CarDetailsPills({
  vehicle,
}: {
  vehicle: booking.AvailableVehicle;
}) {
  const t = useTranslations("booking.results");
  return (
    <div className="flex items-center gap-2 mt-4">
      {CAR_DETAILS_PILLS.map(({ key, icon: Icon, translationKey, image }) => {
        const value =
          vehicle.carDetails[key as keyof typeof vehicle.carDetails];
        if (!value) return null;

        return (
          <div
            key={key}
            className="flex items-center gap-1 bg-[#E7E9F5] px-4 py-1 rounded-full text-sm font-normal"
          >
            {Icon ? (
              <Icon size={16} className="text-black/80" />
            ) : (
              <Image
                src={image}
                alt={`${key} icon`}
                width={16}
                height={16}
                className="w-4"
              />
            )}
            <span className="text-navy">
              {translationKey ? t(`carDetails.${translationKey}`) : value}
            </span>
          </div>
        );
      })}
    </div>
  );
}
