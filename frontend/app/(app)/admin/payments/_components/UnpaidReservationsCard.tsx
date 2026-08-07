"use client";

import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { StatusBadge } from "@/app/(app)/admin/reports/_components/reportTableUtils";
import { formatDate } from "@/shared/utils/formatDate";
import { formatPriceFloat } from "@/shared/utils/formatPrice";
import { LoadPaymentSummaryDialog } from "./LoadPaymentSummaryDialog";
import { MarkAsPaidDialog } from "./MarkAsPaidDialog";
import type { Broker } from "./BrokerCombobox";
import type { UnpaidReservation } from "./unpaid-reservations";

const CHECKBOX_CLASSES =
  "border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand";

const displayDate = (value: string) => {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : formatDate(date);
};

interface UnpaidReservationsCardProps {
  broker: Broker;
  reservations: UnpaidReservation[];
}

export function UnpaidReservationsCard({
  broker,
  reservations,
}: UnpaidReservationsCardProps) {
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [summaryDialogOpen, setSummaryDialogOpen] = useState(false);
  const [markAsPaidDialogOpen, setMarkAsPaidDialogOpen] = useState(false);

  const selectedReservations = useMemo(
    () => reservations.filter((r) => selected.has(r.id)),
    [selected, reservations],
  );

  // Currencies are listed together, so a single sum would be meaningless — total per currency.
  const selectedTotals = useMemo(() => {
    const totals = new Map<string, number>();
    for (const r of selectedReservations) {
      totals.set(
        r.currencyCode,
        (totals.get(r.currencyCode) ?? 0) + r.amountOwed,
      );
    }
    return [...totals.entries()];
  }, [selectedReservations]);

  const allChecked =
    reservations.length > 0 &&
    selectedReservations.length === reservations.length;
  const someChecked = selectedReservations.length > 0 && !allChecked;

  const toggleAll = () => {
    if (allChecked) {
      setSelected(new Set());
    } else {
      setSelected(new Set(reservations.map((r) => r.id)));
    }
  };

  const toggleOne = (id: number) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  // Approved ids come from the summary file, so keep only the ones actually on screen.
  const selectApproved = (reservationIds: number[]) => {
    const onScreen = new Set(reservations.map((r) => r.id));
    setSelected(new Set(reservationIds.filter((id) => onScreen.has(id))));
  };

  return (
    <section className="bg-card rounded-2xl shadow-card overflow-hidden">
      <header className="flex flex-wrap items-center justify-between gap-4 p-6 border-b border-border-light/60">
        <div className="flex flex-col gap-0.5">
          <h2 className="type-h6 text-navy">
            הזמנות שטרם שולמו - {broker.name}
          </h2>
          <p className="type-paragraph text-text-secondary">
            {reservations.length} הזמנות
            {selectedReservations.length > 0 && (
              <span className="text-navy font-medium">
                {" "}
                • {selectedReservations.length} נבחרו
                {selectedTotals.length > 0 && (
                  <>
                    {" "}
                    (סה״כ:{" "}
                    {selectedTotals
                      .map(([currencyCode, total]) =>
                        formatPriceFloat(total, currencyCode),
                      )
                      .join(" • ")}
                    )
                  </>
                )}
              </span>
            )}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          {selectedReservations.length > 0 && (
            <Button
              type="button"
              variant="brand"
              className="h-10 px-6"
              onClick={() => setMarkAsPaidDialogOpen(true)}
            >
              סמן כשולם
            </Button>
          )}
          <Button
            type="button"
            variant="outline"
            className="h-10 px-6"
            onClick={() => setSummaryDialogOpen(true)}
          >
            טען קובץ חיובים
          </Button>
        </div>
      </header>

      <div className="overflow-x-auto">
        <table className="w-full text-sm border-collapse">
          <thead className="bg-background text-text-secondary type-label">
            <tr>
              <th className="px-4 py-3 text-start w-10">
                <Checkbox
                  checked={
                    allChecked ? true : someChecked ? "indeterminate" : false
                  }
                  onCheckedChange={toggleAll}
                  className={CHECKBOX_CLASSES}
                  aria-label="בחר הכל"
                />
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">מזהה</th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מזהה הזמנה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">שם הנהג</th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תאריך איסוף
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מקום איסוף
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                ימי השכרה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סה״כ לתשלום
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סטטוס הזמנה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סטטוס תשלום
              </th>
            </tr>
          </thead>
          <tbody>
            {reservations.map((r) => (
              <tr
                key={r.id}
                className="border-t border-border-light/40 transition-colors hover:bg-background/60"
              >
                <td className="px-4 py-3">
                  <Checkbox
                    checked={selected.has(r.id)}
                    onCheckedChange={() => toggleOne(r.id)}
                    className={CHECKBOX_CLASSES}
                    aria-label={`בחר הזמנה ${r.id}`}
                  />
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-navy font-medium">
                  {r.id}
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  {r.brokerReservationId || "—"}
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  {r.driverName || "—"}
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                  {displayDate(r.pickupDate)}
                </td>
                <td className="px-4 py-3">{r.pickupLocationName || "—"}</td>
                <td className="px-4 py-3 text-start whitespace-nowrap">
                  {r.rentalDays}
                </td>
                <td className="px-4 py-3 text-start whitespace-nowrap font-medium">
                  {formatPriceFloat(r.amountOwed, r.currencyCode)}
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  <StatusBadge status={r.reservationStatus} />
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  <StatusBadge status={r.paymentStatus} />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <LoadPaymentSummaryDialog
        open={summaryDialogOpen}
        onOpenChange={setSummaryDialogOpen}
        onSelectApproved={selectApproved}
      />
      <MarkAsPaidDialog
        open={markAsPaidDialogOpen}
        onOpenChange={setMarkAsPaidDialogOpen}
        brokerCode={broker.code}
        reservationIds={selectedReservations.map((r) => r.id)}
        onPaid={() => setSelected(new Set())}
      />
    </section>
  );
}
