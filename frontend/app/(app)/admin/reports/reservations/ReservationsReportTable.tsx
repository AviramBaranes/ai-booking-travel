"use client";

import { businessReservationReport } from "@/shared/api/reservations-api";
import { reservation } from "@/shared/client";

import { ReportTableShell } from "../_components/ReportTableShell";
import { makeBaseColumns } from "../_components/reportTableUtils";

type ReportRow = reservation.BusinessReservationReportRow;

const columns = makeBaseColumns<ReportRow>();

export default function ReservationsReportTable() {
  return (
    <ReportTableShell
      columns={columns}
      queryKey="reservations-report"
      queryFn={businessReservationReport}
    />
  );
}

