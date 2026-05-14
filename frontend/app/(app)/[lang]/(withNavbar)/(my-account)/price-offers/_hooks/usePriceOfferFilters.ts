"use client";

import { useParams, useSearchParams } from "next/navigation";

export const PRICE_OFFER_STATUSES = ["open", "booked", "declined"] as const;

export type PriceOfferStatus = (typeof PRICE_OFFER_STATUSES)[number];
export type PriceOfferFilterKey = "name" | "status";

export type PriceOfferFilters = {
  name: string;
  status: PriceOfferStatus | null;
};

export type ActivePriceOfferFilter = {
  key: PriceOfferFilterKey;
  value: string;
};

function isPriceOfferStatus(value: string): value is PriceOfferStatus {
  return PRICE_OFFER_STATUSES.includes(value as PriceOfferStatus);
}

export function buildPriceOfferFiltersQuery(filters: PriceOfferFilters) {
  const query = new URLSearchParams();

  if (filters.name) query.set("name", filters.name);
  if (filters.status) query.set("status", filters.status);

  return query;
}

export function usePriceOfferFilters() {
  const searchParams = useSearchParams();
  const { lang } = useParams();

  const name = searchParams.get("name") ?? "";
  const statusParam = searchParams.get("status");
  const status =
    statusParam && isPriceOfferStatus(statusParam) ? statusParam : null;

  const filters: PriceOfferFilters = {
    name,
    status,
  };

  const activeFilters: ActivePriceOfferFilter[] = [];
  if (name) {
    activeFilters.push({ key: "name", value: name });
  }
  if (status) {
    activeFilters.push({ key: "status", value: status });
  }

  const pageParam = searchParams.get("page");
  const page = pageParam ? Math.max(1, parseInt(pageParam, 10) || 1) : 1;

  return {
    lang,
    searchParams,
    filters,
    activeFilters,
    page,
  };
}