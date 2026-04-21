import { getTranslations } from "next-intl/server";
import { NewOrderButton } from "./_components/NewOrderButton";
import { FilterForm } from "./_components/filters/FilterForm";
import { ClearFilterRow } from "./_components/filters/ClearFilterRow";
import { SortDropdown } from "./_components/filters/SortDropdown";
import { Suspense } from "react";
import { ReservationResultsCounter } from "./_components/ReservationResultsCounter";
import { ReservationsGrid } from "./_components/ReservationsGrid";
import { PaginationButtons } from "./_components/filters/PaginationButtons";

export default async function ReservationDetailsPage() {
  const t = await getTranslations("MyAccount.reservations");
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NewOrderButton />
      <div className="flex flex-col gap-6">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <FilterForm />
        <ClearFilterRow />
        <div className="flex items-center gap-4">
          <SortDropdown />
          <Suspense
            fallback={
              <p className="text-xs text-text-secondary">
                {t("showingXResults", {
                  count: "X",
                  total: "X",
                })}
              </p>
            }
          >
            <ReservationResultsCounter />
          </Suspense>
        </div>
        <ReservationsGrid />
        <Suspense>
          <PaginationButtons />
        </Suspense>
      </div>
      <div className="mb-15" />
    </main>
  );
}
