"use client";

import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { ReactNode, useMemo, useState } from "react";

import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { businessReservationReport } from "@/shared/api/reservations-api";
import { reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";
import {
  RESERVATION_REPORT_FILTER_KEYS,
  ReservationReportFilterKey,
  ReservationReportFilters,
  ReservationsReportFilterBar,
  ReportPageSize,
} from "../_components/ReservationsReportFilterBar";

const statusLabels: Record<string, string> = {
  booked: "הוזמן",
  vouchered: "הופק שובר",
  canceled: "בוטל",
};

const statusClasses: Record<string, string> = {
  booked: "bg-navy/5 text-navy ring-navy/15",
  vouchered: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  canceled: "bg-rose-50 text-rose-700 ring-rose-200",
};

type ReportRow = reservation.BusinessReservationReportRow;

type ReportColumn = {
  key: string;
  label: string;
  className: string;
  headerClassName?: string;
  render: (row: ReportRow) => ReactNode;
};

function toNumber(value: string) {
  return value ? Number(value) : undefined;
}

function buildRequest(
  page: number,
  pageSize: number,
  filters: ReservationReportFilters,
): reservation.ReportParams {
  return {
    Page: page,
    PageSize: pageSize,
    PickupDateFrom: filters.pickupFrom || undefined,
    PickupDateTo: filters.pickupTo || undefined,
    CreatedDateFrom: filters.createdFrom || undefined,
    CreatedDateTo: filters.createdTo || undefined,
    VoucheredAtFrom: filters.voucheredFrom || undefined,
    VoucheredAtTo: filters.voucheredTo || undefined,
    Status: filters.status || undefined,
    Broker: filters.broker || undefined,
    Supplier: filters.supplier || undefined,
    OrganizationID: toNumber(filters.organizationId),
    OfficeID: toNumber(filters.officeId),
    AgentID: toNumber(filters.agentId),
    IsBusiness: true,
  };
}

function formatDate(value?: string) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString("he-IL");
}

function formatDateTime(value?: string) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("he-IL", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={cn(
        "inline-flex rounded-full px-2.5 py-1 text-xs font-semibold ring-1",
        statusClasses[status] ?? "bg-gray-50 text-gray-700 ring-gray-200",
      )}
    >
      {statusLabels[status] ?? status}
    </span>
  );
}

function MoneyCell({
  value,
  currency,
  strong = false,
}: {
  value: number;
  currency: string;
  strong?: boolean;
}) {
  return (
    <span
      className={cn(
        "font-mono tabular-nums",
        strong ? "font-bold text-emerald-700" : "text-gray-800",
      )}
    >
      {formatPrice(value, currency)}
    </span>
  );
}

function buildPageNumbers(current: number, total: number): (number | "...")[] {
  if (total <= 7) return Array.from({ length: total }, (_, index) => index + 1);

  const pages: (number | "...")[] = [1];
  if (current > 3) pages.push("...");

  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let page = start; page <= end; page++) pages.push(page);

  if (current < total - 2) pages.push("...");
  pages.push(total);
  return pages;
}

const columns: ReportColumn[] = [
  {
    key: "reservationId",
    label: "מזהה",
    className: "min-w-24 font-mono text-gray-500",
    render: (row) => row.reservationId,
  },
  {
    key: "brokerReservationId",
    label: "מספר הזמנה",
    className: "min-w-44",
    render: (row) => row.brokerReservationId,
  },
  {
    key: "status",
    label: "סטטוס",
    className: "min-w-32",
    render: (row) => <StatusBadge status={row.status} />,
  },
  {
    key: "organizationName",
    label: "רשת",
    className: "min-w-56 font-medium text-gray-900",
    render: (row) => row.organizationName || "-",
  },
  {
    key: "officeName",
    label: "משרד",
    className: "min-w-52",
    render: (row) => row.officeName || "-",
  },
  {
    key: "agentName",
    label: "סוכן",
    className: "min-w-44",
    render: (row) => row.agentName || "-",
  },
  {
    key: "adminName",
    label: "אדמין",
    className: "min-w-40",
    render: (row) => row.adminName || "-",
  },
  {
    key: "brokerName",
    label: "ברוקר",
    className: "min-w-28",
    render: (row) => (
      <span className="rounded bg-slate-100 px-2 py-1 text-xs font-bold uppercase text-slate-700">
        {row.brokerName}
      </span>
    ),
  },
  {
    key: "supplierName",
    label: "ספק",
    className: "min-w-44",
    render: (row) => row.supplierName || "-",
  },
  {
    key: "countryCode",
    label: "מדינה",
    className: "min-w-24 text-center font-semibold text-gray-600",
    render: (row) => row.countryCode || "-",
  },
  {
    key: "pickupDate",
    label: "איסוף",
    className: "min-w-32 font-medium",
    render: (row) => formatDate(row.pickupDate),
  },
  {
    key: "dropoffDate",
    label: "החזרה",
    className: "min-w-32",
    render: (row) => formatDate(row.dropoffDate),
  },
  {
    key: "rentalDays",
    label: "ימים",
    className: "min-w-20 text-center font-semibold",
    render: (row) => row.rentalDays,
  },
  {
    key: "driverName",
    label: "נהג",
    className: "min-w-52 font-medium text-gray-900",
    render: (row) => row.driverName || "-",
  },
  {
    key: "currencyCode",
    label: "מטבע",
    className: "min-w-24 text-center",
    render: (row) => row.currencyCode,
  },
  {
    key: "currencyRate",
    label: "שער",
    className: "min-w-24 font-mono tabular-nums",
    render: (row) => row.currencyRate,
  },
  {
    key: "carSellPriceWithBrokerERP",
    label: "מחיר רכב",
    className: "min-w-36 font-mono tabular-nums",
    render: (row) => formatPrice(row.carSellPriceWithBrokerERP, row.currencyCode),
  },
  {
    key: "carSellPriceWithBrokerERPInILS",
    label: "מחיר רכב בש״ח",
    className: "min-w-44",
    render: (row) => (
      <MoneyCell value={row.carSellPriceWithBrokerERPInILS} currency="ILS" />
    ),
  },
  {
    key: "btERPPrice",
    label: "BT ERP",
    className: "min-w-32 font-mono tabular-nums",
    render: (row) => formatPrice(row.btERPPrice, row.currencyCode),
  },
  {
    key: "btERPPriceInILS",
    label: "BT ERP בש״ח",
    className: "min-w-40",
    render: (row) => <MoneyCell value={row.btERPPriceInILS} currency="ILS" />,
  },
  {
    key: "totalPrice",
    label: "סה״כ",
    className: "min-w-36 font-mono font-semibold tabular-nums",
    render: (row) => formatPrice(row.totalPrice, row.currencyCode),
  },
  {
    key: "totalPriceInILS",
    label: "סה״כ בש״ח",
    className: "min-w-44 bg-emerald-50/70",
    headerClassName: "bg-emerald-100/80 text-emerald-900",
    render: (row) => <MoneyCell value={row.totalPriceInILS} currency="ILS" strong />,
  },
  {
    key: "voucherNumber",
    label: "שובר",
    className: "min-w-36 font-mono text-gray-600",
    render: (row) => row.voucherNumber || "-",
  },
  {
    key: "voucheredAt",
    label: "כורטס בתאריך",
    className: "min-w-44 text-gray-600",
    render: (row) => formatDateTime(row.voucheredAt),
  },
  {
    key: "createdAt",
    label: "נוצר בתאריך",
    className: "min-w-44 text-gray-600",
    render: (row) => formatDateTime(row.createdAt),
  },
];

export default function ReservationsReportTable() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState<ReportPageSize>(25);
  const [filters, setFilters] = useUrlFilters<ReservationReportFilterKey>([
    ...RESERVATION_REPORT_FILTER_KEYS,
  ]);

  const filterSignature = RESERVATION_REPORT_FILTER_KEYS.map(
    (key) => filters[key],
  ).join("|");

  const reportQuery = useQuery({
    queryKey: ["reservations-report", page, pageSize, filterSignature],
    queryFn: () => businessReservationReport(buildRequest(page, pageSize, filters)),
    placeholderData: keepPreviousData,
  });

  const rows = reportQuery.data?.reservations ?? [];
  const total = reportQuery.data?.total ?? 0;
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : 0;
  const isInitialLoading = reportQuery.isLoading && !reportQuery.isPlaceholderData;
  const isRefetching = reportQuery.isFetching && !isInitialLoading;

  const pageNumbers = useMemo(
    () => buildPageNumbers(page, totalPages),
    [page, totalPages],
  );

  function handleFilterSubmit(nextFilters: ReservationReportFilters) {
    setPage(1);
    setFilters(nextFilters);
  }

  function handlePageSizeChange(nextSize: ReportPageSize) {
    setPageSize(nextSize);
    setPage(1);
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
      <div className="border-b border-gray-200 bg-gray-50/80 px-4 py-3">
        <ReservationsReportFilterBar
          key={filterSignature}
          initialFilters={filters}
          onSubmit={handleFilterSubmit}
          pageSize={pageSize}
          onPageSizeChange={handlePageSizeChange}
        />
      </div>

      {reportQuery.isError && (
        <div className="mx-4 mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          לא הצלחנו לטעון את הדוח. נסו שוב בעוד רגע.
        </div>
      )}

      <div className="relative">
        {isRefetching && (
          <div className="absolute top-0 inset-x-0 z-10 h-0.5 overflow-hidden bg-navy/10">
            <div className="h-full w-1/3 animate-[shimmer_1s_ease-in-out_infinite] bg-navy" />
          </div>
        )}

        <div className="overflow-x-auto">
          <table className="min-w-900 w-full table-fixed text-right">
            <thead>
              <tr className="border-b border-gray-200 bg-slate-100 text-xs font-bold text-slate-800">
                {columns.map((column) => (
                  <th
                    key={column.key}
                    className={cn(
                      "whitespace-nowrap px-3 py-3 align-middle",
                      column.className,
                      column.headerClassName,
                    )}
                  >
                    {column.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody
              className={cn(
                "divide-y divide-gray-100 text-sm",
                isRefetching && "opacity-60",
              )}
            >
              {isInitialLoading ? (
                Array.from({ length: 8 }).map((_, rowIndex) => (
                  <tr key={rowIndex}>
                    {columns.map((column) => (
                      <td key={column.key} className="px-3 py-3">
                        <div className="h-4 w-24 animate-pulse rounded bg-gray-200" />
                      </td>
                    ))}
                  </tr>
                ))
              ) : rows.length === 0 ? (
                <tr>
                  <td
                    colSpan={columns.length}
                    className="px-4 py-12 text-center text-sm text-gray-500"
                  >
                    לא נמצאו הזמנות עבור הסינון הנוכחי
                  </td>
                </tr>
              ) : (
                rows.map((row) => (
                  <tr
                    key={row.reservationId}
                    className="bg-white transition-colors odd:bg-white even:bg-slate-50/50 hover:bg-navy/5"
                  >
                    {columns.map((column) => (
                      <td
                        key={column.key}
                        className={cn(
                          "whitespace-nowrap px-3 py-3 align-middle",
                          column.className,
                        )}
                      >
                        {column.render(row)}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <div className="flex items-center justify-between border-t border-gray-200 px-4 py-3">
        <span className="text-sm text-gray-600">
          {total > 0
            ? `עמוד ${page} מתוך ${totalPages} (${total} תוצאות)`
            : "0 תוצאות"}
        </span>
        {totalPages > 1 && (
          <div className="flex items-center gap-1">
            <Button
              type="button"
              variant="outline"
              size="icon"
              disabled={page <= 1}
              onClick={() => setPage((current) => Math.max(1, current - 1))}
            >
              <ChevronRight className="size-4" />
            </Button>
            {pageNumbers.map((pageNumber, index) =>
              pageNumber === "..." ? (
                <span key={`ellipsis-${index}`} className="px-1 text-gray-400">
                  ...
                </span>
              ) : (
                <Button
                  key={pageNumber}
                  type="button"
                  variant={page === pageNumber ? "default" : "outline"}
                  size="sm"
                  onClick={() => setPage(pageNumber)}
                >
                  {pageNumber}
                </Button>
              ),
            )}
            <Button
              type="button"
              variant="outline"
              size="icon"
              disabled={page >= totalPages}
              onClick={() =>
                setPage((current) => Math.min(totalPages, current + 1))
              }
            >
              <ChevronLeft className="size-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
