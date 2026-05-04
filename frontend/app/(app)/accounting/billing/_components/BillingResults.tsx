"use client";

import { useSuspenseQuery } from "@tanstack/react-query";
import { listOpenReservations } from "@/shared/api/reservations-api";
import { CurrencyGroupCard } from "./CurrencyGroupCard";
import type { BillingEntity } from "./BillingEntityCombobox";

interface BillingResultsProps {
  entity: BillingEntity;
}

export function BillingResults({ entity }: BillingResultsProps) {
  const { data } = useSuspenseQuery({
    queryKey: ["open-reservations", entity.kind, entity.id],
    queryFn: () =>
      listOpenReservations({
        OrgID: entity.kind === "org" ? entity.id : undefined,
        OfficeID: entity.kind === "office" ? entity.id : undefined,
      }),
  });

  const groups = data.currencyGroups ?? [];

  if (groups.length === 0) {
    return (
      <div className="bg-card rounded-2xl shadow-card p-12 text-center flex flex-col items-center gap-2">
        <h3 className="type-h6 text-navy">אין הזמנות פתוחות</h3>
        <p className="type-paragraph text-text-secondary">
          לא נמצאו הזמנות פתוחות עבור {entity.name}.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6">
      {groups.map((group) => (
        <CurrencyGroupCard
          key={group.currencyCode}
          entity={entity}
          group={group}
        />
      ))}
    </div>
  );
}
