import { booking } from "@/shared/client";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useSearchParams } from "next/navigation";

export function useSelectedVehicle(params: booking.SearchAvailabilityRequest) {
  const searchParams = useSearchParams();
  const { data } = useAvailableCars(params, { fromCache: true });
  const cid = searchParams.get("cid");
  const cidNumber = cid ? Number(cid) : NaN;

  const selectedVehicle = data?.availableVehicles.find(
    (v) => v.id === cidNumber,
  );

  return selectedVehicle;
}
