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
import { listOffices } from "@/shared/api/accounts-api";

interface OfficeComboboxProps {
  value: number | string;
  onChange: (officeId: number) => void;
  showClear?: boolean;
  placeholder?: string;
}

export function OfficeCombobox({
  value,
  onChange,
  showClear = false,
  placeholder = "בחר משרד...",
}: OfficeComboboxProps) {
  const [search, setSearch] = useState("");

  const { data } = useQuery({
    queryKey: ["offices-combobox", search],
    queryFn: () => listOffices({ Page: 1, Search: search, OrgID: 0 }),
  });

  const offices = data?.offices ?? [];
  const names = offices.map((o) => o.name);
  const selectedName =
    offices.find((o) => o.id === Number(value))?.name ?? null;

  return (
    <Combobox
      items={names}
      value={selectedName}
      onValueChange={(name) => {
        const office = offices.find((o) => o.name === name);
        onChange(office ? office.id : 0);
      }}
    >
      <ComboboxInput
        placeholder={placeholder}
        showClear={showClear}
        className="w-full"
        onChange={(e) => setSearch(e.target.value)}
      />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו משרדים</ComboboxEmpty>
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
