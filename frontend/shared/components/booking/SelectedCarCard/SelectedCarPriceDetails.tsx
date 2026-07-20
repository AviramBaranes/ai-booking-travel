import { availability } from "@/shared/client";
import { PriceDetailRow } from "../PriceDetailRow";
import { useTranslations } from "next-intl";
import { RentalPriceForDays } from "../RentalPriceForDays";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useDirection } from "@/shared/hooks/useDirection";
import clsx from "clsx";

export function SelectedCarPriceDetails({
  vehicle,
  selectedPlanIndex,
  isErpSelected,
  daysCount,
}: {
  vehicle: availability.AvailableVehicle;
  selectedPlanIndex: number;
  isErpSelected: boolean;
  daysCount: number;
}) {
  const t = useTranslations("booking.results");
  const dir = useDirection();

  const { erpFullPrice, erpPrice, fullPrice, price } =
    vehicle.plans[selectedPlanIndex];

  const discountAmount = isErpSelected
    ? fullPrice + erpFullPrice - (price + erpPrice)
    : fullPrice - price;

  const hasDiscount = fullPrice !== price;

  return (
    <div
      className={clsx(
        "flex flex-col gap-2 absolute lg:static bottom-0 mb-18 mx-5 lg:m-0",
        {
          "left-0": dir === "rtl",
          "right-0": dir === "ltr",
        },
      )}
    >
      {hasDiscount && (
        <PriceDetailRow
          altText="coins icon"
          iconSrc="/assets/icons/coins.svg"
          label={t("carDetails.priceBeforeDiscount")}
          price={fullPrice}
          currency={vehicle.priceDetails.currency}
        />
      )}

      {isErpSelected && (
        <PriceDetailRow
          altText="stamp icon"
          iconSrc="/assets/icons/stamp.gif"
          label={t("carDetails.coveragePackage")}
          price={erpFullPrice}
          currency={vehicle.priceDetails.currency}
        />
      )}

      {hasDiscount && (
        <PriceDetailRow
          altText="discount icon"
          iconSrc="/assets/icons/Discount-Green.svg"
          label={t("carDetails.savings")}
          price={discountAmount}
          currency={vehicle.priceDetails.currency}
        />
      )}

      <hr className="mb-6 mt-3 hidden lg:block" />
      <div className="flex justify-between items-start">
        <div>
          <p className="type-label text-brand">{t("carDetails.totalToPay")}</p>
          <RentalPriceForDays daysCount={daysCount} />
        </div>
        <h5 className="type-h5 text-navy">
          {formatPrice(
            price + (isErpSelected ? erpPrice : 0),
            vehicle.priceDetails.currency,
          )}
        </h5>
      </div>
    </div>
  );
}
