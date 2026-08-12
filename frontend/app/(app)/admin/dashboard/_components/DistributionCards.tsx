"use client";

import { useMemo } from "react";
import { Bar, BarChart, Cell, LabelList, Pie, PieChart, XAxis, YAxis } from "recharts";

import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

import { Group, Metric } from "../_lib/aggregate";
import { ChartCard } from "./ChartCard";
import { EmptyState } from "./EntityBreakdownCard";
import {
  CATEGORICAL,
  count,
  ils,
  ilsCompact,
  labelFormatter,
  percent,
} from "../_lib/format";

function metricLabel(metric: Metric): string {
  if (metric === "count") return "הזמנות";
  if (metric === "revenue") return "הכנסות";
  return "רווח";
}

function formatMetric(value: number, metric: Metric): string {
  return metric === "count" ? count(value) : ils(value);
}

function tableView(groups: Group[]) {
  return {
    columns: [
      { label: "ערך" },
      { label: "הזמנות", align: "end" as const },
      { label: "הכנסות", align: "end" as const },
      { label: "רווח", align: "end" as const },
      { label: "% רווח", align: "end" as const },
    ],
    rows: groups.map((group) => [
      group.label,
      count(group.count),
      ils(group.revenue),
      ils(group.profit),
      percent(group.marginPct),
    ]),
  };
}

/**
 * A donut is only honest for a genuine part-to-whole with a handful of slices, so callers
 * cap the input at six. Anything longer belongs in CategoryBarCard.
 */
export function DonutCard({
  title,
  groups,
  metric,
  colors,
}: {
  title: string;
  groups: Group[];
  metric: Metric;
  /** Explicit per-key colours, e.g. the reserved status palette. */
  colors: Record<string, string>;
}) {
  const total = groups.reduce((sum, group) => sum + group.value, 0);

  const config = useMemo(
    () =>
      Object.fromEntries(
        groups.map((group) => [
          group.key,
          { label: group.label, color: colors[group.key] ?? "var(--chart-1)" },
        ]),
      ) satisfies ChartConfig,
    [groups, colors],
  );

  return (
    <ChartCard
      title={title}
      subtitle={`לפי ${metricLabel(metric)}`}
      tableView={tableView(groups)}
    >
      {groups.length === 0 ? (
        <EmptyState />
      ) : (
        <div dir="ltr" className="flex items-center gap-4">
          <ChartContainer config={config} className="h-40 w-40 shrink-0">
            <PieChart>
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    nameKey="label"
                    formatter={(value) => formatMetric(Number(value), metric)}
                  />
                }
              />
              <Pie
                data={groups}
                dataKey="value"
                nameKey="label"
                innerRadius={38}
                outerRadius={68}
                // A 2px surface gap separates the segments — never a drawn border.
                paddingAngle={2}
                stroke="var(--card)"
                strokeWidth={2}
              >
                {groups.map((group) => (
                  <Cell
                    key={group.key}
                    fill={colors[group.key] ?? "var(--chart-1)"}
                  />
                ))}
              </Pie>
            </PieChart>
          </ChartContainer>

          {/* The legend doubles as the direct-label channel, so no value is colour-only. */}
          <ul dir="rtl" className="flex-1 space-y-1.5">
            {groups.map((group) => (
              <li
                key={group.key}
                className="flex items-center justify-between gap-2 text-xs"
              >
                <span className="flex items-center gap-1.5 text-gray-700">
                  <span
                    aria-hidden
                    className="size-2.5 shrink-0 rounded-[2px]"
                    style={{ background: colors[group.key] ?? "var(--chart-1)" }}
                  />
                  {group.label}
                </span>
                <span className="text-gray-500 tabular-nums">
                  {formatMetric(group.value, metric)}
                  {total > 0 && ` · ${percent((group.value / total) * 100, 0)}`}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </ChartCard>
  );
}

/** A two-way split is a stat pair, not a pie. One 100%-wide bar reads it at a glance. */
export function ShareBarCard({
  title,
  groups,
  metric,
}: {
  title: string;
  groups: Group[];
  metric: Metric;
}) {
  const total = groups.reduce((sum, group) => sum + group.value, 0);

  return (
    <ChartCard
      title={title}
      subtitle={`לפי ${metricLabel(metric)}`}
      tableView={tableView(groups)}
    >
      {total === 0 ? (
        <EmptyState />
      ) : (
        <div className="space-y-3">
          <div className="flex h-6 w-full gap-0.5 overflow-hidden rounded">
            {groups.map((group, index) => (
              <div
                key={group.key}
                className="h-full first:rounded-r last:rounded-l"
                style={{
                  width: `${(group.value / total) * 100}%`,
                  background: CATEGORICAL[index % CATEGORICAL.length],
                }}
                title={`${group.label}: ${formatMetric(group.value, metric)}`}
              />
            ))}
          </div>
          <ul className="space-y-1.5">
            {groups.map((group, index) => (
              <li
                key={group.key}
                className="flex items-center justify-between gap-2 text-xs"
              >
                <span className="flex items-center gap-1.5 text-gray-700">
                  <span
                    aria-hidden
                    className="size-2.5 shrink-0 rounded-[2px]"
                    style={{ background: CATEGORICAL[index % CATEGORICAL.length] }}
                  />
                  {group.label}
                </span>
                <span className="text-gray-500 tabular-nums">
                  {formatMetric(group.value, metric)} ·{" "}
                  {percent((group.value / total) * 100, 0)}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </ChartCard>
  );
}

/** Many categories: colour stops distinguishing them, so bar length carries the comparison. */
export function CategoryBarCard({
  title,
  groups,
  metric,
  labelWidth = 110,
}: {
  title: string;
  groups: Group[];
  metric: Metric;
  labelWidth?: number;
}) {
  const config = {
    value: { label: metricLabel(metric), color: "var(--chart-1)" },
  } satisfies ChartConfig;

  return (
    <ChartCard
      title={title}
      subtitle={`לפי ${metricLabel(metric)}`}
      tableView={tableView(groups)}
    >
      {groups.length === 0 ? (
        <EmptyState />
      ) : (
        <div dir="ltr">
          <ChartContainer
            config={config}
            className="w-full"
            style={{ height: Math.max(groups.length * 30 + 24, 140) }}
          >
            <BarChart
              data={groups}
              layout="vertical"
              margin={{ top: 4, right: 52, left: 4, bottom: 4 }}
              barCategoryGap={6}
            >
              <XAxis type="number" hide />
              <YAxis
                type="category"
                dataKey="label"
                tickLine={false}
                axisLine={false}
                width={labelWidth}
                tickMargin={6}
                fontSize={11}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value) => formatMetric(Number(value), metric)}
                  />
                }
              />
              <Bar dataKey="value" fill="var(--color-value)" radius={[0, 4, 4, 0]}>
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

/** Ordered bands, so the bars take steps of one hue rather than eight identities. */
const ORDINAL_RAMP = [
  "#86b6ef",
  "#5598e7",
  "#3987e5",
  "#2a78d6",
  "#256abf",
  "#1c5cab",
];

export function HistogramCard({
  title,
  subtitle,
  bars,
}: {
  title: string;
  subtitle?: string;
  bars: { label: string; count: number }[];
}) {
  const config = {
    count: { label: "הזמנות", color: "var(--chart-1)" },
  } satisfies ChartConfig;

  const total = bars.reduce((sum, bar) => sum + bar.count, 0);

  return (
    <ChartCard
      title={title}
      subtitle={subtitle}
      tableView={{
        columns: [
          { label: "טווח" },
          { label: "הזמנות", align: "end" },
          { label: "חלק מהסך", align: "end" },
        ],
        rows: bars.map((bar) => [
          bar.label,
          count(bar.count),
          total > 0 ? percent((bar.count / total) * 100, 0) : "—",
        ]),
      }}
    >
      {total === 0 ? (
        <EmptyState />
      ) : (
        <div dir="ltr">
          <ChartContainer config={config} className="h-40 w-full">
            <BarChart data={bars} margin={{ top: 16, right: 8, left: 4, bottom: 0 }}>
              <XAxis
                dataKey="label"
                tickLine={false}
                axisLine={{ stroke: "var(--chart-axis)" }}
                tickMargin={8}
                fontSize={11}
              />
              <YAxis hide />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                {bars.map((bar, index) => (
                  <Cell
                    key={bar.label}
                    fill={ORDINAL_RAMP[index % ORDINAL_RAMP.length]}
                  />
                ))}
                <LabelList
                  dataKey="count"
                  position="top"
                  className="fill-foreground"
                  fontSize={11}
                />
              </Bar>
            </BarChart>
          </ChartContainer>
        </div>
      )}
    </ChartCard>
  );
}
