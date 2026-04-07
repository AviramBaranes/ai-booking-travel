"use client";

import { useState } from "react";
import {
  useQuery,
  useMutation,
  useQueryClient,
  keepPreviousData,
} from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import {
  ArrowUp,
  ArrowDown,
  ChevronsUpDown,
  ChevronRight,
  ChevronLeft,
  Trash2,
} from "lucide-react";

import { isAppError, AppError } from "@/shared/api/AppError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";

import { CrudTableProps, SortState } from "./types";
import { EditableRow } from "./EditableRow";
import { CreateRow } from "./CreateRow";
import { TableSkeleton } from "./TableSkeleton";

export function CrudTable<
  TRow = Record<string, unknown>,
  TCreate = unknown,
  TUpdate = unknown,
>({
  columns,
  queryKey,
  queryKeyDeps,
  getId,
  listFn,
  extractList,
  extractTotal,
  createFn,
  updateFn,
  deleteFn,
  createSchema,
  updateSchema,
  bulkActions,
  pageSize = 20,
  filterSlot,
  hideCreate,
}: CrudTableProps<TRow, TCreate, TUpdate>) {
  const tErrors = useTranslations("ApiErrors");
  const queryClient = useQueryClient();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [sort, setSort] = useState<SortState | null>(null);
  const [page, setPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());

  const clearError = () => setError(null);

  function handleMutationError(err: unknown) {
    if (isAppError(err)) {
      setError(tErrors((err as AppError).code));
    } else {
      setError(tErrors("internal_error"));
    }
  }

  function handleSort(key: string) {
    setSort((prev) => {
      if (prev?.key === key) {
        if (prev.dir === "asc") return { key, dir: "desc" };
        return null;
      }
      return { key, dir: "asc" };
    });
    setPage(1);
  }

  function toggleSelect(id: number) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function toggleSelectAll() {
    if (selectedIds.size === rows.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(rows.map(getId)));
    }
  }

  const listQuery = useQuery({
    queryKey: [queryKey, sort, page, ...(queryKeyDeps ?? [])],
    queryFn: () => listFn(sort, page),
    placeholderData: keepPreviousData,
  });

  const rows = listQuery.data ? extractList(listQuery.data) : [];
  const total =
    listQuery.data && extractTotal ? extractTotal(listQuery.data) : 0;
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : 0;

  const createMutation = useMutation({
    mutationFn: (data: TCreate) => createFn?.(data) ?? Promise.resolve(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [queryKey] });
      clearError();
    },
    onError: handleMutationError,
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: TUpdate }) =>
      updateFn(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [queryKey] });
      setEditingId(null);
      clearError();
    },
    onError: handleMutationError,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteFn?.(id) ?? Promise.resolve(),
    onSuccess: (_data, deletedId) => {
      queryClient.invalidateQueries({ queryKey: [queryKey] });
      setSelectedIds((prev) => {
        const next = new Set(prev);
        next.delete(deletedId);
        return next;
      });
      clearError();
    },
    onError: handleMutationError,
  });

  const deleteSelectedMutation = useMutation({
    mutationFn: (ids: number[]) =>
      Promise.all(ids.map((id) => deleteFn?.(id) ?? Promise.resolve())),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [queryKey] });
      setSelectedIds(new Set());
      clearError();
    },
    onError: handleMutationError,
  });

  function handleDeleteSelected() {
    const ids = Array.from(selectedIds);
    if (ids.length === 0) return;
    if (!confirm(`האם למחוק ${ids.length} שורות?`)) return;
    deleteSelectedMutation.mutate(ids);
  }

  const isInitialLoading = listQuery.isLoading && !listQuery.isPlaceholderData;
  const isRefetching = !isInitialLoading && listQuery.isFetching;

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 relative">
      {isRefetching && (
        <div className="absolute top-0 inset-x-0 h-0.5 bg-blue-100 overflow-hidden rounded-t-lg z-10">
          <div className="h-full w-1/3 bg-blue-500 animate-[shimmer_1s_ease-in-out_infinite]" />
        </div>
      )}
      {filterSlot && (
        <div className="px-4 py-3 border-b border-gray-200">{filterSlot}</div>
      )}

      <ErrorDisplay className="px-4 pt-3">{error}</ErrorDisplay>

      {isInitialLoading ? (
        <TableSkeleton columns={columns} />
      ) : (
        <div
          className={
            isRefetching
              ? "opacity-60 pointer-events-none transition-opacity duration-200"
              : "transition-opacity duration-200"
          }
        >
          {selectedIds.size > 0 && bulkActions && (
            <div className="px-4 py-2 bg-blue-50 border-b border-blue-200 flex items-center gap-3">
              <span className="text-sm text-gray-600">
                {selectedIds.size} נבחרו
              </span>
              {bulkActions.map((action, i) => (
                <button
                  key={i}
                  type="button"
                  onClick={() => action.onClick(Array.from(selectedIds))}
                  className={`inline-flex items-center gap-1 px-3 py-1 text-sm rounded cursor-pointer ${
                    action.variant === "danger"
                      ? "bg-red-100 text-red-700 hover:bg-red-200"
                      : "bg-blue-100 text-blue-700 hover:bg-blue-200"
                  }`}
                >
                  {action.icon}
                  {action.label}
                </button>
              ))}
            </div>
          )}

          <div className="overflow-x-auto">
            <table className="w-full text-right">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50">
                  <th className="px-3 py-2 w-10">
                    <input
                      type="checkbox"
                      className="h-4 w-4 accent-blue-600"
                      checked={
                        rows.length > 0 && selectedIds.size === rows.length
                      }
                      onChange={toggleSelectAll}
                    />
                  </th>
                  {columns.map((col) => (
                    <th
                      key={col.key}
                      className={`px-3 py-2 text-xs font-semibold text-gray-600 uppercase ${
                        col.sortable ? "cursor-pointer select-none" : ""
                      }`}
                      onClick={
                        col.sortable ? () => handleSort(col.key) : undefined
                      }
                    >
                      <div className="flex items-center gap-1">
                        {col.label}
                        {col.sortable && (
                          <span className="text-gray-400">
                            {sort?.key === col.key ? (
                              sort.dir === "asc" ? (
                                <ArrowUp size={14} />
                              ) : (
                                <ArrowDown size={14} />
                              )
                            ) : (
                              <ChevronsUpDown size={14} />
                            )}
                          </span>
                        )}
                      </div>
                    </th>
                  ))}
                  <th className="px-3 py-2 w-20">
                    {deleteFn && (
                      <button
                        type="button"
                        disabled={
                          selectedIds.size === 0 ||
                          deleteSelectedMutation.isPending
                        }
                        onClick={handleDeleteSelected}
                        className="p-1 text-gray-400 hover:text-red-600 disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer"
                        title={`מחק נבחרים (${selectedIds.size})`}
                      >
                        <Trash2 size={14} />
                      </button>
                    )}
                  </th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => {
                  const id = getId(row);
                  return (
                    <EditableRow
                      key={id}
                      row={row}
                      columns={columns}
                      isEditing={editingId === id}
                      isPending={updateMutation.isPending}
                      schema={updateSchema}
                      selected={selectedIds.has(id)}
                      onToggleSelect={() => toggleSelect(id)}
                      onEdit={() => {
                        setEditingId(id);
                        clearError();
                      }}
                      onCancel={() => {
                        setEditingId(null);
                        clearError();
                      }}
                      onSave={(data) =>
                        updateMutation.mutate({
                          id,
                          data: data as TUpdate,
                        })
                      }
                      onDelete={
                        deleteFn
                          ? () => {
                              if (confirm("האם למחוק?")) {
                                deleteMutation.mutate(id);
                              }
                            }
                          : undefined
                      }
                    />
                  );
                })}
                {!hideCreate && createSchema && (
                  <CreateRow
                    columns={columns}
                    schema={createSchema}
                    isPending={createMutation.isPending}
                    onSubmit={(data) => createMutation.mutate(data as TCreate)}
                  />
                )}
              </tbody>
            </table>
          </div>

          {totalPages > 1 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200">
              <span className="text-sm text-gray-600">
                עמוד {page} מתוך {totalPages} ({total} תוצאות)
              </span>
              <div className="flex items-center gap-1">
                <button
                  type="button"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                  className="p-1.5 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
                >
                  <ChevronRight size={16} />
                </button>
                {buildPageNumbers(page, totalPages).map((p, i) =>
                  p === "..." ? (
                    <span key={`ellipsis-${i}`} className="px-1 text-gray-400">
                      ...
                    </span>
                  ) : (
                    <button
                      key={p}
                      type="button"
                      onClick={() => setPage(p as number)}
                      className={`px-2.5 py-1 rounded text-sm cursor-pointer ${
                        page === p
                          ? "bg-blue-600 text-white"
                          : "border border-gray-300 text-gray-600 hover:bg-gray-100"
                      }`}
                    >
                      {p}
                    </button>
                  ),
                )}
                <button
                  type="button"
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                  className="p-1.5 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
                >
                  <ChevronLeft size={16} />
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function buildPageNumbers(current: number, total: number): (number | "...")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

  const pages: (number | "...")[] = [1];
  if (current > 3) pages.push("...");

  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) pages.push(i);

  if (current < total - 2) pages.push("...");
  pages.push(total);
  return pages;
}
