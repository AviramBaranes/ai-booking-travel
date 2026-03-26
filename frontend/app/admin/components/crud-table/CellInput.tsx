import { UseFormRegister, FieldValues, Path } from "react-hook-form";
import { ColumnDef } from "./types";

interface CellInputProps<TRow, TForm extends FieldValues> {
  column: ColumnDef<TRow>;
  register: UseFormRegister<TForm>;
  name: Path<TForm>;
}

export function CellInput<TRow, TForm extends FieldValues>({
  column,
  register,
  name,
}: CellInputProps<TRow, TForm>) {
  const base = "border border-gray-300 rounded px-2 py-1 text-sm w-full";

  switch (column.type) {
    case "checkbox":
      return (
        <input
          type="checkbox"
          className="h-4 w-4 accent-blue-600"
          {...register(name)}
        />
      );
    case "select":
      return (
        <select className={base} {...register(name)}>
          {column.options?.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      );
    case "number":
      return (
        <input
          type="number"
          step="any"
          className={base}
          {...register(name, { valueAsNumber: true })}
        />
      );
    case "date":
      return <input type="date" className={base} {...register(name)} />;
    default:
      return <input type="text" className={base} {...register(name)} />;
  }
}
