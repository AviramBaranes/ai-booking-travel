"use client";

import { useMemo } from "react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { ReportDateRangeInput } from "./ReportDateRangeInput";
import {
  ReportAccountEntity,
  ReportAccountEntityCombobox,
} from "./ReportAccountEntityCombobox";
import { ReportSupplierCombobox } from "./ReportSupplierCombobox";

export const RESERVATION_REPORT_FILTER_KEYS = [
  "pickupFrom",
  "pickupTo",
  "createdFrom",
  "createdTo",
  "voucheredFrom",
  "voucheredTo",
  "status",
  "broker",
  "supplier",
  "organizationId",
  "officeId",
  "agentId",
] as const;

export type ReservationReportFilterKey =
  (typeof RESERVATION_REPORT_FILTER_KEYS)[number];

export type ReservationReportFilters = Record<ReservationReportFilterKey, string>;

interface ReservationsReportFilterBarProps {
  initialFilters: ReservationReportFilters;
  onSubmit: (filters: ReservationReportFilters) => void;
}

export const emptyReservationReportFilters: ReservationReportFilters = {
  pickupFrom: "",
  pickupTo: "",
  createdFrom: "",
  createdTo: "",
  voucheredFrom: "",
  voucheredTo: "",
  status: "",
  broker: "",
  supplier: "",
  organizationId: "",
  officeId: "",
  agentId: "",
};

const inputClass =
  "h-9 border border-gray-300 rounded px-2 py-1.5 text-sm w-full bg-white";

function selectedEntity(filters: ReservationReportFilters): ReportAccountEntity | null {
  if (filters.organizationId) {
    return {
      kind: "organization",
      id: Number(filters.organizationId),
      name: `רשת #${filters.organizationId}`,
    };
  }
  if (filters.officeId) {
    return {
      kind: "office",
      id: Number(filters.officeId),
      name: `משרד #${filters.officeId}`,
    };
  }
  if (filters.agentId) {
    return {
      kind: "agent",
      id: Number(filters.agentId),
      name: `סוכן #${filters.agentId}`,
    };
  }
  return null;
}

function entityUpdate(entity: ReportAccountEntity | null) {
  if (!entity) {
    return { organizationId: "", officeId: "", agentId: "" };
  }

  return {
    organizationId: entity.kind === "organization" ? String(entity.id) : "",
    officeId: entity.kind === "office" ? String(entity.id) : "",
    agentId: entity.kind === "agent" ? String(entity.id) : "",
  };
}

export function ReservationsReportFilterBar({
  initialFilters,
  onSubmit,
}: ReservationsReportFilterBarProps) {
  const [filters, setFilters] = useState(initialFilters);
  const entity = useMemo(() => selectedEntity(filters), [filters]);

  function updateFilters(updates: Partial<ReservationReportFilters>) {
    setFilters((current) => ({ ...current, ...updates }));
  }

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onSubmit(filters);
  }

  function handleReset() {
    setFilters(emptyReservationReportFilters);
    onSubmit(emptyReservationReportFilters);
  }

  return (
    <form className="flex items-end gap-3 flex-wrap" onSubmit={handleSubmit}>
      <ReportDateRangeInput
        label="תאריך איסוף"
        fromValue={filters.pickupFrom}
        toValue={filters.pickupTo}
        onChange={(range) =>
          updateFilters({ pickupFrom: range.from, pickupTo: range.to })
        }
      />
      <ReportDateRangeInput
        label="תאריך יצירה"
        fromValue={filters.createdFrom}
        toValue={filters.createdTo}
        onChange={(range) =>
          updateFilters({ createdFrom: range.from, createdTo: range.to })
        }
      />
      <ReportDateRangeInput
        label="תאריך כרטוס"
        fromValue={filters.voucheredFrom}
        toValue={filters.voucheredTo}
        onChange={(range) =>
          updateFilters({ voucheredFrom: range.from, voucheredTo: range.to })
        }
      />
      <div className="min-w-36">
        <label className="block text-xs text-gray-500 mb-1">סטטוס</label>
        <select
          className={inputClass}
          value={filters.status}
          onChange={(event) => updateFilters({ status: event.target.value })}
        >
          <option value="">הכל</option>
          <option value="booked">הוזמן</option>
          <option value="vouchered">הופק שובר</option>
          <option value="canceled">בוטל</option>
        </select>
      </div>
      <div className="min-w-36">
        <label className="block text-xs text-gray-500 mb-1">ברוקר</label>
        <select
          className={inputClass}
          value={filters.broker}
          onChange={(event) => updateFilters({ broker: event.target.value })}
        >
          <option value="">הכל</option>
          <option value="flex">Flex</option>
          <option value="hertz">Hertz</option>
        </select>
      </div>
      <div className="min-w-40">
        <label className="block text-xs text-gray-500 mb-1">ספק</label>
        <ReportSupplierCombobox
          value={filters.supplier}
          onChange={(supplier) => updateFilters({ supplier })}
        />
      </div>
      <div className="min-w-56">
        <label className="block text-xs text-gray-500 mb-1">ישות</label>
        <ReportAccountEntityCombobox
          value={entity}
          onChange={(nextEntity) => updateFilters(entityUpdate(nextEntity))}
        />
      </div>
      <div className="flex gap-2">
        <Button type="submit" className="h-9 bg-navy px-5 hover:bg-dark-navy">
          הפעל סינון
        </Button>
        <Button type="button" variant="outline" className="h-9" onClick={handleReset}>
          נקה
        </Button>
      </div>
    </form>
  );
}
