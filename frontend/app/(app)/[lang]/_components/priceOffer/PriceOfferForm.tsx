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

interface PriceOfferFormProps {
  onSubmit: (data: { name: string; price: number; currency: string }) => void;
  isPending: boolean;
  error: Error | null;
  nameInputPlaceholder: string;
  priceInputPlaceholder: string;
  submitText: string;
  initialName?: string;
  initialPrice?: number;
  initialCurrency?: string;
}

export function PriceOfferForm({
  onSubmit,
  isPending,
  error,
  nameInputPlaceholder,
  priceInputPlaceholder,
  submitText,
  initialName = "",
  initialPrice = 0,
  initialCurrency = "ILS",
}: PriceOfferFormProps) {
  const dir = useDirection();
  const [priceOfferName, setPriceOfferName] = useState(initialName);
  const [price, setPrice] = useState(initialPrice);
  const [currency, setCurrency] = useState(initialCurrency);

  const translatedError = useTranslatedError(error);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit({
          name: priceOfferName,
          price,
          currency,
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

      <Button
        variant="brand"
        className="font-bold py-6 px-8"
        loading={isPending}
        disabled={!priceOfferName || price <= 0}
      >
        {submitText}
      </Button>
      {!!translatedError && <ErrorDisplay>{translatedError}</ErrorDisplay>}
    </form>
  );
}
