"use client";

import { useMemo } from "react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";
import { dashboardReport } from "@/shared/api/reservations-api";
import { cn } from "@/lib/utils";

import { DashboardDateRangePicker } from "./_components/DashboardDateRangePicker";
// import { DashboardHelp } from "./_components/DashboardHelp";
import { SegmentedControl } from "./_components/SegmentedControl";
import { HeroFigure, StatTile } from "./_components/StatTile";
import { TimeSeriesCard, VolumeCard } from "./_components/TimeSeriesCard";
import { EntityBreakdownCard } from "./_components/EntityBreakdownCard";
import {
  CategoryBarCard,
  DonutCard,
  HistogramCard,
  ShareBarCard,
} from "./_components/DistributionCards";
import { PaymentsPanel } from "./_components/PaymentsPanel";
import { CouponPanel } from "./_components/CouponPanel";
import {
  Audience,
  DRIVER_AGE_BANDS,
  LEAD_TIME_BANDS,
  Metric,
  RENTAL_DAYS_BANDS,
  buildLookups,
  buildTimeSeries,
  computePayments,
  computeTotals,
  filterRows,
  groupRows,
  histogram,
  pickGranularity,
} from "./_lib/aggregate";
import {
  BROKER_LABELS,
  GEAR_TYPE_LABELS,
  PAYMENT_STATUS_COLORS,
  PAYMENT_STATUS_LABELS,
  RESERVATION_STATUS_COLORS,
  RESERVATION_STATUS_LABELS,
  count,
  countryLabel,
  decimal,
  ils,
  orEmptyLabel,
  percent,
} from "./_lib/format";
import {
  DateRangeValue,
  defaultDateRange,
  previousPeriod,
} from "./_lib/dateRange";

const FILTER_KEYS = ["from", "to", "audience", "canceled", "metric"] as const;
type FilterKey = (typeof FILTER_KEYS)[number];

const AUDIENCE_OPTIONS = [
  { value: "all", label: "הכל" },
  { value: "business", label: "עסקי" },
  { value: "private", label: "פרטי" },
];

const METRIC_OPTIONS = [
  { value: "count", label: "לפי כמות" },
  { value: "revenue", label: "לפי הכנסות" },
  { value: "profit", label: "לפי רווח" },
];

export default function DashboardShell() {
  const [urlFilters, setUrlFilters] = useUrlFilters<FilterKey>([...FILTER_KEYS]);

  const fallbackRange = useMemo(() => defaultDateRange(), []);
  const range = useMemo<DateRangeValue>(
    () => ({
      from: urlFilters.from || fallbackRange.from,
      to: urlFilters.to || fallbackRange.to,
    }),
    [urlFilters.from, urlFilters.to, fallbackRange],
  );
  const audience = (urlFilters.audience || "all") as Audience;
  const includeCanceled = urlFilters.canceled === "true";
  const metric = (urlFilters.metric || "count") as Metric;

  const current = useQuery({
    queryKey: ["admin-dashboard", range.from, range.to],
    queryFn: () => dashboardReport({ From: range.from, To: range.to }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60_000,
  });

  // The equal-length window immediately before, so every headline number carries a delta.
  const comparison = useMemo(() => previousPeriod(range), [range]);
  const previous = useQuery({
    queryKey: ["admin-dashboard", comparison.from, comparison.to],
    queryFn: () => dashboardReport({ From: comparison.from, To: comparison.to }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60_000,
  });

  const lookups = useMemo(
    () =>
      current.data
        ? buildLookups(current.data)
        : { organizations: new Map(), offices: new Map(), users: new Map() },
    [current.data],
  );

  const rows = useMemo(
    () =>
      filterRows(current.data?.reservations ?? [], { audience, includeCanceled }),
    [current.data, audience, includeCanceled],
  );
  const previousRows = useMemo(
    () =>
      filterRows(previous.data?.reservations ?? [], { audience, includeCanceled }),
    [previous.data, audience, includeCanceled],
  );

  const totals = useMemo(() => computeTotals(rows), [rows]);
  const previousTotals = useMemo(() => computeTotals(previousRows), [previousRows]);
  const payments = useMemo(() => computePayments(rows), [rows]);

  const granularity = useMemo(() => pickGranularity(range), [range]);
  const buckets = useMemo(
    () => buildTimeSeries(rows, range, granularity),
    [rows, range, granularity],
  );

  // Every "all rows" grouping below reads the same filtered slice, so the whole page
  // always describes one consistent population.
  const statusGroups = useMemo(
    () =>
      groupRows(rows, (row) => row.status, (key) => RESERVATION_STATUS_LABELS[key] ?? key, metric),
    [rows, metric],
  );
  const paymentStatusGroups = useMemo(
    () =>
      groupRows(rows, (row) => row.paymentStatus, (key) => PAYMENT_STATUS_LABELS[key] ?? key, metric),
    [rows, metric],
  );
  const gearGroups = useMemo(
    () => groupRows(rows, (row) => row.gearType, (key) => GEAR_TYPE_LABELS[key] ?? key, metric),
    [rows, metric],
  );
  const brokerGroups = useMemo(
    () => groupRows(rows, (row) => row.broker, (key) => BROKER_LABELS[key] ?? key, metric),
    [rows, metric],
  );
  const erpGroups = useMemo(
    () =>
      groupRows(
        rows,
        (row) => (row.hasErp ? "erp" : "no-erp"),
        (key) => (key === "erp" ? "עם ERP" : "בלי ERP"),
        metric,
      ),
    [rows, metric],
  );
  const audienceGroups = useMemo(
    () =>
      groupRows(
        rows,
        (row) => (row.isBusiness ? "business" : "private"),
        (key) => (key === "business" ? "עסקי" : "פרטי"),
        metric,
      ),
    [rows, metric],
  );
  // These go in complete; each card folds its own tail and can be expanded in place.
  const countryGroups = useMemo(
    () => groupRows(rows, (row) => row.countryCode, countryLabel, metric),
    [rows, metric],
  );
  const supplierGroups = useMemo(
    () =>
      groupRows(rows, (row) => orEmptyLabel(row.supplierName), (key) => key, metric),
    [rows, metric],
  );
  const carTypeGroups = useMemo(
    () => groupRows(rows, (row) => orEmptyLabel(row.carType), (key) => key, metric),
    [rows, metric],
  );

  const isInitialLoading = current.isLoading && !current.isPlaceholderData;
  // Hold the previous render at reduced opacity on refetch — no skeleton flash, no jump.
  const isRefetching = current.isFetching && !isInitialLoading;

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="type-h4 text-navy">לוח בקרה</h1>
          <p className="type-paragraph text-text-secondary">
            כל הנתונים מחושבים על הזמנות שנוצרו בטווח שנבחר, בשקלים.
          </p>
        </div>
        {/* <DashboardHelp /> */}
      </header>

      {/* One filter row above everything it scopes. */}
      <div
        data-noprint
        className="flex flex-wrap items-end gap-3 rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
      >
        <DashboardDateRangePicker
          value={range}
          onChange={(next) => setUrlFilters({ from: next.from, to: next.to })}
        />
        <FilterField label="קהל">
          <SegmentedControl
            value={audience}
            onChange={(next) => setUrlFilters({ audience: next })}
            options={AUDIENCE_OPTIONS}
          />
        </FilterField>
        <FilterField label="ביטולים">
          <SegmentedControl
            value={includeCanceled ? "true" : "false"}
            onChange={(next) => setUrlFilters({ canceled: next })}
            options={[
              { value: "false", label: "ללא מבוטלות" },
              { value: "true", label: "כולל מבוטלות" },
            ]}
          />
        </FilterField>
        <FilterField label="פילוחים">
          <SegmentedControl
            value={metric}
            onChange={(next) => setUrlFilters({ metric: next })}
            options={METRIC_OPTIONS}
          />
        </FilterField>
      </div>

      {current.isError && (
        <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          לא הצלחנו לטעון את הנתונים. נסו שוב בעוד רגע.
        </div>
      )}

      {isInitialLoading ? (
        <LoadingSkeleton />
      ) : (
        <div
          className={cn(
            "space-y-6 transition-opacity",
            isRefetching && "opacity-60",
          )}
        >
          <section className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,2fr)]">
            <HeroFigure
              label="סך הרווח"
              value={ils(totals.profit)}
              current={totals.profit}
              previous={previousTotals.profit}
              hint={`${percent(totals.marginPct)} שיעור רווח`}
            />
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              <StatTile
                label="הזמנות"
                value={count(totals.count)}
                current={totals.count}
                previous={previousTotals.count}
              />
              <StatTile
                label="הכנסות"
                value={ils(totals.revenue)}
                current={totals.revenue}
                previous={previousTotals.revenue}
              />
              <StatTile
                label="עלות"
                value={ils(totals.cost)}
                current={totals.cost}
                previous={previousTotals.cost}
                invertDelta
              />
              <StatTile
                label="שיעור ביטולים"
                value={percent(totals.cancellationRate)}
                current={totals.cancellationRate}
                previous={previousTotals.cancellationRate}
                invertDelta
                hint={`${count(totals.canceled)} מבוטלות`}
              />
            </div>
          </section>

          <section className="grid gap-4 xl:grid-cols-[minmax(0,2fr)_minmax(0,1fr)]">
            <TimeSeriesCard buckets={buckets} granularity={granularity} />
            <VolumeCard buckets={buckets} granularity={granularity} />
          </section>

          <EntityBreakdownCard rows={rows} lookups={lookups} />

          <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <DonutCard
              title="סטטוס הזמנה"
              groups={statusGroups}
              metric={metric}
              colors={statusColorsByLabel(RESERVATION_STATUS_COLORS)}
            />
            <DonutCard
              title="סטטוס תשלום"
              groups={paymentStatusGroups}
              metric={metric}
              colors={statusColorsByLabel(PAYMENT_STATUS_COLORS)}
            />
            <DonutCard
              title="תיבת הילוכים"
              groups={gearGroups}
              metric={metric}
              colors={{
                auto: "var(--chart-1)",
                manual: "var(--chart-2)",
                electric: "var(--chart-3)",
              }}
            />
            <ShareBarCard title="ברוקר" groups={brokerGroups} metric={metric} />
            <ShareBarCard title="ERP" groups={erpGroups} metric={metric} />
            <ShareBarCard title="עסקי מול פרטי" groups={audienceGroups} metric={metric} />
            <CategoryBarCard title="מדינות" groups={countryGroups} metric={metric} />
            <CategoryBarCard title="מותג ספק" groups={supplierGroups} metric={metric} />
            <CategoryBarCard title="סוג רכב" groups={carTypeGroups} metric={metric} />
          </section>

          <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <StatTile
              label="ימי השכרה ממוצע"
              value={decimal(totals.avgRentalDays)}
              current={totals.avgRentalDays}
              previous={previousTotals.avgRentalDays}
            />
            <StatTile
              label="גיל נהג ממוצע"
              value={decimal(totals.avgDriverAge, 0)}
              current={totals.avgDriverAge}
              previous={previousTotals.avgDriverAge}
            />
            <StatTile
              label="ערך הזמנה ממוצע"
              value={ils(totals.avgOrderValue)}
              current={totals.avgOrderValue}
              previous={previousTotals.avgOrderValue}
            />
            <StatTile
              label="רווח ממוצע להזמנה"
              value={ils(totals.avgProfit)}
              current={totals.avgProfit}
              previous={previousTotals.avgProfit}
            />
            <StatTile
              label="עלות ממוצעת להזמנה"
              value={ils(totals.avgCost)}
              current={totals.avgCost}
              previous={previousTotals.avgCost}
              invertDelta
            />
            <StatTile
              label="מחיר ERP ממוצע"
              value={ils(totals.avgErpPrice)}
              hint={`${count(totals.withErp)} הזמנות עם ERP`}
            />
            <StatTile
              label="הזמנה מראש ממוצעת"
              value={`${decimal(totals.avgLeadTimeDays, 0)} ימים`}
              current={totals.avgLeadTimeDays}
              previous={previousTotals.avgLeadTimeDays}
            />
            <StatTile
              label="הזמנות עם קופון"
              value={percent(totals.couponRate, 0)}
              hint={`${count(totals.withCoupon)} הזמנות`}
            />
          </section>

          <section className="grid gap-4 md:grid-cols-3">
            <HistogramCard
              title="גיל נהג"
              subtitle="התפלגות לפי טווחי גיל"
              bars={histogram(rows, DRIVER_AGE_BANDS, (row) => row.driverAge)}
            />
            <HistogramCard
              title="ימי השכרה"
              subtitle="התפלגות לפי אורך ההשכרה"
              bars={histogram(rows, RENTAL_DAYS_BANDS, (row) => row.rentalDays)}
            />
            <HistogramCard
              title="הזמנה מראש"
              subtitle="ימים בין ההזמנה לאיסוף"
              bars={histogram(rows, LEAD_TIME_BANDS, (row) => row.leadTimeDays)}
            />
          </section>

          <PaymentsPanel payments={payments} />

          <CouponPanel rows={rows} />
        </div>
      )}
    </div>
  );
}

function FilterField({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <span className="mb-1 block text-xs text-gray-500">{label}</span>
      {children}
    </div>
  );
}

/**
 * Groups are keyed by the raw status value but labelled in Hebrew; the colour map is keyed
 * the same way so a status never changes colour when the label does.
 */
function statusColorsByLabel(colors: Record<string, string>): Record<string, string> {
  return colors;
}

function LoadingSkeleton() {
  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div
            key={index}
            className="h-24 animate-pulse rounded-lg border border-gray-200 bg-white"
          />
        ))}
      </div>
      <div className="h-72 animate-pulse rounded-lg border border-gray-200 bg-white" />
      <div className="grid gap-4 md:grid-cols-3">
        {Array.from({ length: 3 }).map((_, index) => (
          <div
            key={index}
            className="h-56 animate-pulse rounded-lg border border-gray-200 bg-white"
          />
        ))}
      </div>
    </div>
  );
}
