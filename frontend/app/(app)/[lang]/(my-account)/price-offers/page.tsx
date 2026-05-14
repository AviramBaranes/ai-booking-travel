import { getTranslations } from "next-intl/server";
import { Suspense } from "react";
import { NewOrderButton } from "../_components/NewOrderButton";
import { PriceOfferResultsCounter } from "./_components/PriceOfferResultsCounter";
import { ClearFilterRow } from "./_components/filters/ClearFilterRow";
import { FilterForm } from "./_components/filters/FilterForm";

export default async function PriceOffersPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string>>;
}) {
  const t = await getTranslations("MyAccount.priceOffers");
  const resolvedParams = await searchParams;
  const suspenseKey = new URLSearchParams(resolvedParams).toString();

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NewOrderButton btnText={t("newSearch")} />
      <div className="flex flex-col gap-6">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <FilterForm />
        <ClearFilterRow />
        <Suspense
          key={`counter-${suspenseKey}`}
          fallback={
            <p className="text-xs text-text-secondary">
              {t("showingXResults", {
                count: "X",
                total: "X",
              })}
            </p>
          }
        >
          <PriceOfferResultsCounter />
        </Suspense>
      </div>
      <div className="mb-15" />
    </main>
  );
}
