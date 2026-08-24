"use client";

import { useTranslations } from "next-intl";

type RentalPriceForDaysProps = {
  daysCount: number;
  isErp?: boolean;
};

export function RentalPriceForDays({
  daysCount,
  isErp,
}: RentalPriceForDaysProps) {
  const t = useTranslations("booking.shared");

  return (
    <span className="text-[14px] leading-4.5 text-border-muted">
      {t(isErp ? "erpPriceForDays" : "rentalPriceForDays", { daysCount })}
    </span>
  );
}
