import { useTranslations } from "next-intl";
import { reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";
import { OrderSummaryRow } from "./OrderSummaryRow";
import { OrderSummarySubTitle } from "./OrderSummarySubTitle";

export function CostBreakdownSection({
  reservation: res,
}: {
  reservation: reservation.GetReservationResponse;
}) {
  const t = useTranslations("MyAccount.reservation.summary");

  return (
    <>
      <OrderSummarySubTitle title={t("sections.costBreakdown")} />
      <OrderSummaryRow
        label={t("labels.rentalPrice")}
        value={formatPrice(res.priceBefDesc, res.currencyCode)}
      />
      <OrderSummaryRow
        label={t("labels.couponDiscount")}
        value={formatPrice(res.discountAmount, res.currencyCode)}
      />
      {res.erpPrice > 0 && (
        <OrderSummaryRow
          label={t("labels.fullCoverage")}
          value={formatPrice(res.erpPrice, res.currencyCode)}
        />
      )}
      <div className="text-white bg-navy py-3 5 px-5 flex justify-between items-center rounded-xl mt-8">
        <span className="type-paragraph">{t("labels.totalToPay")}</span>
        <h4 className="type-h4">
          {formatPrice(res.totalPrice, res.currencyCode)}
        </h4>
      </div>
    </>
  );
}
