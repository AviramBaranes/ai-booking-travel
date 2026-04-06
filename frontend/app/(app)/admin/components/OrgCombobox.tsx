"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
} from "@/components/ui/combobox";
import { listOrganizations } from "@/shared/api/accounts-api";

interface OrgComboboxProps {
  value: number | string;
  onChange: (orgId: number) => void;
  showClear?: boolean;
  placeholder?: string;
}

export function OrgCombobox({
  value,
  onChange,
  showClear = false,
  placeholder = "בחר רשת...",
}: OrgComboboxProps) {
  const [search, setSearch] = useState("");

  const { data } = useQuery({
    queryKey: ["organizations-combobox", search],
    queryFn: () => listOrganizations({ Page: 1, Search: search }),
  });

  const orgs = data?.organizations ?? [];
  const names = orgs.map((o) => o.name);
  const selectedName = orgs.find((o) => o.id === Number(value))?.name ?? null;

  return (
    <Combobox
      items={names}
      value={selectedName}
      onValueChange={(name) => {
        const org = orgs.find((o) => o.name === name);
        onChange(org ? org.id : 0);
      }}
      onInputValueChange={setSearch}
    >
      <ComboboxInput
        placeholder={placeholder}
        showClear={showClear}
        className="w-full"
      />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו רשתות</ComboboxEmpty>
        <ComboboxList>
          {(name) => (
            <ComboboxItem key={name} value={name}>
              {name}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
