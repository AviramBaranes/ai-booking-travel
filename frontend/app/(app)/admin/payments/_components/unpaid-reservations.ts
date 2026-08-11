import type { supplier_payments } from "@/shared/client";

/**
 * UnpaidReservation is a reservation row as the API returns it: the currency code travels on the
 * row itself, so this screen can list every currency together, ordered by pickup date.
 */
export type UnpaidReservation = supplier_payments.UnpaidSupplierReservation;
