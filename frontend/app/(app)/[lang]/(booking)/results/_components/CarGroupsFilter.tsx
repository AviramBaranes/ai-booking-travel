import Image from "next/image";
import { CAR_GROUPS_FILTERS } from "./carGroupsFilters";
import { useState } from "react";
import clsx from "clsx";

interface CarGroupFiltersProps {
  title: string;
  setAcrissCodes: (codes: Set<string>) => void;
}

export function CarGroupsFilter({
  title,
  setAcrissCodes,
}: CarGroupFiltersProps) {
  const [selectedGroups, setSelectedGroups] = useState<Set<string>>(new Set());

  function toggleGroup(groupName: string) {
    setSelectedGroups((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(groupName)) {
        newSet.delete(groupName);
      } else {
        newSet.add(groupName);
      }

      setAcrissCodes(
        new Set(
          Array.from(newSet).flatMap((groupName) => {
            const group = CAR_GROUPS_FILTERS.find((g) => g.name === groupName);
            return group ? group.acrissCodes : [];
          }),
        ),
      );

      return newSet;
    });
  }

  return (
    <div className="mt-12">
      <h5 className="type-h5 mb-8 text-navy">{title}</h5>
      <div className="flex items-center justify-between">
        {CAR_GROUPS_FILTERS.map((group) => (
          <div
            onClick={() => toggleGroup(group.name)}
            className={clsx(
              "bg-white text-center rounded-lg shadow-card px-2 py-2 cursor-pointer hover:shadow-card-hover",
              {
                "border-brand border": selectedGroups.has(group.name),
              },
            )}
            key={group.name}
          >
            <p
              className={clsx("type-paragraph text-navy", {
                "font-bold": selectedGroups.has(group.name),
              })}
            >
              {group.name}
            </p>
            <Image src={group.image} alt={group.name} width={125} height={90} />
          </div>
        ))}
      </div>
    </div>
  );
}
