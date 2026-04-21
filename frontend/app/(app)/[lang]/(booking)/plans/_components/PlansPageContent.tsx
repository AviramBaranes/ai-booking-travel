"use client";

import { booking } from "@/shared/client";
import { InclusionsDisplay } from "./InclustionsDisplay";
import { useSelectedVehicle } from "../_hooks/useSelectedVehicle";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useParams } from "next/navigation";
import { Loading } from "@/shared/components/Loading";
import { OtherPlansButton } from "./OtherPlansButton";
import { ImportantInfoButton } from "./ImportantInfoButton";
import { SignalsDisplay } from "../../_components/SignalsDisplay";
import { ErpCheckbox } from "./ErpCheckbox";
import { AddOnsDisplay } from "./AddOnsDisplay";
import { SelectedCarCard } from "@/shared/components/booking/SelectedCarCard/SelectedCarCard";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { useSearchParams, useRouter } from "next/navigation";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { useState } from "react";
import { ErpDialog } from "./ErpDialog";
import { FeesNote } from "./FeesNote";

interface PlansPageContentProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function PlansPageContent({ searchRequest }: PlansPageContentProps) {
  const t = useTranslations("booking.plansPage");
  const { lang } = useParams();
  const router = useRouter();
  const currentSearchParams = useSearchParams();

  const selectedPlan = useBookingSessionStore((s) => s.selectedPlanIndex);
  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedAddons = useBookingSessionStore((s) => s.selectedAddons);
  const setSelectedPlan = useBookingSessionStore((s) => s.setSelectedPlanIndex);
  const setIsErpSelected = useBookingSessionStore((s) => s.setIsErpSelected);
  const setSelectedAddons = useBookingSessionStore((s) => s.setSelectedAddons);

  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const [isErpDialogOpen, setIsErpDialogOpen] = useState(false);

  if (!vehicle) {
    return <Loading />;
  }

  return (
    <div className="flex gap-4">
      <div className="w-3/4">
        <div className="flex gap-4 mb-6">
          <div className="w-1/2">
            <InclusionsDisplay
              title={t("inclusionsTitle")}
              inclusions={vehicle.plans[selectedPlan].planInclusions}
            />
          </div>
          <div className="w-1/2">
            <InclusionsDisplay
              title={t("rentalTerms")}
              inclusions={vehicle.plans[selectedPlan].info}
            />
          </div>
        </div>
        <FeesNote vehicle={vehicle} />
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
              selectedAddons={selectedAddons}
              setSelectedAddons={setSelectedAddons}
            />
          </>
        )}
      </div>
      <div className="w-1/4">
        <SelectedCarCard
          isErpSelected={isErpSelected}
          daysCount={data?.daysCount ?? 0}
          vehicle={vehicle}
          selectedPlanIndex={selectedPlan}
        >
          <>
            <FreeCancellationBadge
              pickupDate={searchRequest.PickupDate}
              pickupTime={searchRequest.PickupTime}
              text={t("freeCancellation")}
            />
            <Button
              variant="brand"
              className="mt-4 mx-auto type-paragraph font-bold py-6 px-8 cursor-pointer"
              onClick={() => {
                if (isErpSelected) {
                  router.push(
                    `/${lang}/order?${currentSearchParams.toString()}`,
                  );
                } else {
                  setIsErpDialogOpen(true);
                }
              }}
            >
              {t("continueCta")}
            </Button>
          </>
        </SelectedCarCard>
      </div>
      <ErpDialog
        open={isErpDialogOpen}
        onApprove={() => {
          setIsErpSelected(true);
          router.push(`/${lang}/order?${currentSearchParams.toString()}`);
        }}
        onDecline={() => {
          router.push(`/${lang}/order?${currentSearchParams.toString()}`);
        }}
        erpPrice={vehicle.plans[selectedPlan].erpPrice}
        erpPriceCurrency={vehicle.priceDetails.currency}
      />
    </div>
  );
}
