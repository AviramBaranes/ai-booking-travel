"use client";

import { cn } from "@/lib/utils";

export interface SegmentedOption {
  value: string;
  label: string;
}

export function SegmentedControl({
  value,
  onChange,
  options,
  className,
}: {
  value: string;
  onChange: (value: string) => void;
  options: SegmentedOption[];
  className?: string;
}) {
  return (
    <div
      role="group"
      className={cn(
        "inline-flex rounded-md border border-gray-200 bg-gray-50 p-0.5",
        className,
      )}
    >
      {options.map((option) => {
        const isActive = option.value === value;
        return (
          <button
            key={option.value}
            type="button"
            aria-pressed={isActive}
            onClick={() => onChange(option.value)}
            className={cn(
              "rounded px-2.5 py-1 text-xs transition-colors",
              isActive
                ? "bg-white font-semibold text-navy shadow-sm"
                : "text-gray-600 hover:text-gray-900",
            )}
          >
            {option.label}
          </button>
        );
      })}
    </div>
  );
}
