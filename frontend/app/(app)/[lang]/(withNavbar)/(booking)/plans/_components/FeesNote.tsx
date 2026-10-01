import { availability } from "@/shared/client";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";

export function FeesNote({
  vehicle,
}: {
  vehicle: availability.AvailableVehicle;
}) {
  const t = useTranslations("booking.plansPage");
  const { data } = useBookingSettings();
  const selectedPlan = useBookingSessionStore((s) => s.selectedPlanIndex);
  const plan = vehicle.plans[selectedPlan];

  if (
    !vehicle.priceDetails.fees.dropCharge &&
    !vehicle.priceDetails.fees.youngDriverFee &&
    !vehicle.priceDetails.fees.seniorDriverFee &&
    !plan.deposit &&
    !plan.theftExcess
  ) {
    return null;
  }

  return (
    <div className="border border-warning bg-warning/15 p-6 flex flex-col gap-4 rounded-lg">
      <h6 className="type-h6 text-black">{t("feesNoteTitle")}</h6>
      {!!vehicle.priceDetails.fees.youngDriverFee && (
        <FeeDisplay
          title={data.youngDriverTitle}
          content={data.youngDriverContent}
          amount={vehicle.priceDetails.fees.youngDriverFee}
          currency={vehicle.priceDetails.fees.youngDriverFeeCurrency}
        />
      )}
      {!!vehicle.priceDetails.fees.seniorDriverFee && (
        <FeeDisplay
          title={data.seniorDriverTitle}
          content={data.seniorDriverContent}
          amount={vehicle.priceDetails.fees.seniorDriverFee}
          currency={vehicle.priceDetails.fees.seniorDriverFeeCurrency}
        />
      )}
      {!!vehicle.priceDetails.fees.dropCharge && (
        <FeeDisplay
          title={data.dropoffChargeTitle}
          content={data.dropoffChargeContent}
          amount={vehicle.priceDetails.fees.dropCharge}
          currency={vehicle.priceDetails.fees.dropChargeCurrency}
        />
      )}
      {!!plan.deposit && (
        <FeeDisplay
          title={data.depositChargeTitle}
          content={data.depositChargeContent}
          amount={plan.deposit}
          currency={plan.depositCurrency}
        />
      )}
      {!!plan.theftExcess && (
        <FeeDisplay
          title={data.theftExcessTitle}
          content={data.theftExcessContent}
          amount={plan.theftExcess}
          currency={plan.theftExcessCurrency}
        />
      )}
    </div>
  );
}

interface FeeDisplayProps {
  title: string;
  content: string;
  amount: number;
  currency: string;
}
function FeeDisplay({ title, content, amount, currency }: FeeDisplayProps) {
  return (
    <>
      <h6 className="type-h6 text-black">
        {title}: {formatPrice(amount, currency)}
      </h6>
      <p className="type-h6 font-normal text-black">{content}</p>
    </>
  );
}
