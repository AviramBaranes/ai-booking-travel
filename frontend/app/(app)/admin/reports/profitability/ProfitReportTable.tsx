"use client";

import { profitabilityReport } from "@/shared/api/reservations-api";
import { reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";

import { ReportTableShell } from "../_components/ReportTableShell";
import { makeBaseColumns, MoneyCell, ReportColumn } from "../_components/reportTableUtils";

// The generated ProfitReportRow only contains the 4 extra fields because the Encore
// client generator doesn't flatten Go embedded structs. At runtime the JSON response
// is flat and includes all BusinessReservationReportRow fields, so we model the full
// shape here via an intersection.
type ProfitRow = reservation.BusinessReservationReportRow & reservation.ProfitReportRow;

const profitExtraColumns: ReportColumn<ProfitRow>[] = [
  {
    key: "purchasePrice",
    label: "מחיר קנייה",
    className: "min-w-36 tabular-nums",
    render: (row) => formatPrice(row.purchasePrice, row.currencyCode),
  },
  {
    key: "purchasePriceInILS",
    label: 'מחיר קנייה בש"ח',
    className: "min-w-44",
    render: (row) => <MoneyCell value={row.purchasePriceInILS} currency="ILS" />,
  },
  {
    key: "profit",
    label: "רווח",
    className: "min-w-36 tabular-nums",
    render: (row) => formatPrice(row.profit, row.currencyCode),
  },
  {
    key: "profitInILS",
    label: 'רווח בש"ח',
    className: "min-w-44 bg-emerald-50/70",
    headerClassName: "bg-emerald-100/80 text-emerald-900",
    render: (row) => <MoneyCell value={row.profitInILS} currency="ILS" strong />,
  },
];

const columns: ReportColumn<ProfitRow>[] = [
  ...makeBaseColumns<ProfitRow>(),
  ...profitExtraColumns,
];

export default function ProfitReportTable() {
  return (
    <ReportTableShell
      columns={columns}
      queryKey="profit-report"
      queryFn={
        profitabilityReport as (
          params: reservation.ReportParams,
        ) => Promise<{ reservations: ProfitRow[]; total: number }>
      }
    />
  );
}
