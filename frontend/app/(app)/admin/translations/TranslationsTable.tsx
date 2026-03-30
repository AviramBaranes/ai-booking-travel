"use client";

import { useState } from "react";
import { z } from "zod";
import { BadgeCheck } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { booking } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/components/crud-table/CrudTable";
import {
  ColumnDef,
  SortState,
} from "@/app/(app)/admin/components/crud-table/types";
import {
  listBrokerTranslations,
  deleteBrokerTranslation,
  updateBrokerTranslation,
  verifyBrokerTranslation,
} from "@/shared/api/translations-api";

const updateSchema = z.object({
  target_text: z.string().min(1, "שדה חובה"),
});

interface Filters {
  search: string;
  status: string;
}

function buildListParams(
  sort: SortState | null,
  page: number,
  filters: Filters,
): booking.ListBrokerTranslationsRequest {
  return {
    Page: page,
    Search: filters.search,
    Status: filters.status,
    SortDir: sort?.dir ?? "asc",
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
        <label className="block text-xs text-gray-500 mb-1">חיפוש</label>
        <input
          type="text"
          className={inputClass}
          value={filters.search}
          onChange={(e) => onChange({ ...filters, search: e.target.value })}
          placeholder="חיפוש לפי טקסט"
        />
      </div>
      <div>
        <label className="block text-xs text-gray-500 mb-1">סטטוס</label>
        <select
          className={inputClass}
          value={filters.status}
          onChange={(e) => onChange({ ...filters, status: e.target.value })}
        >
          <option value="">הכל</option>
          <option value="pending">ממתין</option>
          <option value="translated">מתורגם</option>
          <option value="verified">מאומת</option>
        </select>
      </div>
    </div>
  );
}

export default function TranslationsTable() {
  const queryClient = useQueryClient();
  const tStatus = useTranslations("TranslationStatus");
  const [filters, setFilters] = useState<Filters>({
    search: "",
    status: "",
  });

  const columns: ColumnDef<booking.BrokerTranslationRow>[] = [
    { key: "id", label: "מזהה", type: "number", editable: false },
    { key: "source_text", label: "טקסט מקור", type: "text", editable: false },
    { key: "target_text", label: "תרגום", type: "text" },
    {
      key: "status",
      label: "סטטוס",
      type: "text",
      editable: false,
      format: (v) => tStatus(String(v)),
    },
    {
      key: "confidence_score",
      label: "ציון",
      type: "number",
      editable: false,
      sortable: true,
    },
  ];

  async function handleBulkVerify(ids: number[]) {
    await Promise.all(ids.map((id) => verifyBrokerTranslation(id)));
    queryClient.invalidateQueries({ queryKey: ["translations"] });
  }

  return (
    <CrudTable<
      booking.BrokerTranslationRow,
      never,
      booking.UpdateBrokerTranslationRequest
    >
      columns={columns}
      queryKey="translations"
      queryKeyDeps={[filters.search, filters.status]}
      getId={(r) => r.id}
      listFn={(sort, page) =>
        listBrokerTranslations(buildListParams(sort, page, filters))
      }
      extractList={(r) =>
        (r as booking.ListBrokerTranslationsResponse | undefined)
          ?.translations ?? []
      }
      extractTotal={(r) =>
        (r as booking.ListBrokerTranslationsResponse | undefined)?.total ?? 0
      }
      updateFn={(id, data) => updateBrokerTranslation(id, data)}
      deleteFn={deleteBrokerTranslation}
      updateSchema={updateSchema}
      hideCreate
      pageSize={15}
      bulkActions={[
        {
          label: "אמת נבחרים",
          icon: <BadgeCheck size={14} />,
          onClick: (ids) => handleBulkVerify(ids),
        },
      ]}
      filterSlot={<FilterBar filters={filters} onChange={setFilters} />}
    />
  );
}
