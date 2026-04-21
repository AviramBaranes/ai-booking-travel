import { listReservations } from "@/shared/api/reservations";
import { reservation } from "@/shared/client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ReservationFilters } from "./useReservationFilters";

interface UseReservationsParams extends ReservationFilters {
  Page: number;
  SortBy: string;
}

export function useReservations(params: UseReservationsParams) {
  const queryKey = ["reservations", params];
  const suspenseResult = useSuspenseQuery({
    queryKey,
    queryFn: () =>
      listReservations({
        Page: params.Page,
        SortBy: params.SortBy,
        Status: params.status ?? undefined,
        PickupDate: params.pickupDate
          ? `${params.pickupDate.getFullYear()}-${String(params.pickupDate.getMonth() + 1).padStart(2, "0")}-${String(params.pickupDate.getDate()).padStart(2, "0")}`
          : undefined,
        BookingID: params.bookingId || undefined,
        Name: params.driverName || undefined,
      }),
  });

  return { ...suspenseResult, queryKey };
}
