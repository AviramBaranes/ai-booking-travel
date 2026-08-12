import { format } from "date-fns/format";
import { addDays } from "date-fns/addDays";
import { addMonths } from "date-fns/addMonths";
import { startOfWeek } from "date-fns/startOfWeek";
import { startOfMonth } from "date-fns/startOfMonth";
import { he } from "date-fns/locale/he";
import { accounts, reports } from "@/shared/client";

import { DateRangeValue, parseDateValue, rangeLengthInDays } from "./dateRange";

export type DashboardRow = reports.DashboardReservation;

export type Audience = "all" | "business" | "private";

export interface RowFilters {
  audience: Audience;
  includeCanceled: boolean;
}

export const CANCELED = "canceled";

/**
 * The one place rows are narrowed. Every panel reads the same filtered slice, so the whole
 * page always describes a single, consistent population.
 */
export function filterRows(
  rows: DashboardRow[],
  { audience, includeCanceled }: RowFilters,
): DashboardRow[] {
  return rows.filter((row) => {
    if (!includeCanceled && row.status === CANCELED) return false;
    if (audience === "business" && !row.isBusiness) return false;
    if (audience === "private" && row.isBusiness) return false;
    return true;
  });
}

// --- Totals -----------------------------------------------------------------

export interface Totals {
  count: number;
  revenue: number;
  cost: number;
  profit: number;
  erpRevenue: number;
  erpCost: number;
  discount: number;
  marginPct: number;
  canceled: number;
  cancellationRate: number;
  withErp: number;
  withCoupon: number;
  couponRate: number;
  penaltyAmount: number;
  avgRentalDays: number;
  avgDriverAge: number;
  avgLeadTimeDays: number;
  avgOrderValue: number;
  avgProfit: number;
  avgCost: number;
  /** Averaged over ERP bookings only — including the zeros would understate the product. */
  avgErpPrice: number;
}

function mean(total: number, count: number): number {
  return count > 0 ? total / count : 0;
}

export function computeTotals(rows: DashboardRow[]): Totals {
  let revenue = 0;
  let cost = 0;
  let profit = 0;
  let erpRevenue = 0;
  let erpCost = 0;
  let discount = 0;
  let penaltyAmount = 0;
  let canceled = 0;
  let withErp = 0;
  let withCoupon = 0;
  let rentalDays = 0;
  let driverAge = 0;
  let leadTimeDays = 0;

  for (const row of rows) {
    revenue += row.revenueIls;
    cost += row.costIls;
    profit += row.profitIls;
    erpRevenue += row.erpRevenueIls;
    erpCost += row.erpCostIls;
    discount += row.discountIls;
    penaltyAmount += row.penaltyAmountIls;
    rentalDays += row.rentalDays;
    driverAge += row.driverAge;
    leadTimeDays += row.leadTimeDays;
    if (row.status === CANCELED) canceled += 1;
    if (row.hasErp) withErp += 1;
    if (row.couponName) withCoupon += 1;
  }

  const count = rows.length;

  return {
    count,
    revenue,
    cost,
    profit,
    erpRevenue,
    erpCost,
    discount,
    marginPct: revenue > 0 ? (profit / revenue) * 100 : 0,
    canceled,
    cancellationRate: count > 0 ? (canceled / count) * 100 : 0,
    withErp,
    withCoupon,
    couponRate: count > 0 ? (withCoupon / count) * 100 : 0,
    penaltyAmount,
    avgRentalDays: mean(rentalDays, count),
    avgDriverAge: mean(driverAge, count),
    avgLeadTimeDays: mean(leadTimeDays, count),
    avgOrderValue: mean(revenue, count),
    avgProfit: mean(profit, count),
    avgCost: mean(cost, count),
    avgErpPrice: mean(erpRevenue, withErp),
  };
}

// --- Time series ------------------------------------------------------------

export type Granularity = "day" | "week" | "month";

export function pickGranularity(range: DateRangeValue): Granularity {
  const days = rangeLengthInDays(range);
  if (days <= 31) return "day";
  if (days <= 120) return "week";
  return "month";
}

export interface TimeBucket {
  key: string;
  label: string;
  count: number;
  revenue: number;
  cost: number;
  profit: number;
  business: number;
  private: number;
}

function bucketStart(date: Date, granularity: Granularity): Date {
  if (granularity === "week") return startOfWeek(date, { weekStartsOn: 0 });
  if (granularity === "month") return startOfMonth(date);
  return date;
}

function bucketKey(date: Date, granularity: Granularity): string {
  return format(bucketStart(date, granularity), "yyyy-MM-dd");
}

function bucketLabel(key: string, granularity: Granularity): string {
  const date = parseDateValue(key);
  if (granularity === "month") return format(date, "MMM yy", { locale: he });
  return format(date, "d MMM", { locale: he });
}

function emptyBucket(key: string, granularity: Granularity): TimeBucket {
  return {
    key,
    label: bucketLabel(key, granularity),
    count: 0,
    revenue: 0,
    cost: 0,
    profit: 0,
    business: 0,
    private: 0,
  };
}

/**
 * Buckets rows over the full range, including empty periods — a line that silently skips
 * quiet days misrepresents the trend.
 */
export function buildTimeSeries(
  rows: DashboardRow[],
  range: DateRangeValue,
  granularity: Granularity,
): TimeBucket[] {
  const buckets = new Map<string, TimeBucket>();

  const start = bucketStart(parseDateValue(range.from), granularity);
  const end = parseDateValue(range.to);
  for (
    let cursor = start;
    cursor <= end;
    cursor =
      granularity === "month"
        ? addMonths(cursor, 1)
        : addDays(cursor, granularity === "week" ? 7 : 1)
  ) {
    const key = format(cursor, "yyyy-MM-dd");
    buckets.set(key, emptyBucket(key, granularity));
  }

  for (const row of rows) {
    const key = bucketKey(new Date(row.createdAt), granularity);
    const bucket = buckets.get(key) ?? emptyBucket(key, granularity);
    bucket.count += 1;
    bucket.revenue += row.revenueIls;
    bucket.cost += row.costIls;
    bucket.profit += row.profitIls;
    if (row.isBusiness) {
      bucket.business += row.profitIls;
    } else {
      bucket.private += row.profitIls;
    }
    buckets.set(key, bucket);
  }

  return [...buckets.values()].sort((a, b) => a.key.localeCompare(b.key));
}

// --- Grouping ---------------------------------------------------------------

export type Metric = "count" | "revenue" | "profit";

export interface Group {
  key: string;
  label: string;
  count: number;
  revenue: number;
  cost: number;
  profit: number;
  marginPct: number;
  /** The value of the metric currently being displayed, so charts can read one field. */
  value: number;
}

export const OTHER_KEY = "__other__";

export function groupRows(
  rows: DashboardRow[],
  keyOf: (row: DashboardRow) => string | null,
  labelOf: (key: string) => string,
  metric: Metric = "count",
): Group[] {
  const groups = new Map<string, Group>();

  for (const row of rows) {
    const key = keyOf(row);
    if (key === null) continue;

    const group = groups.get(key) ?? {
      key,
      label: labelOf(key),
      count: 0,
      revenue: 0,
      cost: 0,
      profit: 0,
      marginPct: 0,
      value: 0,
    };
    group.count += 1;
    group.revenue += row.revenueIls;
    group.cost += row.costIls;
    group.profit += row.profitIls;
    groups.set(key, group);
  }

  return [...groups.values()]
    .map((group) => ({
      ...group,
      marginPct: group.revenue > 0 ? (group.profit / group.revenue) * 100 : 0,
      value: group[metric],
    }))
    .sort((a, b) => b.value - a.value);
}

/**
 * Keeps the leading groups and folds the tail into "אחר". Never solve too many categories
 * by generating more colours — past the palette's slots they stop being distinguishable.
 */
export function topN(groups: Group[], n: number): Group[] {
  if (groups.length <= n) return groups;

  const head = groups.slice(0, n);
  const tail = groups.slice(n);
  const other = tail.reduce(
    (acc, group) => ({
      ...acc,
      count: acc.count + group.count,
      revenue: acc.revenue + group.revenue,
      cost: acc.cost + group.cost,
      profit: acc.profit + group.profit,
      value: acc.value + group.value,
    }),
    {
      key: OTHER_KEY,
      label: `אחר (${tail.length})`,
      count: 0,
      revenue: 0,
      cost: 0,
      profit: 0,
      marginPct: 0,
      value: 0,
    } as Group,
  );

  return [
    ...head,
    { ...other, marginPct: other.revenue > 0 ? (other.profit / other.revenue) * 100 : 0 },
  ];
}

// --- Histograms -------------------------------------------------------------

export interface Band {
  label: string;
  min: number;
  max: number;
}

export const DRIVER_AGE_BANDS: Band[] = [
  { label: "18-24", min: 18, max: 24 },
  { label: "25-29", min: 25, max: 29 },
  { label: "30-39", min: 30, max: 39 },
  { label: "40-49", min: 40, max: 49 },
  { label: "50-59", min: 50, max: 59 },
  { label: "60+", min: 60, max: Infinity },
];

export const RENTAL_DAYS_BANDS: Band[] = [
  { label: "1-2", min: 1, max: 2 },
  { label: "3-5", min: 3, max: 5 },
  { label: "6-9", min: 6, max: 9 },
  { label: "10-14", min: 10, max: 14 },
  { label: "15+", min: 15, max: Infinity },
];

export const LEAD_TIME_BANDS: Band[] = [
  { label: "0-2", min: 0, max: 2 },
  { label: "3-7", min: 3, max: 7 },
  { label: "8-30", min: 8, max: 30 },
  { label: "31-90", min: 31, max: 90 },
  { label: "90+", min: 91, max: Infinity },
];

export interface HistogramBar {
  label: string;
  count: number;
}

export function histogram(
  rows: DashboardRow[],
  bands: Band[],
  valueOf: (row: DashboardRow) => number,
): HistogramBar[] {
  return bands.map((band) => ({
    label: band.label,
    count: rows.filter((row) => {
      const value = valueOf(row);
      return value >= band.min && value <= band.max;
    }).length,
  }));
}

// --- Payments ---------------------------------------------------------------

export interface Payments {
  openToSuppliers: number;
  paidToSuppliers: number;
  openFromCustomers: number;
  collected: number;
  aging: { label: string; amount: number }[];
}

const AGING_BANDS: Band[] = [
  { label: "0-30 ימים", min: 0, max: 30 },
  { label: "31-60 ימים", min: 31, max: 60 },
  { label: "61-90 ימים", min: 61, max: 90 },
  { label: "90+ ימים", min: 91, max: Infinity },
];

/** Mirrors the backend's own definitions of what is outstanding — see the queries in
 *  reservation_query.sql that drive the supplier-payments and collections screens. */
export function isOpenToSupplier(row: DashboardRow): boolean {
  return !row.supplierPaid && row.status !== CANCELED;
}

export function isOpenFromCustomer(row: DashboardRow): boolean {
  return (
    (row.status === "vouchered" && row.paymentStatus === "unpaid") ||
    (row.status === CANCELED && row.paymentStatus === "refund_pending")
  );
}

export function computePayments(rows: DashboardRow[]): Payments {
  let openToSuppliers = 0;
  let paidToSuppliers = 0;
  let openFromCustomers = 0;
  let collected = 0;

  const now = Date.now();
  const aging = AGING_BANDS.map((band) => ({ label: band.label, amount: 0 }));

  for (const row of rows) {
    const supplierOwed = row.costIls + (row.penaltySupplierPaid ? 0 : row.penaltyAmountIls);

    if (isOpenToSupplier(row)) {
      openToSuppliers += supplierOwed;
    } else if (row.supplierPaid) {
      paidToSuppliers += row.costIls;
    }

    if (isOpenFromCustomer(row)) {
      // A refund pending is money flowing the other way, so it reduces what we are owed.
      const owed =
        row.status === CANCELED ? -row.revenueIls : row.revenueIls;
      openFromCustomers += owed;

      const ageDays = Math.floor(
        (now - new Date(row.createdAt).getTime()) / (1000 * 60 * 60 * 24),
      );
      const bandIndex = AGING_BANDS.findIndex(
        (band) => ageDays >= band.min && ageDays <= band.max,
      );
      if (bandIndex >= 0) aging[bandIndex].amount += owed;
    }

    if (row.paymentStatus === "paid") {
      collected += row.revenueIls;
    }
  }

  return { openToSuppliers, paidToSuppliers, openFromCustomers, collected, aging };
}

// --- Entity lookups ---------------------------------------------------------

export type EntityDimension =
  | "organization"
  | "office"
  | "agent"
  | "customer";

export interface Lookups {
  organizations: Map<number, string>;
  offices: Map<number, string>;
  users: Map<number, string>;
}

export function buildLookups(response: reports.DashboardResponse): Lookups {
  const toMap = (names: accounts.AccountName[]) =>
    new Map(names.map((name) => [name.id, name.name]));

  return {
    organizations: toMap(response.organizations),
    offices: toMap(response.offices),
    users: toMap(response.users),
  };
}

export function entityGrouping(
  dimension: EntityDimension,
  lookups: Lookups,
): {
  keyOf: (row: DashboardRow) => string | null;
  labelOf: (key: string) => string;
} {
  switch (dimension) {
    case "organization":
      return {
        keyOf: (row) => (row.organizationId ? String(row.organizationId) : null),
        labelOf: (key) => lookups.organizations.get(Number(key)) ?? `רשת #${key}`,
      };
    case "office":
      return {
        keyOf: (row) => (row.officeId ? String(row.officeId) : null),
        labelOf: (key) => lookups.offices.get(Number(key)) ?? `משרד #${key}`,
      };
    case "agent":
      return {
        keyOf: (row) => (row.isBusiness ? String(row.userId) : null),
        labelOf: (key) => lookups.users.get(Number(key)) ?? `סוכן #${key}`,
      };
    case "customer":
      return {
        keyOf: (row) => (row.isBusiness ? null : String(row.userId)),
        labelOf: (key) => lookups.users.get(Number(key)) ?? `לקוח #${key}`,
      };
  }
}
