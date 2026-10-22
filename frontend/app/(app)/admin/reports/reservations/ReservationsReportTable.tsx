"use client";

import LoginAsUserButton from "@/app/(app)/admin/_components/LoginAsUserButton";
import { businessReservationReport } from "@/shared/api/reservations-api";
import { reports } from "@/shared/client";

import { ReportTableShell } from "../_components/ReportTableShell";
import { ReservationReportFilters } from "../_components/ReservationsReportFilterBar";
import {
  LIMITED_COLUMNS_KEYS,
  ReportColumn,
  makeBaseColumns,
} from "../_components/reportTableUtils";

type ReportRow = reports.BusinessReservationReportRow;

const loginAsUserColumn: ReportColumn<ReportRow> = {
  key: "loginAsUser",
  label: "",
  className: "w-8 px-1 text-center",
  skipExport: true,
  render: (row) => (
    // Keep the click from also opening the row's reservation dialog.
    <span onClick={(e) => e.stopPropagation()}>
      <LoginAsUserButton userId={row.userId} label="התחבר כמשתמש" compact />
    </span>
  ),
};

const baseColumns = makeBaseColumns<ReportRow>();
const userNameIndex = baseColumns.findIndex((col) => col.key === "userName");
const columns = [
  ...baseColumns.slice(0, userNameIndex + 1),
  loginAsUserColumn,
  ...baseColumns.slice(userNameIndex + 1),
];
const limitedColumns = columns.filter(
  (col) => col === loginAsUserColumn || LIMITED_COLUMNS_KEYS.includes(col.key),
);

interface ReservationsReportTableProps {
  showFilters?: boolean;
  fixedFilters?: Partial<ReservationReportFilters>;
}

export default function ReservationsReportTable({
  showFilters = true,
  fixedFilters,
}: ReservationsReportTableProps) {
  return (
    <ReportTableShell
      fullColumns={columns}
      limitedColumns={limitedColumns}
      queryKey="reservations-report"
      queryFn={businessReservationReport}
      showFilters={showFilters}
      fixedFilters={fixedFilters}
    />
  );
}

