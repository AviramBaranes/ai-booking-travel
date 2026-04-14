"use client";

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import type {
  PlansCountFilter,
  AddOnsFilter,
} from "../../_hooks/useDevFilters";
import { Bug } from "lucide-react";
import { useState } from "react";

interface DevFiltersProps {
  plansCountFilter: PlansCountFilter;
  addOnsFilter: AddOnsFilter;
  onPlansCountChange: (value: Exclude<PlansCountFilter, null>) => void;
  onAddOnsChange: (value: Exclude<AddOnsFilter, null>) => void;
}

const PLANS_OPTIONS = [
  { value: "2", label: "2 תוכניות" },
  { value: "3", label: "3 תוכניות" },
] as const;

const ADDONS_OPTIONS = [
  { value: "has", label: "יש תוספות" },
  { value: "not", label: "אין תוספות" },
] as const;

function FilterGroup({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
        {title}
      </p>
      <div className="flex gap-2">{children}</div>
    </div>
  );
}

function ToggleChip({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "rounded-full border px-3 py-1 text-sm font-medium transition-colors",
        active
          ? "border-brand bg-brand text-white"
          : "border-border bg-background text-foreground hover:border-brand/50",
      )}
    >
      {label}
    </button>
  );
}

export function DevFilters({
  plansCountFilter,
  addOnsFilter,
  onPlansCountChange,
  onAddOnsChange,
}: DevFiltersProps) {
  const [open, setOpen] = useState(false);
  const hasActive = plansCountFilter !== null || addOnsFilter !== null;

  return (
    <div className="fixed bottom-20 z-50 right-4">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <button
            type="button"
            aria-label="פילטרים לפיתוח"
            className={cn(
              "flex size-11 items-center justify-center rounded-full shadow-lg transition-colors",
              hasActive
                ? "bg-brand text-white"
                : "bg-white text-brand border border-brand/30",
            )}
          >
            <Bug size={20} />
          </button>
        </PopoverTrigger>

        <PopoverContent
          side="top"
          align="end"
          sideOffset={10}
          className="w-60 p-4"
        >
          <div className="mb-3 flex items-center gap-2 border-b pb-3">
            <Bug size={15} className="text-brand" />
            <span className="text-sm font-semibold text-navy">
              פילטרים לפיתוח
            </span>
          </div>

          <div className="flex flex-col gap-4">
            <FilterGroup title="מספר תוכניות">
              {PLANS_OPTIONS.map(({ value, label }) => (
                <ToggleChip
                  key={value}
                  label={label}
                  active={plansCountFilter === value}
                  onClick={() => onPlansCountChange(value)}
                />
              ))}
            </FilterGroup>

            <FilterGroup title="תוספות">
              {ADDONS_OPTIONS.map(({ value, label }) => (
                <ToggleChip
                  key={value}
                  label={label}
                  active={addOnsFilter === value}
                  onClick={() => onAddOnsChange(value)}
                />
              ))}
            </FilterGroup>
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
}
