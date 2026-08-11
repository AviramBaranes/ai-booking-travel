"use client";

import { useState } from "react";
import { z } from "zod";
import { markup_rate } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/_components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/_components/crud-table/types";
import {
  listMarkupRates,
  createMarkupRate,
  updateMarkupRate,
  deleteMarkupRate,
} from "@/shared/api/markup-rates-api";

const columns: ColumnDef<markup_rate.MarkupRateResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "countryCode", label: "קוד מדינה", type: "text", sortable: true },
  {
    key: "broker",
    label: "ספק",
    type: "select",
    sortable: true,
    options: [
      { label: "Flex", value: "flex" },
      { label: "Hertz", value: "hertz" },
    ],
  },
  { key: "markUpGross", label: "מרווח ברוטו", type: "number", sortable: true },
  { key: "markUpNet", label: "מרווח נטו", type: "number", sortable: true },
];

const rateSchema = z.object({
  countryCode: z.string().min(1, "שדה חובה"),
  broker: z.string().min(1, "שדה חובה"),
  markUpGross: z.number({ error: "מספר נדרש" }),
  markUpNet: z.number({ error: "מספר נדרש" }),
});

// Map camelCase frontend keys to snake_case backend sort fields
const sortKeyMap: Record<string, string> = {
  country: "country",
  broker: "broker",
  markUpGross: "mark_up_gross",
  markUpNet: "mark_up_net",
};

interface Filters {
  country: string;
  broker: string;
}

function buildListParams(
  sort: SortState | null,
  page: number,
  filters: Filters,
): markup_rate.ListMarkupRatesParams {
  return {
    Country: filters.country,
    broker: filters.broker,
    SortBy: sort ? (sortKeyMap[sort.key] ?? "") : "",
    SortDir: sort?.dir ?? "",
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
        <label className="block text-xs text-gray-500 mb-1">מדינה</label>
        <input
          type="text"
          className={inputClass}
          value={filters.country}
          onChange={(e) => onChange({ ...filters, country: e.target.value })}
          placeholder="סינון לפי מדינה"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">מותג</label>
        <input
          type="text"
          className={inputClass}
          value={filters.broker}
          onChange={(e) => onChange({ ...filters, broker: e.target.value })}
          placeholder="סינון לפי ספק"
        />
      </div>
    </div>
  );
}

export default function PricingTable() {
  const [filters, setFilters] = useState<Filters>({
    country: "",
    broker: "",
  });

  return (
    <CrudTable<
      markup_rate.MarkupRateResponse,
      markup_rate.CreateMarkupRateParams,
      markup_rate.UpdateMarkupRateParams
    >
      columns={columns}
      queryKey="-markup-rates"
      queryKeyDeps={[filters.country, filters.broker]}
      getId={(r) => r.id}
      listFn={(sort, page) =>
        listMarkupRates(buildListParams(sort, page, filters))
      }
      extractList={(r) =>
        (r as markup_rate.ListMarkupRatesResponse | undefined)?.rates ?? []
      }
      extractTotal={(r) =>
        (r as markup_rate.ListMarkupRatesResponse | undefined)?.total ?? 0
      }
      createFn={createMarkupRate}
      updateFn={updateMarkupRate}
      deleteFn={deleteMarkupRate}
      createSchema={rateSchema}
      updateSchema={rateSchema}
      pageSize={15}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
