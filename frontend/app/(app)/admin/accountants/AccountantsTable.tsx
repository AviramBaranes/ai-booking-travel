"use client";

import { z } from "zod";
import { accounts } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import { ColumnDef } from "@/app/(app)/admin/_components/crud-table/types";
import {
  listAccountants,
  createAccountant,
  updateUser,
} from "@/shared/api/accounts-api";

interface AccountantUpdateData {
  firstName?: string;
  lastName?: string;
  email: string;
  password: string;
}

const columns: ColumnDef<accounts.AccountantResponse>[] = [
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

export default function AccountantsTable() {
  return (
    <CrudTable<
      accounts.AccountantResponse,
      accounts.CreateAccountantRequest,
      AccountantUpdateData
    >
      columns={[
        ...columns,
        {
          key: "password" as keyof accounts.AccountantResponse,
          label: "סיסמה",
          type: "password",
        },
      ]}
      queryKey="accountants"
      getId={(r) => r.id}
      listFn={() => listAccountants()}
      extractList={(r) =>
        (r as accounts.ListAccountantsResponse | undefined)?.accountants ?? []
      }
      createFn={createAccountant}
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
