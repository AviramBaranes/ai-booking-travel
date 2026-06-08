import { price_offer } from "@/shared/client";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { getTranslations } from "next-intl/server";
import { ApproveButton } from "./ApproveButton";

export async function ClientOfferCarCard({
  offer,
}: {
  offer: price_offer.GetPriceOfferResponse;
}) {
  const t = await getTranslations("MyAccount");

  return (
    <div className="sticky top-4">
      <SelectedCarCardWrapper>
        <SelectedCarHeader carDetails={offer.carDetails} />
        <FreeCancellationBadge
          pickupDate={offer.pickupDate}
          pickupTime={offer.pickupTime}
          text={t("reservation.freeCancellation")}
        />
        {offer.status === "open" && (
          <ApproveButton
            id={offer.id}
            text={t("priceOffer.clientOrderCTA")}
            successText={t("priceOffer.orderApproveSuccess")}
          />
        )}
      </SelectedCarCardWrapper>
    </div>
  );
}
