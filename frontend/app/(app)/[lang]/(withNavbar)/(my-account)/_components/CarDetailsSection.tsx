"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { SummaryRow } from "./SummaryRow";

interface CarDetailsSectionProps {
  rentalDays: number;
  carType: string;
  model: string;
  brand: string;
  isAutomatic: boolean;
}

export function CarDetailsSection({
  rentalDays,
  carType,
  model,
  brand,
  isAutomatic,
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
    </>
  );
}
