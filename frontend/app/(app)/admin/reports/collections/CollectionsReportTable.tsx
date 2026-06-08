"use client";

import { useQuery } from "@tanstack/react-query";

import { cn } from "@/lib/utils";
import { collectionsReport } from "@/shared/api/reservations-api";
import { reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";

import { MoneyCell, ReportColumn } from "../_components/reportTableUtils";
import { CollectionsDetailDialog } from "../_components/CollectionsDetailDialog";
import { BillingEntity } from "@/app/(app)/accounting/billing/_components/BillingEntityCombobox";
import { useState } from "react";

type CollectionRow = reservation.BusinessesBalancesReportRow;

const columns: ReportColumn<CollectionRow>[] = [
  {
    key: "billingEntityName",
    label: "גורם משלם",
    className: "font-medium text-gray-900",
    render: (row) => row.billingEntityName || "-",
  },
  {
    key: "totalInEuro",
    label: "סהכ הזמנות ביורו",
    className: "text-center",
    render: (row) => formatPrice(row.totalOpenBalanceInEuro, "EUR"),
  },
  {
    key: "totalInDollar",
    label: "סהכ הזמנות בדולר",
    className: "text-center",
    render: (row) => formatPrice(row.totalOpenBalanceInDollar, "USD"),
  },
  {
    key: "totalInOtherCurrency",
    label: "סהכ הזמנות שאינן בדולר או ביורו",
    className: "",
    render: (row) => formatPrice(row.totalInOtherCurrency, "ILS"),
  },
];

export default function CollectionsReportTable() {
  const [entity, setEntity] = useState<BillingEntity | null>(null);
  const reportQuery = useQuery({
    queryKey: ["collections-report"],
    queryFn: collectionsReport,
  });

  const rows = reportQuery.data?.businesses ?? [];
  const total = reportQuery.data?.total ?? 0;

  return (
    <>
      <CollectionsDetailDialog
        entity={entity}
        onClose={() => setEntity(null)}
      />
      <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
        {reportQuery.isError && (
          <div className="mx-4 mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
            לא הצלחנו לטעון את הדוח. נסו שוב בעוד רגע.
          </div>
        )}

        <div className="overflow-x-auto">
          <table className="w-full table-fixed text-right">
            <thead>
              <tr className="border-b border-gray-200 bg-slate-100 text-xs font-bold text-slate-800">
                {columns.map((column) => (
                  <th
                    key={column.key}
                    className={cn(
                      "whitespace-normal px-3 py-3 align-middle leading-snug",
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
            <tbody className="divide-y divide-gray-100 text-sm">
              {reportQuery.isLoading ? (
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
                    לא נמצאו עסקים עבור הדוח
                  </td>
                </tr>
              ) : (
                rows.map((row) => (
                  <tr
                    key={`${row.billingEntityType}-${row.billingEntityId}`}
                    className="bg-white transition-colors odd:bg-white even:bg-slate-50/50 hover:bg-navy/5"
                    onClick={() =>
                      setEntity({
                        id: row.billingEntityId,
                        kind:
                          row.billingEntityType === "organization"
                            ? "org"
                            : "office",
                        name: row.billingEntityName,
                      })
                    }
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

        <div className="border-t border-gray-200 px-4 py-3 text-sm text-gray-600">
          {total > 0 ? `${total} תוצאות` : "0 תוצאות"}
        </div>
      </div>
    </>
  );
}
