import { useState } from "react";
import { Copy, Check, Edit, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import { PriceOfferStatus } from "@/app/(app)/[lang]/_components/priceOffer/PriceOfferForm";
import { EditPriceOfferDialog } from "./EditPriceOfferDialog";
import { useParams } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { renewPriceOffer } from "@/shared/api/price-offers-api";
import { usePriceOffer } from "../../_hooks/usePriceOffer";

const CLIENT_PRICE_OFFER_LINK_PREFIX = "/offers/";

export function PriceOfferActions({ priceOfferId }: { priceOfferId: number }) {
  const { lang } = useParams();
  const t = useTranslations("MyAccount.priceOffer.summary");

  const { data: priceOffer, refetch } = usePriceOffer(priceOfferId);

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const handleCopy = () => {
    navigator.clipboard.writeText(
      `${window.location.origin}/${lang}${CLIENT_PRICE_OFFER_LINK_PREFIX}${priceOffer.token}`,
    );
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  const { mutate: renewOffer, isPending } = useMutation({
    mutationFn: (id: number) => renewPriceOffer(id),
    onSuccess: () => {
      refetch();
    },
  });

  return (
    <>
      <div className="flex gap-1 items-center w-1/4 justify-end">
        <Button
          variant="ghost"
          className="py-6 text-border-muted font-semibold flex gap-4"
          onClick={handleCopy}
          disabled={isCopied}
        >
          {isCopied ? (
            <Check className="w-6 h-6 text-green-500" />
          ) : (
            <Copy className="w-6 h-6" />
          )}
          {t("clientLink")}
        </Button>
        <Button
          variant="ghost"
          className="py-6 text-border-muted font-semibold flex gap-4"
          onClick={() => setIsEditDialogOpen(true)}
        >
          <Edit className="w-6 h-6" />
          {t("editOffer")}
        </Button>
        <Button
          variant="ghost"
          loading={isPending}
          className="py-6 text-border-muted font-semibold flex gap-4"
          onClick={() => renewOffer(priceOfferId)}
        >
          <RefreshCw className="w-6 h-6" />
          {t("renewOffer")}
        </Button>
      </div>
      <EditPriceOfferDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        priceOfferId={priceOfferId}
        initialName={priceOffer.name}
        initialPrice={priceOffer.offeredPrice}
        initialCurrency={priceOffer.offeredCurrencyCode}
        initialStatus={priceOffer.status as PriceOfferStatus}
      />
    </>
  );
}
