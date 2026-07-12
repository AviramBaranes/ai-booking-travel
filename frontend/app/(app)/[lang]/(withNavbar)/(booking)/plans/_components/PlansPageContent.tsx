"use client";

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
import { useMemo, useState } from "react";
import { ErpDialog } from "./ErpDialog";
import { FeesNote } from "./FeesNote";
import { PriceOfferDialog } from "./PriceOfferDialog";
import { useSearchRequest } from "../../_hooks/useSearchRequest";
import useAuthStore from "@/shared/auth/authStore";

export function PlansPageContent() {
  const t = useTranslations("booking.plansPage");
  const { lang } = useParams();
  const router = useRouter();
  const currentSearchParams = useSearchParams();
  const user = useAuthStore((s) => s.user);
  const isAgent = user?.role === "agent";

  const selectedPlan = useBookingSessionStore((s) => s.selectedPlanIndex);
  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedAddons = useBookingSessionStore((s) => s.selectedAddons);
  const setSelectedPlan = useBookingSessionStore((s) => s.setSelectedPlanIndex);
  const setIsErpSelected = useBookingSessionStore((s) => s.setIsErpSelected);
  const setSelectedAddons = useBookingSessionStore((s) => s.setSelectedAddons);

  const { searchRequest } = useSearchRequest();

  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const [isErpDialogOpen, setIsErpDialogOpen] = useState(false);
  const [isPriceOfferDialogOpen, setIsPriceOfferDialogOpen] = useState(false);

  const { addOns, planInclusions } = useMemo(() => {
    if (!data) {
      return { planInclusions: [], addOns: [] };
    }
    const supplier = data.suppliersInfo.find(
      (s) => s.name === vehicle?.plans[selectedPlan].supplierName,
    );

    const selectedPlanName = vehicle?.plans[selectedPlan].planName;
    const selectedPlanInclusions = supplier?.inclusions.find(
      (inc) => inc.productName === selectedPlanName,
    )?.productInclusions;

    return {
      planInclusions: selectedPlanInclusions ?? [],
      addOns: supplier?.addOns ?? [],
    };
  }, [data, selectedPlan, vehicle]);

  if (!vehicle) {
    return <Loading />;
  }

  return (
    <div className="flex gap-4">
      <div className="w-3/4">
        <div className="flex gap-4 mb-6">
          {planInclusions.length > 0 && (
            <div className="w-1/2">
              <InclusionsDisplay
                title={t("inclusionsTitle")}
                inclusions={planInclusions}
              />
            </div>
          )}
          {vehicle.plans[selectedPlan].info.length > 0 && (
            <div className="w-1/2">
              <InclusionsDisplay
                title={t("rentalTerms")}
                inclusions={vehicle.plans[selectedPlan].info}
              />
            </div>
          )}
        </div>
        <FeesNote vehicle={vehicle} />
        <div className="flex justify-between items-center my-6">
          <div className="flex gap-4">
            {vehicle.plans.length > 1 && (
              <OtherPlansButton
                plans={vehicle.plans}
                suppliersInfo={data?.suppliersInfo ?? []}
                selectedPlan={selectedPlan}
                onSelectPlan={setSelectedPlan}
                currency={vehicle.priceDetails.currency}
                daysCount={data?.daysCount ?? 0}
              />
            )}
            <ImportantInfoButton
              plans={vehicle.plans}
              suppliersInfo={data?.suppliersInfo ?? []}
              selectedPlanIndex={selectedPlan}
            />
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
        {!!addOns.length && (
          <>
            <hr className="mt-10 mb-6" />
            <AddOnsDisplay
              addons={addOns}
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
              className="mt-4 type-paragraph font-bold py-6 px-8 cursor-pointer"
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
            {isAgent && (
              <Button
                variant="brand"
                className="type-paragraph font-bold py-6 px-8 cursor-pointer bg-navy"
                onClick={() => {
                  setIsPriceOfferDialogOpen(true);
                }}
              >
                {t("createPriceOffer")}
              </Button>
            )}
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
      {isAgent && (
        <PriceOfferDialog
          open={isPriceOfferDialogOpen}
          onOpenChange={setIsPriceOfferDialogOpen}
          searchRequest={searchRequest}
        />
      )}
    </div>
  );
}
