"use client";

import { useTranslations } from "next-intl";

type RentalPriceForDaysProps = {
  daysCount: number;
};

export function RentalPriceForDays({ daysCount }: RentalPriceForDaysProps) {
  const t = useTranslations("booking.shared");

  return (
    <span className="text-[14px] leading-4.5 text-border-muted">
      {t("rentalPriceForDays", { daysCount })}
    </span>
  );
}
