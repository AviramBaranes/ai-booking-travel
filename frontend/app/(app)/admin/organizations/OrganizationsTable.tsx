"use client";

import { useState } from "react";
import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/components/crud-table/types";
import {
  listOrganizations,
  createOrganization,
  updateOrganization,
} from "@/shared/api/accounts-api";

const ORGANIC_OPTIONS = [
  { value: "true", label: "אורגני" },
  { value: "false", label: "לא אורגני" },
];

const columns: ColumnDef<accounts.ListOrganizationsRow>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "name", label: "שם", type: "text" },
  {
    key: "isOrganic",
    label: "סוג רשת",
    type: "select",
    options: ORGANIC_OPTIONS,
    format: (v) => (v ? "אורגני" : "לא אורגני"),
  },
  { key: "phone", label: "טלפון", type: "text" },
  { key: "address", label: "כתובת", type: "text" },
  { key: "obligo", label: "אובליגו", type: "number" },
  {
    key: "officeCount",
    label: "משרדים",
    type: "link",
    editable: false,
    href: (row) => `/admin/offices?orgId=${row.id}`,
  },
  {
    key: "contactCount",
    label: "אנשי קשר",
    type: "link",
    editable: false,
    href: (row) => `/admin/contacts?orgId=${row.id}`,
  },
  {
    key: "agentCount",
    label: "סוכנים",
    type: "link",
    editable: false,
    href: (row) => `/admin/agents?orgId=${row.id}`,
  },
];

const booleanFromSelect = z.preprocess(
  (v) => v === "true" || v === true,
  z.boolean(),
);

const obligoField = z.preprocess(
  (v) =>
    v === "" || v === undefined || (typeof v === "number" && isNaN(v))
      ? undefined
      : Number(v),
  z.number().min(0, "ערך מינימלי 0").optional(),
);

const createSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  isOrganic: booleanFromSelect,
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
  obligo: obligoField,
});

const updateSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  isOrganic: booleanFromSelect,
  phone: z.string().optional().default(""),
  address: z.string().optional().default(""),
  obligo: obligoField,
});

interface Filters {
  search: string;
  isOrganic: string;
}

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: Filters,
): accounts.ListOrganizationsRequest {
  return {
    Search: filters.search,
    IsOrganic: filters.isOrganic,
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
      <div>
        <label className="block text-xs text-gray-500 mb-1">סוג רשת</label>
        <select
          className={inputClass}
          value={filters.isOrganic}
          onChange={(e) => onChange({ ...filters, isOrganic: e.target.value })}
        >
          <option value="">הכל</option>
          <option value="true">אורגני</option>
          <option value="false">לא אורגני</option>
        </select>
      </div>
    </div>
  );
}

type OrgUpdateData = {
  name: string;
  isOrganic: boolean;
  phone?: string;
  address?: string;
  obligo?: number;
};

type OrgCreateData = OrgUpdateData;

export default function OrganizationsTable() {
  const [filters, setFilters] = useState<Filters>({
    search: "",
    isOrganic: "",
  });

  return (
    <CrudTable<accounts.ListOrganizationsRow, OrgCreateData, OrgUpdateData>
      columns={columns}
      queryKey="organizations"
      queryKeyDeps={[filters.search, filters.isOrganic]}
      getId={(r) => r.id}
      listFn={(sort, page) =>
        listOrganizations(buildRequest(sort, page, filters))
      }
      extractList={(r) =>
        (r as accounts.ListOrganizationsResponse | undefined)?.organizations ??
        []
      }
      extractTotal={(r) =>
        (r as accounts.ListOrganizationsResponse | undefined)?.total ?? 0
      }
      createFn={(data) =>
        createOrganization({
          name: data.name,
          isOrganic: data.isOrganic,
          phone: data.phone || undefined,
          address: data.address || undefined,
          obligo: data.obligo || undefined,
        })
      }
      updateFn={(id, data) =>
        updateOrganization(id, {
          name: data.name,
          isOrganic: data.isOrganic,
          phone: data.phone,
          address: data.address,
          obligo: data.obligo,
        })
      }
      createSchema={createSchema}
      updateSchema={updateSchema}
      pageSize={15}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
