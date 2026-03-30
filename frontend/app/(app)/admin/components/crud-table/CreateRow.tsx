import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus } from "lucide-react";
import { ZodType } from "zod";
import { FieldValues } from "react-hook-form";
import { ColumnDef } from "./types";
import { CellInput } from "./CellInput";

interface CreateRowProps<TRow> {
  columns: ColumnDef<TRow>[];
  schema: ZodType<FieldValues>;
  isPending: boolean;
  onSubmit: (data: Record<string, unknown>) => void;
}

export function CreateRow<TRow>({
  columns,
  schema,
  isPending,
  onSubmit,
}: CreateRowProps<TRow>) {
  const editableColumns = columns.filter((c) => c.editable !== false);

  const defaults: Record<string, unknown> = {};
  for (const col of editableColumns) {
    switch (col.type) {
      case "number":
        defaults[col.key] = 0;
        break;
      case "checkbox":
        defaults[col.key] = false;
        break;
      default:
        defaults[col.key] = "";
    }
  }

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(schema as any),
    defaultValues: defaults as Record<string, string>,
  });

  const submit = (data: FieldValues) => {
    onSubmit(data as Record<string, unknown>);
    reset();
  };

  return (
    <tr className="bg-green-50/50 border-t-2 border-gray-200">
      <td className="px-3 py-2 w-10" />
      {columns.map((col) => (
        <td key={col.key} className="px-3 py-2">
          {col.editable === false ? (
            <span className="text-sm text-gray-400">—</span>
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
        <button
          type="button"
          disabled={isPending}
          onClick={handleSubmit(submit)}
          className="p-1 text-green-600 hover:text-green-800 disabled:opacity-50 cursor-pointer"
        >
          <Plus size={16} />
        </button>
      </td>
    </tr>
  );
}
