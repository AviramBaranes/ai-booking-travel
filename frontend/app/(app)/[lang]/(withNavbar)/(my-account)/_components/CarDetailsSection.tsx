"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { SummaryRow } from "./SummaryRow";
import { formatPrice } from "@/shared/utils/formatPrice";

interface CarDetailsSectionProps {
  rentalDays: number;
  carType: string;
  model: string;
  brand: string;
  isAutomatic: boolean;
  excess?: number;
  excessCurrency?: string;
}

export function CarDetailsSection({
  rentalDays,
  carType,
  model,
  brand,
  isAutomatic,
  excess,
  excessCurrency,
}: CarDetailsSectionProps) {
  const t = useTranslations("MyAccount.summary");

  return (
    <>
      <SummarySubTitle title={t("sections.carDetails")} />
      <SummaryRow
        label={t("labels.rentalDays")}
        value={rentalDays.toString()}
      />
      <SummaryRow
        label={t("labels.carType")}
        value={carType}
      />
      <SummaryRow label={t("labels.model")} value={model} />
      <SummaryRow label={t("labels.gear")} value={isAutomatic ? t("labels.auto") : t("labels.manual")} />
      <SummaryRow
        label={t("labels.brand")}
        value={brand}
      />
      {!!excess && excessCurrency && (
        <SummaryRow
          label={t("labels.excess")}
          value={formatPrice(excess, excessCurrency)}
        />
      )}
    </>
  );
}
