import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useSelectedVehicle } from "../_hooks/useSelectedVehicle";
import { CheckIcon, ShieldCheck, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { useDirection } from "@/shared/hooks/useDirection";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { createPriceOffer } from "@/shared/api/price-offers-api";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useParams } from "next/navigation";

interface PriceOfferDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  searchRequest: booking.SearchAvailabilityRequest;
}

const PRICE_OFFER_CURRENCIES = [
  {
    code: "ILS",
    symbol: "₪",
  },
  {
    code: "USD",
    symbol: "$",
  },
  {
    code: "EUR",
    symbol: "€",
  },
];

const PRICE_OFFER_URL_PREFIX = "/price-offers/"; // + docNumber

export function PriceOfferDialog({
  open,
  onOpenChange,
  searchRequest,
}: PriceOfferDialogProps) {
  const { lang } = useParams();
  const dir = useDirection();
  const t = useTranslations("booking.plansPage");
  const { data } = useAvailableCars(searchRequest, { fromCache: true });
  const vehicle = useSelectedVehicle(searchRequest);

  const [priceOfferName, setPriceOfferName] = useState("");
  const [price, setPrice] = useState(0);
  const [currency, setCurrency] = useState("ILS");

  const [priceOfferId, setPriceOfferId] = useState<number | null>(null);

  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedPlanIndex = useBookingSessionStore((s) => s.selectedPlanIndex);
  const selectedPlan = vehicle?.plans[selectedPlanIndex];

  const { mutate, error, isPending } = useMutation({
    mutationFn: () =>
      createPriceOffer({
        includeERP: isErpSelected,
        name: priceOfferName,
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
        <h5 className="type-h5 text-navy mx-0.5">{t("priceOfferSubtitle")}</h5>
        {!priceOfferId ? (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              mutate();
            }}
            className="bg-white shadow-card rounded-xl p-3 my-3 flex flex-col justify-center gap-3"
          >
            <Input
              type="text"
              className="bg-background py-6"
              placeholder={t("enterPriceOfferName")}
              value={priceOfferName}
              onChange={(e) => setPriceOfferName(e.target.value)}
            />
            <div className="flex gap-3">
              <Input
                type="number"
                className="bg-background py-6"
                placeholder={t("enterPrice")}
                value={price > 0 ? price : ""}
                onChange={(e) => setPrice(Number(e.target.value))}
              />
              <DropdownMenu dir={dir}>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="default"
                    className="bg-background text-navy py-6 px-6 border-border"
                  >
                    {PRICE_OFFER_CURRENCIES.find((c) => c.code === currency)
                      ?.symbol ?? ""}
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  {PRICE_OFFER_CURRENCIES.map((c) => (
                    <DropdownMenuItem
                      key={c.code}
                      onClick={() => {
                        setCurrency(c.code);
                      }}
                      className="gap-2 flex"
                    >
                      <span>{c.symbol}</span>
                      {c.code === currency && (
                        <CheckIcon className="ms-auto size-4" />
                      )}
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>

            <Button
              variant="brand"
              className="font-bold py-6 px-8"
              loading={isPending}
              disabled={!priceOfferName || price <= 0}
            >
              {t("createPriceOffer")}
            </Button>
          </form>
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
