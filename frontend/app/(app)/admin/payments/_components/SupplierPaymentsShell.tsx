"use client";

import { Suspense, useState } from "react";

import { BrokerCombobox, type Broker } from "./BrokerCombobox";
import { SupplierPaymentsResults } from "./SupplierPaymentsResults";

export function SupplierPaymentsSkeleton() {
  return (
    <div className="bg-card rounded-2xl shadow-card p-6 animate-pulse h-40" />
  );
}

export function SupplierPaymentsShell() {
  const [broker, setBroker] = useState<Broker | null>(null);

  return (
    <>
      <div className="bg-card rounded-2xl shadow-card p-6 flex flex-col gap-2">
        <label className="type-label text-navy">ספק</label>
        <div className="max-w-md">
          <BrokerCombobox value={broker} onChange={setBroker} />
        </div>
      </div>

      {broker === null ? (
        <div className="bg-card rounded-2xl shadow-card p-12 text-center flex flex-col items-center gap-2">
          <h3 className="type-h6 text-navy">עדיין לא נבחר ספק</h3>
          <p className="type-paragraph text-text-secondary">
            בחר ספק מהרשימה כדי להציג את ההזמנות שטרם שולמו לו.
          </p>
        </div>
      ) : (
        <Suspense key={broker.code} fallback={<SupplierPaymentsSkeleton />}>
          <SupplierPaymentsResults broker={broker} />
        </Suspense>
      )}
    </>
  );
}
