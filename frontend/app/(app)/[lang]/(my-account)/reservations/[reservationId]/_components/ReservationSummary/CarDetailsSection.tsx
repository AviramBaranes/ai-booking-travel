import { useTranslations } from "next-intl";
import { reservation } from "@/shared/client";
import { OrderSummaryRow } from "./OrderSummaryRow";
import { OrderSummarySubTitle } from "./OrderSummarySubTitle";

export function CarDetailsSection({
  reservation: res,
}: {
  reservation: reservation.GetReservationResponse;
}) {
  const t = useTranslations("MyAccount.reservation.summary");

  return (
    <>
      <OrderSummarySubTitle title={t("sections.carDetails")} />
      <OrderSummaryRow
        label={t("labels.rentalDays")}
        value={res.rentalDays.toString()}
      />
      <OrderSummaryRow
        label={t("labels.carType")}
        value={res.carDetails.carType}
      />
      <OrderSummaryRow label={t("labels.model")} value={res.carDetails.model} />
      <OrderSummaryRow
        label={t("labels.brand")}
        value={res.carDetails.supplierName}
      />
    </>
  );
}
