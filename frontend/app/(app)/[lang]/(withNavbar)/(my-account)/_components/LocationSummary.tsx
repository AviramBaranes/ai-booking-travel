import { useParams } from "next/navigation";
import { useState } from "react";
import {
  hasStationInfo,
  SupplierInfoData,
  SupplierInfoDialog,
  SupplierInfoTab,
} from "@/shared/components/booking/SupplierInfoDialog";

interface LocationDateTimeSummaryProps {
  title: string;
  locationName: string;
  date: string;
  time: string;
  linkText: string;
  /** Omitted for reservations booked before supplier info was captured; hides the station link. */
  supplierInfo?: SupplierInfoData;
  /** Which station this summary describes, so the dialog opens on the matching tab. */
  stationTab: Extract<SupplierInfoTab, "pickupDetails" | "dropoffDetails">;
}
export function LocationDateTimeSummary({
  title,
  locationName,
  date,
  time,
  linkText,
  supplierInfo,
  stationTab,
}: LocationDateTimeSummaryProps) {
  const { lang } = useParams();
  const [open, setOpen] = useState(false);

  const stationInfo =
    stationTab === "pickupDetails"
      ? supplierInfo?.pickupDetails
      : supplierInfo?.dropoffDetails;
  const showStationLink = hasStationInfo(stationInfo);

  return (
    <div className="flex flex-col gap-2 mt-2">
      <p className="type-label text-navy">{title}</p>
      <p className="type-paragraph text-text-secondary">{locationName}</p>
      <p className="type-paragraph text-text-secondary">
        {new Date(date).toLocaleDateString(lang)} | {time}
      </p>
      {showStationLink && (
        <>
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="text-brand underline type-label print:hidden cursor-pointer w-fit"
          >
            {linkText}
          </button>
          <SupplierInfoDialog
            open={open}
            onOpenChange={setOpen}
            initialTab={stationTab}
            termsAndConditions={supplierInfo?.termsAndConditions}
            pickupDetails={supplierInfo?.pickupDetails}
            dropoffDetails={supplierInfo?.dropoffDetails}
          />
        </>
      )}
    </div>
  );
}
