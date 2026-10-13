"use client";

import { useState, type FormEvent } from "react";
import { useMutation } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { findCanceledReservationForBilling } from "@/shared/api/reservations-api";
import type { reservation_pricing } from "@/shared/client";
import type { BillingEntity } from "./BillingEntityCombobox";

const ALREADY_LISTED_ERROR = "ההזמנה כבר מופיעה ברשימה";

interface AddCanceledReservationFormProps {
  entity: BillingEntity;
  listedIds: Set<number>;
  onAdd: (reservation: reservation_pricing.BillingReservation) => void;
}

// An entity sometimes pays for a reservation that was canceled before it was paid for. Such a
// reservation is not open for billing, so the accountant brings it in by its supplier number to
// record the payment; once billed it turns into a refund owed back to the entity.
export function AddCanceledReservationForm({
  entity,
  listedIds,
  onAdd,
}: AddCanceledReservationFormProps) {
  const [brokerReservationId, setBrokerReservationId] = useState("");
  const [localError, setLocalError] = useState<string | null>(null);

  const { mutate, isPending, error, reset } = useMutation({
    mutationFn: (id: string) =>
      findCanceledReservationForBilling({
        BrokerReservationID: id,
        OrgID: entity.kind === "org" ? entity.id : undefined,
        OfficeID: entity.kind === "office" ? entity.id : undefined,
      }),
    onSuccess: ({ reservation }) => {
      if (listedIds.has(reservation.id)) {
        setLocalError(ALREADY_LISTED_ERROR);
        return;
      }
      onAdd(reservation);
      setBrokerReservationId("");
    },
  });
  const translatedError = useTranslatedError(error);

  const onSubmit = (e: FormEvent) => {
    e.preventDefault();
    const id = brokerReservationId.trim();
    if (!id) return;
    setLocalError(null);
    mutate(id);
  };

  return (
    <form
      onSubmit={onSubmit}
      className="bg-card rounded-2xl shadow-card p-6 flex flex-col gap-2"
    >
      <Label
        htmlFor="canceled-broker-reservation-id"
        className="type-label text-navy"
      >
        הוספת הזמנה מבוטלת לפי מס׳ ספק
      </Label>
      <div className="flex flex-wrap items-center gap-3">
        <Input
          id="canceled-broker-reservation-id"
          value={brokerReservationId}
          onChange={(e) => {
            setBrokerReservationId(e.target.value);
            setLocalError(null);
            reset();
          }}
          placeholder="מס׳ הזמנה אצל הספק"
          className="max-w-xs h-10"
        />
        <Button
          type="submit"
          variant="brand"
          className="h-10 px-6"
          disabled={isPending || !brokerReservationId.trim()}
        >
          הוסף
        </Button>
      </div>
      <ErrorDisplay>{localError ?? translatedError}</ErrorDisplay>
    </form>
  );
}
