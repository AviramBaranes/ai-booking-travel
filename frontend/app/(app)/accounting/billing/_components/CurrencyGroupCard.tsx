"use client";

import { useMemo, useState } from "react";
import { format } from "date-fns/format";
import { he } from "date-fns/locale/he";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { reservation } from "@/shared/client";
import { BillDialog } from "./BillDialog";
import type { BillingEntity } from "./BillingEntityCombobox";

const CHECKBOX_CLASSES =
  "border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand";

interface CurrencyGroupCardProps {
  entity: BillingEntity;
  group: reservation.CurrencyGroup;
}

export function CurrencyGroupCard({ entity, group }: CurrencyGroupCardProps) {
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [dialogOpen, setDialogOpen] = useState(false);
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

  const selectedIds = useMemo(
    () => group.reservations.filter((r) => selected.has(r.id)).map((r) => r.id),
    [selected, group.reservations],
  );

  const allChecked =
    group.reservations.length > 0 &&
    selectedIds.length === group.reservations.length;
  const someChecked = selectedIds.length > 0 && !allChecked;

  const toggleAll = () => {
    if (allChecked) {
      setSelected(new Set());
    } else {
      setSelected(new Set(group.reservations.map((r) => r.id)));
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
            {group.reservations.length} הזמנות פתוחות
            {selectedIds.length > 0 && (
              <span className="text-navy font-medium">
                {" "}
                • {selectedIds.length} נבחרו
              </span>
            )}
          </p>
        </div>
        <Button
          type="button"
          variant="brand"
          className="h-10 px-6"
          disabled={selectedIds.length === 0}
          onClick={() => setDialogOpen(true)}
        >
          חייב נבחרים
        </Button>
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
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מס׳ הזמנה
              </th>
              <th className="px-4 py-3 text-start whitespace-nowrap">
                מס׳ ספק
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
                    "border-t border-border-light/40 transition-colors " +
                    (isRefundPending
                      ? "bg-brand/10 hover:bg-brand/15"
                      : "hover:bg-background/60")
                  }
                >
                  <td className="px-4 py-3">
                    <Checkbox
                      checked={checked}
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
                      (r.profitOnCar >= 0 ? "text-success" : "text-destructive")
                    }
                  >
                    {currencyFormatter.format(r.profitOnCar)}
                  </td>
                  <td className="px-4 py-3 text-start whitespace-nowrap font-medium">
                    {currencyFormatter.format(r.totalPrice)}
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
          </tbody>
        </table>
      </div>

      {dialogOpen && (
        <BillDialog
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          entity={entity}
          currencyCode={group.currencyCode}
          selectedIds={selectedIds}
          onSuccess={() => setSelected(new Set())}
        />
      )}
    </section>
  );
}
