"use client";

import { useTranslations } from "next-intl";
import { SummarySubTitle } from "./SummarySubTitle";
import { LocationDateTimeSummary } from "./LocationSummary";
import { SummaryRow } from "./SummaryRow";
import { SupplierInfoData } from "@/shared/components/booking/SupplierInfoDialog";

interface RentalSummaryProps {
  pickupDate: string;
  pickupTime: string;
  pickupLocationName: string;
  dropoffDate: string;
  dropoffTime: string;
  dropoffLocationName: string;
  flightNumber?: string;
  supplierInfo?: SupplierInfoData;
}

export function RentalSummary({
  pickupDate,
  pickupTime,
  pickupLocationName,
  dropoffDate,
  dropoffTime,
  dropoffLocationName,
  flightNumber,
  supplierInfo,
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
            supplierInfo={supplierInfo}
            stationTab="pickupDetails"
          />
        </div>
        <div className="w-1/2">
          <LocationDateTimeSummary
            title={t("rentalSummary.dropoffDetails")}
            date={dropoffDate}
            time={dropoffTime}
            locationName={dropoffLocationName}
            linkText={t("rentalSummary.stationDetails")}
            supplierInfo={supplierInfo}
            stationTab="dropoffDetails"
          />
        </div>
      </div>
    </>
  );
}
