"use client";

import { booking, broker } from "@/shared/client";
import { InclusionsDisplay } from "./InclustionsDisplay";
import { useState } from "react";
import { useSelectedVehicle } from "../_hooks/useSelectedVehicle";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { redirect, useParams } from "next/navigation";
import { Loading } from "@/shared/components/Loading";
import { OtherPlansButton } from "./OtherPlansButton";
import { ImportantInfoButton } from "./ImportantInfoButton";
import { SignalsDisplay } from "../../_components/SignalsDisplay";
import { ErpCheckbox } from "./ErpCheckbox";
import { AddOnsDisplay } from "./AddOnsDisplay";
import { AddonsGallery } from "@/payload-types";

interface PlansPageContentProps {
  addonsGallery: AddonsGallery;
  searchRequest: booking.SearchAvailabilityRequest;
}

export function PlansPageContent({
  addonsGallery,
  searchRequest,
}: PlansPageContentProps) {
  const { lang } = useParams();
  const [selectedPlan, setSelectedPlan] = useState(0);
  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest);
  const [isErpSelected, setIsErpSelected] = useState(false);
  const [selectedAddons, setSelectedAddons] = useState<broker.SelectAddOn[]>(
    [],
  );

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
        <div className="flex justify-between items-center my-6">
          <div className="flex gap-4">
            {vehicle.plans.length > 1 && (
              <OtherPlansButton
                plans={vehicle.plans}
                selectedPlan={selectedPlan}
                onSelectPlan={setSelectedPlan}
                currency={vehicle.priceDetails.currency}
                daysCount={data?.daysCount ?? 0}
              />
            )}
            <ImportantInfoButton />
          </div>
          {vehicle.signals && (
            <div className="flex items-center gap-2">
              <SignalsDisplay
                remainingCount={vehicle.signals.remainingCount}
                liveViewers={vehicle.signals.liveViewers}
              />
            </div>
          )}
        </div>
        <hr />
        <ErpCheckbox
          isSelected={isErpSelected}
          setSelected={setIsErpSelected}
          vehicle={vehicle}
          selectedPlan={selectedPlan}
          daysCount={data?.daysCount ?? 0}
        />
        {!!vehicle.addOns?.length && (
          <>
            <hr className="mt-10 mb-6" />
            <AddOnsDisplay
              addons={vehicle.addOns}
              addOnsGallery={addonsGallery}
              selectedAddons={selectedAddons}
              setSelectedAddons={setSelectedAddons}
            />
          </>
        )}
      </div>
      <div className="w-1/4"></div>
    </div>
  );
}
