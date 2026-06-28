"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { SummaryRow } from "./SummaryRow";
import { broker, reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";
import { AddonsGallery } from "@/payload-types";
import { useMemo } from "react";
import { useParams } from "next/navigation";

interface PayAtPickupSectionProps {
  currency: string;
  payAtPickup: reservation.PayAtPickup;
  addonsGallery: AddonsGallery;
}

export function PayAtPickupSection({
  currency,
  payAtPickup,
  addonsGallery,
}: PayAtPickupSectionProps) {
  const { lang } = useParams();
  const { fees, deposit, depositCurrency } = payAtPickup;
  const selectedAddons = useMemo(
    () =>
      payAtPickup.selectedAddons?.map((addon) => {
        const addonGalleryItem = addonsGallery.addons?.find(
          (item) => item.addonId === addon.id.toString(),
        );
        return {
          ...addon,
          name: addonGalleryItem
            ? lang === "he"
              ? addonGalleryItem.hebrewName
              : addonGalleryItem.englishName
            : addon.name,
        };
      }),
    [addonsGallery, lang, payAtPickup.selectedAddons],
  );
  const t = useTranslations("MyAccount.summary");

  if (
    fees.dropCharge === 0 &&
    fees.youngDriverFee === 0 &&
    deposit === 0 &&
    (!selectedAddons || selectedAddons.length === 0)
  ) {
    return null;
  }

  return (
    <>
      <SummarySubTitle title={t("payAtPickup.title")} />
      {fees.dropCharge > 0 && (
        <SummaryRow
          label={t("payAtPickup.dropCharge")}
          value={formatPrice(fees.dropCharge, fees.dropChargeCurrency)}
        />
      )}
      {fees.youngDriverFee > 0 && (
        <SummaryRow
          label={t("payAtPickup.youngDriverFee")}
          value={formatPrice(fees.youngDriverFee, fees.youngDriverFeeCurrency)}
        />
      )}
      {deposit > 0 && (
        <SummaryRow
          label={t("payAtPickup.deposit")}
          value={formatPrice(deposit, depositCurrency)}
        />
      )}

      {selectedAddons?.map((addon) => (
        <div key={addon.id}>
          <SummaryRow
            label={addon.name}
            value={formatPrice(addon.price, currency)}
          />
          <p>{t("payAtPickup.quantity", { quantity: addon.quantity })}</p>
        </div>
      ))}
    </>
  );
}
