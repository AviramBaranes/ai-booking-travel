import { price_offer } from "@/shared/client";
import { AccountCardLabelValue } from "../../_components/AccountCardLabelValue";
import { formatPrice } from "@/shared/utils/formatPrice";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useParams } from "next/navigation";
import { statusToBg, statusToColor } from "../_utils/statusesStyles";

export function PriceOfferCard({
  priceOffer,
}: {
  priceOffer: price_offer.PriceOfferSummary;
}) {
  const { lang } = useParams();
  const tLabels = useTranslations("MyAccount.priceOffers.labels");
  const tStatuses = useTranslations("MyAccount.priceOffer.summary.status");

  return (
    <Link
      href={`/${lang}/price-offers/${priceOffer.id}`}
      className="p-6 flex flex-col gap-4 rounded-xl bg-white shadow-card hover:shadow-card-hover hover:border hover:border-brand"
    >
      <AccountCardLabelValue
        valClassName="font-semibold"
        label={tLabels("name")}
        value={priceOffer.name}
      />
      <AccountCardLabelValue
        label={tLabels("pickupDate")}
        value={new Date(priceOffer.pickupDate).toLocaleDateString(lang)}
      />
      <AccountCardLabelValue
        label={tLabels("pickupLocation")}
        value={priceOffer.pickupLocationName}
      />
      <AccountCardLabelValue
        label={tLabels("totalPrice")}
        value={formatPrice(priceOffer.totalPrice, priceOffer.currencyCode)}
      />
      <AccountCardLabelValue
        valClassName="font-semibold"
        label={tLabels("offeredPrice")}
        value={formatPrice(
          priceOffer.offeredPrice,
          priceOffer.offeredCurrencyCode,
        )}
      />
      <div className="px-6 py-1 flex flex-col">
        <p className="text-xs text-muted">{tLabels("status")}</p>
        <p
          className={`rounded-md py-1 mt-2 px-2 w-fit text-sm ${statusToBg(priceOffer.status)} ${statusToColor(priceOffer.status)}`}
        >
          {tStatuses(priceOffer.status)}
        </p>
      </div>
    </Link>
  );
}
