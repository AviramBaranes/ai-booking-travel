"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";

import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
  ComboboxGroup,
  ComboboxLabel,
  ComboboxCollection,
} from "@/components/ui/combobox";
import {
  listInorganicOffices,
  listOrganicOrganizations,
} from "@/shared/api/accounts-api";

export type BillingEntity =
  | { kind: "org"; id: number; name: string }
  | { kind: "office"; id: number; name: string };

interface BillingEntityComboboxProps {
  value: BillingEntity | null;
  onChange: (value: BillingEntity | null) => void;
}

type EntityOption = {
  key: string;
  label: string;
  entity: BillingEntity;
};

type EntityGroup = {
  value: string;
  items: EntityOption[];
};

export function BillingEntityCombobox({
  value,
  onChange,
}: BillingEntityComboboxProps) {
  const { data: orgsData } = useQuery({
    queryKey: ["billing-entity-orgs"],
    queryFn: () => listOrganicOrganizations(),
  });

  const { data: officesData } = useQuery({
    queryKey: ["billing-entity-offices"],
    queryFn: () => listInorganicOffices(),
  });

  const groups = useMemo<EntityGroup[]>(() => {
    const result: EntityGroup[] = [];
    const orgs = orgsData?.organizations ?? [];
    const offices = officesData?.offices ?? [];
    if (orgs.length > 0) {
      result.push({
        value: "רשתות",
        items: orgs.map((o) => ({
          key: String(o.id),
          label: o.name,
          entity: { kind: "org", id: o.id, name: o.name },
        })),
      });
    }
    if (offices.length > 0) {
      result.push({
        value: "סוכנויות",
        items: offices.map((o) => ({
          key: String(o.id),
          label: o.name,
          entity: { kind: "office", id: o.id, name: o.name },
        })),
      });
    }
    return result;
  }, [orgsData, officesData]);

  const selectedOption = useMemo<EntityOption | null>(() => {
    if (!value) return null;
    for (const group of groups) {
      const found = group.items.find(
        (i) => i.entity.kind === value.kind && i.entity.id === value.id,
      );
      if (found) return found;
    }
    return { key: String(value.id), label: value.name, entity: value };
  }, [groups, value]);

  return (
    <Combobox
      items={groups}
      value={selectedOption}
      onValueChange={(option: EntityOption | null) =>
        onChange(option ? option.entity : null)
      }
      itemToStringLabel={(item: EntityOption) => item.label}
      itemToStringValue={(item: EntityOption) => item.key}
      isItemEqualToValue={(a: EntityOption, b: EntityOption) => a.key === b.key}
    >
      <ComboboxInput
        placeholder="בחר ישות לחיוב..."
        showClear
        className="w-full"
      />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו תוצאות</ComboboxEmpty>
        <ComboboxList>
          {(group: EntityGroup) => (
            <ComboboxGroup key={group.value} items={group.items}>
              <ComboboxLabel>{group.value}</ComboboxLabel>
              <ComboboxCollection>
                {(item: EntityOption) => (
                  <ComboboxItem key={item.key} value={item}>
                    {item.label}
                  </ComboboxItem>
                )}
              </ComboboxCollection>
            </ComboboxGroup>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
