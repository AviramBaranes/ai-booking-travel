import { reservation } from "@/shared/client";
import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { LocationDateTimeSummary } from "./LocationSummary";

interface RentalSummaryProps {
  pickupDate: string;
  pickupTime: string;
  pickupLocationName: string;
  dropoffDate: string;
  dropoffTime: string;
  dropoffLocationName: string;
}

export function RentalSummary({
  pickupDate,
  pickupTime,
  pickupLocationName,
  dropoffDate,
  dropoffTime,
  dropoffLocationName,
}: RentalSummaryProps) {
  const t = useTranslations("MyAccount.summary");

  return (
    <>
      <SummarySubTitle title={t("sections.rentalInfo")} />
      <div className="flex">
        <div className="w-1/2">
          <LocationDateTimeSummary
            title={t("rentalSummary.pickupDetails")}
            date={pickupDate}
            time={pickupTime}
            locationName={pickupLocationName}
            linkText={t("rentalSummary.stationDetails")}
          />
        </div>
        <div className="w-1/2">
          <LocationDateTimeSummary
            title={t("rentalSummary.returnDetails")}
            date={dropoffDate}
            time={dropoffTime}
            locationName={dropoffLocationName}
            linkText={t("rentalSummary.stationDetails")}
          />
        </div>
      </div>
    </>
  );
}
