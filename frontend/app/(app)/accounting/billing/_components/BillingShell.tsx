"use client";

import { Suspense, useState } from "react";
import {
  BillingEntityCombobox,
  type BillingEntity,
} from "./BillingEntityCombobox";
import { BillingResults } from "./BillingResults";

function ReservationsSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      {[0, 1].map((i) => (
        <div
          key={i}
          className="bg-card rounded-2xl shadow-card p-6 animate-pulse h-40"
        />
      ))}
    </div>
  );
}

export function BillingShell() {
  const [entity, setEntity] = useState<BillingEntity | null>(null);

  return (
    <>
      <div className="bg-card rounded-2xl shadow-card p-6 flex flex-col gap-2">
        <label className="type-label text-navy">ישות לחיוב</label>
        <div className="max-w-md">
          <BillingEntityCombobox value={entity} onChange={setEntity} />
        </div>
      </div>

      {entity === null ? (
        <div className="bg-card rounded-2xl shadow-card p-12 text-center flex flex-col items-center gap-2">
          <h3 className="type-h6 text-navy">עדיין לא נבחרה ישות</h3>
          <p className="type-paragraph text-text-secondary">
            בחר רשת או סוכנות מהרשימה כדי להציג הזמנות פתוחות.
          </p>
        </div>
      ) : (
        <Suspense key={`${entity.kind}-${entity.id}`} fallback={<ReservationsSkeleton />}>
          <BillingResults entity={entity} />
        </Suspense>
      )}
    </>
  );
}
