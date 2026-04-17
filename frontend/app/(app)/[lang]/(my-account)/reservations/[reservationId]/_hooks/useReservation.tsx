import { getReservationById } from "@/shared/api/reservations";
import { useSuspenseQuery } from "@tanstack/react-query";

export function useReservation(reservationId: number) {
  return useSuspenseQuery({
    queryKey: ["reservation", reservationId],
    queryFn: () => getReservationById(reservationId),
  });
}
