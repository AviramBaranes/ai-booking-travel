import { getReservationById } from "@/shared/api/reservations-api";
import { useQuery } from "@tanstack/react-query";

export function useReservation(reservationId: number) {
  const queryKey = ["reservation", reservationId];
  const suspenseResult = useQuery({
    queryKey,
    queryFn: () => getReservationById(reservationId),
  });

  return { ...suspenseResult, queryKey };
}
