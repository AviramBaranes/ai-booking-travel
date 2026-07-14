"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { useDirection } from "@/shared/hooks/useDirection";
import { ChevronDown } from "lucide-react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import {
  buildPriceOfferFiltersQuery,
  PRICE_OFFER_STATUSES,
  PriceOfferStatus,
  usePriceOfferFilters,
} from "../../_hooks/usePriceOfferFilters";
import { statusToColor } from "../../_utils/statusesStyles";

export function FilterForm() {
  const { searchParams } = usePriceOfferFilters();

  return <FilterFormFields key={searchParams.toString()} />;
}

function FilterFormFields() {
  const router = useRouter();
  const { lang, filters } = usePriceOfferFilters();
  const t = useTranslations("MyAccount.priceOffers");

  const [name, setName] = useState(filters.name);
  const [status, setStatus] = useState<PriceOfferStatus | null>(
    filters.status,
  );

  function handleSearch(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = buildPriceOfferFiltersQuery({
      name,
      status,
    });

    const queryString = query.toString();
    const basePath = `/${lang}/price-offers`;
    router.push(queryString ? `${basePath}?${queryString}` : basePath);
  }

  return (
    <form
      className="flex flex-col lg:flex-row gap-4 justify-between items-start lg:items-center"
      onSubmit={handleSearch}
    >
      <legend className="type-label text-navy w-50">{t("filtersLabel")}</legend>
      <div className="flex justify-between gap-4 w-full">
        <Input
          className="bg-white border w-1/2 border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary"
          placeholder={t("namePlaceholder")}
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <StatusDropdown
          value={status}
          onChange={setStatus}
          placeholder={t("statusPlaceholder")}
        />
      </div>
      <Button variant="brand" className="py-6 w-full lg:w-40 font-semibold">
        {t("searchButton")}
      </Button>
    </form>
  );
}

function StatusDropdown({
  value,
  onChange,
  placeholder,
}: {
  value: PriceOfferStatus | null;
  placeholder: string;
  onChange: (value: PriceOfferStatus | null) => void;
}) {
  const t = useTranslations("MyAccount.priceOffer.summary.status");
  const dir = useDirection();

  return (
    <DropdownMenu dir={dir}>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className={cn(
            "w-1/2 flex items-center justify-between bg-white border rounded-lg px-4 h-12 cursor-pointer",
            "text-sm font-normal font-[inherit]",
          )}
        >
          <span
            className={cn(
              "text-sm font-normal text-muted-foreground",
              value ? statusToColor(value) : "",
            )}
          >
            {value ? t(value) : placeholder}
          </span>
          <ChevronDown className="w-4 h-4 text-muted shrink-0" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="start"
        className="w-(--radix-dropdown-menu-trigger-width) items-start"
      >
        {PRICE_OFFER_STATUSES.map((status) => (
          <DropdownMenuItem
            className={`font-semibold ${statusToColor(status)}`}
            key={status}
            onClick={() => onChange(status)}
          >
            {t(status)}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}