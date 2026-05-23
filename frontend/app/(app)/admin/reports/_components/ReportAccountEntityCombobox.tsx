"use client";

import { useMemo, useState } from "react";
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
  listAgents,
  listOffices,
  listOrganizations,
} from "@/shared/api/accounts-api";

export type ReportAccountEntity =
  | { kind: "organization"; id: number; name: string }
  | { kind: "office"; id: number; name: string }
  | { kind: "agent"; id: number; name: string };

interface ReportAccountEntityComboboxProps {
  value: ReportAccountEntity | null;
  onChange: (value: ReportAccountEntity | null) => void;
}

type EntityOption = {
  key: string;
  label: string;
  entity: ReportAccountEntity;
};

type EntityGroup = {
  value: string;
  items: EntityOption[];
};

function fullName(firstName: string, lastName: string) {
  return `${firstName} ${lastName}`.trim();
}

export function ReportAccountEntityCombobox({
  value,
  onChange,
}: ReportAccountEntityComboboxProps) {
  const [search, setSearch] = useState("");

  const { data: orgsData } = useQuery({
    queryKey: ["report-entity-orgs", search],
    queryFn: () => listOrganizations({ Page: 1, Search: search }),
  });

  const { data: officesData } = useQuery({
    queryKey: ["report-entity-offices", search],
    queryFn: () => listOffices({ Page: 1, Search: search, OrgID: 0 }),
  });

  const { data: agentsData } = useQuery({
    queryKey: ["report-entity-agents", search],
    queryFn: () => listAgents({ Page: 1, Search: search, OrgID: 0, OfficeID: 0 }),
  });

  const groups = useMemo<EntityGroup[]>(() => {
    const result: EntityGroup[] = [];
    const orgs = orgsData?.organizations ?? [];
    const offices = officesData?.offices ?? [];
    const agents = agentsData?.agents ?? [];

    if (orgs.length > 0) {
      result.push({
        value: "רשתות",
        items: orgs.map((org) => ({
          key: `organization:${org.id}`,
          label: org.name,
          entity: { kind: "organization", id: org.id, name: org.name },
        })),
      });
    }

    if (offices.length > 0) {
      result.push({
        value: "משרדים",
        items: offices.map((office) => ({
          key: `office:${office.id}`,
          label: office.name,
          entity: { kind: "office", id: office.id, name: office.name },
        })),
      });
    }

    if (agents.length > 0) {
      result.push({
        value: "סוכנים",
        items: agents.map((agent) => ({
          key: `agent:${agent.id}`,
          label: fullName(agent.firstName, agent.lastName) || agent.email,
          entity: {
            kind: "agent",
            id: agent.id,
            name: fullName(agent.firstName, agent.lastName) || agent.email,
          },
        })),
      });
    }

    return result;
  }, [agentsData, officesData, orgsData]);

  const selectedOption = useMemo<EntityOption | null>(() => {
    if (!value) return null;
    for (const group of groups) {
      const found = group.items.find(
        (item) =>
          item.entity.kind === value.kind && item.entity.id === value.id,
      );
      if (found) return found;
    }
    return {
      key: `${value.kind}:${value.id}`,
      label: value.name,
      entity: value,
    };
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
        placeholder="כל הישויות"
        showClear
        className="w-full"
        onChange={(event) => setSearch(event.target.value)}
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
