import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { availability, broker } from "@/shared/client";
import { SupplierInfoDialog } from "@/shared/components/booking/SupplierInfoDialog";

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

  const selectedPlan = plans[selectedPlanIndex];
  const supplierInfo = suppliersInfo.find(
    (s) => s.name === selectedPlan.supplierName,
  );

  return (
    <>
      <Button
        onClick={() => setOpen(true)}
        variant="outline"
        className="px-4 lg:px-8 py-6 type-paragraph text-navy rounded-lg"
      >
        {t("importantInfo")}
      </Button>
      <SupplierInfoDialog
        open={open}
        onOpenChange={setOpen}
        termsAndConditions={supplierInfo?.termsAndConditions}
        pickupDetails={supplierInfo?.pickupDetails}
        dropoffDetails={supplierInfo?.dropoffDetails}
      />
    </>
  );
}
