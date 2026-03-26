"use client";

import { useState } from "react";
import { z } from "zod";
import { booking } from "@/shared/client";
import { CrudTable } from "@/app/admin/components/crud-table/CrudTable";
import { ColumnDef, SortState } from "@/app/admin/components/crud-table/types";
import {
  listHertzMarkupRates,
  createHertzMarkupRate,
  updateHertzMarkupRate,
  deleteHertzMarkupRate,
} from "@/shared/api/hertz-rates-api";

const columns: ColumnDef<booking.HertzMarkupRateResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "country", label: "מדינה", type: "text", sortable: true },
  { key: "brand", label: "מותג", type: "text", sortable: true },
  {
    key: "pickupDateFrom",
    label: "מתאריך איסוף",
    type: "date",
    sortable: true,
  },
  { key: "pickupDateTo", label: "עד תאריך איסוף", type: "date" },
  { key: "carGroup", label: "קבוצת רכב", type: "text", sortable: true },
  { key: "numOfRentalDaysFrom", label: "מימים", type: "number" },
  { key: "numOfRentalDaysTo", label: "עד ימים", type: "number" },
  { key: "markUpGross", label: "מרווח ברוטו", type: "number" },
  { key: "markUpNet", label: "מרווח נטו", type: "number" },
];

const hertzRateSchema = z.object({
  country: z.string().min(1, "שדה חובה"),
  brand: z.string().min(1, "שדה חובה"),
  pickupDateFrom: z.string().min(1, "שדה חובה"),
  pickupDateTo: z.string().min(1, "שדה חובה"),
  carGroup: z.string().min(1, "שדה חובה"),
  numOfRentalDaysFrom: z
    .number({ error: "מספר נדרש" })
    .int("מספר שלם נדרש")
    .min(0, "ערך מינימלי 0"),
  numOfRentalDaysTo: z
    .number({ error: "מספר נדרש" })
    .int("מספר שלם נדרש")
    .min(0, "ערך מינימלי 0"),
  markUpGross: z.number({ error: "מספר נדרש" }),
  markUpNet: z.number({ error: "מספר נדרש" }),
});

// Map camelCase frontend keys to snake_case backend sort fields
const sortKeyMap: Record<string, string> = {
  country: "country",
  brand: "brand",
  pickupDateFrom: "pickup_date_from",
  carGroup: "car_group",
  numOfRentalDaysFrom: "num_of_rental_days_from",
};

interface Filters {
  country: string;
  brand: string;
  carGroup: string;
}

function buildListParams(
  sort: SortState | null,
  page: number,
  filters: Filters,
): booking.ListHertzMarkupRatesRequest {
  return {
    Country: filters.country,
    Brand: filters.brand,
    CarGroup: filters.carGroup,
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
          value={filters.brand}
          onChange={(e) => onChange({ ...filters, brand: e.target.value })}
          placeholder="סינון לפי מותג"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">קבוצת רכב</label>
        <input
          type="text"
          className={inputClass}
          value={filters.carGroup}
          onChange={(e) => onChange({ ...filters, carGroup: e.target.value })}
          placeholder="סינון לפי קבוצה"
        />
      </div>
    </div>
  );
}

export default function PricingTable() {
  const [filters, setFilters] = useState<Filters>({
    country: "",
    brand: "",
    carGroup: "",
  });

  return (
    <CrudTable<
      booking.HertzMarkupRateResponse,
      booking.CreateHertzMarkupRateRequest,
      booking.UpdateHertzMarkupRateRequest
    >
      columns={columns}
      queryKey="hertz-markup-rates"
      queryKeyDeps={[filters.country, filters.brand, filters.carGroup]}
      getId={(r) => r.id}
      listFn={(sort, page) =>
        listHertzMarkupRates(buildListParams(sort, page, filters))
      }
      extractList={(r) =>
        (r as booking.ListHertzMarkupRatesResponse | undefined)?.rates ?? []
      }
      extractTotal={(r) =>
        (r as booking.ListHertzMarkupRatesResponse | undefined)?.total ?? 0
      }
      createFn={createHertzMarkupRate}
      updateFn={updateHertzMarkupRate}
      deleteFn={deleteHertzMarkupRate}
      createSchema={hertzRateSchema}
      updateSchema={hertzRateSchema}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
