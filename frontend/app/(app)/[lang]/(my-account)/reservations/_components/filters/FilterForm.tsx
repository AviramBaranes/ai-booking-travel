"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { useDirection } from "@/shared/hooks/useDirection";
import { ChevronDown } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";
import { statusToColor } from "../../[reservationId]/_components/ReservationSummary/HeaderSection";
import { cn } from "@/lib/utils";
import { CalendarInput } from "@/app/(app)/[lang]/_components/home/SearchForm/CalendarInput";
import { useRouter } from "next/navigation";
import {
  buildReservationFiltersQuery,
  RESERVATION_STATUSES,
  ReservationStatus,
  useReservationFilters,
} from "../../_hooks/useReservationFilters";

export function FilterForm() {
  const router = useRouter();
  const { lang, searchParams } = useReservationFilters();
  const t = useTranslations("MyAccount.reservations");

  const [status, setStatus] = useState<ReservationStatus | null>(null);
  const [pickupDate, setPickupDate] = useState<Date | undefined>(undefined);
  const [bookingId, setBookingId] = useState("");
  const [driverName, setDriverName] = useState("");

  useEffect(() => {
    const statusParam = searchParams.get("status");
    setStatus(
      statusParam &&
        RESERVATION_STATUSES.includes(statusParam as ReservationStatus)
        ? (statusParam as ReservationStatus)
        : null,
    );

    const pickupDateParam = searchParams.get("pd");
    if (!pickupDateParam) {
      setPickupDate(undefined);
    } else {
      const parsed = new Date(pickupDateParam);
      setPickupDate(Number.isNaN(parsed.getTime()) ? undefined : parsed);
    }

    setBookingId(searchParams.get("bid") ?? "");
    setDriverName(searchParams.get("dn") ?? "");
  }, [searchParams]);

  function handleSearch(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = buildReservationFiltersQuery({
      status,
      pickupDate,
      bookingId,
      driverName,
    });

    router.push(`/${lang}/reservations?${query.toString()}`);
  }

  return (
    <form
      className="flex gap-4 justify-between items-center"
      onSubmit={handleSearch}
    >
      <legend className="type-label text-navy w-40">{t("filtersLabel")}</legend>
      <div className="flex justify-between gap-4 w-full">
        <Input
          className="bg-white border w-1/4 border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary"
          placeholder={t("bookingIdPlaceholder")}
          value={bookingId}
          onChange={(e) => setBookingId(e.target.value)}
        />
        <StatusDropdown
          value={status}
          onChange={setStatus}
          placeholder={t("statusPlaceholder")}
        />
        <div className="w-1/4">
          <CalendarInput
            placeholder={t("pickupDatePlaceholder")}
            value={pickupDate}
            onSelect={(e) => {
              setPickupDate(e);
            }}
            showIcon={false}
            numberOfMonths={1}
            align="center"
          />
        </div>
        <Input
          className="bg-white border w-1/4 border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary"
          placeholder={t("driverNamePlaceholder")}
          value={driverName}
          onChange={(e) => setDriverName(e.target.value)}
        />
      </div>
      <Button variant="brand" className="py-6 w-40 font-semibold">
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
  value: ReservationStatus | null;
  placeholder: string;
  onChange: (value: ReservationStatus | null) => void;
}) {
  const t = useTranslations("MyAccount.reservation.summary.status");
  const dir = useDirection();

  return (
    <DropdownMenu dir={dir}>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className={cn(
            "w-1/4 flex items-center justify-between bg-white border rounded-lg px-4 h-12 cursor-pointer",
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
        {RESERVATION_STATUSES.map((status) => (
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
