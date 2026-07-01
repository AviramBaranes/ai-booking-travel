"use client";

import { DayPickerLocale, he } from "react-day-picker/locale";
import { useEffect, useState } from "react";
import { DateRange } from "react-day-picker";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { X, Calendar as CalendarIcon } from "lucide-react";
import {
  Sheet,
  SheetContent,
  SheetClose,
  SheetTitle,
} from "@/components/ui/sheet";
import { Calendar } from "@/components/ui/calendar";
import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { FieldError } from "react-hook-form";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";

interface CalendarInputTriggerProps {
  placeholder: string;
  value?: Date;
  error?: FieldError;
  onClick: () => void;
}

// Display-only field used as a trigger to open the shared CalendarSheet.
export function CalendarInputTrigger({
  placeholder,
  value,
  error,
  onClick,
}: CalendarInputTriggerProps) {
  const { lang } = useParams();
  const locale = lang === "he" ? he : undefined;
  const displayValue = value ? value.toLocaleDateString(locale?.code) : "";

  return (
    <div className="flex flex-col items-start">
      <div onClick={onClick} className="w-full cursor-pointer">
        <Field>
          <InputGroup className="search-form-input px-0">
            <InputGroupInput
              aria-invalid={!!error}
              value={displayValue}
              placeholder={placeholder}
              className="text-start px-2 text-sm pointer-events-none"
              readOnly
              tabIndex={-1}
            />
            <InputGroupAddon align="inline-start" className="pl-1 pr-0">
              <CalendarIcon className="size-5 mr-2 text-brand" />
            </InputGroupAddon>
          </InputGroup>
        </Field>
      </div>
      <ErrorDisplay>{error?.message}</ErrorDisplay>
    </div>
  );
}

interface CalendarSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  value?: DateRange;
  onConfirm: (range: DateRange | undefined) => void;
}

const formatRangeDate = (locale: DayPickerLocale | undefined, date: Date) => {
  const parts = new Intl.DateTimeFormat(locale?.code ?? "en-GB", {
    weekday: "short",
    day: "numeric",
    month: "short",
  }).formatToParts(date);

  const get = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((p) => p.type === type)?.value ?? "";

  return `${get("weekday")} ${get("day")} ${get("month")}`;
};

export function CalendarSheet({
  open,
  onOpenChange,
  value,
  onConfirm,
}: CalendarSheetProps) {
  const { lang } = useParams();
  const locale = lang === "he" ? he : undefined;
  const t = useTranslations("SearchForm");
  const [range, setRange] = useState<DateRange | undefined>(value);

  // Sync the local selection with the current form values whenever the sheet opens.
  useEffect(() => {
    if (open) setRange(value);
  }, [open, value]);

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const maxDate = new Date(Date.now() + 365 * 24 * 60 * 60 * 1000);

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="top"
        showCloseButton={false}
        className="p-0 rounded-none border-0 w-full bottom-0 flex flex-col"
      >
        <SheetTitle className="sr-only">{t("selectDates")}</SheetTitle>
        {/* Top bar */}
        <div className="flex h-15 items-center justify-between px-4 border-b border-border-light shrink-0">
          <SheetClose asChild>
            <button aria-label="Close menu">
              <X className="size-5 text-navy" />
            </button>
          </SheetClose>
        </div>

        {/* Two-month range calendar stacked in a column */}
        <div className="flex-1 overflow-y-scroll flex justify-center">
          <Calendar
            mode="range"
            showOutsideDays={false}
            locale={locale}
            numberOfMonths={2}
            selected={range}
            onSelect={setRange}
            className="min-h-75 w-full px-8 bg-white"
            classNames={{
              months: "relative flex flex-col gap-4 items-center",
              today: "bg-brand/35",
              day_button:
                "text-navy data-[selected-single=true]:bg-brand data-[selected-single=true]:text-white",
            }}
            disabled={(date) => date < today || date > maxDate}
          />
        </div>

        {/* Confirm button fixed to the bottom */}
        <div className="p-4 border-t border-border-light shrink-0">
          <p className="type-paragraph text-center">
            {range?.from && range?.to
              ? `${formatRangeDate(locale, range.from)} - ${formatRangeDate(locale, range.to)}`
              : t("selectDates")}
          </p>
          <Button
            type="button"
            variant="brand"
            className="w-full py-6 my-2 type-paragraph font-bold cursor-pointer"
            onClick={() => {
              onConfirm(range);
              onOpenChange(false);
            }}
          >
            {t("confirmDates")}
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  );
}
