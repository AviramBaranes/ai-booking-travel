"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { SummaryRow } from "./SummaryRow";
import { broker, reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";

interface PayAtPickupSectionProps {
  currency: string;
  fees: broker.Fees;
  selectedAddons: reservation.PayAtPickup["selectedAddons"];
}

export function PayAtPickupSection({
  currency,
  fees,
  selectedAddons,
}: PayAtPickupSectionProps) {
  const t = useTranslations("MyAccount.summary");

  if (
    fees.dropCharge === 0 &&
    fees.youngDriverFee === 0 &&
    selectedAddons.length === 0
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

      {selectedAddons.map((addon) => (
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
