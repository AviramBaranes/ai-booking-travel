import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import { useParams } from "next/navigation";
import { SummaryRow } from "../../../../_components/SummaryRow";
import { Copy, Check, Edit } from "lucide-react";
import { useState } from "react";
import { statusToColor } from "../../../_utils/statusesStyles";
import { booking } from "@/shared/client";
import { EditPriceOfferDialog } from "./EditPriceOfferDialog";
import { PriceOfferStatus } from "@/app/(app)/[lang]/_components/priceOffer/PriceOfferForm";

const CLIENT_PRICE_OFFER_LINK_PREFIX = "/offers/";

export function HeaderSection({
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

  return (
    <>
      <div className="flex items-center justify-between">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
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
        </div>
      </div>
      <hr />
      <SummaryRow label={t("labels.name")} value={priceOffer.name} />
      <SummaryRow
        label={t("labels.status")}
        value={t(`status.${priceOffer.status}`)}
        valClassName={statusToColor(priceOffer.status)}
      />
      <SummaryRow
        label={t("labels.createdAt")}
        value={new Date(priceOffer.createdAt).toLocaleDateString(lang)}
      />
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
