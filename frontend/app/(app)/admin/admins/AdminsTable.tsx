"use client";

import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import { ColumnDef } from "@/app/(app)/admin/_components/crud-table/types";
import { listAdmins, createAdmin, updateUser } from "@/shared/api/accounts-api";

interface AdminUpdateData {
  email: string;
  password: string;
}

const columns: ColumnDef<accounts.AdminResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "email", label: "אימייל", type: "text" },
  { key: "lastLogin", label: "כניסה אחרונה", type: "text", editable: false },
  {
    key: "createdAt",
    label: "נוצר בתאריך",
    type: "text",
    editable: false,
    format: (v) => {
      const d = new Date(v as string);
      return d.toLocaleString("he-IL", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    },
  },
];

const createColumns: ColumnDef<accounts.CreateAdminRequest>[] = [
  { key: "email", label: "אימייל", type: "text" },
  { key: "password", label: "סיסמה", type: "password" },
];

const createSchema = z.object({
  email: z.string().email("אימייל לא תקין"),
  password: z.string().min(8, "סיסמה חייבת להכיל לפחות 8 תווים"),
});

const updateSchema = z.object({
  email: z.string().email("אימייל לא תקין"),
  password: z.string().optional().default(""),
});

export default function AdminsTable() {
  return (
    <CrudTable<
      accounts.AdminResponse,
      accounts.CreateAdminRequest,
      AdminUpdateData
    >
      columns={[
        ...columns,
        {
          key: "password" as keyof accounts.AdminResponse,
          label: "סיסמה",
          type: "password",
        },
      ]}
      queryKey="admins"
      getId={(r) => r.id}
      listFn={() => listAdmins()}
      extractList={(r) =>
        (r as accounts.ListAdminsResponse | undefined)?.admins ?? []
      }
      createFn={createAdmin}
      updateFn={(id, data) =>
        updateUser(id, {
          email: data.email,
          password: data.password,
          phoneNumber: "",
          officeId: 0,
        })
      }
      createSchema={createSchema}
      updateSchema={updateSchema}
    />
  );
}
