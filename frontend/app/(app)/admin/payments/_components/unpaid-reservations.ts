import type { supplier_payments } from "@/shared/client";

/**
 * UnpaidReservation is a reservation row as the API returns it: the currency code travels on the
 * row itself, so this screen can list every currency together, ordered by pickup date.
 */
export type UnpaidReservation = supplier_payments.UnpaidSupplierReservation;

/** UnpaidPenalty is a cancellation or no-show fee row, listed in the same table. */
export type UnpaidPenalty = supplier_payments.UnpaidSupplierPenalty;

/** UnpaidRow is one line of the table, which lists reservations and fees together. */
export type UnpaidRow = UnpaidReservation | UnpaidPenalty;

/** Only a fee carries a type, which is what tells the two rows apart. */
export const isPenalty = (row: UnpaidRow): row is UnpaidPenalty =>
  "type" in row;

/**
 * toUnpaidRows lists reservations and fees together. Each list arrives ordered by pickup date, so
 * they are merged on the same key to keep the lines of one booking side by side.
 */
export function toUnpaidRows(
  reservations: UnpaidReservation[],
  penalties: UnpaidPenalty[],
): UnpaidRow[] {
  return [...reservations, ...penalties].sort((a, b) =>
    a.pickupDate.localeCompare(b.pickupDate),
  );
}

/** Labels for the תשלום עבור column. */
export const PAYMENT_FOR_RESERVATION = "הזמנה";

const PENALTY_TYPE_LABELS: Record<string, string> = {
  cancellation: "ביטול מאוחר",
  no_show: "אי הגעה",
};

export const penaltyTypeLabel = (type: string) =>
  PENALTY_TYPE_LABELS[type] ?? type;
