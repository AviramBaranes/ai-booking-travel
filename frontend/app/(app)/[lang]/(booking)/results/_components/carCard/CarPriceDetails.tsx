import { booking } from "@/shared/client";
import { useDirection } from "@/shared/hooks/useDirection";
import { formatPrice } from "@/shared/utils/formatPrice";
import clsx from "clsx";
import Image from "next/image";
import { useTranslations } from "next-intl";
import { RentalPriceForDays } from "../../../../../../../shared/components/booking/RentalPriceForDays";
import { ContinueToPlansLink } from "./ContinueToPlansLink";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";

export function CarPriceDetails({
  vehicle,
  searchRequest,
  daysCount,
}: {
  daysCount: number;
  vehicle: booking.AvailableVehicle;
  searchRequest: booking.SearchAvailabilityRequest;
}) {
  const t = useTranslations("booking.results.carDetails");
  const dir = useDirection();
  const firstPlan = vehicle.plans[0];

  return (
    <div
      className={clsx("absolute bottom-0 flex flex-col items-end gap-2", {
        "right-0": dir === "ltr",
        "left-0": dir === "rtl",
      })}
    >
      {firstPlan.fullPrice !== firstPlan.price && (
        <>
          <p className="type-label mx-4 flex items-center gap-2 text-navy">
            <Image
              src="/assets/icons/coins.svg"
              alt="Coins Icon"
              width={24}
              height={24}
              className="w-6 h-6"
            />
            {t("priceBeforeDiscount")}{" "}
            {formatPrice(firstPlan.fullPrice, vehicle.priceDetails.currency)}
          </p>
          <p className="type-label mx-4 bg-success rounded-md px-2 py-1 text-white flex items-center gap-2">
            <Image
              src="/assets/icons/discount.svg"
              alt="Discount Icon"
              width={24}
              height={24}
              className="w-6 h-6"
            />
            {t("savings")}{" "}
            {formatPrice(
              firstPlan.fullPrice - firstPlan.price,
              vehicle.priceDetails.currency,
            )}
          </p>
        </>
      )}
      <p className="type-h4 mx-4 text-navy flex items-center gap-2">
        <span className="type-paragraph text-text-secondary">{t("sum")}</span>
        {formatPrice(firstPlan.price, vehicle.priceDetails.currency)}
      </p>

      <div className="mx-4">
        <RentalPriceForDays daysCount={daysCount} />
      </div>
      <FreeCancellationBadge
        pickupDate={searchRequest.PickupDate}
        pickupTime={searchRequest.PickupTime}
        text={t("freeCancellation")}
      />
      <ContinueToPlansLink
        carIndex={vehicle.id}
        className={clsx("bg-brand type-label p-6 text-white", {
          "rounded-tr-2xl": dir === "rtl",
          "rounded-tl-2xl": dir === "ltr",
        })}
      >
        {t("continueCta")}
      </ContinueToPlansLink>
    </div>
  );
}
