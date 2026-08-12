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
import { ExpandToggle } from "./ExpandToggle";
import { SegmentedControl } from "./SegmentedControl";
import { count, ils, ilsCompact, labelFormatter, percent } from "../_lib/format";

const DIMENSIONS: { value: EntityDimension; label: string }[] = [
  { value: "organization", label: "רשתות" },
  { value: "office", label: "משרדים" },
  { value: "agent", label: "סוכנים" },
  { value: "customer", label: "לקוחות פרטיים" },
];

/** How many entities the chart shows before folding the tail; overridable per card. */
const LIMIT = 10;

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
  const [expanded, setExpanded] = useState(false);

  const groups = useMemo(() => {
    const { keyOf, labelOf } = entityGrouping(dimension, lookups);
    return groupRows(rows, keyOf, labelOf, metric);
  }, [rows, dimension, lookups, metric]);

  const visible = useMemo(
    () => (expanded ? groups : topN(groups, LIMIT)),
    [groups, expanded],
  );
  const formatValue = (value: number) =>
    metric === "count" ? count(value) : ils(value);
  const metricLabel = METRICS.find((m) => m.value === metric)?.label;

  return (
    <ChartCard
      title="ביצועים לפי ישות"
      subtitle={
        expanded
          ? `כל הישויות לפי ${metricLabel}`
          : `${LIMIT} המובילים לפי ${metricLabel}`
      }
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
          <ExpandToggle
            expanded={expanded}
            onToggle={() => setExpanded((current) => !current)}
            hiddenCount={groups.length - LIMIT}
            limit={LIMIT}
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
      {visible.length === 0 ? (
        <EmptyState />
      ) : (
        <div
          dir="ltr"
          className={expanded ? "max-h-[36rem] overflow-y-auto" : undefined}
        >
          <ChartContainer
            config={config}
            className="w-full"
            style={{ height: Math.max(visible.length * 34 + 32, 160) }}
          >
            <BarChart
              data={visible}
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
