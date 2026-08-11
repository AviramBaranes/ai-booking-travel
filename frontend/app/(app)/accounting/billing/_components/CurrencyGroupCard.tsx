"use client";

import { useMemo, useState } from "react";
import { format } from "date-fns/format";
import { he } from "date-fns/locale/he";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { ReservationDetailDialog } from "@/app/(app)/admin/_components/ReservationDetailDialog";
import {
  PAYMENT_FOR_RESERVATION,
  penaltyTypeLabel,
} from "@/app/(app)/admin/payments/_components/unpaid-reservations";
import { queries } from "@/shared/client";
import { BillDialog } from "./BillDialog";
import type { BillingEntity } from "./BillingEntityCombobox";
import { PromptBillDialog } from "./PromptBillDialog";

const CHECKBOX_CLASSES =
  "border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand";

// Fees are tinted so they read as an exception among the rental charges.
const PENALTY_ROW_CLASSES = "bg-destructive/10 hover:bg-destructive/20";

interface CurrencyGroupCardProps {
  entity: BillingEntity;
  group: queries.CurrencyGroup;
  showActions?: boolean;
}

export function CurrencyGroupCard({
  entity,
  group,
  showActions = true,
}: CurrencyGroupCardProps) {
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [promptBillDialogOpen, setPromptBillDialogOpen] = useState(false);
  const [billDialogOpen, setBillDialogOpen] = useState(false);
  // Reservations and fees are numbered independently, so each keeps its own selection.
  const [selectedPenaltyIds, setSelectedPenaltyIds] = useState<Set<number>>(
    new Set(),
  );
  const [detailReservationId, setDetailReservationId] = useState<number | null>(
    null,
  );
  const tStatus = useTranslations("MyAccount.reservation.summary.status");

  const currencyFormatter = useMemo(
    () =>
      new Intl.NumberFormat("he-IL", {
        style: "currency",
        currency: group.currencyCode,
        maximumFractionDigits: 2,
      }),
    [group.currencyCode],
  );

  const selectedReservations = useMemo(
    () => group.reservations.filter((r) => selected.has(r.id)),
    [selected, group.reservations],
  );
  const selectedIds = useMemo(
    () => selectedReservations.map((r) => r.id),
    [selectedReservations],
  );

  const selectedPenalties = useMemo(
    () => group.penalties.filter((p) => selectedPenaltyIds.has(p.id)),
    [selectedPenaltyIds, group.penalties],
  );
  const selectedPenaltyIdList = useMemo(
    () => selectedPenalties.map((p) => p.id),
    [selectedPenalties],
  );

  const selectedTotal = useMemo(
    () =>
      selectedReservations.reduce((sum, r) => sum + r.totalPrice, 0) +
      selectedPenalties.reduce((sum, p) => sum + p.amount, 0),
    [selectedReservations, selectedPenalties],
  );

  const rowCount = group.reservations.length + group.penalties.length;
  const selectedCount = selectedIds.length + selectedPenaltyIdList.length;

  const allChecked = rowCount > 0 && selectedCount === rowCount;
  const someChecked = selectedCount > 0 && !allChecked;

  const toggleAll = () => {
    if (allChecked) {
      setSelected(new Set());
      setSelectedPenaltyIds(new Set());
    } else {
      setSelected(new Set(group.reservations.map((r) => r.id)));
      setSelectedPenaltyIds(new Set(group.penalties.map((p) => p.id)));
    }
  };

  const toggle = (setIds: typeof setSelected, id: number) => {
    setIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleOne = (id: number) => toggle(setSelected, id);
  const togglePenalty = (id: number) => toggle(setSelectedPenaltyIds, id);

  const clearSelection = () => {
    setSelected(new Set());
    setSelectedPenaltyIds(new Set());
  };

  const formatDate = (s: string) => {
    if (!s) return "—";
    const d = new Date(s);
    if (Number.isNaN(d.getTime())) return s;
    return format(d, "d בMMM yyyy", { locale: he });
  };

  return (
    <section className="bg-card rounded-2xl shadow-card overflow-hidden">
      <header className="flex flex-wrap items-center justify-between gap-4 p-6 border-b border-border-light/60">
        <div className="flex flex-col gap-0.5">
          <h2 className="type-h6 text-navy">מטבע {group.currencyCode}</h2>
          <p className="type-paragraph text-text-secondary">
            {rowCount} שורות פתוחות
            {selectedCount > 0 && (
              <span className="text-navy font-medium">
                {" "}
                • {selectedCount} נבחרו
                {selectedTotal > 0 && (
                  <> (סה"כ לתשלום: {currencyFormatter.format(selectedTotal)})</>
                )}
              </span>
            )}
          </p>
        </div>
        {showActions && (
          <Button
            type="button"
            variant="brand"
            className="h-10 px-6"
            disabled={selectedCount === 0}
            onClick={() => setPromptBillDialogOpen(true)}
          >
            חייב נבחרים
          </Button>
        )}
      </header>

      <div className="overflow-x-auto">
        <table className="w-full text-sm border-collapse">
          <thead className="bg-background text-text-secondary type-label">
            <tr>
              {showActions && (
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
              )}
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מס׳ הזמנה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מס׳ ספק
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תשלום עבור
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סטטוס תשלום
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סטטוס הזמנה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מחיר רכש
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מחיר מכירה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מחיר כיסוי
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">רווח</th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                סה״כ לתשלום
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תאריך כרטוס
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תאריך יצירה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                תאריך איסוף
              </th>
            </tr>
          </thead>
          <tbody>
            {group.reservations.map((r) => {
              const checked = selected.has(r.id);
              const isRefundPending = r.paymentStatus === "refund_pending";
              return (
                <tr
                  key={r.id}
                  className={
                    "border-t border-border-light/40 transition-colors cursor-pointer " +
                    (isRefundPending
                      ? "bg-brand/10 hover:bg-brand/15"
                      : "hover:bg-background/60")
                  }
                  onClick={() => setDetailReservationId(r.id)}
                >
                  {showActions && (
                    /* Stops the checkbox from also opening the reservation behind it. */
                    <td
                      className="px-4 py-3"
                      onClick={(e) => e.stopPropagation()}
                    >
                      <Checkbox
                        checked={checked}
                        onCheckedChange={() => toggleOne(r.id)}
                        className={CHECKBOX_CLASSES}
                        aria-label={`בחר הזמנה ${r.id}`}
                      />
                    </td>
                  )}
                  <td className="px-4 py-3 whitespace-nowrap text-navy font-medium">
                    {r.id}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {r.brokerReservationId || "—"}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {PAYMENT_FOR_RESERVATION}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {tStatus(r.paymentStatus)}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {tStatus(r.reservationStatus)}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap">
                    {currencyFormatter.format(r.carPurchasePrice)}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap">
                    {currencyFormatter.format(r.carSellingPrice)}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap">
                    {currencyFormatter.format(r.erpSellingPrice)}
                  </td>
                  <td
                    className={
                      "px-4 py-3 text-start whitespace-nowrap font-medium " +
                      (r.totalProfit >= 0 ? "text-success" : "text-destructive")
                    }
                  >
                    {currencyFormatter.format(r.totalProfit)}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap font-medium">
                    {currencyFormatter.format(r.totalPrice)}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                    {formatDate(r.voucheredAt)}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                    {formatDate(r.createdAt)}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                    {formatDate(r.pickupDate)}
                  </td>
                </tr>
              );
            })}

            {/* Fees close the table: the rows above are ordered by voucher date, which a fee has
                no equivalent of. */}
            {group.penalties.map((p) => (
              <tr
                key={`penalty-${p.id}`}
                className={`border-t border-border-light/40 transition-colors cursor-pointer ${PENALTY_ROW_CLASSES}`}
                onClick={() => setDetailReservationId(p.reservationId)}
              >
                {showActions && (
                  /* Stops the checkbox from also opening the reservation behind it. */
                  <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
                    <Checkbox
                      checked={selectedPenaltyIds.has(p.id)}
                      onCheckedChange={() => togglePenalty(p.id)}
                      className={CHECKBOX_CLASSES}
                      aria-label={`בחר חיוב ${p.id}`}
                    />
                  </td>
                )}
                <td className="px-4 py-3 whitespace-nowrap text-navy font-medium">
                  {p.id}
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  {p.brokerReservationId || "—"}
                </td>
                <td className="px-4 py-3 whitespace-nowrap">
                  {penaltyTypeLabel(p.type)}
                </td>
                {/* A fee carries no status, and no price breakdown: it is passed on as charged. */}
                <td className="px-4 py-3 whitespace-nowrap">—</td>
                <td className="px-4 py-3 whitespace-nowrap">—</td>
                <td className="px-4 py-3 text-start whitespace-nowrap">—</td>
                <td className="px-4 py-3 text-start whitespace-nowrap">—</td>
                <td className="px-4 py-3 text-start whitespace-nowrap">—</td>
                <td className="px-4 py-3 text-start whitespace-nowrap">—</td>
                <td className="px-4 py-3 text-start whitespace-nowrap font-medium">
                  {currencyFormatter.format(p.amount)}
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                  —
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                  {formatDate(p.createdAt)}
                </td>
                <td className="px-4 py-3 whitespace-nowrap text-text-secondary">
                  —
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <BillDialog
        open={billDialogOpen}
        onOpenChange={setBillDialogOpen}
        entity={entity}
        currencyCode={group.currencyCode}
        selectedIds={selectedIds}
        selectedPenaltyIds={selectedPenaltyIdList}
        totalDue={selectedTotal}
        onSuccess={clearSelection}
      />
      <PromptBillDialog
        open={promptBillDialogOpen}
        onOpenChange={setPromptBillDialogOpen}
        entity={entity}
        selectedIds={selectedIds}
        selectedPenaltyIds={selectedPenaltyIdList}
        onSuccess={clearSelection}
        onContinueToInvoiceCreation={() => setBillDialogOpen(true)}
      />
      <ReservationDetailDialog
        reservationId={detailReservationId}
        onClose={() => setDetailReservationId(null)}
      />
    </section>
  );
}
