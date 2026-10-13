"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { listOpenReservations } from "@/shared/api/reservations-api";
import type { queries, reservation_pricing } from "@/shared/client";
import { AddCanceledReservationForm } from "./AddCanceledReservationForm";
import { CurrencyGroupCard } from "./CurrencyGroupCard";
import type { BillingEntity } from "./BillingEntityCombobox";

interface BillingResultsProps {
  entity: BillingEntity;
  showActions?: boolean;
}

export function BillingResults({
  entity,
  showActions = true,
}: BillingResultsProps) {
  const { data } = useQuery({
    queryKey: ["open-reservations", entity.kind, entity.id],
    queryFn: () =>
      listOpenReservations({
        OrgID: entity.kind === "org" ? entity.id : undefined,
        OfficeID: entity.kind === "office" ? entity.id : undefined,
      }),
  });

  // Canceled, unpaid reservations brought in by hand. They are not open, so the server does not
  // list them until they are billed and turn into refunds.
  const [added, setAdded] = useState<reservation_pricing.BillingReservation[]>(
    [],
  );

  const { groups, listedIds } = useMemo(() => {
    const fetched = data?.currencyGroups ?? [];
    const listedIds = new Set(
      fetched.flatMap((g) => g.reservations.map((r) => r.id)),
    );

    const groups: queries.CurrencyGroup[] = fetched.map((g) => ({
      ...g,
      reservations: [...g.reservations],
    }));
    for (const r of added) {
      // Once billed it comes back from the server as a refund, which replaces it.
      if (listedIds.has(r.id)) continue;
      listedIds.add(r.id);

      let group = groups.find((g) => g.currencyCode === r.currencyCode);
      if (!group) {
        group = {
          currencyCode: r.currencyCode,
          reservations: [],
          penalties: [],
        };
        groups.push(group);
      }
      group.reservations.push(r);
    }

    return { groups, listedIds };
  }, [data, added]);

  return (
    <div className="flex flex-col gap-6">
      {showActions && (
        <AddCanceledReservationForm
          entity={entity}
          listedIds={listedIds}
          onAdd={(r) => setAdded((prev) => [...prev, r])}
        />
      )}

      {groups.length === 0 ? (
        <div className="bg-card rounded-2xl shadow-card p-12 text-center flex flex-col items-center gap-2">
          <h3 className="type-h6 text-navy">אין הזמנות פתוחות</h3>
          <p className="type-paragraph text-text-secondary">
            לא נמצאו הזמנות פתוחות עבור {entity.name}.
          </p>
        </div>
      ) : (
        groups.map((group) => (
          <CurrencyGroupCard
            key={group.currencyCode}
            entity={entity}
            group={group}
            showActions={showActions}
          />
        ))
      )}
    </div>
  );
}
