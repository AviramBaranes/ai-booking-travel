import { availability } from "@/shared/client";
import { useTranslations } from "next-intl";
import { SelectedCarPriceDetails } from "./SelectedCarPriceDetails";
import { SelectedCarCardWrapper } from "./SelectedCarCardWrapper";
import { SelectedCarHeader } from "./SelectedCarHeader";

interface SelectedCarCardProps {
  daysCount: number;
  selectedPlanIndex: number;
  vehicle: availability.AvailableVehicle;
  isErpSelected: boolean;
  children?: React.ReactNode;
  headerClassName?: string;
}

export function SelectedCarCard({
  children,
  vehicle,
  daysCount,
  selectedPlanIndex,
  isErpSelected,
  headerClassName,
}: SelectedCarCardProps) {
  const t = useTranslations("booking.results");

  return (
    <SelectedCarCardWrapper>
      <SelectedCarHeader
        className={headerClassName}
        carDetails={vehicle.carDetails}
      />

      <SelectedCarPriceDetails
        daysCount={daysCount}
        isErpSelected={isErpSelected}
        selectedPlanIndex={selectedPlanIndex}
        vehicle={vehicle}
      />

      {children}
    </SelectedCarCardWrapper>
  );
}
