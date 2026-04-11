import { SuppliersGallery } from "@/payload-types";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { SupplierLogo } from "./SupplierLogo";
import { CarTags } from "./CarTags";
import { CarModel } from "./CarModel";
import { CarDetailsPills } from "./CarDetailsPills";
import { CarChecks } from "./CarChecks";
import { CarSignals } from "./CarSignals";

interface CarCardProps {
  supplierGallery: SuppliersGallery;
  vehicle: booking.AvailableVehicle;
}

export function CarCard({ vehicle, supplierGallery }: CarCardProps) {
  const t = useTranslations("booking.results");

  return (
    <div className="shadow-[0_4px_12px_0_rgba(63,63,63,0.10)] rounded-2xl border-cars-border relative pr-4 flex gap-2 bg-white overflow-hidden">
      <div className="flex flex-col items-center gap-6 mt-10 mb-12">
        <SupplierLogo
          supplierName={vehicle.carDetails.supplierName}
          supplierGallery={supplierGallery}
        />
        <Image
          src={vehicle.carDetails.imageUrl}
          alt={vehicle.carDetails.model}
          width={175}
          height={150}
        />
      </div>
      <div className="flex flex-col items-start gap-2 mt-5">
        <CarTags vehicle={vehicle} />
        <CarModel
          model={vehicle.carDetails.model}
          orSimilarText={t("carDetails.orSimilar")}
        />
        <CarDetailsPills vehicle={vehicle} />
        <CarChecks
          checks={[
            vehicle.locationDetails.locationType === "shuttle"
              ? t(`carDetails.shuttlePickup`)
              : t(`carDetails.terminalPickup`),
          ]}
        />
      </div>
      <CarSignals vehicle={vehicle} />
    </div>
  );
}
