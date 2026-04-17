import { withErrorHandler } from "./_api";

export function getReservationById(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.GetReservation(reservationId),
  );
}
