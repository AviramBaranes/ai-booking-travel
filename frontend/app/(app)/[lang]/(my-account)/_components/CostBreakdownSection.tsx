import { useTranslations } from "next-intl";
import { formatPrice } from "@/shared/utils/formatPrice";
import { SummarySubTitle } from "./SummarySubTitle";
import { SummaryRow } from "./SummaryRow";

interface CostBreakdownSectionProps {
  priceBefDesc: number;
  discountAmount: number;
  erpPrice: number;
  totalPrice: number;
  currencyCode: string;
}

export function CostBreakdownSection({
  priceBefDesc,
  discountAmount,
  erpPrice,
  totalPrice,
  currencyCode,
}: CostBreakdownSectionProps) {
  const t = useTranslations("MyAccount.summary");

  return (
    <>
      <SummarySubTitle title={t("sections.costBreakdown")} />
      <SummaryRow
        label={t("labels.rentalPrice")}
        value={formatPrice(priceBefDesc, currencyCode)}
      />
      {discountAmount > 0 && (
        <SummaryRow
          label={t("labels.couponDiscount")}
          value={formatPrice(discountAmount, currencyCode)}
        />
      )}
      {erpPrice > 0 && (
        <SummaryRow
          label={t("labels.fullCoverage")}
          value={formatPrice(erpPrice, currencyCode)}
        />
      )}
      <div className="text-white bg-navy py-3 5 px-5 flex justify-between items-center rounded-xl mt-8">
        <span className="type-paragraph">{t("labels.totalToPay")}</span>
        <h4 className="type-h4">{formatPrice(totalPrice, currencyCode)}</h4>
      </div>
    </>
  );
}
