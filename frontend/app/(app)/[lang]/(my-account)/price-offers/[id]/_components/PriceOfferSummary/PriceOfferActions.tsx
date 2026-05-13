import { useState } from "react";
import { Copy, Check, Edit, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { PriceOfferStatus } from "@/app/(app)/[lang]/_components/priceOffer/PriceOfferForm";
import { EditPriceOfferDialog } from "./EditPriceOfferDialog";
import { useParams } from "next/navigation";

const CLIENT_PRICE_OFFER_LINK_PREFIX = "/offers/";

export function PriceOfferActions({
  priceOffer,
}: {
  priceOffer: booking.GetAgentPriceOfferResponse;
}) {
  const { lang } = useParams();
  const t = useTranslations("MyAccount.priceOffer.summary");

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const handleCopy = () => {
    navigator.clipboard.writeText(
      `${window.location.origin}/${lang}${CLIENT_PRICE_OFFER_LINK_PREFIX}${priceOffer.token}`,
    );
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  function researchOffer() {
    const urlParams = new URLSearchParams();

    urlParams.set("pl", priceOffer.pickupLocationId.toString());
    urlParams.set("rl", priceOffer.dropoffLocationId.toString());
    urlParams.set("pd", priceOffer.pickupDate);
    urlParams.set("pt", priceOffer.pickupTime);
    urlParams.set("rd", priceOffer.returnDate);
    urlParams.set("rt", priceOffer.dropoffTime);
    urlParams.set("da", priceOffer.driverAge.toString());

    const href = `/${lang}/results?${urlParams.toString()}`;
    window.open(href, "_blank");
  }

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
          className="py-6 text-border-muted font-semibold flex gap-4"
          onClick={researchOffer}
        >
          <Search className="w-6 h-6" />
          {t("renewOffer")}
        </Button>
      </div>
      <EditPriceOfferDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        priceOfferId={priceOffer.id}
        initialName={priceOffer.name}
        initialPrice={priceOffer.offeredPrice}
        initialCurrency={priceOffer.offeredCurrencyCode}
        initialStatus={priceOffer.status as PriceOfferStatus}
      />
    </>
  );
}
