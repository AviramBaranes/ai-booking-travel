"use client";

import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
} from "@/components/ui/combobox";
import { location } from "@/shared/client";
import { searchLocationsForMerge } from "@/shared/api/locations-api";
import LocationSummary from "./LocationSummary";

interface MergeLocationSearchProps {
  label: string;
  value: location.MergeLocation | null;
  onChange: (l: location.MergeLocation | null) => void;
}

export default function MergeLocationSearch({
  label,
  value,
  onChange,
}: MergeLocationSearchProps) {
  const [search, setSearch] = useState("");
  const [debounced, setDebounced] = useState("");

  useEffect(() => {
    const t = setTimeout(() => setDebounced(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  const { data, isFetching } = useQuery({
    queryKey: ["location-merge-search", debounced],
    queryFn: () => searchLocationsForMerge(debounced),
    enabled: debounced.length >= 2,
  });

  const items = debounced.length >= 2 ? (data?.locations ?? []) : [];

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-gray-700">{label}</label>
      <Combobox
        items={items}
        filter={null}
        value={value}
        onValueChange={(l) => onChange(l as location.MergeLocation | null)}
        onInputValueChange={setSearch}
        itemToStringLabel={(l: location.MergeLocation) => l.name}
        isItemEqualToValue={(x: location.MergeLocation, y: location.MergeLocation) =>
          x.id === y.id
        }
      >
        <ComboboxInput
          placeholder="חיפוש לפי שם, עיר, IATA או מזהה"
          showClear
          className="w-full"
        />
        <ComboboxContent>
          <ComboboxEmpty>
            {debounced.length < 2
              ? "הקלידו לפחות 2 תווים"
              : isFetching
                ? "מחפש..."
                : "לא נמצאו מיקומים"}
          </ComboboxEmpty>
          <ComboboxList>
            {(l: location.MergeLocation) => (
              <ComboboxItem key={l.id} value={l}>
                <div className="flex flex-col">
                  <span>
                    {l.name}{" "}
                    <span className="text-muted-foreground">
                      · {l.country_code} · #{l.id}
                    </span>
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {l.broker_codes.map((c) => c.broker).join(", ")}
                  </span>
                </div>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
      {value && <LocationSummary location={value} />}
    </div>
  );
}
