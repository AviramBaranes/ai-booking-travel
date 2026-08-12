import { format } from "date-fns/format";
import { startOfWeek } from "date-fns/startOfWeek";
import { startOfMonth } from "date-fns/startOfMonth";
import { endOfMonth } from "date-fns/endOfMonth";
import { startOfYear } from "date-fns/startOfYear";
import { endOfYear } from "date-fns/endOfYear";
import { subDays } from "date-fns/subDays";
import { subMonths } from "date-fns/subMonths";
import { subYears } from "date-fns/subYears";
import { differenceInCalendarDays } from "date-fns/differenceInCalendarDays";

export interface DateRangeValue {
  from: string;
  to: string;
}

/** The backend reads from/to as calendar dates, so ranges are plain yyyy-MM-dd strings. */
export function toDateValue(date: Date): string {
  return format(date, "yyyy-MM-dd");
}

export function parseDateValue(value: string): Date {
  return new Date(`${value}T00:00:00`);
}

/**
 * "הכל" needs a concrete lower bound because the endpoint always takes a window. The
 * business has no reservations before this, so it reads as all-time without a special case.
 */
const ALL_TIME_START = new Date(2020, 0, 1);

export interface DateRangePreset {
  id: string;
  label: string;
  getRange: () => DateRangeValue;
}

function range(from: Date, to: Date): DateRangeValue {
  return { from: toDateValue(from), to: toDateValue(to) };
}

export const PRESETS: DateRangePreset[] = [
  {
    id: "today",
    label: "היום",
    getRange: () => range(new Date(), new Date()),
  },
  {
    id: "yesterday",
    label: "אתמול",
    getRange: () => range(subDays(new Date(), 1), subDays(new Date(), 1)),
  },
  {
    id: "last7",
    label: "7 ימים אחרונים",
    getRange: () => range(subDays(new Date(), 6), new Date()),
  },
  {
    id: "last30",
    label: "30 ימים אחרונים",
    getRange: () => range(subDays(new Date(), 29), new Date()),
  },
  {
    id: "thisWeek",
    label: "השבוע",
    getRange: () => range(startOfWeek(new Date(), { weekStartsOn: 0 }), new Date()),
  },
  {
    id: "thisMonth",
    label: "החודש",
    getRange: () => range(startOfMonth(new Date()), new Date()),
  },
  {
    id: "lastMonth",
    label: "החודש הקודם",
    getRange: () => {
      const previous = subMonths(new Date(), 1);
      return range(startOfMonth(previous), endOfMonth(previous));
    },
  },
  {
    id: "last3Months",
    label: "3 חודשים אחרונים",
    getRange: () => range(startOfMonth(subMonths(new Date(), 2)), new Date()),
  },
  {
    id: "yearToDate",
    label: "מתחילת השנה",
    getRange: () => range(startOfYear(new Date()), new Date()),
  },
  {
    id: "lastYear",
    label: "שנה קודמת",
    getRange: () => {
      const previous = subYears(new Date(), 1);
      return range(startOfYear(previous), endOfYear(previous));
    },
  },
  {
    id: "all",
    label: "הכל",
    getRange: () => range(ALL_TIME_START, new Date()),
  },
];

/** The dashboard opens on the current week. */
export function defaultDateRange(): DateRangeValue {
  return range(startOfWeek(new Date(), { weekStartsOn: 0 }), new Date());
}

/**
 * The window of the same length immediately before the selected one, used to show how
 * every headline number moved.
 */
export function previousPeriod(value: DateRangeValue): DateRangeValue {
  const from = parseDateValue(value.from);
  const to = parseDateValue(value.to);
  const days = differenceInCalendarDays(to, from) + 1;

  return range(subDays(from, days), subDays(to, days));
}

export function rangeLengthInDays(value: DateRangeValue): number {
  return (
    differenceInCalendarDays(parseDateValue(value.to), parseDateValue(value.from)) + 1
  );
}
