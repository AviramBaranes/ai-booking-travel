"use client";

import { booking } from "@/shared/client";
import { InclusionsDisplay } from "./InclustionsDisplay";
import { useState } from "react";
import { useSelectedVehicle } from "../_hooks/useSelectedVehicle";
import { redirect, useParams } from "next/navigation";
import { Loading } from "@/shared/components/Loading";

interface PlansPageContentProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function PlansPageContent({ searchRequest }: PlansPageContentProps) {
  const { lang } = useParams();
  const [selectedPlan, setSelectedPlan] = useState(0);
  const vehicle = useSelectedVehicle(searchRequest);

  if (!vehicle) {
    return <Loading />;
  }

  return (
    <div className="flex">
      <div className="w-3/4">
        <div className="flex gap-4">
          <div className="w-1/2">
            <InclusionsDisplay
              title="מה התוכנית כוללת?"
              inclusions={vehicle.plans[selectedPlan].planInclusions}
            />
          </div>
          <div className="w-1/2">
            <InclusionsDisplay
              title="תנאי התוכנית"
              inclusions={vehicle.plans[selectedPlan].info}
            />
          </div>
        </div>
      </div>
      <div className="w-1/4"></div>
    </div>
  );
}
