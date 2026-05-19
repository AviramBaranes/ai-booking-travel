import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { availability } from "@/shared/client";
import { useTranslations } from "next-intl";
import { useSelectedVehicle } from "../_hooks/useSelectedVehicle";
import { ShieldCheck, X } from "lucide-react";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { createPriceOffer } from "@/shared/api/price-offers-api";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useParams } from "next/navigation";
import { PriceOfferForm } from "../../../_components/priceOffer/PriceOfferForm";

interface PriceOfferDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  searchRequest: availability.SearchAvailabilityParams;
}

const PRICE_OFFER_URL_PREFIX = "/price-offers/"; // + docNumber

export function PriceOfferDialog({
  open,
  onOpenChange,
  searchRequest,
}: PriceOfferDialogProps) {
  const { lang } = useParams();
  const t = useTranslations("booking.plansPage");
  const { data } = useAvailableCars(searchRequest, { fromCache: true });
  const vehicle = useSelectedVehicle(searchRequest);

  const [priceOfferId, setPriceOfferId] = useState<number | null>(null);

  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedPlanIndex = useBookingSessionStore((s) => s.selectedPlanIndex);
  const selectedPlan = vehicle?.plans[selectedPlanIndex];

  const { mutate, error, isPending } = useMutation({
    mutationFn: ({
      name,
      price,
      currency,
    }: {
      name: string;
      price: number;
      currency: string;
    }) =>
      createPriceOffer({
        includeERP: isErpSelected,
        name,
        offeredCurrencyCode: currency,
        offeredPrice: price,
        snapshotId: data?.snapshotId ?? 0,
        rateQualifier: selectedPlan?.rateQualifier ?? "",
        supplierCode: selectedPlan?.supplierCode ?? "",
      }),
    onSuccess: (data) => {
      setPriceOfferId(data.id);
    },
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="min-w-1/3 max-w-md py-6 px-10 flex flex-col gap-4 bg-background border-border-light/50 rounded-2xl shadow-modal"
        showCloseButton={false}
      >
        <div className="flex items-center justify-between p-3 pb-0">
          <DialogTitle className="flex items-center gap-4">
            <ShieldCheck className="w-8 h-8 text-success" />
            <span className="type-h5 text-navy">{t("priceOfferTitle")}</span>
          </DialogTitle>
          <button
            onClick={() => onOpenChange(false)}
            className="p-2 cursor-pointer"
          >
            <X className="w-6 h-6 text-navy" />
          </button>
        </div>
        <hr />
        {!priceOfferId ? (
          <>
            <h5 className="type-h5 text-navy mx-0.5">
              {t("priceOfferSubtitle")}
            </h5>
            <PriceOfferForm
              error={error}
              isPending={isPending}
              nameInputPlaceholder={t("enterPriceOfferName")}
              priceInputPlaceholder={t("enterPrice")}
              submitText={t("createPriceOffer")}
              onSubmit={({ name, price, currency }) => {
                mutate({ name, price, currency });
              }}
            />
          </>
        ) : (
          <SuccessBadge>
            {t("priceOfferCreatedSuccess")}{" "}
            <a
              href={`/${lang}/${PRICE_OFFER_URL_PREFIX}${priceOfferId}`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-brand-blue underline"
            >
              {t("viewPriceOffer")}
            </a>
          </SuccessBadge>
        )}
      </DialogContent>
    </Dialog>
  );
}
