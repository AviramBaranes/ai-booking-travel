import { getTranslations } from "next-intl/server";
import { NewOrderButton } from "../_components/NewOrderButton";

export default async function PriceOffersPage() {
  const t = await getTranslations("MyAccount.priceOffers");
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NewOrderButton btnText={t("newSearch")} />
      <div className="flex flex-col gap-6">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
      </div>
      <div className="mb-15" />
    </main>
  );
}
