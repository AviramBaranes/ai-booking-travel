import { useState, Fragment } from "react";
import { useTranslations } from "next-intl";
import DOMPurify from "isomorphic-dompurify";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { availability, broker } from "@/shared/client";
import { X, Info, MapPin, Phone, Clock } from "lucide-react";

interface ImportantInfoButtonProps {
  plans: availability.Plan[];
  suppliersInfo: broker.SupplierInfo[];
  selectedPlanIndex: number;
}
export function ImportantInfoButton({
  plans,
  suppliersInfo,
  selectedPlanIndex,
}: ImportantInfoButtonProps) {
  const t = useTranslations("booking.infoDialog");
  const [open, setOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<
    "terms" | "pickupDetails" | "dropoffDetails"
  >("terms");

  const selectedPlan = plans[selectedPlanIndex];
  const supplierInfo = suppliersInfo.find(
    (s) => s.name === selectedPlan.supplierName,
  );
  const isDropoffDifferent =
    supplierInfo?.pickupDetails.address !==
    supplierInfo?.dropoffDetails.address;

  return (
    <>
      <Button
        onClick={() => setOpen(true)}
        variant="outline"
        className="px-4 lg:px-8 py-6 type-paragraph text-navy rounded-lg"
      >
        {t("importantInfo")}
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          className="max-w-242.5! flex max-h-[90vh] flex-col overflow-hidden rounded-3xl border-none bg-background p-0 shadow-[0px_4px_24px_0px_rgba(0,0,0,0.25)] ring-0 max-sm:mx-auto max-sm:w-11/12"
          showCloseButton={false}
        >
          <div className="flex items-center justify-between px-10 pt-4 pb-0">
            <DialogTitle className="flex items-center gap-4">
              <Info className="w-8 h-8 text-navy" />
              <span className="type-h5 text-navy">{t("importantInfo")}</span>
            </DialogTitle>
            <button
              onClick={() => setOpen(false)}
              className="p-2 cursor-pointer"
            >
              <X className="w-6 h-6 text-navy" />
            </button>
          </div>

          {/* Divider + Tabs */}
          <div className="flex flex-col items-center gap-3 px-12">
            <div className="w-full border-t border-border-light" />
            <div className="flex gap-6 items-center text-[22px] leading-[30.8px] text-brand-blue">
              <button
                className={`cursor-pointer underline ${
                  activeTab === "terms" ? "font-bold" : "font-normal"
                }`}
                onClick={() => setActiveTab("terms")}
              >
                {t("terms")}
              </button>
              <button
                className={`cursor-pointer underline ${
                  activeTab === "pickupDetails" ? "font-bold" : "font-normal"
                }`}
                onClick={() => setActiveTab("pickupDetails")}
              >
                {t("pickupDetails")}
              </button>
              {isDropoffDifferent && (
                <button
                  className={`cursor-pointer underline ${
                    activeTab === "dropoffDetails" ? "font-bold" : "font-normal"
                  }`}
                  onClick={() => setActiveTab("dropoffDetails")}
                >
                  {t("dropoffDetails")}
                </button>
              )}
            </div>
          </div>

          <div className="min-h-0 flex-1 overflow-y-auto">
            {/* Terms */}
            {activeTab === "terms" && (
              <div
                dir="ltr"
                className="flex flex-col items-stretch justify-center gap-6 px-12 pb-6"
              >
                {supplierInfo?.termsAndConditions.map(
                  ({ title, htmlContent }, i) => (
                    <div key={title+i} className="flex flex-col gap-3">
                      <h6 className="type-h6 text-navy">{title}</h6>
                      <div
                        className="type-paragraph text-navy"
                        dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(htmlContent) }}
                      />
                    </div>
                  ),
                )}
              </div>
            )}
            {activeTab === "pickupDetails" && (
              <StationInfoTab stationInfo={supplierInfo?.pickupDetails!} />
            )}
            {activeTab === "dropoffDetails" && (
              <StationInfoTab stationInfo={supplierInfo?.dropoffDetails!} />
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

function StationInfoTab({ stationInfo }: { stationInfo: broker.StationInfo }) {
  const t = useTranslations("booking.infoDialog");
  const { address, locationInfo, openingHours, phoneNumber } = stationInfo;

  return (
    <div className="flex flex-col gap-6 px-12 pb-8 pt-2">
      {/* Address */}
      {address && (
        <div className="flex items-start gap-3">
          <MapPin className="mt-0.5 h-5 w-5 shrink-0 text-brand-blue" />
          <div className="flex flex-col gap-0.5">
            <span className="type-label text-text-secondary">
              {t("station.address")}
            </span>
            <span dir="ltr" className="type-paragraph text-navy text-left">
              {address}
            </span>
          </div>
        </div>
      )}

      {/* Location info */}
      {locationInfo && (
        <div className="flex items-start gap-3">
          <Info className="mt-0.5 h-5 w-5 shrink-0 text-brand-blue" />
          <div className="flex flex-col gap-0.5">
            <span className="type-label text-text-secondary">
              {t("station.locationInfo")}
            </span>
            <span
              dir="ltr"
              className="type-paragraph text-navy text-left"
              dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(locationInfo) }}
            />
          </div>
        </div>
      )}

      {/* Phone */}
      {phoneNumber && (
        <div className="flex items-start gap-3">
          <Phone className="mt-0.5 h-5 w-5 shrink-0 text-brand-blue" />
          <div className="flex flex-col gap-0.5">
            <span className="type-label text-text-secondary">
              {t("station.phone")}
            </span>
            <a
              href={`tel:${phoneNumber}`}
              className="type-paragraph text-navy hover:underline"
            >
              {phoneNumber}
            </a>
          </div>
        </div>
      )}

      {/* Opening hours */}
      {openingHours && openingHours.length > 0 && (
        <div className="flex items-start gap-3">
          <Clock className="mt-0.5 h-5 w-5 shrink-0 text-brand-blue" />
          <div className="flex flex-col gap-1.5">
            <span className="type-label text-text-secondary">
              {t("station.openingHours")}
            </span>
            <div className="grid grid-cols-[auto_1fr] gap-x-6 gap-y-1">
              {openingHours.map(({ day, openTime, closeTime }) => (
                <Fragment key={day}>
                  <span className="type-paragraph font-medium text-navy">
                    {t(`station.days.${day}`)}
                  </span>
                  {openTime === "Closed" ? (
                    <span className="type-paragraph text-red-500">
                      {t("station.closed")}
                    </span>
                  ) : (
                    <span className="type-paragraph text-text-secondary">
                      {openTime} – {closeTime}
                    </span>
                  )}
                </Fragment>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
