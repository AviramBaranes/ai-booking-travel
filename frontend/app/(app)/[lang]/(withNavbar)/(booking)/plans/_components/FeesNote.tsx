import { availability } from "@/shared/client";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";

export function FeesNote({ vehicle }: { vehicle: availability.AvailableVehicle }) {
  const t = useTranslations("booking.plansPage");
  const { data } = useBookingSettings();

  if (
    !vehicle.priceDetails.fees.dropCharge &&
    !vehicle.priceDetails.fees.youngDriverFee
  ) {
    return null;
  }

  return (
    <div className="border border-destructive bg-destructive/15 p-6 flex flex-col gap-4 rounded-lg">
      <h6 className="type-h6 text-navy">{t("feesNoteTitle")}</h6>
      {!!vehicle.priceDetails.fees.youngDriverFee && (
        <FeeDisplay
          title={data.youngDriverTitle}
          content={data.youngDriverContent}
          amount={vehicle.priceDetails.fees.youngDriverFee}
          currency={vehicle.priceDetails.fees.youngDriverFeeCurrency}
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
      <h6 className="type-h6 text-navy">
        {title}: {formatPrice(amount, currency)}
      </h6>
      <p className="type-h6 font-normal text-navy">{content}</p>
    </>
  );
}
