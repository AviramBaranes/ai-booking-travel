import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function listCoupons() {
  return withErrorHandler((client) => client.booking.ListCoupons());
}

export function createCoupon(data: booking.CreateCouponRequest) {
  return withErrorHandler((client) => client.booking.CreateCoupon(data));
}

export function updateCoupon(id: number, data: booking.UpdateCouponRequest) {
  return withErrorHandler((client) => client.booking.UpdateCoupon(id, data));
}

export function deleteCoupon(id: number) {
  return withErrorHandler((client) => client.booking.DeleteCoupon(id));
}
