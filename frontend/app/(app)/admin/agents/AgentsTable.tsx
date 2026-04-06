"use client";

import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/components/crud-table/types";
import { listAgents, createAgent, updateUser } from "@/shared/api/accounts-api";
import { ContactsFilterBar } from "@/app/(app)/admin/components/ContactsFilterBar";
import { OfficeCombobox } from "@/app/(app)/admin/components/OfficeCombobox";
import { useUrlFilters } from "@/app/(app)/admin/hooks/useUrlFilters";
import LoginAsAgentButton from "./LoginAsAgentButton";

const formatDate = (v: unknown) => {
  const d = new Date(v as string);
  return d.toLocaleString("he-IL", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const columns: ColumnDef<accounts.AgentResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "email", label: "אימייל", type: "text" },
  { key: "phoneNumber", label: "טלפון", type: "text" },
  {
    key: "officeId",
    label: "משרד",
    type: "number",
    renderCell: (_v, row) => row.officeName || "",
    renderEditCell: ({ value, onChange }) => (
      <OfficeCombobox
        value={value as number}
        onChange={onChange}
        placeholder="בחר משרד..."
      />
    ),
  },
  {
    key: "organizationName",
    label: "רשת",
    type: "text",
    editable: false,
  },
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
];

const createSchema = z.object({
  email: z.string().email("אימייל לא תקין"),
  password: z.string().min(8, "סיסמה חייבת להכיל לפחות 8 תווים"),
  phoneNumber: z.string().min(1, "שדה חובה"),
  officeId: z.coerce.number().min(1, "יש לבחור משרד"),
});

const updateSchema = z.object({
  email: z.string().email("אימייל לא תקין"),
  phoneNumber: z.string().optional(),
  officeId: z.coerce.number().optional(),
  password: z.string().optional().default(""),
});

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: { search: string; orgId: string; officeId: string },
): accounts.ListAgentsRequest {
  return {
    Search: filters.search,
    OrgID: filters.orgId ? Number(filters.orgId) : 0,
    OfficeID: filters.officeId ? Number(filters.officeId) : 0,
    Page: page,
  };
}

export default function AgentsTable() {
  const [filters, setFilters] = useUrlFilters(["search", "orgId", "officeId"]);

  return (
    <CrudTable<
      accounts.AgentResponse,
      accounts.CreateAgentRequest,
      accounts.UpdateUserRequest
    >
      columns={[
        ...columns,
        {
          key: "password" as keyof accounts.AgentResponse,
          label: "סיסמה",
          type: "password",
        },
        {
          key: "id" as keyof accounts.AgentResponse,
          label: "פעולות",
          type: "text",
          editable: false,
          renderCell: (_v, row) => <LoginAsAgentButton agentId={row.id} />,
        },
      ]}
      queryKey="agents"
      queryKeyDeps={[filters.search, filters.orgId, filters.officeId]}
      getId={(r) => r.id}
      listFn={(sort, page) => listAgents(buildRequest(sort, page, filters))}
      extractList={(r) =>
        (r as accounts.ListAgentsResponse | undefined)?.agents ?? []
      }
      extractTotal={(r) =>
        (r as accounts.ListAgentsResponse | undefined)?.total ?? 0
      }
      createFn={createAgent}
      updateFn={(id, data) => updateUser(id, data)}
      createSchema={createSchema as never}
      updateSchema={updateSchema as never}
      pageSize={15}
      filterSlot={<ContactsFilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
