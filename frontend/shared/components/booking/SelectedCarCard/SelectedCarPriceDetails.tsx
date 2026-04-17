import { booking } from "@/shared/client";
import { PriceDetailRow } from "../PriceDetailRow";
import { useTranslations } from "next-intl";
import { RentalPriceForDays } from "../RentalPriceForDays";
import { formatPrice } from "@/shared/utils/formatPrice";

export function SelectedCarPriceDetails({
  vehicle,
  selectedPlanIndex,
  isErpSelected,
  daysCount,
}: {
  vehicle: booking.AvailableVehicle;
  selectedPlanIndex: number;
  isErpSelected: boolean;
  daysCount: number;
}) {
  const t = useTranslations("booking.results");

  const selectedPlan = vehicle.plans[selectedPlanIndex];
  return (
    <>
      {selectedPlan.fullPrice !== selectedPlan.price && (
        <>
          <PriceDetailRow
            altText="coins icon"
            iconSrc="/assets/icons/coins.svg"
            label={t("carDetails.priceBeforeDiscount")}
            price={selectedPlan.fullPrice}
            currency={vehicle.priceDetails.currency}
          />
          <PriceDetailRow
            altText="discount icon"
            iconSrc="/assets/icons/Discount-Green.svg"
            label={t("carDetails.savings")}
            price={selectedPlan.fullPrice - selectedPlan.price}
            currency={vehicle.priceDetails.currency}
          />
        </>
      )}

      {isErpSelected && (
        <PriceDetailRow
          altText="stamp icon"
          iconSrc="/assets/icons/stamp.gif"
          label={t("carDetails.coveragePackage")}
          price={selectedPlan.erpPrice}
          currency={vehicle.priceDetails.currency}
        />
      )}

      <hr className="mb-6 mt-3" />

      <div className="flex justify-between items-start">
        <div>
          <p className="type-label text-brand">{t("carDetails.totalToPay")}</p>
          <RentalPriceForDays daysCount={daysCount} />
        </div>
        <h5 className="type-h5 text-navy">
          {formatPrice(
            selectedPlan.price + (isErpSelected ? selectedPlan.erpPrice : 0),
            vehicle.priceDetails.currency,
          )}
        </h5>
      </div>
    </>
  );
}
