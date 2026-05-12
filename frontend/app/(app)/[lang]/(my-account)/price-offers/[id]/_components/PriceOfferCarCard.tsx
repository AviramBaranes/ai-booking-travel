"use client";

import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { useTranslations } from "next-intl";
import { usePriceOffer } from "../_hooks/usePriceOffer";

export function PriceOfferCarCard({ priceOfferId }: { priceOfferId: number }) {
  const { data: priceOffer } = usePriceOffer(priceOfferId);
  const t = useTranslations("MyAccount.reservation");

  return (
    <div className="sticky top-24">
      <SelectedCarCardWrapper>
        <SelectedCarHeader carDetails={priceOffer.carDetails} />
        <FreeCancellationBadge
          pickupDate={priceOffer.pickupDate}
          pickupTime={priceOffer.pickupTime}
          text={t("freeCancellation")}
        />
      </SelectedCarCardWrapper>
    </div>
  );
}
