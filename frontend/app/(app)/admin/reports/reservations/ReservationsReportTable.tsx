"use client";

import { businessReservationReport } from "@/shared/api/reservations-api";
import { reports } from "@/shared/client";

import { ReportTableShell } from "../_components/ReportTableShell";
import { ReservationReportFilters } from "../_components/ReservationsReportFilterBar";
import { limitColumns, makeBaseColumns } from "../_components/reportTableUtils";

type ReportRow = reports.BusinessReservationReportRow;

const columns = makeBaseColumns<ReportRow>();
const limitedColumns = limitColumns(columns)

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

