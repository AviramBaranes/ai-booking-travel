"use client";

import { useState } from "react";
import { z } from "zod";
import { Power, PowerOff } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { booking } from "@/shared/client";
import { CrudTable } from "@/app/admin/components/crud-table/CrudTable";
import { ColumnDef, SortState } from "@/app/admin/components/crud-table/types";
import {
  listLocations,
  insertLocation,
  deleteLocation,
  toggleLocation,
  bulkToggleLocations,
} from "@/shared/api/locations-api";

const columns: ColumnDef<booking.LocationRow>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "name", label: "שם", type: "text", editable: false },
  { key: "country_code", label: "קוד מדינה", type: "text", editable: false },
  { key: "country", label: "מדינה", type: "text", editable: false },
  { key: "city", label: "עיר", type: "text", editable: false },
  { key: "iata", label: "IATA", type: "text", editable: false },
  { key: "enabled", label: "פעיל", type: "checkbox" },
  {
    key: "broker_location_id",
    label: "קוד ספק",
    type: "text",
    editable: false,
  },
];

const createSchema = z.object({
  broker: z.string().min(1, "שדה חובה"),
  id: z.string().min(1, "שדה חובה"),
  name: z.string().min(1, "שדה חובה"),
  country: z.string().min(1, "שדה חובה"),
  country_code: z.string().min(1, "שדה חובה"),
  city: z.string(),
  iata: z.string(),
});

const updateSchema = z.object({
  enabled: z.boolean(),
});

interface Filters {
  countryCode: string;
  broker: string;
  name: string;
  iata: string;
  enabled: string;
}

function buildListParams(
  _sort: SortState | null,
  page: number,
  filters: Filters,
): booking.ListLocationsRequest {
  return {
    CountryCode: filters.countryCode,
    Broker: filters.broker,
    Name: filters.name,
    Iata: filters.iata,
    Enabled: filters.enabled,
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
        <label className="block text-xs text-gray-500 mb-1">שם</label>
        <input
          type="text"
          className={inputClass}
          value={filters.name}
          onChange={(e) => onChange({ ...filters, name: e.target.value })}
          placeholder="סינון לפי שם"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">קוד מדינה</label>
        <input
          type="text"
          className={inputClass}
          value={filters.countryCode}
          onChange={(e) =>
            onChange({ ...filters, countryCode: e.target.value })
          }
          placeholder="סינון לפי קוד מדינה"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">IATA</label>
        <input
          type="text"
          className={inputClass}
          value={filters.iata}
          onChange={(e) => onChange({ ...filters, iata: e.target.value })}
          placeholder="סינון לפי IATA"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">ספק</label>
        <input
          type="text"
          className={inputClass}
          value={filters.broker}
          onChange={(e) => onChange({ ...filters, broker: e.target.value })}
          placeholder="סינון לפי ספק"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">סטטוס</label>
        <select
          className={inputClass}
          value={filters.enabled}
          onChange={(e) => onChange({ ...filters, enabled: e.target.value })}
        >
          <option value="">הכל</option>
          <option value="true">פעיל</option>
          <option value="false">לא פעיל</option>
        </select>
      </div>
    </div>
  );
}

export default function LocationsTable() {
  const queryClient = useQueryClient();
  const [filters, setFilters] = useState<Filters>({
    countryCode: "",
    broker: "",
    name: "",
    iata: "",
    enabled: "",
  });

  async function handleBulkToggle(ids: number[], enabled: boolean) {
    await bulkToggleLocations(ids, enabled);
    queryClient.invalidateQueries({ queryKey: ["locations"] });
  }

  return (
    <CrudTable<
      booking.LocationRow,
      booking.InsertLocationParams,
      { enabled: boolean }
    >
      columns={columns}
      queryKey="locations"
      queryKeyDeps={[
        filters.countryCode,
        filters.broker,
        filters.name,
        filters.iata,
        filters.enabled,
      ]}
      getId={(r) => r.id}
      listFn={(sort, page) =>
        listLocations(buildListParams(sort, page, filters))
      }
      extractList={(r) =>
        (r as booking.ListLocationsResponse | undefined)?.locations ?? []
      }
      extractTotal={(r) =>
        (r as booking.ListLocationsResponse | undefined)?.total ?? 0
      }
      createFn={insertLocation}
      updateFn={(id, data) => toggleLocation(id, data.enabled)}
      deleteFn={deleteLocation}
      createSchema={createSchema}
      updateSchema={updateSchema}
      pageSize={15}
      bulkActions={[
        {
          label: "הפעל נבחרים",
          icon: <Power size={14} />,
          onClick: (ids) => handleBulkToggle(ids, true),
        },
        {
          label: "השבת נבחרים",
          icon: <PowerOff size={14} />,
          onClick: (ids) => handleBulkToggle(ids, false),
          variant: "danger",
        },
      ]}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
