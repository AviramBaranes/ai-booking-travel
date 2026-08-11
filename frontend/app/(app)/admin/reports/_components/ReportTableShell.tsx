"use client";

import { keepPreviousData, useMutation, useQuery } from "@tanstack/react-query";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";

import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { reports } from "@/shared/client";

import {
  RESERVATION_REPORT_FILTER_KEYS,
  ReservationReportFilterKey,
  ReservationReportFilters,
  ReservationsReportFilterBar,
  emptyReservationReportFilters,
  initialReservationReportFilters,
  ReportPageSize,
} from "./ReservationsReportFilterBar";
import {
  ReportColumn,
  buildPageNumbers,
  buildRequest,
} from "./reportTableUtils";
import { formatPriceFloat } from "@/shared/utils/formatPrice";
import { ReservationDetailDialog } from "@/app/(app)/admin/_components/ReservationDetailDialog";
import { Input } from "@/components/ui/input";
import clsx from "clsx";

export type ReportQueryFn<T> = (params: reports.ReportParams) => Promise<{
  reservations: T[];
  count: number;
  totalSales: number;
  totalProfit?: number;
  profitPercentage?: number;
}>;

interface ReportTableShellProps<T extends { reservationId: number }> {
  fullColumns: ReportColumn<T>[];
  limitedColumns: ReportColumn<T>[];
  queryKey: string;
  queryFn: ReportQueryFn<T>;
  showStatusFilter?: boolean;
  showFilters?: boolean;
  fixedFilters?: Partial<ReservationReportFilters>;
  isProfitReport?: boolean;
}

export function ReportTableShell<T extends { reservationId: number }>({
  fullColumns,
  limitedColumns,
  queryKey,
  queryFn,
  showStatusFilter = true,
  showFilters = true,
  fixedFilters,
  isProfitReport = false,
}: ReportTableShellProps<T>) {
  const [page, setPage] = useState(1);
  const [showLimitedColumns, setShowLimitedColumns] = useState(true);
  const [isExporting, setIsExporting] = useState(false);
  const [selectedReservationId, setSelectedReservationId] = useState<
    number | null
  >(null);
  const [pageSize, setPageSize] = useState<ReportPageSize>(10);
  const [urlFilters, setUrlFilters] = useUrlFilters<ReservationReportFilterKey>(
    [...RESERVATION_REPORT_FILTER_KEYS],
  );

  const didInitDates = useRef(false);
  useEffect(() => {
    if (didInitDates.current) return;
    didInitDates.current = true;
    if (!showFilters) return;
    if (!urlFilters.createdFrom && !urlFilters.createdTo) {
      setUrlFilters({
        createdFrom: initialReservationReportFilters.createdFrom,
        createdTo: initialReservationReportFilters.createdTo,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const urlFiltersWithDefaults: ReservationReportFilters = {
    ...urlFilters,
    skipCanceled: urlFilters.skipCanceled || emptyReservationReportFilters.skipCanceled,
  };
  const filters = showFilters ? urlFiltersWithDefaults : emptyReservationReportFilters;
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
  const count = reportQuery.data?.count ?? 0;
  const totalSales = reportQuery.data?.totalSales ?? 0;
  const totalProfit = reportQuery.data?.totalProfit ?? 0;
  const profitPercentage = reportQuery.data?.profitPercentage ?? 0;
  const totalPages = count > 0 ? Math.ceil(count / pageSize) : 0;
  const isInitialLoading =
    reportQuery.isLoading && !reportQuery.isPlaceholderData;
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

  const { mutateAsync: exportMutate } = useMutation({
    mutationKey: [queryKey, "export", filterSignature],
    mutationFn: () =>
      queryFn(buildRequest(page, pageSize, effectiveFilters, true)),
  });

  const columns = showLimitedColumns ? limitedColumns : fullColumns;

  async function handleExport() {
    setIsExporting(true);
    const { reservations } = await exportMutate();
    const headers = columns.map((col) => col.label).join("\t");
    const rows = reservations
      .map((row) =>
        columns
          .map((col) => {
            const cell = col.exportValue
              ? col.exportValue(row)
              : col.render(row);
            if (typeof cell === "string" || typeof cell === "number") {
              return String(cell).replace(/\t/g, " ");
            }
            return "";
          })
          .join("\t"),
      )
      .join("\n");
    const tsvContent = `${headers}\n${rows}`;
    const blob = new Blob([tsvContent], { type: "text/tab-separated-values" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    const baseName = isProfitReport ? "דוח-רווח" : "דוח-הזמנות";
    const filtersDisplay = Object.entries(effectiveFilters)
      .filter(([_, value]) => value)
      .map(([key, value]) => `${key}-${value}`)
      .join("-");

    link.download = `${baseName}-${filtersDisplay}-${Date.now()}.tsv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    setIsExporting(false);
  }

  return (
    <>
      <ReservationDetailDialog
        reservationId={selectedReservationId}
        onClose={() => setSelectedReservationId(null)}
      />
      <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
        {showFilters && (
          <>
              <Button
                variant="outline"
                className="m-4 p-3 font-normal rounded-lg border-navy"
                onClick={handleExport}
                loading={isExporting}
              >
                ייצוא לאקסל
              </Button>
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
            <div className="w-36 m-3">
              <select
                className="h-9 border border-gray-300 rounded px-2 py-1.5 text-sm w-full bg-white"
                value={showLimitedColumns ? "limited" : "full"}
                onChange={(e) =>
                  setShowLimitedColumns(e.target.value === "limited")
                }
              >
                <option value="limited">תצוגה מוגבלת</option>
                <option value="full">תצוגה מלאה</option>
              </select>
            </div>
          </>
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
            <table className={clsx("w-full table-fixed text-right",{
              "min-w-900":!showLimitedColumns,
              "min-w-350 md:min-w-0":showLimitedColumns
            })}>
              <thead>
                <tr className="border-b border-gray-200 bg-slate-100 text-xs font-bold text-slate-800">
                  {columns.map((column) => (
                    <th
                      key={column.key}
                      className={cn(
                        "whitespace-nowrap px-1 py-3 align-middle",
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
                      onClick={() =>
                        setSelectedReservationId(row.reservationId)
                      }
                      className="cursor-pointer bg-white transition-colors odd:bg-white even:bg-slate-50/50 hover:bg-navy/5"
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

        <div className="flex flex-col lg:flex-row items-center justify-between border-t border-gray-200 px-4 py-3">
          <div className="flex items-start gap-4">
            <span className="text-sm text-gray-600">
              {count > 0
                ? `עמוד ${page} מתוך ${totalPages} (${count} תוצאות)`
                : "0 תוצאות"}
            </span>

            <span className="text-sm text-gray-600">
              סה"כ מכירות: {formatPriceFloat(totalSales, "ILS")}
            </span>
            {totalProfit > 0 && (
              <span className="text-sm text-gray-600">
                סה"כ רווח: {formatPriceFloat(totalProfit, "ILS")} (
                {profitPercentage.toFixed(2)}%)
              </span>
            )}
          </div>

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
                  <span
                    key={`ellipsis-${index}`}
                    className="px-1 text-gray-400"
                  >
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
    </>
  );
}
