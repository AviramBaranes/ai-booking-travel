import { ReactNode } from "react";
import { ZodType } from "zod";
import { FieldValues } from "react-hook-form";

export type ColumnType = "text" | "number" | "date" | "checkbox" | "select";

export interface SelectOption {
  label: string;
  value: string;
}

export interface ColumnDef<TRow> {
  key: keyof TRow & string;
  label: string;
  type: ColumnType;
  options?: SelectOption[];
  editable?: boolean;
  sortable?: boolean;
}

export interface SortState {
  key: string;
  dir: "asc" | "desc";
}

export interface BulkAction {
  label: string;
  icon?: ReactNode;
  onClick: (ids: number[]) => void | Promise<void>;
  variant?: "danger" | "default";
}

export interface CrudTableProps<TRow, TCreate, TUpdate> {
  columns: ColumnDef<TRow>[];
  queryKey: string;
  queryKeyDeps?: unknown[];
  getId: (row: TRow) => number;
  listFn: (sort: SortState | null, page: number) => Promise<unknown>;
  extractList: (response: unknown) => TRow[];
  extractTotal?: (response: unknown) => number;
  createFn: (data: TCreate) => Promise<unknown>;
  updateFn: (id: number, data: TUpdate) => Promise<unknown>;
  deleteFn: (id: number) => Promise<unknown>;
  createSchema: ZodType<TCreate & FieldValues>;
  updateSchema: ZodType<TUpdate & FieldValues>;
  selectable?: boolean;
  bulkActions?: BulkAction[];
  pageSize?: number;
  filterSlot?: ReactNode;
}
