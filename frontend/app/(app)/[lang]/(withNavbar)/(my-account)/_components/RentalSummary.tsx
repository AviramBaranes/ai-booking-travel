"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { LocationDateTimeSummary } from "./LocationSummary";
import { SummaryRow } from "./SummaryRow";

interface RentalSummaryProps {
  pickupDate: string;
  pickupTime: string;
  pickupLocationName: string;
  dropoffDate: string;
  dropoffTime: string;
  dropoffLocationName: string;
  flightNumber?: string;
}

export function RentalSummary({
  pickupDate,
  pickupTime,
  pickupLocationName,
  dropoffDate,
  dropoffTime,
  dropoffLocationName,
  flightNumber,
}: RentalSummaryProps) {
  const t = useTranslations("MyAccount.summary");

  return (
    <>
      <SummarySubTitle title={t("sections.rentalInfo")} />
      {flightNumber && (
        <div className="w-1/4">
          <SummaryRow label={t("labels.flightNumber")} value={flightNumber} />
        </div>
      )}
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
            title={t("rentalSummary.dropoffDetails")}
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
