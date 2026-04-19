import { withErrorHandler } from "./_api";

export function getReservationById(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.GetReservation(reservationId),
  );
}

export function applyVoucher(reservationId: number, voucherCode: string) {
  return withErrorHandler((client) =>
    client.reservation.ApplyVoucher(reservationId, {
      voucher: voucherCode,
    }),
  );
}

export function cancelReservation(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.CancelReservation(reservationId),
  );
}
