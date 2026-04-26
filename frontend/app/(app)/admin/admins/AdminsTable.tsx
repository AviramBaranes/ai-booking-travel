"use client";

import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import { ColumnDef } from "@/app/(app)/admin/_components/crud-table/types";
import { listAdmins, createAdmin, updateUser } from "@/shared/api/accounts-api";

interface AdminUpdateData {
  firstName?: string;
  lastName?: string;
  email: string;
  password: string;
}

const columns: ColumnDef<accounts.AdminResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "firstName", label: "שם פרטי", type: "text" },
  { key: "lastName", label: "שם משפחה", type: "text" },
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

const createSchema = z.object({
  firstName: z.string().min(1, "שדה חובה"),
  lastName: z.string().min(1, "שדה חובה"),
  email: z.string().email("אימייל לא תקין"),
  password: z.string().min(8, "סיסמה חייבת להכיל לפחות 8 תווים"),
});

const updateSchema = z.object({
  firstName: z.string().min(1, "שדה חובה").optional(),
  lastName: z.string().min(1, "שדה חובה").optional(),
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
          firstName: data.firstName,
          lastName: data.lastName,
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
