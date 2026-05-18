"use client";

import { z } from "zod";
import { currency } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import { ColumnDef } from "@/app/(app)/admin/_components/crud-table/types";
import {
  listCurrencies,
  createCurrency,
  updateCurrency,
  deleteCurrency,
} from "@/shared/api/currencies-api";

const columns: ColumnDef<currency.CurrencyResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "currencyCode", label: "קוד מטבע", type: "text" },
  { key: "currencyISOName", label: "שם ISO", type: "text" },
  { key: "rate", label: "שער", type: "number" },
];

const currencySchema = z.object({
  currencyCode: z.string().min(1, "שדה חובה"),
  currencyISOName: z.string().min(1, "שדה חובה"),
  rate: z.number({ error: "מספר נדרש" }).min(0, "ערך מינימלי 0"),
});

export default function CurrenciesTable() {
  return (
    <CrudTable<
      currency.CurrencyResponse,
      currency.CreateCurrencyParams,
      currency.UpdateCurrencyParams
    >
      columns={columns}
      queryKey="currencies"
      getId={(r) => r.id}
      listFn={() => listCurrencies()}
      extractList={(r) =>
        (r as currency.ListCurrenciesResponse | undefined)?.currencies ?? []
      }
      createFn={createCurrency}
      updateFn={updateCurrency}
      deleteFn={deleteCurrency}
      createSchema={currencySchema}
      updateSchema={currencySchema}
    />
  );
}
