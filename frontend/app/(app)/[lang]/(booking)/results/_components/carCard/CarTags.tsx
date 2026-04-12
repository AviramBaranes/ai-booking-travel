import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";

export function CarTags({ vehicle }: { vehicle: booking.AvailableVehicle }) {
  const t = useTranslations("booking.results");

  return (
    <div className="flex gap-2">
      <TagPill content={vehicle.carDetails.acriss} />
      <TagPill content={vehicle.carDetails.carType} />
      {vehicle.signals?.tags.map((tag) => (
        <TagPill key={tag} content={t(`carDetails.${tag}`)} />
      ))}
    </div>
  );
}

function TagPill({ content }: { content: string }) {
  return (
    <span className="text-xs font-normal bg-navy text-white rounded-full px-3 py-1">
      {content}
    </span>
  );
}
