"use client";

import { OrgCombobox } from "@/app/(app)/admin/_components/OrgCombobox";

interface OfficesFilters {
  search: string;
  orgId: string;
}

interface OfficesFilterBarProps {
  filters: OfficesFilters;
  onChange: (updates: Partial<OfficesFilters>) => void;
}

export function OfficesFilterBar({ filters, onChange }: OfficesFilterBarProps) {
  const inputClass =
    "border border-gray-300 rounded px-2 py-1.5 text-sm w-full";

  return (
    <div className="flex items-end gap-3 flex-wrap">
      <div>
        <label className="block text-xs text-gray-500 mb-1">חיפוש</label>
        <input
          type="text"
          className={inputClass}
          value={filters.search}
          onChange={(e) => onChange({ search: e.target.value })}
          placeholder="חיפוש לפי שם"
        />
      </div>
      <div className="min-w-48">
        <label className="block text-xs text-gray-500 mb-1">רשת</label>
        <OrgCombobox
          value={filters.orgId ? Number(filters.orgId) : 0}
          onChange={(v) => onChange({ orgId: v ? String(v) : "" })}
          showClear
          placeholder="כל הרשתות"
        />
      </div>
    </div>
  );
}
