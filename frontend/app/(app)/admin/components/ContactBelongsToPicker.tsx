"use client";

import { useState } from "react";
import { OrgCombobox } from "@/app/(app)/admin/components/OrgCombobox";
import { OfficeCombobox } from "@/app/(app)/admin/components/OfficeCombobox";

type AssociationType = "office" | "org";

export function parseAssociation(value: unknown): {
  type: AssociationType;
  id: number;
} {
  if (typeof value === "string" && value.includes(":")) {
    const [type, idStr] = value.split(":");
    return {
      type: type === "org" ? "org" : "office",
      id: Number(idStr) || 0,
    };
  }
  return { type: "office", id: 0 };
}

export function encodeAssociation(type: AssociationType, id: number): string {
  return `${type}:${id}`;
}

interface ContactBelongsToPickerProps {
  value: unknown;
  onChange: (value: unknown) => void;
  initialType?: AssociationType;
}

export function ContactBelongsToPicker({
  value,
  onChange,
  initialType,
}: ContactBelongsToPickerProps) {
  const parsed = parseAssociation(value);
  const [type, setType] = useState<AssociationType>(
    parsed.id > 0 ? parsed.type : (initialType ?? "office"),
  );

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex gap-3 text-sm">
        <label className="flex items-center gap-1 cursor-pointer">
          <input
            type="radio"
            name="belongsTo"
            value="office"
            checked={type === "office"}
            onChange={() => {
              setType("office");
              onChange(encodeAssociation("office", 0));
            }}
          />
          משרד
        </label>
        <label className="flex items-center gap-1 cursor-pointer">
          <input
            type="radio"
            name="belongsTo"
            value="org"
            checked={type === "org"}
            onChange={() => {
              setType("org");
              onChange(encodeAssociation("org", 0));
            }}
          />
          רשת
        </label>
      </div>
      {type === "office" ? (
        <OfficeCombobox
          value={parsed.type === "office" ? parsed.id : 0}
          onChange={(id) => onChange(encodeAssociation("office", id))}
        />
      ) : (
        <OrgCombobox
          value={parsed.type === "org" ? parsed.id : 0}
          onChange={(id) => onChange(encodeAssociation("org", id))}
        />
      )}
    </div>
  );
}
