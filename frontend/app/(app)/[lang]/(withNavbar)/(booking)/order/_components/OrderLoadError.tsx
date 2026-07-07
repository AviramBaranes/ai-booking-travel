"use client";

import { useTranslations } from "next-intl";

export function OrderLoadError() {
  const t = useTranslations("booking.orderPage");

  return (
    <div className="p-20 text-center">
      <h4 className="type-h4 text-navy">{t("loadError")}</h4>
    </div>
  );
}
