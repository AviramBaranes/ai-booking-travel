import { booking } from "@/shared/client";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";

export function FeesNote({ vehicle }: { vehicle: booking.AvailableVehicle }) {
  const t = useTranslations("booking.plansPage");
  const { data } = useBookingSettings();

  if (
    !vehicle.priceDetails.dropCharge &&
    !vehicle.priceDetails.youngDriverFee
  ) {
    return null;
  }

  return (
    <div className="border border-destructive bg-destructive/15 p-6 flex flex-col gap-4 rounded-lg">
      <h6 className="type-h6 text-navy">{t("feesNoteTitle")}</h6>
      {!!vehicle.priceDetails.youngDriverFee && (
        <FeeDisplay
          title={data.youngDriverTitle}
          content={data.youngDriverContent}
          amount={vehicle.priceDetails.youngDriverFee}
          currency={vehicle.priceDetails.youngDriverFeeCurrency}
        />
      )}
      {!!vehicle.priceDetails.dropCharge && (
        <FeeDisplay
          title={data.dropoffChargeTitle}
          content={data.dropoffChargeContent}
          amount={vehicle.priceDetails.dropCharge}
          currency={vehicle.priceDetails.dropChargeCurrency}
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
