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
import { CarPriceDetails } from "./CarPriceDetails";

interface CarCardProps {
  vehicle: booking.AvailableVehicle;
  daysCount: number;
  supplierGallery: SuppliersGallery;
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarCard({
  vehicle,
  daysCount,
  supplierGallery,
  searchRequest,
}: CarCardProps) {
  const t = useTranslations("booking.results");

  return (
    <div className="shadow-[0_4px_12px_0_rgba(63,63,63,0.10)] rounded-2xl h-70 border-cars-border relative pr-4 flex gap-2 bg-white overflow-hidden">
      <div className="flex flex-col items-center gap-6 mt-10 mb-12">
        <SupplierLogo
          supplierName={vehicle.carDetails.supplierName}
          supplierGallery={supplierGallery}
        />
        <Image
          src={vehicle.carDetails.imageUrl}
          alt={vehicle.carDetails.model}
          width={176}
          height={100}
          className="w-44 h-25 object-cover"
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
            vehicle.locationDetails.locationType === "Shuttle"
              ? t(`carDetails.shuttlePickup`)
              : t(`carDetails.terminalPickup`),
          ]}
        />
      </div>
      <CarSignals vehicle={vehicle} />
      <CarPriceDetails
        vehicle={vehicle}
        searchRequest={searchRequest}
        daysCount={daysCount}
      />
    </div>
  );
}
