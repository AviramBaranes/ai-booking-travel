"use client";

import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { booking } from "@/shared/client";

interface CarResultsProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function CarResults({ searchRequest }: CarResultsProps) {
  const { data, status } = useAvailableCars(searchRequest);

  if (status === "pending") return null;
  if (status === "error") return null;

  return (
    <div>
      {data.availableVehicles.map((vehicle, i) => (
        <div key={i}>{vehicle.carDetails.model}</div>
      ))}
    </div>
  );
}
