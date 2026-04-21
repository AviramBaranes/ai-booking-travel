import { getReservationById } from "@/shared/api/reservations";
import { useSuspenseQuery } from "@tanstack/react-query";

export function useReservation(reservationId: number) {
  const queryKey = ["reservation", reservationId];
  const suspenseResult = useSuspenseQuery({
    queryKey,
    queryFn: () => getReservationById(reservationId),
  });

  return { ...suspenseResult, queryKey };
}
