"use client";

import { useQuery } from "@tanstack/react-query";

import { listUnpaidSupplierReservations } from "@/shared/api/reservations-api";
import type { Broker } from "./BrokerCombobox";
import { UnpaidReservationsCard } from "./UnpaidReservationsCard";

interface SupplierPaymentsResultsProps {
  broker: Broker;
}

export function SupplierPaymentsResults({
  broker,
}: SupplierPaymentsResultsProps) {
  const { data } = useQuery({
    queryKey: ["unpaid-supplier-reservations", broker.code],
    queryFn: () => listUnpaidSupplierReservations({ Broker: broker.code }),
  });

  const reservations = data?.reservations ?? [];
  const penalties = data?.penalties ?? [];

  if (reservations.length === 0 && penalties.length === 0) {
    return (
      <div className="bg-card rounded-2xl shadow-card p-12 text-center flex flex-col items-center gap-2">
        <h3 className="type-h6 text-navy">אין הזמנות לתשלום</h3>
        <p className="type-paragraph text-text-secondary">
          כל ההזמנות של {broker.name} כבר שולמו.
        </p>
      </div>
    );
  }

  return (
    <UnpaidReservationsCard
      broker={broker}
      reservations={reservations}
      penalties={penalties}
    />
  );
}
