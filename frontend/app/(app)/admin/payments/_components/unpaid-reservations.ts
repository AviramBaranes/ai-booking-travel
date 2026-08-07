import type { supplier_payments } from "@/shared/client";

/**
 * UnpaidReservation is a reservation row with its currency folded back in. The API groups rows
 * by currency; this screen lists every currency together, so the code travels on the row.
 */
export type UnpaidReservation = supplier_payments.UnpaidSupplierReservation & {
  currencyCode: string;
};

/** flattenCurrencyGroups merges the per-currency groups into one list ordered by pickup date. */
export function flattenCurrencyGroups(
  groups: supplier_payments.UnpaidSupplierCurrencyGroup[],
): UnpaidReservation[] {
  return groups
    .flatMap((group) =>
      group.reservations.map((reservation) => ({
        ...reservation,
        currencyCode: group.currencyCode,
      })),
    )
    .sort((a, b) => a.pickupDate.localeCompare(b.pickupDate));
}
