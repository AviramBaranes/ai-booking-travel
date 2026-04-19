import { reservation } from "@/shared/client";
import { useTranslations } from "next-intl";
import { OrderSummarySubTitle } from "./OrderSummarySubTitle";
import { LocationDateTimeSummary } from "./LocationSummary";

export function RentalSummary({
  reservation,
}: {
  reservation: reservation.GetReservationResponse;
}) {
  const t = useTranslations("MyAccount.reservation.summary");

  return (
    <>
      <OrderSummarySubTitle title={t("rental.title")} />
      <div className="flex">
        <div className="w-1/2">
          <LocationDateTimeSummary
            title={t("rental.pickupDetails")}
            date={reservation.pickupDate}
            time={reservation.pickupTime}
            locationName={reservation.pickupLocationName}
            linkText={t("rental.stationDetails")}
          />
        </div>
        <div className="w-1/2">
          <LocationDateTimeSummary
            title={t("rental.returnDetails")}
            date={reservation.returnDate}
            time={reservation.dropoffTime}
            locationName={reservation.dropoffLocationName}
            linkText={t("rental.stationDetails")}
          />
        </div>
      </div>
    </>
  );
}
