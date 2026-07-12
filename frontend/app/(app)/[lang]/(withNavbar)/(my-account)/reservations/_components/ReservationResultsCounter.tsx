"use client";

import { useTranslations } from "next-intl";
import { useReservationFilters } from "../_hooks/useReservationFilters";
import { useReservations } from "../_hooks/useReservations";

export function ReservationResultsCounter() {
  const t = useTranslations("MyAccount.reservations");
  const { sortBy, filters, page } = useReservationFilters();
  const { data, isLoading } = useReservations({
    Page: page,
    SortBy: sortBy,
    ...filters,
  });

  if (isLoading || !data) {
    return (
      <p className="text-xs text-text-secondary">
        {t("showingXResults", {
          count: "X",
          total: "X",
        })}
      </p>
    );
  }

  const {
    total,
    reservations: { length },
  } = data;

  return (
    <p className="text-xs text-text-secondary">
      {t("showingXResults", {
        count: length,
        total: total,
      })}
    </p>
  );
}
