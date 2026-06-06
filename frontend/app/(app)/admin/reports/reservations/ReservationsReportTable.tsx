"use client";

import { businessReservationReport } from "@/shared/api/reservations-api";
import { reservation } from "@/shared/client";

import { ReportTableShell } from "../_components/ReportTableShell";
import { ReservationReportFilters } from "../_components/ReservationsReportFilterBar";
import { makeBaseColumns } from "../_components/reportTableUtils";

type ReportRow = reservation.BusinessReservationReportRow;

const columns = makeBaseColumns<ReportRow>();

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
      columns={columns}
      queryKey="reservations-report"
      queryFn={businessReservationReport}
      showFilters={showFilters}
      fixedFilters={fixedFilters}
    />
  );
}

