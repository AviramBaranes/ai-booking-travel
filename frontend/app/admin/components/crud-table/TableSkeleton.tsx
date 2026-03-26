import { ColumnDef } from "./types";

interface TableSkeletonProps<TRow> {
  columns: ColumnDef<TRow>[];
  rows?: number;
}

export function TableSkeleton<TRow>({
  columns,
  rows = 16,
}: TableSkeletonProps<TRow>) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-right">
        <thead>
          <tr className="border-b border-gray-200 bg-gray-50">
            {columns.map((col) => (
              <th
                key={col.key}
                className="px-3 py-2 text-xs font-semibold text-gray-600 uppercase"
              >
                {col.label}
              </th>
            ))}
            <th className="px-3 py-2 w-20" />
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }, (_, i) => (
            <tr key={i} className="border-b border-gray-100">
              {columns.map((col) => (
                <td key={col.key} className="px-3 py-2.5">
                  <div
                    className={`h-4 rounded bg-gray-200 animate-pulse ${
                      col.type === "checkbox" ? "w-4" : "w-full max-w-30"
                    }`}
                    style={{ animationDelay: `${i * 75}ms` }}
                  />
                </td>
              ))}
              <td className="px-3 py-2.5">
                <div className="flex items-center gap-1">
                  <div className="h-4 w-4 rounded bg-gray-200 animate-pulse" />
                  <div className="h-4 w-4 rounded bg-gray-200 animate-pulse" />
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
