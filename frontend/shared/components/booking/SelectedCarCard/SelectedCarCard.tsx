import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { SelectedCarPriceDetails } from "./SelectedCarPriceDetails";
import { SelectedCarCardWrapper } from "./SelectedCarCardWrapper";
import { SelectedCarHeader } from "./SelectedCarHeader";

interface SelectedCarCardProps {
  daysCount: number;
  selectedPlanIndex: number;
  vehicle: booking.AvailableVehicle;
  isErpSelected: boolean;
  children?: React.ReactNode;
}

export function SelectedCarCard({
  children,
  vehicle,
  daysCount,
  selectedPlanIndex,
  isErpSelected,
}: SelectedCarCardProps) {
  const t = useTranslations("booking.results");

  return (
    <SelectedCarCardWrapper>
      <SelectedCarHeader carDetails={vehicle.carDetails} />

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
