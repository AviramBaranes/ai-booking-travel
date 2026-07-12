"use client";

import { PaginationButtons } from "../../../_components/PaginationButtons";
import { useTranslations } from "next-intl";
import { useReservationFilters } from "../../_hooks/useReservationFilters";
import { useReservations } from "../../_hooks/useReservations";

export function ReservationPaginationButtons() {
  const t = useTranslations("MyAccount.reservations");
  const { lang, searchParams, sortBy, filters, page } = useReservationFilters();
  const {
    data,
    isLoading
  } = useReservations({ Page: page, SortBy: sortBy, ...filters });

  if (isLoading || !data) {
    return null
  }
  
  return (
    <PaginationButtons
      total={data.total}
      page={page}
      searchParams={searchParams}
      basePath={`/${lang}/reservations`}
      previousLabel={t("pagination.prev")}
      nextLabel={t("pagination.next")}
    />
  );
}