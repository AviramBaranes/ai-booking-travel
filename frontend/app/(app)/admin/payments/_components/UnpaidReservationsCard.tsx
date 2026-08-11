"use client";

import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { StatusBadge } from "@/app/(app)/admin/reports/_components/reportTableUtils";
import { ReservationDetailDialog } from "@/app/(app)/admin/_components/ReservationDetailDialog";
import { formatDate } from "@/shared/utils/formatDate";
import { formatPriceFloat } from "@/shared/utils/formatPrice";
import { LoadPaymentSummaryDialog } from "./LoadPaymentSummaryDialog";
import { MarkAsPaidDialog } from "./MarkAsPaidDialog";
import type { Broker } from "./BrokerCombobox";
import {
  isPenalty,
  PAYMENT_FOR_RESERVATION,
  penaltyTypeLabel,
  toUnpaidRows,
  type UnpaidPenalty,
  type UnpaidReservation,
} from "./unpaid-reservations";

const CHECKBOX_CLASSES =
  "border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand";

// Fees are tinted so they read as an exception among the rental charges.
const PENALTY_ROW_CLASSES = "bg-destructive/10 hover:bg-destructive/20";

const displayDate = (value: string) => {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : formatDate(date);
};

interface UnpaidReservationsCardProps {
  broker: Broker;
  reservations: UnpaidReservation[];
  penalties: UnpaidPenalty[];
}

export function UnpaidReservationsCard({
  broker,
  reservations,
  penalties,
}: UnpaidReservationsCardProps) {
  // Reservations and fees are numbered independently, so each keeps its own selection.
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [selectedPenaltyIds, setSelectedPenaltyIds] = useState<Set<number>>(
    new Set(),
  );
  const [summaryDialogOpen, setSummaryDialogOpen] = useState(false);
  const [markAsPaidDialogOpen, setMarkAsPaidDialogOpen] = useState(false);
  // A fee has no details of its own, so it opens the reservation it was charged on.
  const [detailReservationId, setDetailReservationId] = useState<number | null>(
    null,
  );

  const selectedReservations = useMemo(
    () => reservations.filter((r) => selected.has(r.id)),
    [selected, reservations],
  );

  const selectedPenalties = useMemo(
    () => penalties.filter((p) => selectedPenaltyIds.has(p.id)),
    [selectedPenaltyIds, penalties],
  );

  const rows = useMemo(
    () => toUnpaidRows(reservations, penalties),
    [reservations, penalties],
  );

  const rowCount = rows.length;
  const selectedCount = selectedReservations.length + selectedPenalties.length;

  // Currencies are listed together, so a single sum would be meaningless — total per currency.
  const selectedTotals = useMemo(() => {
    const totals = new Map<string, number>();
    for (const row of [...selectedReservations, ...selectedPenalties]) {
      totals.set(
        row.currencyCode,
        (totals.get(row.currencyCode) ?? 0) + row.amountOwed,
      );
    }
    return [...totals.entries()];
  }, [selectedReservations, selectedPenalties]);

  const allChecked = rowCount > 0 && selectedCount === rowCount;
  const someChecked = selectedCount > 0 && !allChecked;

  const toggleAll = () => {
    if (allChecked) {
      setSelected(new Set());
      setSelectedPenaltyIds(new Set());
    } else {
      setSelected(new Set(reservations.map((r) => r.id)));
      setSelectedPenaltyIds(new Set(penalties.map((p) => p.id)));
    }
  };

  const toggle = (
    setIds: typeof setSelected,
    id: number,
  ) => {
    setIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleOne = (id: number) => toggle(setSelected, id);
  const togglePenalty = (id: number) => toggle(setSelectedPenaltyIds, id);

  // Approved ids come from the summary file, so keep only the ones actually on screen.
  const selectApproved = (reservationIds: number[], penaltyIds: number[]) => {
    const reservationsOnScreen = new Set(reservations.map((r) => r.id));
    setSelected(
      new Set(reservationIds.filter((id) => reservationsOnScreen.has(id))),
    );

    const penaltiesOnScreen = new Set(penalties.map((p) => p.id));
    setSelectedPenaltyIds(
      new Set(penaltyIds.filter((id) => penaltiesOnScreen.has(id))),
    );
  };

  const clearSelection = () => {
    setSelected(new Set());
    setSelectedPenaltyIds(new Set());
  };

  return (
    <section className="bg-card rounded-2xl shadow-card overflow-hidden">
      <header className="flex flex-wrap items-center justify-between gap-4 p-6 border-b border-border-light/60">
        <div className="flex flex-col gap-0.5">
          <h2 className="type-h6 text-navy">
            הזמנות שטרם שולמו - {broker.name}
          </h2>
          <p className="type-paragraph text-text-secondary">
            {rowCount} שורות
            {selectedCount > 0 && (
              <span className="text-navy font-medium">
                {" "}
                • {selectedCount} נבחרו
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
          {selectedCount > 0 && (
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
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תשלום עבור
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
            {rows.map((row) => {
              const penalty = isPenalty(row);
              const checked = penalty
                ? selectedPenaltyIds.has(row.id)
                : selected.has(row.id);

              return (
                <tr
                  key={`${penalty ? "penalty" : "reservation"}-${row.id}`}
                  className={`border-t border-border-light/40 transition-colors cursor-pointer ${
                    penalty ? PENALTY_ROW_CLASSES : "hover:bg-background/60"
                  }`}
                  onClick={() =>
                    setDetailReservationId(
                      isPenalty(row) ? row.reservationId : row.id,
                    )
                  }
                >
                  {/* Stops the checkbox from also opening the reservation behind it. */}
                  <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
                    <Checkbox
                      checked={checked}
                      onCheckedChange={() =>
                        penalty ? togglePenalty(row.id) : toggleOne(row.id)
                      }
                      className={CHECKBOX_CLASSES}
                      aria-label={`בחר ${penalty ? "חיוב" : "הזמנה"} ${row.id}`}
                    />
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-navy font-medium">
                    {row.id}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {row.brokerReservationId || "—"}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {isPenalty(row)
                      ? penaltyTypeLabel(row.type)
                      : PAYMENT_FOR_RESERVATION}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {row.driverName || "—"}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                    {displayDate(row.pickupDate)}
                  </td>
                  <td className="px-4 py-3">{row.pickupLocationName || "—"}</td>
                  {/* A fee has no rental days and no status of its own. */}
                  <td className="px-4 py-3 text-start whitespace-nowrap">
                    {isPenalty(row) ? "—" : row.rentalDays}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap font-medium">
                    {formatPriceFloat(row.amountOwed, row.currencyCode)}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {isPenalty(row) ? (
                      "—"
                    ) : (
                      <StatusBadge status={row.reservationStatus} />
                    )}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {isPenalty(row) ? (
                      "—"
                    ) : (
                      <StatusBadge status={row.paymentStatus} />
                    )}
                  </td>
                </tr>
              );
            })}
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
        penaltyIds={selectedPenalties.map((p) => p.id)}
        onPaid={clearSelection}
      />
      <ReservationDetailDialog
        reservationId={detailReservationId}
        onClose={() => setDetailReservationId(null)}
      />
    </section>
  );
}
