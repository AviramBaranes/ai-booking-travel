"use client";

import { useMemo, useState } from "react";
import { Bar, BarChart, LabelList, XAxis, YAxis } from "recharts";

import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

import { DashboardRow, groupRows, topN } from "../_lib/aggregate";
import { ChartCard } from "./ChartCard";
import { EmptyState } from "./EntityBreakdownCard";
import { ExpandToggle } from "./ExpandToggle";
import { count, ils, percent } from "../_lib/format";

const config = {
  value: { label: "שימושים", color: "var(--chart-1)" },
} satisfies ChartConfig;

const LIMIT = 8;

export function CouponPanel({ rows }: { rows: DashboardRow[] }) {
  const [expanded, setExpanded] = useState(false);

  const groups = useMemo(
    () =>
      groupRows(
        rows,
        (row) => row.couponName || null,
        (key) => key,
        "count",
      ),
    [rows],
  );
  const visible = useMemo(
    () => (expanded ? groups : topN(groups, LIMIT)),
    [groups, expanded],
  );

  const discountByCoupon = useMemo(() => {
    const totals = new Map<string, number>();
    for (const row of rows) {
      if (!row.couponName) continue;
      totals.set(
        row.couponName,
        (totals.get(row.couponName) ?? 0) + row.discountIls,
      );
    }
    return totals;
  }, [rows]);

  return (
    <ChartCard
      title="קופונים"
      subtitle="שימושים והשפעה על הרווח"
      actions={
        <ExpandToggle
          expanded={expanded}
          onToggle={() => setExpanded((current) => !current)}
          hiddenCount={groups.length - LIMIT}
          limit={LIMIT}
        />
      }
      tableView={{
        columns: [
          { label: "קופון" },
          { label: "שימושים", align: "end" },
          { label: "הנחה", align: "end" },
          { label: "הכנסות", align: "end" },
          { label: "רווח", align: "end" },
          { label: "רווח ממוצע", align: "end" },
        ],
        rows: groups.map((group) => [
          group.label,
          count(group.count),
          ils(discountByCoupon.get(group.key) ?? 0),
          ils(group.revenue),
          ils(group.profit),
          ils(group.count > 0 ? group.profit / group.count : 0),
        ]),
      }}
    >
      {groups.length === 0 ? (
        <EmptyState message="לא נעשה שימוש בקופונים בטווח שנבחר" />
      ) : (
        <div
          dir="ltr"
          className={expanded ? "max-h-[32rem] overflow-y-auto" : undefined}
        >
          <ChartContainer
            config={config}
            className="w-full"
            style={{ height: Math.max(visible.length * 32 + 24, 140) }}
          >
            <BarChart
              data={visible}
              layout="vertical"
              margin={{ top: 4, right: 44, left: 4, bottom: 4 }}
              barCategoryGap={6}
            >
              <XAxis type="number" hide />
              <YAxis
                type="category"
                dataKey="label"
                tickLine={false}
                axisLine={false}
                width={120}
                fontSize={11}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value) => count(Number(value))}
                  />
                }
              />
              <Bar dataKey="value" fill="var(--color-value)" radius={[0, 4, 4, 0]}>
                <LabelList
                  dataKey="value"
                  position="right"
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

export function couponUsageHint(rows: DashboardRow[]): string {
  const withCoupon = rows.filter((row) => row.couponName).length;
  if (rows.length === 0) return "";
  return percent((withCoupon / rows.length) * 100, 0);
}
