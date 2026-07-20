import { getTranslations } from "next-intl/server";
import { NewOrderButton } from "../_components/NewOrderButton";
import { PriceOfferResultsCounter } from "./_components/PriceOfferResultsCounter";
import { PriceOffersGrid } from "./_components/PriceOffersGrid";
import { ClearFilterRow } from "./_components/filters/ClearFilterRow";
import { FilterForm } from "./_components/filters/FilterForm";
import { PriceOfferPaginationButtons } from "./_components/filters/PriceOfferPaginationButtons";

export default async function PriceOffersPage() {
  const t = await getTranslations("MyAccount.priceOffers");

  return (
    <main className="lg:w-2/3 mx-5 lg:mx-auto lg:pt-15 pb-6">
      <NewOrderButton btnText={t("newSearch")} />
      <div className="flex flex-col gap-6">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <FilterForm />
        <ClearFilterRow />
        <PriceOfferResultsCounter />
        <PriceOffersGrid />
        <PriceOfferPaginationButtons />
      </div>
      <div className="mb-15" />
    </main>
  );
}
