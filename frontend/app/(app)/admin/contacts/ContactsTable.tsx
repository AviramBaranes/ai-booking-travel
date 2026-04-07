"use client";

import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/_components/crud-table/types";
import {
  listContacts,
  createContact,
  updateContact,
  deleteContact,
} from "@/shared/api/accounts-api";
import { ContactsFilterBar } from "@/app/(app)/admin/_components/ContactsFilterBar";
import {
  ContactBelongsToPicker,
  parseAssociation,
  encodeAssociation,
} from "@/app/(app)/admin/_components/ContactBelongsToPicker";
import { useUrlFilters } from "@/app/(app)/admin/_hooks/useUrlFilters";

const columns: ColumnDef<accounts.ContactResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "firstName", label: "שם פרטי", type: "text" },
  { key: "lastName", label: "שם משפחה", type: "text" },
  { key: "role", label: "תפקיד", type: "text" },
  { key: "cellphone", label: "טלפון", type: "text" },
  { key: "email", label: "אימייל", type: "text" },
  {
    key: "officeId",
    label: "שייך ל",
    type: "text",
    renderCell: (_value, row) => {
      if (row.officeName) return `משרד: ${row.officeName}`;
      if (row.organizationName) return `רשת: ${row.organizationName}`;
      return "";
    },
    renderEditCell: ({ value, onChange, row }) => {
      const contact = row as accounts.ContactResponse | undefined;
      const initialType =
        contact && contact.organizationId > 0 ? "org" : "office";
      const encoded =
        typeof value === "number" ||
        (typeof value === "string" && !value.includes(":"))
          ? encodeAssociation(
              initialType,
              initialType === "office"
                ? (contact?.officeId ?? 0)
                : (contact?.organizationId ?? 0),
            )
          : value;
      return (
        <ContactBelongsToPicker
          value={encoded}
          onChange={onChange}
          initialType={initialType}
        />
      );
    },
  },
];

const associationField = z
  .string()
  .refine((v) => parseAssociation(v).id > 0, "יש לבחור משרד או רשת");

const createSchema = z.object({
  firstName: z.string().min(1, "שדה חובה"),
  lastName: z.string().min(1, "שדה חובה"),
  role: z.string().min(1, "שדה חובה"),
  cellphone: z.string().min(1, "שדה חובה"),
  email: z.string().email("אימייל לא תקין"),
  officeId: associationField,
});

const updateSchema = z.object({
  firstName: z.string().optional(),
  lastName: z.string().optional(),
  role: z.string().optional(),
  cellphone: z.string().optional(),
  email: z.string().email("אימייל לא תקין").optional().or(z.literal("")),
  officeId: associationField,
});

function formDataToCreatePayload(
  data: Record<string, unknown>,
): accounts.CreateContactRequest {
  const { type, id } = parseAssociation(data.officeId);
  return {
    firstName: data.firstName as string,
    lastName: data.lastName as string,
    role: data.role as string,
    cellphone: data.cellphone as string,
    email: data.email as string,
    officeId: type === "office" ? id : undefined,
    organizationId: type === "org" ? id : undefined,
  };
}

function formDataToUpdatePayload(
  data: Record<string, unknown>,
): accounts.UpdateContactRequest {
  const { type, id } = parseAssociation(data.officeId);
  return {
    firstName: data.firstName as string,
    lastName: data.lastName as string,
    role: data.role as string,
    cellphone: data.cellphone as string,
    email: data.email as string,
    officeId: type === "office" ? id : undefined,
    organizationId: type === "org" ? id : undefined,
  };
}

function buildRequest(
  _sort: SortState | null,
  page: number,
  filters: { search: string; orgId: string; officeId: string },
): accounts.ListContactsRequest {
  return {
    Search: filters.search,
    OrgID: filters.orgId ? Number(filters.orgId) : 0,
    OfficeID: filters.officeId ? Number(filters.officeId) : 0,
    Page: page,
  };
}

export default function ContactsTable() {
  const [filters, setFilters] = useUrlFilters(["search", "orgId", "officeId"]);

  return (
    <CrudTable<
      accounts.ContactResponse,
      accounts.CreateContactRequest,
      accounts.UpdateContactRequest
    >
      columns={columns}
      queryKey="contacts"
      queryKeyDeps={[filters.search, filters.orgId, filters.officeId]}
      getId={(r) => r.id}
      listFn={(sort, page) => listContacts(buildRequest(sort, page, filters))}
      extractList={(r) =>
        (r as accounts.ListContactsResponse | undefined)?.contacts ?? []
      }
      extractTotal={(r) =>
        (r as accounts.ListContactsResponse | undefined)?.total ?? 0
      }
      createFn={(data) => createContact(formDataToCreatePayload(data as never))}
      updateFn={(id, data) =>
        updateContact(id, formDataToUpdatePayload(data as never))
      }
      deleteFn={deleteContact}
      createSchema={createSchema as never}
      updateSchema={updateSchema as never}
      pageSize={15}
      filterSlot={<ContactsFilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
