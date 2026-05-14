"use client";

import { useTranslations } from "next-intl";
import { useReservationFilters } from "../_hooks/useReservationFilters";
import { useReservations } from "../_hooks/useReservations";
import { ReservationCard } from "./ReservationCard";

export function ReservationsGrid() {
  const t = useTranslations("MyAccount.reservations");
  const { sortBy, filters, page } = useReservationFilters();
  const {
    data: { total, reservations },
    refetch,
  } = useReservations({
    Page: page,
    SortBy: sortBy,
    ...filters,
  });

  if (total === 0 || reservations.length === 0) {
    return (
      <div className="p-10 pb-60 mx-auto text-center">
        <h4 className="type-h4 text-navy">{t("noReservationsFound")}</h4>
      </div>
    );
  }
  return (
    <div className="grid grid-cols-4 gap-6">
      {reservations.map((reservation) => (
        <ReservationCard
          refetchReservations={refetch}
          reservation={reservation}
          key={reservation.id}
        />
      ))}
    </div>
  );
}
