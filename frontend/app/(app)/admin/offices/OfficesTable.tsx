"use client";

import { z } from "zod";
import { office } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/_components/crud-table/types";
import {
  listOffices,
  createOffice,
  updateOffice,
} from "@/shared/api/accounts-api";
import { OrgCombobox } from "@/app/(app)/admin/_components/OrgCombobox";
import { OfficesFilterBar } from "@/app/(app)/admin/_components/OfficesFilterBar";
import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";

const columns: ColumnDef<office.OfficeResponse>[] = [
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
  {
    key: "icountClientId",
    label: "מזהה לקוח iCount",
    type: "number",
  },
  { key: "phone", label: "טלפון", type: "text" },
  { key: "address", label: "כתובת", type: "text" },
  {
    key: "contactCount",
    label: "אנשי קשר",
    type: "link",
    editable: false,
    href: (row) => `/admin/contacts?officeId=${row.id}`,
  },
  {
    key: "agentCount",
    label: "סוכנים",
    type: "link",
    editable: false,
    href: (row) => `/admin/agents?officeId=${row.id}`,
  },
];

const icountField = z.preprocess(
  (v) =>
    v === "" || v === undefined || (typeof v === "number" && isNaN(v))
      ? undefined
      : Number(v),
  z.number().optional(),
);

const createSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  organizationId: z.preprocess(
    (v) => (typeof v === "string" ? Number(v) : v),
    z.number().min(1, "שדה חובה"),
  ),
  icountClientId: icountField,
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
});

const updateSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  organizationId: z.preprocess(
    (v) => (typeof v === "string" ? Number(v) : v),
    z.number().min(1, "שדה חובה"),
  ),
  icountClientId: icountField,
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
});

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: { search: string; orgId: string },
): office.ListOfficesParams {
  return {
    Search: filters.search,
    OrgID: filters.orgId ? Number(filters.orgId) : 0,
    Page: page,
  };
}

export default function OfficesTable() {
  const [filters, setFilters] = useUrlFilters(["search", "orgId"]);

  return (
    <CrudTable<
      office.OfficeResponse,
      office.CreateOfficeParams,
      office.UpdateOfficeParams
    >
      columns={columns}
      queryKey="offices"
      queryKeyDeps={[filters.search, filters.orgId]}
      getId={(r) => r.id}
      listFn={(sort, page) => listOffices(buildRequest(sort, page, filters))}
      extractList={(r) =>
        (r as office.ListOfficesResponse | undefined)?.offices ?? []
      }
      extractTotal={(r) =>
        (r as office.ListOfficesResponse | undefined)?.total ?? 0
      }
      createFn={createOffice}
      updateFn={updateOffice}
      createSchema={createSchema}
      updateSchema={updateSchema}
      pageSize={15}
      filterSlot={<OfficesFilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
