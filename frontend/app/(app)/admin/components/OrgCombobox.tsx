"use client";

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

type Org = { id: number; name: string };

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
  const { data } = useQuery({
    queryKey: ["organizations-all"],
    queryFn: () => listOrganizations({ Page: 1, Search: "" }),
  });

  const orgs: Org[] =
    (data as { organizations?: Org[] } | undefined)?.organizations ?? [];

  const selected = orgs.find((o) => o.id === Number(value)) ?? null;

  return (
    <Combobox
      items={orgs}
      itemToStringValue={(org) => org.name}
      value={selected}
      onValueChange={(org) => {
        onChange(org ? org.id : 0);
      }}
    >
      <ComboboxInput
        placeholder={placeholder}
        showClear={showClear}
        className="w-full"
      />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו רשתות</ComboboxEmpty>
        <ComboboxList>
          {(org) => (
            <ComboboxItem key={org.id} value={org}>
              {org.name}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
