import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { ShieldCheck, X } from "lucide-react";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { updatePriceOffer } from "@/shared/api/price-offers-api";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { PriceOfferForm } from "@/app/(app)/[lang]/_components/priceOffer/PriceOfferForm";
import { usePriceOffer } from "../../_hooks/usePriceOffer";

interface EditPriceOfferDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  priceOfferId: number;
  initialName: string;
  initialPrice: number;
  initialCurrency: string;
}

export function EditPriceOfferDialog({
  open,
  onOpenChange,
  priceOfferId,
  initialName,
  initialPrice,
  initialCurrency,
}: EditPriceOfferDialogProps) {
  const { refetch } = usePriceOffer(priceOfferId);
  const t = useTranslations("MyAccount.priceOffer.editDialog");

  const [isSuccess, setIsSuccess] = useState(false);

  const { mutate, isPending, error } = useMutation({
    mutationFn: (params: booking.UpdatePriceOfferParams) =>
      updatePriceOffer(priceOfferId, params),
    onSuccess: () => {
      setIsSuccess(true);
      refetch();
      setTimeout(() => {
        setIsSuccess(false);
        onOpenChange(false);
      }, 2000);
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
            <span className="type-h5 text-navy">{t("title")}</span>
          </DialogTitle>
          <button
            onClick={() => onOpenChange(false)}
            className="p-2 cursor-pointer"
          >
            <X className="w-6 h-6 text-navy" />
          </button>
        </div>
        <hr />
        {!isSuccess ? (
          <>
            <h5 className="type-h5 text-navy mx-0.5">{t("subtitle")}</h5>
            <PriceOfferForm
              error={error}
              isPending={isPending}
              nameInputPlaceholder={t("priceOfferName")}
              priceInputPlaceholder={t("priceOfferPrice")}
              submitText={t("saveButton")}
              onSubmit={({ name, price, currency }) => {
                mutate({
                  name,
                  offeredPrice: price,
                  offeredCurrencyCode: currency,
                });
              }}
              initialCurrency={initialCurrency}
              initialName={initialName}
              initialPrice={initialPrice}
            />
          </>
        ) : (
          <SuccessBadge>{t("priceOfferUpdatedSuccess")}</SuccessBadge>
        )}
      </DialogContent>
    </Dialog>
  );
}
