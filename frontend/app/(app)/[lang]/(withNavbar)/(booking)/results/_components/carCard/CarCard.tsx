import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { SupplierLogo } from "../../../../../../../shared/components/booking/SupplierLogo";
import { CarTags } from "./CarTags";
import { CarModel } from "./CarModel";
import { CarDetailsPills } from "../../../../../../../shared/components/booking/CarDetailsPills";
import { CarChecks } from "./CarChecks";
import { CarSignals } from "./CarSignals";
import { CarPriceDetails } from "./CarPriceDetails";

interface CarCardProps {
  vehicle: booking.AvailableVehicle;
  daysCount: number;
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarCard({ vehicle, daysCount, searchRequest }: CarCardProps) {
  const t = useTranslations("booking.results");

  return (
    <div className="shadow-[0_4px_12px_0_rgba(63,63,63,0.10)] rounded-2xl h-80 border-cars-border relative pr-4 flex gap-2 bg-white overflow-hidden">
      <div className="flex flex-col items-center gap-6 mt-10 mb-12">
        <SupplierLogo supplierName={vehicle.carDetails.supplierName} />
        <Image
          src={vehicle.carDetails.imageUrl}
          alt={vehicle.carDetails.model}
          width={176}
          height={100}
          className="w-44 h-25 object-cover"
        />
      </div>
      <div className="flex flex-col items-start gap-2 my-auto">
        <CarTags vehicle={vehicle} />
        <CarModel
          model={vehicle.carDetails.model}
          orSimilarText={t("carDetails.orSimilar")}
        />
        <CarDetailsPills carDetails={vehicle.carDetails} />
        <CarChecks
          checks={[
            {
              text:
                vehicle.locationDetails.locationType === "Shuttle"
                  ? t(`carDetails.shuttlePickup`)
                  : t(`carDetails.terminalPickup`),
              image: "/assets/icons/V.svg",
            },
            {
              text: t(`carDetails.erpRecommendation`),
              image: "/assets/icons/stamp.gif",
            },
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
