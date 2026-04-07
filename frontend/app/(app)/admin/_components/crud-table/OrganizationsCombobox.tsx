"use client";

import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox";
import { listOrganizations } from "@/shared/api/accounts-api";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

interface OrganizationsComboboxProps {
  showClear?: boolean;
  onSelect: (org: { id: number; name: string }) => void;
}
export function OrganizationsCombobox({
  showClear,
  onSelect,
}: OrganizationsComboboxProps) {
  const [search, setSearch] = useState("");
  const { data } = useQuery({
    queryKey: ["organization", search],
    queryFn: () => listOrganizations({ Search: search, Page: 1 }),
  });

  return (
    <Combobox items={data?.organizations ?? []}>
      <ComboboxInput
        showClear={showClear}
        placeholder="Select a framework"
        onChange={(e) => setSearch(e.target.value)}
      />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו רשתות</ComboboxEmpty>
        <ComboboxList>
          {(item) => (
            <ComboboxItem
              key={item.id}
              value={item.name}
              onSelect={() => onSelect(item)}
            >
              {item.name}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
