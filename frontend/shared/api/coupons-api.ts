import { withErrorHandler } from "./_api";

export function listCoupons() {
  return withErrorHandler((client) => client.booking.ListCoupons());
}
