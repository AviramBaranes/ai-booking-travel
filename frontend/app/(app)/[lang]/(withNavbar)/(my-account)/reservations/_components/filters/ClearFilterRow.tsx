"use client";

import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { X } from "lucide-react";
import {
  ReservationFilterKey,
  useReservationFilters,
} from "../../_hooks/useReservationFilters";

export function ClearFilterRow() {
  const router = useRouter();
  const { lang, searchParams, activeFilters } = useReservationFilters();
  const t = useTranslations("MyAccount.reservations");
  const tStatus = useTranslations("MyAccount.reservation.summary.status");
  const labelByKey: Record<ReservationFilterKey, string> = {
    bid: t("bookingIdPlaceholder"),
    status: t("statusPlaceholder"),
    pd: t("pickupDatePlaceholder"),
    dn: t("driverNamePlaceholder"),
  };

  const locale = lang === "he" ? "he-IL" : "en-US";

  const formatPickupDate = (value: string) => {
    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) {
      return value;
    }

    return new Intl.DateTimeFormat(locale, {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    }).format(parsed);
  };

  const clearFilter = (key: ReservationFilterKey) => {
    const nextQuery = new URLSearchParams(searchParams.toString());
    nextQuery.delete(key);

    const queryString = nextQuery.toString();
    const basePath = `/${lang}/reservations`;
    router.push(queryString ? `${basePath}?${queryString}` : basePath);
  };

  if (activeFilters.length === 0) {
    return null;
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {activeFilters.map((filter) => {
        const value =
          filter.key === "status"
            ? tStatus(filter.value)
            : filter.key === "pd"
              ? formatPickupDate(filter.value)
              : filter.value;

        return (
          <Button
            key={filter.key}
            variant="outline"
            type="button"
            onClick={() => clearFilter(filter.key)}
            className="type-paragraph font-normal flex items-center gap-1"
          >
            <X className="size-4" />

            <span className="font-semibold">
              {labelByKey[filter.key]}
              {": "}
            </span>
            {value}
          </Button>
        );
      })}
    </div>
  );
}
