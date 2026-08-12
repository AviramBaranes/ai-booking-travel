"use client";

import { ReactNode, useState } from "react";
import { Table2, ChartColumnBig } from "lucide-react";

import { cn } from "@/lib/utils";

export interface TableColumn {
  label: string;
  align?: "start" | "end";
}

export interface TableView {
  columns: TableColumn[];
  rows: (string | number)[][];
}

interface ChartCardProps {
  title: string;
  subtitle?: string;
  /** Rendered in the header, e.g. a metric toggle that scopes only this card's encoding. */
  actions?: ReactNode;
  children: ReactNode;
  /**
   * Every chart ships a table twin. Three slots in the palette sit below 3:1 contrast on
   * white, so the table is the WCAG-clean way to read any value — not an optional extra.
   */
  tableView: TableView;
  className?: string;
}

export function ChartCard({
  title,
  subtitle,
  actions,
  children,
  tableView,
  className,
}: ChartCardProps) {
  const [showTable, setShowTable] = useState(false);

  return (
    <section
      className={cn(
        "flex flex-col rounded-lg border border-gray-200 bg-white shadow-sm",
        className,
      )}
    >
      <header className="flex flex-wrap items-start justify-between gap-3 border-b border-gray-100 px-4 py-3">
        <div>
          <h2 className="text-sm font-semibold text-navy">{title}</h2>
          {subtitle && (
            <p className="mt-0.5 text-xs text-gray-500">{subtitle}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
          {actions}
          <button
            type="button"
            onClick={() => setShowTable((current) => !current)}
            aria-pressed={showTable}
            className="flex items-center gap-1 rounded border border-gray-200 px-2 py-1 text-xs text-gray-600 transition-colors hover:bg-gray-50"
            title={showTable ? "הצגת תרשים" : "הצגת טבלה"}
          >
            {showTable ? (
              <ChartColumnBig className="size-3.5" />
            ) : (
              <Table2 className="size-3.5" />
            )}
            {showTable ? "תרשים" : "טבלה"}
          </button>
        </div>
      </header>

      <div className="flex-1 p-4">
        {showTable ? <DataTable {...tableView} /> : children}
      </div>
    </section>
  );
}

function DataTable({ columns, rows }: TableView) {
  if (rows.length === 0) {
    return <p className="py-8 text-center text-sm text-gray-500">אין נתונים</p>;
  }

  return (
    <div className="max-h-72 overflow-auto">
      <table className="w-full text-right text-sm">
        <thead className="sticky top-0 bg-slate-100 text-xs font-bold text-slate-800">
          <tr>
            {columns.map((column) => (
              <th
                key={column.label}
                className={cn(
                  "whitespace-nowrap px-3 py-2",
                  column.align === "end" && "text-left",
                )}
              >
                {column.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {rows.map((row, rowIndex) => (
            <tr key={rowIndex} className="even:bg-slate-50/50">
              {row.map((cell, cellIndex) => (
                <td
                  key={cellIndex}
                  className={cn(
                    "whitespace-nowrap px-3 py-2",
                    columns[cellIndex]?.align === "end" &&
                      "text-left tabular-nums",
                  )}
                >
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
