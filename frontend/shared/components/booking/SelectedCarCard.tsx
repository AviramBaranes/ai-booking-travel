import { RentalPriceForDays } from "@/shared/components/booking/RentalPriceForDays";
import { CarDetailsPills } from "@/shared/components/booking/CarDetailsPills";
import { SupplierLogo } from "@/shared/components/booking/SupplierLogo";
import { booking } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { PriceDetailRow } from "./PriceDetailRow";

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

  const selectedPlan = vehicle.plans[selectedPlanIndex];
  return (
    <div className="bg-white shadow-card p-6 flex rounded-2xl flex-col gap-2 justify-between border border-cars-border">
      <div className="flex-col flex items-center">
        <div className="mb-12">
          <SupplierLogo supplierName={vehicle.carDetails.supplierName} />
        </div>

        <Image
          src={vehicle.carDetails.imageUrl}
          alt={vehicle.carDetails.model}
          width={176}
          height={100}
          className="w-44 h-25 object-cover"
        />
      </div>
      <div className="flex gap-2 flex-col items-start">
        <h5 className="type-h5 text-navy">
          {vehicle.carDetails.model} ({vehicle.carDetails.acriss})
        </h5>
        <span className="type-paragraph text-navy">
          {t("carDetails.orSimilar")}
        </span>
      </div>
      <div className="mb-6">
        <CarDetailsPills vehicle={vehicle} />
      </div>

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
      {children}
    </div>
  );
}
