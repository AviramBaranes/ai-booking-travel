"use client";

import * as React from "react";
import {
  DayPicker,
  getDefaultClassNames,
  type ChevronProps,
  type DayButton,
  type Locale,
  type RootProps,
  type WeekNumberProps,
} from "react-day-picker";

import { cn } from "@/lib/utils";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
} from "lucide-react";
import { isAfter } from "date-fns/isAfter";
import { isBefore } from "date-fns/isBefore";
import { isSameDay } from "date-fns/isSameDay";

const defaultClassNames = getDefaultClassNames();

type CalendarPreview = {
  locale?: Partial<Locale>;
  previewFrom?: Date;
  previewTo?: Date;
  onPreviewDayEnter?: (date: Date) => void;
  onPreviewDayLeave?: () => void;
};

/**
 * The hover preview reaches the custom components through context instead of
 * through props on components declared inside `Calendar`.
 *
 * react-day-picker keeps `components` in a memo and renders each day with
 * `createElement(components.DayButton)`. A component declared during render
 * gets a fresh identity every time the preview moves, which React treats as a
 * different element type: it unmounts every day button and mounts a new one.
 * The button that received `mousedown` is then gone before `mouseup`, so the
 * browser never fires `click` and the day looks enabled but cannot be picked.
 */
const CalendarPreviewContext = React.createContext<CalendarPreview>({});

function CalendarRoot({
  className,
  rootRef,
  onMouseLeave,
  ...props
}: RootProps) {
  const { onPreviewDayLeave } = React.useContext(CalendarPreviewContext);

  return (
    <div
      data-slot="calendar"
      ref={rootRef}
      className={cn(className)}
      onMouseLeave={(e) => {
        onMouseLeave?.(e);
        onPreviewDayLeave?.();
      }}
      {...props}
    />
  );
}

function CalendarChevron({ className, orientation, ...props }: ChevronProps) {
  if (orientation === "left") {
    return <ChevronLeftIcon className={cn("size-4", className)} {...props} />;
  }

  if (orientation === "right") {
    return <ChevronRightIcon className={cn("size-4", className)} {...props} />;
  }

  return <ChevronDownIcon className={cn("size-4", className)} {...props} />;
}

function CalendarWeekNumber({ children, week, ...props }: WeekNumberProps) {
  return (
    <td {...props}>
      <div className="flex size-(--cell-size) items-center justify-center text-center">
        {children}
      </div>
    </td>
  );
}

function CalendarDayButtonSlot(props: React.ComponentProps<typeof DayButton>) {
  const preview = React.useContext(CalendarPreviewContext);

  return <CalendarDayButton {...preview} {...props} />;
}

const calendarComponents = {
  Root: CalendarRoot,
  Chevron: CalendarChevron,
  DayButton: CalendarDayButtonSlot,
  WeekNumber: CalendarWeekNumber,
};

function Calendar({
  className,
  classNames,
  showOutsideDays = true,
  captionLayout = "label",
  buttonVariant = "ghost",
  locale,
  formatters,
  components,
  previewFrom,
  previewTo,
  onPreviewDayEnter,
  onPreviewDayLeave,
  ...props
}: React.ComponentProps<typeof DayPicker> & {
  buttonVariant?: React.ComponentProps<typeof Button>["variant"];
  previewFrom?: Date;
  previewTo?: Date;
  onPreviewDayEnter?: (date: Date) => void;
  onPreviewDayLeave?: () => void;
}) {
  const showWeekNumber = props.showWeekNumber;

  // Every object below has to keep a stable identity across renders, or
  // react-day-picker rebuilds the calendar instead of updating it.
  const mergedComponents = React.useMemo(
    () =>
      components ? { ...calendarComponents, ...components } : calendarComponents,
    [components],
  );

  const mergedFormatters = React.useMemo(
    () => ({
      formatMonthDropdown: (date: Date) =>
        date.toLocaleString(locale?.code, { month: "short" }),
      ...formatters,
    }),
    [locale?.code, formatters],
  );

  const preview = React.useMemo(
    () => ({
      locale,
      previewFrom,
      previewTo,
      onPreviewDayEnter,
      onPreviewDayLeave,
    }),
    [locale, previewFrom, previewTo, onPreviewDayEnter, onPreviewDayLeave],
  );

  const mergedClassNames = React.useMemo(
    () => ({
      root: cn("w-fit", defaultClassNames.root),
      months: cn(
        "relative flex flex-col gap-4 md:flex-row",
        defaultClassNames.months,
      ),
      month: cn("flex w-full flex-col gap-4", defaultClassNames.month),
      nav: cn(
        "absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1",
        defaultClassNames.nav,
      ),
      button_previous: cn(
        buttonVariants({ variant: buttonVariant }),
        "size-(--cell-size) p-0 select-none aria-disabled:opacity-50",
        defaultClassNames.button_previous,
      ),
      button_next: cn(
        buttonVariants({ variant: buttonVariant }),
        "size-(--cell-size) p-0 select-none aria-disabled:opacity-50",
        defaultClassNames.button_next,
      ),
      month_caption: cn(
        "flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)",
        defaultClassNames.month_caption,
      ),
      dropdowns: cn(
        "flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium",
        defaultClassNames.dropdowns,
      ),
      dropdown_root: cn(
        "relative rounded-(--cell-radius)",
        defaultClassNames.dropdown_root,
      ),
      dropdown: cn(
        "absolute inset-0 bg-popover opacity-0",
        defaultClassNames.dropdown,
      ),
      caption_label: cn(
        "font-medium select-none",
        captionLayout === "label"
          ? "text-sm"
          : "flex items-center gap-1 rounded-(--cell-radius) text-sm [&>svg]:size-3.5 [&>svg]:text-muted-foreground",
        defaultClassNames.caption_label,
      ),
      table: "w-full border-collapse",
      weekdays: cn("flex", defaultClassNames.weekdays),
      weekday: cn(
        "flex-1 rounded-(--cell-radius) text-[0.8rem] font-normal text-muted-foreground select-none",
        defaultClassNames.weekday,
      ),
      week: cn("mt-2 flex w-full", defaultClassNames.week),
      week_number_header: cn(
        "w-(--cell-size) select-none",
        defaultClassNames.week_number_header,
      ),
      week_number: cn(
        "text-[0.8rem] text-muted-foreground select-none",
        defaultClassNames.week_number,
      ),
      day: cn(
        "group/day relative aspect-square h-full w-full rounded-(--cell-radius) p-0 text-center select-none ltr:[&:last-child[data-selected=true]_button]:rounded-r-(--cell-radius) rtl:[&:last-child[data-selected=true]_button]:rounded-l-(--cell-radius)",
        showWeekNumber
          ? "ltr:[&:nth-child(2)[data-selected=true]_button]:rounded-l-(--cell-radius) rtl:[&:nth-child(2)[data-selected=true]_button]:rounded-r-(--cell-radius)"
          : "ltr:[&:first-child[data-selected=true]_button]:rounded-l-(--cell-radius) rtl:[&:first-child[data-selected=true]_button]:rounded-r-(--cell-radius)",
        defaultClassNames.day,
      ),
      range_middle: cn("rounded-none", defaultClassNames.range_middle),
      range_start: cn(
        "relative isolate z-0 ltr:rounded-l-(--cell-radius) rtl:rounded-r-(--cell-radius) bg-brand/55 after:absolute after:inset-y-0 ltr:after:right-0 rtl:after:left-0 after:w-4 after:bg-brand/55",
        defaultClassNames.range_start,
      ),
      range_end: cn(
        "relative isolate z-0 ltr:rounded-r-(--cell-radius) rtl:rounded-l-(--cell-radius) bg-brand/55 after:absolute after:inset-y-0 ltr:after:left-0 rtl:after:right-0 after:w-4 after:bg-brand/55",
        defaultClassNames.range_end,
      ),
      today: cn(
        "rounded-(--cell-radius) bg-muted text-foreground data-[selected=true]:rounded-none",
        defaultClassNames.today,
      ),
      outside: cn(
        "text-muted-foreground aria-selected:text-muted-foreground",
        defaultClassNames.outside,
      ),
      disabled: cn(
        "text-muted-foreground opacity-50",
        defaultClassNames.disabled,
      ),
      hidden: cn("invisible", defaultClassNames.hidden),
      ...classNames,
    }),
    [classNames, captionLayout, buttonVariant, showWeekNumber],
  );

  return (
    <CalendarPreviewContext.Provider value={preview}>
      <DayPicker
        showOutsideDays={showOutsideDays}
        className={cn(
          "group/calendar bg-background p-2 [--cell-radius:var(--radius-md)] [--cell-size:--spacing(7)] in-data-[slot=card-content]:bg-transparent in-data-[slot=popover-content]:bg-transparent",
          String.raw`rtl:**:[.rdp-button\_next>svg]:rotate-180`,
          String.raw`rtl:**:[.rdp-button\_previous>svg]:rotate-180`,
          className,
        )}
        captionLayout={captionLayout}
        locale={locale}
        formatters={mergedFormatters}
        classNames={mergedClassNames}
        components={mergedComponents}
        {...props}
      />
    </CalendarPreviewContext.Provider>
  );
}

function CalendarDayButton({
  className,
  day,
  modifiers,
  locale,
  previewFrom,
  previewTo,
  onPreviewDayEnter,
  onPreviewDayLeave,
  onMouseEnter,
  onMouseLeave,
  ...props
}: React.ComponentProps<typeof DayButton> & {
  locale?: Partial<Locale>;
  previewFrom?: Date;
  previewTo?: Date;
  onPreviewDayEnter?: (date: Date) => void;
  onPreviewDayLeave?: () => void;
}) {
  const ref = React.useRef<HTMLButtonElement>(null);
  React.useEffect(() => {
    if (modifiers.focused) ref.current?.focus();
  }, [modifiers.focused]);

  let previewStart = false;
  let previewMiddle = false;
  let previewEnd = false;

  if (previewFrom && previewTo) {
    const start = isAfter(previewFrom, previewTo) ? previewTo : previewFrom;
    const end = isAfter(previewFrom, previewTo) ? previewFrom : previewTo;
    const current = day.date;

    previewStart = isSameDay(current, start);
    previewEnd = isSameDay(current, end);
    previewMiddle =
      !previewStart &&
      !previewEnd &&
      !isBefore(current, start) &&
      !isAfter(current, end);
  }

  const previewOnlyStart = previewStart && !modifiers.range_start;
  const previewOnlyEnd = previewEnd && !modifiers.range_end;

  return (
    <Button
      ref={ref}
      type="button"
      variant="ghost"
      size="icon"
      data-day={day.date.toLocaleDateString(locale?.code)}
      data-selected-single={
        modifiers.selected &&
        !modifiers.range_start &&
        !modifiers.range_end &&
        !modifiers.range_middle
      }
      data-range-start={modifiers.range_start || previewStart}
      data-range-end={modifiers.range_end || previewEnd}
      data-range-middle={modifiers.range_middle || previewMiddle}
      data-preview-start={previewOnlyStart}
      data-preview-end={previewOnlyEnd}
      onMouseEnter={(e) => {
        onMouseEnter?.(e);
        onPreviewDayEnter?.(day.date);
      }}
      className={cn(
        "relative isolate z-10 flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 border-0 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-(--cell-radius) ltr:data-[range-end=true]:rounded-r-(--cell-radius) rtl:data-[range-end=true]:rounded-l-(--cell-radius) data-[range-end=true]:bg-brand data-[range-end=true]:text-white data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-brand/55 data-[range-middle=true]:text-white data-[range-start=true]:rounded-(--cell-radius) ltr:data-[range-start=true]:rounded-l-(--cell-radius) rtl:data-[range-start=true]:rounded-r-(--cell-radius) data-[range-start=true]:bg-brand data-[range-start=true]:text-white data-[selected-single=true]:bg-brand data-[selected-single=true]:text-white data-[preview-start=true]:before:absolute data-[preview-start=true]:before:inset-y-0 ltr:data-[preview-start=true]:before:right-0 rtl:data-[preview-start=true]:before:left-0 data-[preview-start=true]:before:w-1/2 data-[preview-start=true]:before:bg-brand/55 data-[preview-start=true]:before:-z-10 data-[preview-end=true]:before:absolute data-[preview-end=true]:before:inset-y-0 ltr:data-[preview-end=true]:before:left-0 rtl:data-[preview-end=true]:before:right-0 data-[preview-end=true]:before:w-1/2 data-[preview-end=true]:before:bg-brand/55 data-[preview-end=true]:before:-z-10 dark:hover:text-foreground [&>span]:text-xs [&>span]:opacity-70",
        defaultClassNames.day,
        className,
      )}
      {...props}
    />
  );
}

export { Calendar, CalendarDayButton };
