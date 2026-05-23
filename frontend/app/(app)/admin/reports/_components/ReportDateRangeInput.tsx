"use client";

import { useMemo, useState } from "react";
import { CalendarIcon } from "lucide-react";
import { DateRange } from "react-day-picker";
import { format } from "date-fns/format";
import { isSameDay } from "date-fns/isSameDay";
import { he } from "date-fns/locale/he";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface ReportDateRangeInputProps {
  label: string;
  fromValue: string;
  toValue: string;
  onChange: (range: { from: string; to: string }) => void;
}

function parseDate(value: string): Date | undefined {
  if (!value) return undefined;
  const date = new Date(`${value}T00:00:00`);
  return Number.isNaN(date.getTime()) ? undefined : date;
}

function toDateValue(date: Date | undefined): string {
  if (!date) return "";
  return format(date, "yyyy-MM-dd");
}

function formatDisplayDate(date: Date | undefined): string {
  if (!date) return "";
  return format(date, "d בMMM yyyy", { locale: he });
}

export function ReportDateRangeInput({
  label,
  fromValue,
  toValue,
  onChange,
}: ReportDateRangeInputProps) {
  const [open, setOpen] = useState(false);
  const [hoverDate, setHoverDate] = useState<Date | undefined>(undefined);
  const selected = useMemo<DateRange>(
    () => ({ from: parseDate(fromValue), to: parseDate(toValue) }),
    [fromValue, toValue],
  );

  const displayValue = selected.from
    ? selected.to
      ? `${formatDisplayDate(selected.from)} - ${formatDisplayDate(selected.to)}`
      : formatDisplayDate(selected.from)
    : "הכל";

  function handleSelect(range: DateRange | undefined) {
    onChange({
      from: toDateValue(range?.from),
      to: toDateValue(range?.to),
    });
    setHoverDate(undefined);
    if (range?.from && range.to && !isSameDay(range.from, range.to)) {
      setOpen(false);
    }
  }

  return (
    <div className="min-w-56">
      <label className="block text-xs text-gray-500 mb-1">{label}</label>
      <Popover
        open={open}
        onOpenChange={(nextOpen) => {
          setOpen(nextOpen);
          if (!nextOpen) setHoverDate(undefined);
        }}
      >
        <PopoverTrigger asChild>
          <button
            type="button"
            className="flex h-8 w-full items-center justify-between gap-2 rounded border border-gray-300 bg-white px-2 py-1.5 text-sm text-gray-700"
          >
            <span className="truncate">{displayValue}</span>
            <CalendarIcon className="size-4 shrink-0 text-gray-400" />
          </button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-auto p-2" dir="rtl">
          <Calendar
            mode="range"
            selected={selected}
            onSelect={handleSelect}
            numberOfMonths={2}
            locale={he}
            previewFrom={selected.from && !selected.to ? selected.from : undefined}
            previewTo={selected.from && !selected.to ? hoverDate : undefined}
            onPreviewDayEnter={(date) => setHoverDate(date)}
            onPreviewDayLeave={() => setHoverDate(undefined)}
          />
          <div className="flex items-center justify-between border-t border-gray-200 pt-2">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => onChange({ from: "", to: "" })}
            >
              נקה
            </Button>
            <Button type="button" size="sm" onClick={() => setOpen(false)}>
              סגור
            </Button>
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
}
