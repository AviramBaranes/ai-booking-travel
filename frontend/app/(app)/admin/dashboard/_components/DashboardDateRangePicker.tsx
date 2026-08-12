"use client";

import { useMemo, useState } from "react";
import { CalendarIcon, Check } from "lucide-react";
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
import { cn } from "@/lib/utils";

import { DateRangeValue, PRESETS, toDateValue } from "../_lib/dateRange";

interface DashboardDateRangePickerProps {
  value: DateRangeValue;
  onChange: (range: DateRangeValue) => void;
}

function parseDate(value: string): Date | undefined {
  if (!value) return undefined;
  const date = new Date(`${value}T00:00:00`);
  return Number.isNaN(date.getTime()) ? undefined : date;
}

function formatDisplayDate(date: Date | undefined): string {
  return date ? format(date, "dd/MM/yyyy") : "";
}

export function DashboardDateRangePicker({
  value,
  onChange,
}: DashboardDateRangePickerProps) {
  const [open, setOpen] = useState(false);
  const [hoverDate, setHoverDate] = useState<Date | undefined>(undefined);

  const selected = useMemo<DateRange>(
    () => ({ from: parseDate(value.from), to: parseDate(value.to) }),
    [value.from, value.to],
  );

  // A preset is "active" when its range is exactly what is selected, so navigating the
  // calendar away from it clears the highlight rather than lying about the selection.
  const activePreset = useMemo(
    () =>
      PRESETS.find((preset) => {
        const range = preset.getRange();
        return range.from === value.from && range.to === value.to;
      })?.id,
    [value.from, value.to],
  );

  const displayValue = selected.from
    ? `${formatDisplayDate(selected.from)} ← ${formatDisplayDate(selected.to ?? selected.from)}`
    : "בחרו טווח תאריכים";

  function handleSelect(range: DateRange | undefined) {
    if (!range?.from) return;

    onChange({
      from: toDateValue(range.from),
      to: toDateValue(range.to ?? range.from),
    });
    setHoverDate(undefined);

    if (range.from && range.to && !isSameDay(range.from, range.to)) {
      setOpen(false);
    }
  }

  return (
    <div className="min-w-64">
      <label className="mb-1 block text-xs text-gray-500">תאריכים</label>
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
            className="flex h-9 w-full items-center justify-between gap-2 rounded-md border border-gray-300 bg-white px-3 text-sm text-gray-700 transition-colors hover:border-navy/40"
          >
            <span className="truncate tabular-nums">{displayValue}</span>
            <CalendarIcon className="size-4 shrink-0 text-gray-400" />
          </button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-auto p-0" dir="rtl">
          <div className="flex flex-col-reverse sm:flex-row-reverse">
            <ul className="flex max-h-72 shrink-0 flex-col gap-0.5 overflow-y-auto border-gray-200 p-2 sm:max-h-none sm:w-44 sm:border-r">
              {PRESETS.map((preset) => {
                const isActive = activePreset === preset.id;
                return (
                  <li key={preset.id}>
                    <button
                      type="button"
                      onClick={() => {
                        onChange(preset.getRange());
                        setOpen(false);
                      }}
                      className={cn(
                        "flex w-full items-center justify-between rounded px-2 py-1.5 text-right text-sm transition-colors",
                        isActive
                          ? "font-semibold text-navy"
                          : "text-gray-700 hover:bg-gray-100",
                      )}
                    >
                      {preset.label}
                      {isActive && <Check className="size-4 shrink-0" />}
                    </button>
                  </li>
                );
              })}
            </ul>
            <div className="p-2">
              <Calendar
                mode="range"
                selected={selected}
                onSelect={handleSelect}
                numberOfMonths={2}
                locale={he}
                captionLayout="dropdown"
                startMonth={new Date(2020, 0)}
                endMonth={new Date(new Date().getFullYear() + 1, 11)}
                defaultMonth={selected.from}
                previewFrom={
                  selected.from && !selected.to ? selected.from : undefined
                }
                previewTo={
                  selected.from && !selected.to ? hoverDate : undefined
                }
                onPreviewDayEnter={(date) => setHoverDate(date)}
                onPreviewDayLeave={() => setHoverDate(undefined)}
              />
              <div className="flex items-center justify-end border-t border-gray-200 pt-2">
                <Button type="button" size="sm" onClick={() => setOpen(false)}>
                  סגור
                </Button>
              </div>
            </div>
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
}
