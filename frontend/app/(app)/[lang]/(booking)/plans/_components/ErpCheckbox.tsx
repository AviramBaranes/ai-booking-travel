import { Checkbox } from "@/components/ui/checkbox";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { useMemo, useState } from "react";
import { RentalPriceForDays } from "../../../../../../shared/components/booking/RentalPriceForDays";
import { formatPrice } from "@/shared/utils/formatPrice";

interface ErpCheckboxProps {
  vehicle: booking.AvailableVehicle;
  daysCount: number;
  selectedPlan: number;
  isSelected: boolean;
  setSelected: (selected: boolean) => void;
}

export function ErpCheckbox({
  vehicle,
  daysCount,
  selectedPlan,
  isSelected,
  setSelected,
}: ErpCheckboxProps) {
  const t = useTranslations("booking.erpCheckbox");

  const [isReadMore, setIsReadMore] = useState(false);
  const description = useMemo(() => {
    const fullText = t("description");
    if (isReadMore) {
      return fullText;
    }

    const words = fullText.split(" ");
    if (words.length <= 24) {
      return fullText;
    }
    return words.slice(0, 24).join(" ");
  }, [t, isReadMore]);

  return (
    <>
      <h5 className="type-h5 text-navy my-6">{t("title")}</h5>
      <div className="bg-white border-brand border rounded-lg flex justify-between p-6">
        <div className="w-3/4">
          <label className="flex items-center gap-2 cursor-pointer">
            <Checkbox
              checked={isSelected}
              onCheckedChange={setSelected}
              className="border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand"
            />
            <span className="type-label text-navy ">{t("label")}</span>
          </label>
          <p className="type-paragraph text-text-secondary ">
            {description}
            <button
              className="text-navy type-label underline mx-2 cursor-pointer"
              onClick={() => setIsReadMore((s) => !s)}
            >
              {isReadMore ? t("readLess") : t("readMore")}
            </button>
          </p>
        </div>
        <div className="w-1/4 mt-4 flex items-end flex-col">
          <RentalPriceForDays daysCount={daysCount} />
          <h4 className="type-h4 text-navy">
            {formatPrice(
              vehicle.plans[selectedPlan].erpPrice,
              vehicle.priceDetails.currency,
            )}
          </h4>
        </div>
      </div>
    </>
  );
}
