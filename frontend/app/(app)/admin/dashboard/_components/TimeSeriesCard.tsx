"use client";

import { useMemo, useState } from "react";
import {
  CartesianGrid,
  Line,
  LineChart,
  XAxis,
  YAxis,
  Bar,
  BarChart,
} from "recharts";

import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

import { Granularity, TimeBucket, toCumulative } from "../_lib/aggregate";
import { ChartCard } from "./ChartCard";
import { SegmentedControl } from "./SegmentedControl";
import { count, ils, ilsCompact } from "../_lib/format";

const GRANULARITY_LABELS: Record<Granularity, string> = {
  day: "לפי יום",
  week: "לפי שבוע",
  month: "לפי חודש",
};

const moneyConfig = {
  revenue: { label: "הכנסות", color: "var(--chart-1)" },
  cost: { label: "עלות", color: "var(--chart-2)" },
  profit: { label: "רווח", color: "var(--chart-3)" },
} satisfies ChartConfig;

const splitConfig = {
  business: { label: "עסקי", color: "var(--chart-4)" },
  private: { label: "פרטי", color: "var(--chart-5)" },
} satisfies ChartConfig;

const countConfig = {
  count: { label: "הזמנות", color: "var(--chart-1)" },
} satisfies ChartConfig;

type View = "combined" | "split";

const MONEY_KEYS = ["revenue", "cost", "profit"] as const;

/**
 * The axis tick is abbreviated to fit, which makes a weekly or monthly bucket look like a
 * single day. The tooltip has room, so it names the whole span the point covers.
 */
function periodTooltipLabel(_label: unknown, payload?: readonly { payload?: TimeBucket }[]) {
  return payload?.[0]?.payload?.periodLabel ?? "";
}

/**
 * Pins the axis floor to the data's real minimum, or to zero when nothing is negative.
 *
 * Left to its own devices recharts rounds the domain out to "nice" ticks, so a single
 * slightly-negative period drags the floor thousands below zero and squashes everything
 * that matters into the top half of the plot.
 */
function moneyDomain(
  buckets: TimeBucket[],
  keys: readonly (keyof TimeBucket)[],
): [number, "auto"] {
  let min = 0;
  for (const bucket of buckets) {
    for (const key of keys) {
      min = Math.min(min, bucket[key] as number);
    }
  }
  return [min, "auto"];
}

/**
 * Dots mark the periods that were actually measured and stay visible at every density.
 *
 * What made a dense series look like a dashed line was the 2px surface ring, not the dots:
 * packed tighter than the rings are wide, it cut the stroke between every pair. So past a
 * comfortable density the ring comes off and the dot shrinks, which keeps every period
 * marked without chopping up the line that connects them.
 */
const COMFORTABLE_DOT_DENSITY = 31;

function dotProps(pointCount: number, color: string) {
  // `fill` has to be set explicitly: recharts does not inherit the line's colour for a
  // custom dot, so a dot with only a surface-coloured stroke renders white on white.
  return pointCount <= COMFORTABLE_DOT_DENSITY
    ? { r: 3, strokeWidth: 2, stroke: "var(--card)", fill: color }
    : { r: 1.5, strokeWidth: 0, fill: color };
}

export function TimeSeriesCard({
  buckets,
  granularity,
}: {
  buckets: TimeBucket[];
  granularity: Granularity;
}) {
  const [view, setView] = useState<View>("combined");
  const config = view === "combined" ? moneyConfig : splitConfig;
  const series =
    view === "combined"
      ? (["revenue", "cost", "profit"] as const)
      : (["business", "private"] as const);

  // The line is a running total from zero, so each point is "everything up to here".
  const running = useMemo(() => toCumulative(buckets), [buckets]);

  return (
    <ChartCard
      title={
        view === "combined"
          ? "הכנסות, עלות ורווח — מצטבר"
          : "רווח מצטבר: עסקי מול פרטי"
      }
      subtitle={`נצבר מתחילת הטווח · ${GRANULARITY_LABELS[granularity]}`}
      actions={
        <SegmentedControl
          value={view}
          onChange={(next) => setView(next as View)}
          options={[
            { value: "combined", label: "הכל" },
            { value: "split", label: "עסקי מול פרטי" },
          ]}
        />
      }
      // The table keeps the per-period figures beside the running ones, so a single day or
      // week can still be reconciled against the row-level reports.
      tableView={{
        columns: [
          { label: "תקופה" },
          { label: "הזמנות", align: "end" },
          { label: "הכנסות בתקופה", align: "end" },
          { label: "עלות בתקופה", align: "end" },
          { label: "רווח בתקופה", align: "end" },
          { label: "הכנסות מצטבר", align: "end" },
          { label: "רווח מצטבר", align: "end" },
        ],
        rows: buckets.map((bucket, index) => [
          bucket.periodLabel,
          count(bucket.count),
          ils(bucket.revenue),
          ils(bucket.cost),
          ils(bucket.profit),
          ils(running[index].revenue),
          ils(running[index].profit),
        ]),
      }}
    >
      {/* Recharts lays out left-to-right; the card around it stays RTL. */}
      <div dir="ltr">
        <ChartContainer config={config} className="h-64 w-full">
          <LineChart data={running} margin={{ top: 8, right: 12, left: 4, bottom: 0 }}>
            <CartesianGrid vertical={false} stroke="var(--chart-grid)" />
            <XAxis
              dataKey="label"
              tickLine={false}
              axisLine={{ stroke: "var(--chart-axis)" }}
              tickMargin={8}
              minTickGap={16}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={4}
              width={56}
              domain={moneyDomain(running, series)}
              tickFormatter={ilsCompact}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  labelFormatter={periodTooltipLabel}
                  formatter={(value) => ils(Number(value))}
                />
              }
            />
            <ChartLegend content={<ChartLegendContent />} />
            {series.map((key) => (
              <Line
                key={key}
                dataKey={key}
                type="monotone"
                stroke={`var(--color-${key})`}
                strokeWidth={2}
                dot={false}
                activeDot={{ r: 4, strokeWidth: 2, stroke: "var(--card)" }}
              />
            ))}
          </LineChart>
        </ChartContainer>
      </div>
    </ChartCard>
  );
}

/**
 * The volume counterpart to TimeSeriesCard: the same running-total reading, but for how
 * many reservations were booked rather than how much money they carried. Kept on its own
 * axis in its own card because a count and a shekel figure never share a scale.
 */
export function CumulativeVolumeCard({
  buckets,
  granularity,
}: {
  buckets: TimeBucket[];
  granularity: Granularity;
}) {
  const running = useMemo(() => toCumulative(buckets), [buckets]);

  return (
    <ChartCard
      title="כמות הזמנות — מצטבר"
      subtitle={`נצבר מתחילת הטווח · ${GRANULARITY_LABELS[granularity]}`}
      tableView={{
        columns: [
          { label: "תקופה" },
          { label: "הזמנות בתקופה", align: "end" },
          { label: "הזמנות מצטבר", align: "end" },
        ],
        rows: buckets.map((bucket, index) => [
          bucket.periodLabel,
          count(bucket.count),
          count(running[index].count),
        ]),
      }}
    >
      <div dir="ltr">
        {/* A single series, so the title names it and no legend box is needed. */}
        <ChartContainer config={countConfig} className="h-64 w-full">
          <LineChart data={running} margin={{ top: 8, right: 12, left: 4, bottom: 0 }}>
            <CartesianGrid vertical={false} stroke="var(--chart-grid)" />
            <XAxis
              dataKey="label"
              tickLine={false}
              axisLine={{ stroke: "var(--chart-axis)" }}
              tickMargin={8}
              minTickGap={16}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={4}
              width={40}
              allowDecimals={false}
            />
            <ChartTooltip
              content={<ChartTooltipContent labelFormatter={periodTooltipLabel} />}
            />
            <Line
              dataKey="count"
              type="monotone"
              stroke="var(--color-count)"
              strokeWidth={2}
              dot={false}
              activeDot={{ r: 4, strokeWidth: 2, stroke: "var(--card)" }}
            />
          </LineChart>
        </ChartContainer>
      </div>
    </ChartCard>
  );
}

/**
 * Money bucket by bucket, so a strong or weak period stands on its own rather than being
 * absorbed into a running total.
 *
 * Straight segments between marked points, not a smoothed curve: the whole reason to look
 * here is to spot a jump or a drop, and interpolation invents values between periods that
 * were never measured.
 */
export function MoneyPerPeriodCard({
  buckets,
  granularity,
}: {
  buckets: TimeBucket[];
  granularity: Granularity;
}) {
  return (
    <ChartCard
      title="הכנסות, עלות ורווח לפי תקופה"
      subtitle={GRANULARITY_LABELS[granularity]}
      tableView={{
        columns: [
          { label: "תקופה" },
          { label: "הכנסות", align: "end" },
          { label: "עלות", align: "end" },
          { label: "רווח", align: "end" },
        ],
        rows: buckets.map((bucket) => [
          bucket.periodLabel,
          ils(bucket.revenue),
          ils(bucket.cost),
          ils(bucket.profit),
        ]),
      }}
    >
      <div dir="ltr">
        <ChartContainer config={moneyConfig} className="h-64 w-full">
          <LineChart data={buckets} margin={{ top: 8, right: 12, left: 4, bottom: 0 }}>
            <CartesianGrid vertical={false} stroke="var(--chart-grid)" />
            <XAxis
              dataKey="label"
              tickLine={false}
              axisLine={{ stroke: "var(--chart-axis)" }}
              tickMargin={8}
              minTickGap={16}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={4}
              width={56}
              domain={moneyDomain(buckets, MONEY_KEYS)}
              tickFormatter={ilsCompact}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  labelFormatter={periodTooltipLabel}
                  formatter={(value) => ils(Number(value))}
                />
              }
            />
            <ChartLegend content={<ChartLegendContent />} />
            {MONEY_KEYS.map((key) => (
              <Line
                key={key}
                dataKey={key}
                type="linear"
                stroke={`var(--color-${key})`}
                strokeWidth={2}
                dot={dotProps(buckets.length, `var(--color-${key})`)}
                activeDot={{ r: 5, strokeWidth: 2, stroke: "var(--card)" }}
              />
            ))}
          </LineChart>
        </ChartContainer>
      </div>
    </ChartCard>
  );
}

/** Volume is a different unit from money, so it gets its own axis in its own chart. */
export function VolumeCard({
  buckets,
  granularity,
}: {
  buckets: TimeBucket[];
  granularity: Granularity;
}) {
  return (
    <ChartCard
      title="כמות הזמנות לפי תקופה"
      subtitle={GRANULARITY_LABELS[granularity]}
      tableView={{
        columns: [{ label: "תקופה" }, { label: "הזמנות", align: "end" }],
        rows: buckets.map((bucket) => [bucket.periodLabel, count(bucket.count)]),
      }}
    >
      <div dir="ltr">
        <ChartContainer config={countConfig} className="h-40 w-full">
          <BarChart data={buckets} margin={{ top: 8, right: 12, left: 4, bottom: 0 }}>
            <CartesianGrid vertical={false} stroke="var(--chart-grid)" />
            <XAxis
              dataKey="label"
              tickLine={false}
              axisLine={{ stroke: "var(--chart-axis)" }}
              tickMargin={8}
              minTickGap={16}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={4}
              width={40}
              allowDecimals={false}
            />
            <ChartTooltip
              content={<ChartTooltipContent labelFormatter={periodTooltipLabel} />}
            />
            <Bar dataKey="count" fill="var(--color-count)" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ChartContainer>
      </div>
    </ChartCard>
  );
}
