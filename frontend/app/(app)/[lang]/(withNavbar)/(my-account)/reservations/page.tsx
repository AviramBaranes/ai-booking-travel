import { getTranslations } from "next-intl/server";
import { NewOrderButton } from "../_components/NewOrderButton";
import { FilterForm } from "./_components/filters/FilterForm";
import { ClearFilterRow } from "./_components/filters/ClearFilterRow";
import { SortDropdown } from "./_components/filters/SortDropdown";
import { Suspense } from "react";
import { ReservationResultsCounter } from "./_components/ReservationResultsCounter";
import { ReservationsGrid } from "./_components/ReservationsGrid";
import { ReservationPaginationButtons } from "./_components/filters/ReservationPaginationButtons";

export default async function ReservationDetailsPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string>>;
}) {
  const t = await getTranslations("MyAccount.reservations");
  const resolvedParams = await searchParams;
  const suspenseKey = new URLSearchParams(resolvedParams).toString();
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NewOrderButton btnText={t("newOrder")} />
      <div className="flex flex-col gap-6">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <FilterForm />
        <ClearFilterRow />
        <div className="flex items-center gap-4">
          <SortDropdown />
          <ReservationResultsCounter />
        </div>
        <ReservationsGrid />
        <ReservationPaginationButtons />
      </div>
      <div className="mb-15" />
    </main>
  );
}
