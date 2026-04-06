"use client";

import { OrgCombobox } from "@/app/(app)/admin/components/OrgCombobox";
import { OfficeCombobox } from "@/app/(app)/admin/components/OfficeCombobox";

interface ContactsFilters {
  search: string;
  orgId: string;
  officeId: string;
}

interface ContactsFilterBarProps {
  filters: ContactsFilters;
  onChange: (updates: Partial<ContactsFilters>) => void;
}

export function ContactsFilterBar({
  filters,
  onChange,
}: ContactsFilterBarProps) {
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
      <div className="min-w-48">
        <label className="block text-xs text-gray-500 mb-1">משרד</label>
        <OfficeCombobox
          value={filters.officeId ? Number(filters.officeId) : 0}
          onChange={(v) => onChange({ officeId: v ? String(v) : "" })}
          showClear
          placeholder="כל המשרדים"
        />
      </div>
    </div>
  );
}
