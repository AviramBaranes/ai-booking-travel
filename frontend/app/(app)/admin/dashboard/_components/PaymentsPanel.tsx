"use client";

import { Bar, BarChart, Cell, LabelList, XAxis, YAxis } from "recharts";
import { Banknote, Wallet } from "lucide-react";

import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

import { Payments } from "../_lib/aggregate";
import { ChartCard } from "./ChartCard";
import { EmptyState } from "./EntityBreakdownCard";
import { Meter, StatTile } from "./StatTile";
import { ils, ilsCompact, labelFormatter } from "../_lib/format";

const ORDINAL_RAMP = ["#86b6ef", "#3987e5", "#256abf", "#104281"];

const config = {
  amount: { label: "יתרה פתוחה", color: "var(--chart-1)" },
} satisfies ChartConfig;

export function PaymentsPanel({ payments }: { payments: Payments }) {
  const supplierTotal = payments.openToSuppliers + payments.paidToSuppliers;
  const customerTotal = payments.openFromCustomers + payments.collected;
  const agingTotal = payments.aging.reduce((sum, band) => sum + band.amount, 0);

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <div className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2">
          <StatTile
            label="פתוח לספקים"
            value={ils(payments.openToSuppliers)}
            hint="טרם שולם"
            icon={<Wallet className="size-4" />}
          />
          <StatTile
            label="שולם לספקים"
            value={ils(payments.paidToSuppliers)}
            icon={<Wallet className="size-4" />}
          />
          <StatTile
            label="פתוח מלקוחות"
            value={ils(payments.openFromCustomers)}
            hint="טרם נגבה"
            icon={<Banknote className="size-4" />}
          />
          <StatTile
            label="נגבה מלקוחות"
            value={ils(payments.collected)}
            icon={<Banknote className="size-4" />}
          />
        </div>

        <div className="space-y-4 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <Meter
            label="שולם לספקים מתוך סך העלות"
            value={payments.paidToSuppliers}
            total={supplierTotal}
            valueLabel={ils(payments.paidToSuppliers)}
          />
          <Meter
            label="נגבה מלקוחות מתוך סך ההכנסות"
            value={payments.collected}
            total={customerTotal}
            valueLabel={ils(payments.collected)}
          />
        </div>
      </div>

      <ChartCard
        title="גיל החוב הפתוח מלקוחות"
        subtitle="לפי הזמן שחלף מיצירת ההזמנה"
        tableView={{
          columns: [
            { label: "טווח" },
            { label: "יתרה פתוחה", align: "end" },
          ],
          rows: payments.aging.map((band) => [band.label, ils(band.amount)]),
        }}
      >
        {agingTotal === 0 ? (
          <EmptyState message="אין חוב פתוח בטווח שנבחר" />
        ) : (
          <div dir="ltr">
            <ChartContainer config={config} className="h-52 w-full">
              <BarChart
                data={payments.aging}
                layout="vertical"
                margin={{ top: 4, right: 60, left: 4, bottom: 4 }}
                barCategoryGap={8}
              >
                <XAxis type="number" hide />
                <YAxis
                  type="category"
                  dataKey="label"
                  tickLine={false}
                  axisLine={false}
                  width={92}
                  fontSize={11}
                />
                <ChartTooltip
                  content={
                    <ChartTooltipContent
                      formatter={(value) => ils(Number(value))}
                    />
                  }
                />
                <Bar dataKey="amount" radius={[0, 4, 4, 0]}>
                  {payments.aging.map((band, index) => (
                    <Cell key={band.label} fill={ORDINAL_RAMP[index]} />
                  ))}
                  <LabelList
                    dataKey="amount"
                    position="right"
                    className="fill-foreground"
                    fontSize={11}
                    formatter={labelFormatter(ilsCompact)}
                  />
                </Bar>
              </BarChart>
            </ChartContainer>
          </div>
        )}
      </ChartCard>
    </div>
  );
}
