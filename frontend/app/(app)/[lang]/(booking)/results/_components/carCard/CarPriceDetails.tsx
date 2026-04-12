import { booking } from "@/shared/client";
import { useDirection } from "@/shared/hooks/useDirection";
import { formatPrice } from "@/shared/utils/formatPrice";
import { isFutureWithinHours } from "@/shared/utils/isFutureWithinHours";
import clsx from "clsx";
import Link from "next/link";
import Image from "next/image";
import { useTranslations } from "next-intl";

export const HOURS_BEFORE_PICKUP_TO_ALLOW_CANCELLATION = 48;

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
      className={clsx("absolute bottom-0 flex flex-col items-end gap-3", {
        "right-0": dir === "ltr",
        "left-0": dir === "rtl",
      })}
    >
      <p className="type-paragraph mx-4">
        {t("rentalPriceForDays", { daysCount })}
      </p>
      {firstPlan.fullPrice !== firstPlan.price && (
        <p className="type-label mx-4">
          {t("priceBeforeDiscount")}{" "}
          {formatPrice(firstPlan.fullPrice, vehicle.priceDetails.currency)}
        </p>
      )}
      <p className="type-h4 mx-4 text-navy">
        {formatPrice(firstPlan.price, vehicle.priceDetails.currency)}
      </p>

      {isFutureWithinHours(
        new Date(searchRequest.PickupDate),
        searchRequest.PickupTime,
        HOURS_BEFORE_PICKUP_TO_ALLOW_CANCELLATION,
      ) && (
        <div className="flex gap-2 items-center mx-4">
          <Image
            src="/assets/icons/V.svg"
            alt="Checked Icon"
            width={16}
            height={4}
            className="w-4"
          />
          <span className="type-label text-[#16A34A]">
            {t("freeCancellation")}
          </span>
        </div>
      )}

      <Link
        href={`/add-ons/?sid=`}
        className={clsx("bg-brand type-label p-6 text-white", {
          "rounded-tr-2xl": dir === "rtl",
          "rounded-tl-2xl": dir === "ltr",
        })}
      >
        {t("continueCta")}
      </Link>
    </div>
  );
}
