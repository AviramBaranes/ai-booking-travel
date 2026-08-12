"use client";

import { ReactNode } from "react";
import { ArrowDown, ArrowUp } from "lucide-react";

import { cn } from "@/lib/utils";
import { delta, percent } from "../_lib/format";

interface StatTileProps {
  label: string;
  value: string;
  /** Current and previous raw values; the tile derives the change itself. */
  current?: number;
  previous?: number;
  /** For rates where a fall is the good outcome (cancellations, open debt). */
  invertDelta?: boolean;
  hint?: string;
  icon?: ReactNode;
}

export function StatTile({
  label,
  value,
  current,
  previous,
  invertDelta = false,
  hint,
  icon,
}: StatTileProps) {
  const change =
    current !== undefined && previous !== undefined
      ? delta(current, previous)
      : null;

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs text-gray-500">{label}</span>
        {icon && <span className="text-gray-300">{icon}</span>}
      </div>
      <p className="mt-1.5 text-2xl font-semibold text-navy">{value}</p>
      <div className="mt-1 flex items-center gap-2">
        {change !== null && <DeltaBadge change={change} invert={invertDelta} />}
        {hint && <span className="text-xs text-gray-400">{hint}</span>}
      </div>
    </div>
  );
}

export function HeroFigure({
  label,
  value,
  current,
  previous,
  hint,
}: {
  label: string;
  value: string;
  current?: number;
  previous?: number;
  hint?: string;
}) {
  const change =
    current !== undefined && previous !== undefined
      ? delta(current, previous)
      : null;

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm">
      <span className="text-xs text-gray-500">{label}</span>
      {/* Proportional figures, not tabular: equal-width digits read loose at display sizes. */}
      <p className="mt-1 text-5xl leading-tight font-semibold text-navy">
        {value}
      </p>
      <div className="mt-2 flex flex-wrap items-center gap-2">
        {change !== null && <DeltaBadge change={change} />}
        {hint && <span className="text-xs text-gray-400">{hint}</span>}
      </div>
    </div>
  );
}

function DeltaBadge({
  change,
  invert = false,
}: {
  change: number;
  invert?: boolean;
}) {
  const rising = change >= 0;
  const good = invert ? !rising : rising;
  const Icon = rising ? ArrowUp : ArrowDown;

  return (
    // The arrow and the sign carry the meaning; colour only reinforces it.
    <span
      className={cn(
        "flex items-center gap-0.5 text-xs font-medium tabular-nums",
        good ? "text-chart-good" : "text-chart-critical",
      )}
    >
      <Icon className="size-3" />
      {percent(Math.abs(change))}
      <span className="text-gray-400">מול התקופה הקודמת</span>
    </span>
  );
}

export function Meter({
  label,
  value,
  total,
  valueLabel,
}: {
  label: string;
  value: number;
  total: number;
  valueLabel: string;
}) {
  const ratio = total > 0 ? Math.min(Math.max(value / total, 0), 1) : 0;

  return (
    <div>
      <div className="flex items-baseline justify-between gap-2">
        <span className="text-xs text-gray-500">{label}</span>
        <span className="text-xs font-medium text-gray-700 tabular-nums">
          {valueLabel} · {percent(ratio * 100, 0)}
        </span>
      </div>
      <div
        className="mt-1.5 h-2 w-full overflow-hidden rounded-full bg-gray-100"
        role="meter"
        aria-valuenow={Math.round(ratio * 100)}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label={label}
      >
        <div
          className="h-full rounded-full bg-chart-1"
          style={{ width: `${ratio * 100}%` }}
        />
      </div>
    </div>
  );
}
