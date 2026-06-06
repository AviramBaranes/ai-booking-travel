"use client";

import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";

import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { reservation } from "@/shared/client";

import {
  RESERVATION_REPORT_FILTER_KEYS,
  ReservationReportFilterKey,
  ReservationReportFilters,
  ReservationsReportFilterBar,
  emptyReservationReportFilters,
  ReportPageSize,
} from "./ReservationsReportFilterBar";
import { ReportColumn, buildPageNumbers, buildRequest } from "./reportTableUtils";

interface ReportTableShellProps<T extends { reservationId: number }> {
  columns: ReportColumn<T>[];
  queryKey: string;
  queryFn: (
    params: reservation.ReportParams,
  ) => Promise<{ reservations: T[]; total: number }>;
  showStatusFilter?: boolean;
  showFilters?: boolean;
  fixedFilters?: Partial<ReservationReportFilters>;
}

export function ReportTableShell<T extends { reservationId: number }>({
  columns,
  queryKey,
  queryFn,
  showStatusFilter = true,
  showFilters = true,
  fixedFilters,
}: ReportTableShellProps<T>) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState<ReportPageSize>(25);
  const [urlFilters, setUrlFilters] = useUrlFilters<ReservationReportFilterKey>([
    ...RESERVATION_REPORT_FILTER_KEYS,
  ]);
  const filters = showFilters ? urlFilters : emptyReservationReportFilters;
  const effectiveFilters = {
    ...filters,
    ...fixedFilters,
  } satisfies ReservationReportFilters;

  const filterSignature = RESERVATION_REPORT_FILTER_KEYS.map(
    (key) => effectiveFilters[key],
  ).join("|");

  const reportQuery = useQuery({
    queryKey: [queryKey, page, pageSize, filterSignature],
    queryFn: () => queryFn(buildRequest(page, pageSize, effectiveFilters)),
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
    setUrlFilters(nextFilters);
  }

  function handlePageSizeChange(nextSize: ReportPageSize) {
    setPageSize(nextSize);
    setPage(1);
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
      {showFilters && (
        <div className="border-b border-gray-200 bg-gray-50/80 px-4 py-3">
          <ReservationsReportFilterBar
            key={filterSignature}
            initialFilters={filters}
            onSubmit={handleFilterSubmit}
            pageSize={pageSize}
            onPageSizeChange={handlePageSizeChange}
            showStatusFilter={showStatusFilter}
          />
        </div>
      )}

      {reportQuery.isError && (
        <div className="mx-4 mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          לא הצלחנו לטעון את הדוח. נסו שוב בעוד רגע.
        </div>
      )}

      <div className="relative">
        {isRefetching && (
          <div className="absolute inset-x-0 top-0 z-10 h-0.5 overflow-hidden bg-navy/10">
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
                      "font-bold",
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
