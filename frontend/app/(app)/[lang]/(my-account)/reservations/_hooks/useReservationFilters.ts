"use client";

import { useParams, useSearchParams } from "next/navigation";

export const SORT_OPTIONS = ["created_at", "pickup_date"] as const;

export const RESERVATION_STATUSES = [
  "booked",
  "vouchered",
  "canceled",
] as const;

export type ReservationStatus = (typeof RESERVATION_STATUSES)[number];
export type ReservationFilterKey = "status" | "pd" | "bid" | "dn";

export type ReservationFilters = {
  status: ReservationStatus | null;
  pickupDate: Date | undefined;
  bookingId: string;
  driverName: string;
};

export type ActiveReservationFilter = {
  key: ReservationFilterKey;
  value: string;
};

function isReservationStatus(value: string): value is ReservationStatus {
  return RESERVATION_STATUSES.includes(value as ReservationStatus);
}

function parsePickupDate(value: string | null): Date | undefined {
  if (!value) {
    return undefined;
  }

  const parsedDate = new Date(value);
  if (Number.isNaN(parsedDate.getTime())) {
    return undefined;
  }

  return parsedDate;
}

function toLocalDateString(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

export function buildReservationFiltersQuery(filters: ReservationFilters) {
  const query = new URLSearchParams();

  if (filters.status) query.set("status", filters.status);
  if (filters.pickupDate)
    query.set("pd", toLocalDateString(filters.pickupDate));
  if (filters.bookingId) query.set("bid", filters.bookingId);
  if (filters.driverName) query.set("dn", filters.driverName);

  return query;
}

export function useReservationFilters() {
  const searchParams = useSearchParams();
  const { lang } = useParams();

  const statusParam = searchParams.get("status");
  const status =
    statusParam && isReservationStatus(statusParam) ? statusParam : null;

  const bookingId = searchParams.get("bid") ?? "";
  const driverName = searchParams.get("dn") ?? "";
  const pickupDateRaw = searchParams.get("pd");
  const pickupDate = parsePickupDate(pickupDateRaw);

  const filters: ReservationFilters = {
    status,
    pickupDate,
    bookingId,
    driverName,
  };

  const activeFilters: ActiveReservationFilter[] = [];
  if (bookingId) {
    activeFilters.push({ key: "bid", value: bookingId });
  }
  if (status) {
    activeFilters.push({ key: "status", value: status });
  }
  if (pickupDateRaw) {
    activeFilters.push({ key: "pd", value: pickupDateRaw });
  }
  if (driverName) {
    activeFilters.push({ key: "dn", value: driverName });
  }

  const pageParam = searchParams.get("page");
  const page = pageParam ? Math.max(1, parseInt(pageParam, 10) || 1) : 1;

  return {
    sortBy: searchParams.get("sortBy") || SORT_OPTIONS[0],
    lang,
    searchParams,
    filters,
    activeFilters,
    page,
  };
}
