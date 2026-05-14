import { booking } from "@/shared/client";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { getTranslations } from "next-intl/server";

export async function ClientOfferCarCard({
  offer,
}: {
  offer: booking.GetPriceOfferResponse;
}) {
  const t = await getTranslations("MyAccount.reservation");

  return (
    <div className="sticky top-24">
      <SelectedCarCardWrapper>
        <SelectedCarHeader carDetails={offer.carDetails} />
        <FreeCancellationBadge
          pickupDate={offer.pickupDate}
          pickupTime={offer.pickupTime}
          text={t("freeCancellation")}
        />
      </SelectedCarCardWrapper>
    </div>
  );
}
