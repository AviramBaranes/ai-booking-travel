import { useTranslations } from "next-intl";
import { SummaryRow } from "../../../../_components/SummaryRow";
import { statusToColor } from "../../../_utils/statusesStyles";
import { booking } from "@/shared/client";
import { PriceOfferActions } from "./PriceOfferActions";
import { useParams } from "next/navigation";

export function HeaderSection({
  priceOffer,
}: {
  priceOffer: booking.GetAgentPriceOfferResponse;
}) {
  const { lang } = useParams();
  const t = useTranslations("MyAccount.priceOffer.summary");

  return (
    <>
      <div className="flex items-center justify-between">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        {priceOffer.status !== "unavailable" && (
          <PriceOfferActions priceOfferId={priceOffer.id} />
        )}
      </div>
      <hr />
      <SummaryRow label={t("labels.name")} value={priceOffer.name} />
      <SummaryRow
        label={t("labels.status")}
        value={t(`status.${priceOffer.status}`)}
        valClassName={statusToColor(priceOffer.status)}
      />
      <SummaryRow
        label={t("labels.createdAt")}
        value={new Date(priceOffer.createdAt).toLocaleDateString(lang)}
      />
    </>
  );
}
