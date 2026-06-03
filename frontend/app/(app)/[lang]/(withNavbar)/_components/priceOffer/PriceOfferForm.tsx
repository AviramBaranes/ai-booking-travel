import { useState } from "react";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useDirection } from "@/shared/hooks/useDirection";
import { Button } from "@/components/ui/button";
import { CheckIcon } from "lucide-react";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { useTranslations } from "next-intl";
import { statusToColor } from "../../(my-account)/price-offers/_utils/statusesStyles";
import { Label } from "@/components/ui/label";

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
  {
    code: "GBP",
    symbol: "£",
  },
];

const PRICE_OFFER_STATUSES: PriceOfferStatus[] = ["open", "booked", "declined"];

export type PriceOfferStatus = "open" | "booked" | "declined";

interface PriceOfferFormProps {
  onSubmit: (data: {
    name: string;
    price: number;
    currency: string;
    status?: PriceOfferStatus;
  }) => void;
  isPending: boolean;
  error: Error | null;
  nameInputPlaceholder: string;
  priceInputPlaceholder: string;
  submitText: string;
  statusLabel?: string;
  showStatusSelect?: boolean;
  initialName?: string;
  initialPrice?: number;
  initialCurrency?: string;
  initialStatus?: PriceOfferStatus;
}

export function PriceOfferForm({
  onSubmit,
  isPending,
  error,
  nameInputPlaceholder,
  priceInputPlaceholder,
  submitText,
  initialStatus,
  statusLabel,
  showStatusSelect = false,
  initialName = "",
  initialPrice = 0,
  initialCurrency = "ILS",
}: PriceOfferFormProps) {
  const dir = useDirection();
  const t = useTranslations("MyAccount.priceOffer.summary");
  const [priceOfferName, setPriceOfferName] = useState(initialName);
  const [price, setPrice] = useState(initialPrice);
  const [currency, setCurrency] = useState(initialCurrency);
  const [status, setStatus] = useState<PriceOfferStatus | null>(
    initialStatus ?? null,
  );

  const translatedError = useTranslatedError(error);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit({
          name: priceOfferName,
          price,
          currency,
          status: showStatusSelect && status ? status : undefined,
        });
      }}
      className="bg-white shadow-card rounded-xl p-3 my-3 flex flex-col justify-center gap-3"
    >
      <Input
        type="text"
        className="bg-background py-6"
        placeholder={nameInputPlaceholder}
        value={priceOfferName}
        onChange={(e) => setPriceOfferName(e.target.value)}
      />
      <div className="flex gap-3">
        <Input
          type="number"
          className="bg-background py-6"
          placeholder={priceInputPlaceholder}
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

      {showStatusSelect && (
        <>
          <Label className="text-sm mt-1">{statusLabel}</Label>
          <DropdownMenu dir={dir}>
            <DropdownMenuTrigger asChild>
              <Button
                type="button"
                variant="default"
                className="bg-background text-navy py-6 px-4 border-border w-full justify-between"
              >
                <span
                  className={`font-medium ${status ? statusToColor(status) : "text-navy"}`}
                >
                  {status ? t(`status.${status}`) : t("labels.status")}
                </span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-(--radix-dropdown-menu-trigger-width)">
              {PRICE_OFFER_STATUSES.map((s) => (
                <DropdownMenuItem
                  key={s}
                  onClick={() => {
                    setStatus(s);
                  }}
                  className="gap-2 flex"
                >
                  <span className={statusToColor(s)}>{t(`status.${s}`)}</span>
                  {s === status && <CheckIcon className="ms-auto size-4" />}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </>
      )}

      <Button
        variant="brand"
        className="font-bold py-6 px-8"
        loading={isPending}
        disabled={
          !priceOfferName ||
          price <= 0 ||
          !PRICE_OFFER_CURRENCIES.some((c) => c.code === currency)
        }
      >
        {submitText}
      </Button>
      {!!translatedError && <ErrorDisplay>{translatedError}</ErrorDisplay>}
    </form>
  );
}
