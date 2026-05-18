import { price_offer } from "@/shared/client";
import { getTranslations } from "next-intl/server";
import { SummaryRow } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/_components/SummaryRow";
import { statusToColor } from "@/app/(app)/[lang]/(withNavbar)/(my-account)/price-offers/_utils/statusesStyles";

export async function HeaderSection({
  offer,
  lang,
}: {
  offer: price_offer.GetPriceOfferResponse;
  lang: string;
}) {
  const t = await getTranslations("MyAccount.priceOffer.summary");

  return (
    <>
      <div className="flex items-center justify-between">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
      </div>
      <hr />
      <SummaryRow label={t("labels.name")} value={offer.name} />
      <SummaryRow
        label={t("labels.status")}
        value={t(`status.${offer.status}`)}
        valClassName={statusToColor(offer.status)}
      />
      <SummaryRow
        label={t("labels.createdAt")}
        value={new Date(offer.createdAt).toLocaleDateString(lang)}
      />
    </>
  );
}
