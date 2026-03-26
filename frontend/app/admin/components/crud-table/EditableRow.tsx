import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Pencil, Trash2, Check, X } from "lucide-react";
import { ZodType } from "zod";
import { FieldValues } from "react-hook-form";
import { ColumnDef } from "./types";
import { CellInput } from "./CellInput";

interface EditableRowProps<TRow> {
  row: TRow;
  columns: ColumnDef<TRow>[];
  isEditing: boolean;
  isPending: boolean;
  onEdit: () => void;
  onCancel: () => void;
  onSave: (data: Record<string, unknown>) => void;
  onDelete: () => void;
  schema: ZodType<FieldValues>;
  selected?: boolean;
  onToggleSelect?: () => void;
}

function formatCellValue<TRow>(row: TRow, column: ColumnDef<TRow>): string {
  const value = (row as Record<string, unknown>)[column.key];
  if (value == null) return "";
  if (column.type === "checkbox") return "";
  return String(value);
}

export function EditableRow<TRow>({
  row,
  columns,
  isEditing,
  isPending,
  onEdit,
  onCancel,
  onSave,
  onDelete,
  schema,
  selected,
  onToggleSelect,
}: EditableRowProps<TRow>) {
  const editableColumns = columns.filter((c) => c.editable !== false);

  const defaults: Record<string, unknown> = {};
  for (const col of editableColumns) {
    defaults[col.key] = (row as Record<string, unknown>)[col.key] ?? "";
  }

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(schema as any),
    defaultValues: defaults as Record<string, string>,
  });

  if (isEditing) {
    return (
      <tr className="bg-blue-50/50">
        <td className="px-3 py-2 w-10">
          <input
            type="checkbox"
            className="h-4 w-4 accent-blue-600"
            checked={selected}
            onChange={onToggleSelect}
          />
        </td>
        {columns.map((col) => (
          <td key={col.key} className="px-3 py-2">
            {col.editable === false ? (
              <span className="text-sm text-gray-500">
                {formatCellValue(row, col)}
              </span>
            ) : (
              <div className="relative pb-4">
                <CellInput
                  column={col}
                  register={register}
                  name={col.key as string}
                />
                {errors[col.key] && (
                  <span
                    className="absolute right-0 bottom-0 text-red-500 text-xs whitespace-nowrap"
                    dir="rtl"
                  >
                    {errors[col.key]?.message as string}
                  </span>
                )}
              </div>
            )}
          </td>
        ))}
        <td className="px-3 py-2">
          <div className="flex items-center gap-1">
            <button
              type="button"
              disabled={isPending}
              onClick={handleSubmit((data) =>
                onSave(data as Record<string, unknown>),
              )}
              className="p-1 text-green-600 hover:text-green-800 disabled:opacity-50 cursor-pointer"
            >
              <Check size={16} />
            </button>
            <button
              type="button"
              disabled={isPending}
              onClick={onCancel}
              className="p-1 text-gray-500 hover:text-gray-700 disabled:opacity-50 cursor-pointer"
            >
              <X size={16} />
            </button>
          </div>
        </td>
      </tr>
    );
  }

  return (
    <tr className="hover:bg-gray-50 border-b border-gray-100">
      <td className="px-3 py-2 w-10">
        <input
          type="checkbox"
          className="h-4 w-4 accent-blue-600"
          checked={selected}
          onChange={onToggleSelect}
        />
      </td>
      {columns.map((col) => (
        <td key={col.key} className="px-3 py-2 text-sm">
          {col.type === "checkbox" ? (
            <input
              type="checkbox"
              checked={!!(row as Record<string, unknown>)[col.key]}
              disabled
              className="h-4 w-4 accent-blue-600"
            />
          ) : (
            formatCellValue(row, col)
          )}
        </td>
      ))}
      <td className="px-3 py-2">
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={onEdit}
            className="p-1 text-gray-400 hover:text-blue-600 cursor-pointer"
          >
            <Pencil size={14} />
          </button>
          <button
            type="button"
            onClick={onDelete}
            className="p-1 text-gray-400 hover:text-red-600 cursor-pointer"
          >
            <Trash2 size={14} />
          </button>
        </div>
      </td>
    </tr>
  );
}
