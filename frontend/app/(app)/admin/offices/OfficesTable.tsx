"use client";

import { useState, useEffect } from "react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";
import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/components/crud-table/types";
import {
  listOffices,
  createOffice,
  updateOffice,
} from "@/shared/api/accounts-api";
import { OrgCombobox } from "@/app/(app)/admin/components/OrgCombobox";

const columns: ColumnDef<accounts.OfficeResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "name", label: "שם", type: "text" },
  {
    key: "organizationId",
    label: "רשת",
    type: "number",
    renderEditCell: ({ value, onChange }) => (
      <OrgCombobox value={value as number} onChange={(v) => onChange(v)} />
    ),
    renderCell: (_value, row) => row.organizationName ?? "",
  },
  { key: "phone", label: "טלפון", type: "text" },
  { key: "address", label: "כתובת", type: "text" },
  {
    key: "contactCount",
    label: "אנשי קשר",
    type: "link",
    editable: false,
    href: (row) => `/admin/contacts?orgId=${row.organizationId}`,
  },
  {
    key: "agentCount",
    label: "סוכנים",
    type: "link",
    editable: false,
    href: (row) => `/admin/agents?orgId=${row.organizationId}`,
  },
];

const createSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  organizationId: z.preprocess(
    (v) => (typeof v === "string" ? Number(v) : v),
    z.number().min(1, "שדה חובה"),
  ),
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
});

const updateSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  organizationId: z.preprocess(
    (v) => (typeof v === "string" ? Number(v) : v),
    z.number().min(1, "שדה חובה"),
  ),
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
});

interface Filters {
  search: string;
  orgId: string;
}

function useUrlFilters(): [Filters, (f: Filters) => void] {
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();

  const [filters, setFiltersState] = useState<Filters>({
    search: searchParams.get("search") ?? "",
    orgId: searchParams.get("orgId") ?? "",
  });

  useEffect(() => {
    const params = new URLSearchParams();
    if (filters.search) params.set("search", filters.search);
    if (filters.orgId) params.set("orgId", filters.orgId);
    const qs = params.toString();
    router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
  }, [filters, router, pathname]);

  return [filters, setFiltersState];
}

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: Filters,
): accounts.ListOfficesRequest {
  return {
    Search: filters.search,
    OrgID: filters.orgId ? Number(filters.orgId) : 0,
    Page: page,
  };
}

function FilterBar({
  filters,
  onChange,
}: {
  filters: Filters;
  onChange: (f: Filters) => void;
}) {
  const inputClass =
    "border border-gray-300 rounded px-2 py-1.5 text-sm w-full";

  return (
    <div className="flex items-end gap-3 flex-wrap">
      <div>
        <label className="block text-xs text-gray-500 mb-1">חיפוש</label>
        <input
          type="text"
          className={inputClass}
          value={filters.search}
          onChange={(e) => onChange({ ...filters, search: e.target.value })}
          placeholder="חיפוש לפי שם"
        />
      </div>
      <div className="min-w-48">
        <label className="block text-xs text-gray-500 mb-1">רשת</label>
        <OrgCombobox
          value={filters.orgId ? Number(filters.orgId) : 0}
          onChange={(v) => onChange({ ...filters, orgId: v ? String(v) : "" })}
          showClear
          placeholder="כל הרשתות"
        />
      </div>
    </div>
  );
}

export default function OfficesTable() {
  const [filters, setFilters] = useUrlFilters();

  return (
    <CrudTable<
      accounts.OfficeResponse,
      accounts.CreateOfficeRequest,
      accounts.UpdateOfficeRequest
    >
      columns={columns}
      queryKey="offices"
      queryKeyDeps={[filters.search, filters.orgId]}
      getId={(r) => r.id}
      listFn={(sort, page) => listOffices(buildRequest(sort, page, filters))}
      extractList={(r) =>
        (r as accounts.ListOfficesResponse | undefined)?.offices ?? []
      }
      extractTotal={(r) =>
        (r as accounts.ListOfficesResponse | undefined)?.total ?? 0
      }
      createFn={createOffice}
      updateFn={updateOffice}
      createSchema={createSchema}
      updateSchema={updateSchema}
      pageSize={15}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
