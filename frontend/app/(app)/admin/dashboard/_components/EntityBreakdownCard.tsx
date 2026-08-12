"use client";

import { useMemo, useState } from "react";
import { Bar, BarChart, CartesianGrid, LabelList, XAxis, YAxis } from "recharts";

import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

import {
  DashboardRow,
  EntityDimension,
  Lookups,
  Metric,
  entityGrouping,
  groupRows,
  topN,
} from "../_lib/aggregate";
import { ChartCard } from "./ChartCard";
import { SegmentedControl } from "./SegmentedControl";
import { count, ils, ilsCompact, labelFormatter, percent } from "../_lib/format";

const DIMENSIONS: { value: EntityDimension; label: string }[] = [
  { value: "organization", label: "רשתות" },
  { value: "office", label: "משרדים" },
  { value: "agent", label: "סוכנים" },
  { value: "customer", label: "לקוחות פרטיים" },
];

const METRICS: { value: Metric; label: string }[] = [
  { value: "profit", label: "רווח" },
  { value: "revenue", label: "הכנסות" },
  { value: "count", label: "הזמנות" },
];

// A single series, so it takes slot 1 for every bar. Shading bars by their own value would
// double-encode the bar length and burn the only free channel.
const config = {
  value: { label: "ערך", color: "var(--chart-1)" },
} satisfies ChartConfig;

export function EntityBreakdownCard({
  rows,
  lookups,
}: {
  rows: DashboardRow[];
  lookups: Lookups;
}) {
  const [dimension, setDimension] = useState<EntityDimension>("organization");
  const [metric, setMetric] = useState<Metric>("profit");

  const groups = useMemo(() => {
    const { keyOf, labelOf } = entityGrouping(dimension, lookups);
    return groupRows(rows, keyOf, labelOf, metric);
  }, [rows, dimension, lookups, metric]);

  const top = useMemo(() => topN(groups, 10), [groups]);
  const formatValue = (value: number) =>
    metric === "count" ? count(value) : ils(value);

  return (
    <ChartCard
      title="ביצועים לפי ישות"
      subtitle={`10 המובילים לפי ${METRICS.find((m) => m.value === metric)?.label}`}
      actions={
        <div className="flex flex-wrap items-center gap-2">
          <SegmentedControl
            value={dimension}
            onChange={(next) => setDimension(next as EntityDimension)}
            options={DIMENSIONS}
          />
          <SegmentedControl
            value={metric}
            onChange={(next) => setMetric(next as Metric)}
            options={METRICS}
          />
        </div>
      }
      tableView={{
        columns: [
          { label: "ישות" },
          { label: "הזמנות", align: "end" },
          { label: "הכנסות", align: "end" },
          { label: "עלות", align: "end" },
          { label: "רווח", align: "end" },
          { label: "% רווח", align: "end" },
        ],
        rows: groups.map((group) => [
          group.label,
          count(group.count),
          ils(group.revenue),
          ils(group.cost),
          ils(group.profit),
          percent(group.marginPct),
        ]),
      }}
    >
      {top.length === 0 ? (
        <EmptyState />
      ) : (
        <div dir="ltr">
          <ChartContainer
            config={config}
            className="w-full"
            style={{ height: Math.max(top.length * 34 + 32, 160) }}
          >
            <BarChart
              data={top}
              layout="vertical"
              margin={{ top: 4, right: 56, left: 8, bottom: 4 }}
              barCategoryGap={6}
            >
              <CartesianGrid horizontal={false} stroke="var(--chart-grid)" />
              <XAxis
                type="number"
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) =>
                  metric === "count" ? count(value) : ilsCompact(value)
                }
              />
              <YAxis
                type="category"
                dataKey="label"
                tickLine={false}
                axisLine={false}
                width={150}
                tickMargin={6}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value) => formatValue(Number(value))}
                  />
                }
              />
              <Bar dataKey="value" fill="var(--color-value)" radius={[0, 4, 4, 0]}>
                {/* Direct labels: three palette slots are sub-3:1 on white, so values must
                    be legible without relying on the fill. */}
                <LabelList
                  dataKey="value"
                  position="right"
                  className="fill-foreground"
                  fontSize={11}
                  formatter={labelFormatter((value) =>
                    metric === "count" ? count(value) : ilsCompact(value),
                  )}
                />
              </Bar>
            </BarChart>
          </ChartContainer>
        </div>
      )}
    </ChartCard>
  );
}

export function EmptyState({ message = "אין נתונים בטווח שנבחר" }: { message?: string }) {
  return (
    <p className="py-12 text-center text-sm text-gray-500">{message}</p>
  );
}
