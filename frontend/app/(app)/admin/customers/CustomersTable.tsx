"use client";

import { customer } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/_components/crud-table/types";
import { listCustomers } from "@/shared/api/accounts-api";
import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";
import LoginAsUserButton from "@/app/(app)/admin/_components/LoginAsUserButton";

const formatDate = (value: unknown) => {
  const date = new Date(value as string);
  return date.toLocaleString("he-IL", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const columns: ColumnDef<customer.CustomerResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "firstName", label: "שם פרטי", type: "text", editable: false },
  { key: "lastName", label: "שם משפחה", type: "text", editable: false },
  { key: "email", label: "אימייל", type: "text", editable: false },
  { key: "phoneNumber", label: "טלפון", type: "text", editable: false },
  {
    key: "lastLogin",
    label: "כניסה אחרונה",
    type: "text",
    editable: false,
    format: formatDate,
  },
  {
    key: "createdAt",
    label: "נוצר בתאריך",
    type: "text",
    editable: false,
    format: formatDate,
  },
  {
    key: "updatedAt",
    label: "עודכן בתאריך",
    type: "text",
    editable: false,
    format: formatDate,
  },
  {
    key: "actions" as keyof customer.CustomerResponse,
    label: "פעולות",
    type: "text",
    editable: false,
    renderCell: (_value, row) => (
      <LoginAsUserButton userId={row.id} label="התחבר כלקוח" />
    ),
  },
];

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: { search: string },
): customer.ListCustomersParams {
  return {
    Search: filters.search,
    Page: page,
  };
}

function FilterBar({
  search,
  onChange,
}: {
  search: string;
  onChange: (search: string) => void;
}) {
  return (
    <div className="flex items-end gap-3 flex-wrap">
      <div>
        <label className="block text-xs text-gray-500 mb-1">חיפוש</label>
        <input
          type="text"
          className="border border-gray-300 rounded px-2 py-1.5 text-sm w-full"
          value={search}
          onChange={(e) => onChange(e.target.value)}
          placeholder="חיפוש לפי שם, אימייל או טלפון"
        />
      </div>
    </div>
  );
}

export default function CustomersTable() {
  const [filters, setFilters] = useUrlFilters(["search"]);

  return (
    <CrudTable<customer.CustomerResponse, never, never>
      columns={columns}
      queryKey="customers"
      queryKeyDeps={[filters.search]}
      getId={(row) => row.id}
      listFn={(sort, page) => listCustomers(buildRequest(sort, page, filters))}
      extractList={(response) =>
        (response as customer.ListCustomersResponse | undefined)?.customers ?? []
      }
      extractTotal={(response) =>
        (response as customer.ListCustomersResponse | undefined)?.total ?? 0
      }
      pageSize={15}
      readOnly
      filterSlot={
        <FilterBar
          search={filters.search}
          onChange={(search) => setFilters({ search })}
        />
      }
    />
  );
}