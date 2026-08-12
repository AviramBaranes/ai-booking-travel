"use client";

import { useState } from "react";
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

import { Granularity, TimeBucket } from "../_lib/aggregate";
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

  return (
    <ChartCard
      title={view === "combined" ? "הכנסות, עלות ורווח לאורך זמן" : "רווח: עסקי מול פרטי"}
      subtitle={GRANULARITY_LABELS[granularity]}
      actions={
        <SegmentedControl
          value={view}
          onChange={(next) => setView(next as View)}
          options={[
            { value: "combined", label: "מצטבר" },
            { value: "split", label: "עסקי מול פרטי" },
          ]}
        />
      }
      tableView={{
        columns: [
          { label: "תקופה" },
          { label: "הזמנות", align: "end" },
          { label: "הכנסות", align: "end" },
          { label: "עלות", align: "end" },
          { label: "רווח", align: "end" },
        ],
        rows: buckets.map((bucket) => [
          bucket.label,
          count(bucket.count),
          ils(bucket.revenue),
          ils(bucket.cost),
          ils(bucket.profit),
        ]),
      }}
    >
      {/* Recharts lays out left-to-right; the card around it stays RTL. */}
      <div dir="ltr">
        <ChartContainer config={config} className="h-64 w-full">
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
              tickFormatter={ilsCompact}
            />
            <ChartTooltip
              content={<ChartTooltipContent formatter={(value) => ils(Number(value))} />}
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
      title="כמות הזמנות לאורך זמן"
      subtitle={GRANULARITY_LABELS[granularity]}
      tableView={{
        columns: [{ label: "תקופה" }, { label: "הזמנות", align: "end" }],
        rows: buckets.map((bucket) => [bucket.label, count(bucket.count)]),
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
            <ChartTooltip content={<ChartTooltipContent />} />
            <Bar dataKey="count" fill="var(--color-count)" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ChartContainer>
      </div>
    </ChartCard>
  );
}
