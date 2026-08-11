"use client";

import { useQuery } from "@tanstack/react-query";
import Image from "next/image";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { getFullReservation } from "@/shared/api/reservations-api";
import { queries } from "@/shared/client";
import { formatPrice, formatPriceFloat } from "@/shared/utils/formatPrice";
import { CreatePenaltyForm } from "./CreatePenaltyForm";

const RESERVATION_STATUS_CANCELED = "canceled";

// ---------------------------------------------------------------------------
// Primitive display helpers
// ---------------------------------------------------------------------------

function Row({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start gap-2 py-1 text-sm">
      <span className="w-52 shrink-0 font-bold text-slate-700">{label}</span>
      <span className="whitespace-nowrap text-slate-600">{value ?? "—"}</span>
    </div>
  );
}

function formatIsraeliDateTime(iso: string): string {
  if (!iso) return "—";
  try {
    return new Intl.DateTimeFormat("he-IL", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      timeZone: "Asia/Jerusalem",
    }).format(new Date(iso));
  } catch {
    return iso;
  }
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mt-4 first:mt-0">
      <h3 className="mb-1 border-b border-slate-200 pb-1 text-xs font-bold uppercase tracking-wider text-slate-400">
        {title}
      </h3>
      {children}
    </div>
  );
}

function yesNo(val: boolean) {
  return val ? "כן" : "לא";
}

function pct(val: number) {
  return `${val.toFixed(2)}%`;
}

// ---------------------------------------------------------------------------
// Pay-at-pickup inline section (mirrors PayAtPickupSection logic)
// ---------------------------------------------------------------------------

function PayAtPickupSection({
  r,
}: {
  r: queries.GetFullReservationResponse;
}) {
  const { fees, selectedAddons, currencyCode } = r;
  const hasContent =
    fees.dropCharge > 0 ||
    fees.youngDriverFee > 0 ||
    (selectedAddons && selectedAddons.length > 0);

  if (!hasContent) return <Row label="תוכן" value="אין" />;

  return (
    <>
      {fees.dropCharge > 0 && (
        <Row
          label="תוספת החזרה בנקודה אחרת"
          value={formatPrice(fees.dropCharge, fees.dropChargeCurrency)}
        />
      )}
      {fees.youngDriverFee > 0 && (
        <Row
          label="תוספת נהג צעיר"
          value={formatPrice(fees.youngDriverFee, fees.youngDriverFeeCurrency)}
        />
      )}
      {selectedAddons?.map((addon) => (
        <Row
          key={addon.id}
          label={addon.name}
          value={`${formatPrice(addon.price, currencyCode)} × ${addon.quantity}`}
        />
      ))}
    </>
  );
}

// ---------------------------------------------------------------------------
// Main dialog
// ---------------------------------------------------------------------------

interface ReservationDetailDialogProps {
  reservationId: number | null;
  onClose: () => void;
}

export function ReservationDetailDialog({
  reservationId,
  onClose,
}: ReservationDetailDialogProps) {
  const { data: r, isLoading } = useQuery({
    queryKey: ["reservation-detail", reservationId],
    queryFn: () => getFullReservation(reservationId!),
    enabled: reservationId != null,
  });

  return (
    <Dialog open={reservationId != null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[90vh] w-[90vw] max-w-4xl! overflow-y-auto" dir="rtl">
        <DialogHeader>
          {/* ps-8 keeps the heading clear of the close button in the top start corner. */}
          <DialogTitle className="ps-8">
            {r ? `הזמנה #${r.id} — ${r.brokerReservationId}` : "פרטי הזמנה"}
          </DialogTitle>
        </DialogHeader>

        {isLoading && (
          <div className="space-y-2 py-4">
            {Array.from({ length: 10 }).map((_, i) => (
              <div key={i} className="h-4 w-full animate-pulse rounded bg-slate-200" />
            ))}
          </div>
        )}

        {r && (
          <div className="pb-2">
            {/* Car image */}
            {r.imageUrl && (
              <div className="mb-4 flex justify-center">
                <Image
                  src={r.imageUrl}
                  alt={r.model}
                  width={280}
                  height={160}
                  className="rounded-lg object-contain"
                  unoptimized
                />
              </div>
            )}

            <Section title="זיהוי">
              <Row label="מזהה הזמנה" value={r.id} />
              <Row label="מזהה ספק" value={r.brokerReservationId} />
              <Row label="מזהה משתמש" value={r.userId} />
              {r.officeId && <Row label="מזהה משרד" value={r.officeId} />}
              {r.organizationId && <Row label="מזהה ארגון" value={r.organizationId} />}
              {r.adminRefId && <Row label="מזהה אדמין" value={r.adminRefId} />}
            </Section>

            <Section title="סטטוס">
              <Row label="סטטוס הזמנה" value={r.reservationStatus} />
              <Row label="סטטוס תשלום" value={r.paymentStatus} />
              {r.voucherNumber && <Row label="מספר שובר" value={r.voucherNumber} />}
              {r.voucheredAt && <Row label="תאריך כרטוס" value={formatIsraeliDateTime(r.voucheredAt)} />}
            </Section>

            <Section title="רכב">
              <Row label="דגם" value={r.model} />
              <Row label="קבוצה" value={r.carGroup} />
              <Row label="סוג" value={r.carType} />
              <Row label="קוד ACRISS" value={r.acriss} />
              <Row label="ספק" value={r.supplierName} />
              <Row label="מושבים" value={r.seats} />
              <Row label="תיקים" value={r.bags} />
              <Row label="דלתות" value={r.doors} />
              <Row label="מיזוג" value={yesNo(r.hasAC)} />
              <Row label="גיר אוטומטי" value={yesNo(r.isAutoGear)} />
              <Row label="חשמלי" value={yesNo(r.isElectric)} />
            </Section>

            <Section title="נהג">
              <Row label="שם" value={`${r.driverTitle} ${r.driverFirstName} ${r.driverLastName}`} />
              <Row label="גיל" value={r.driverAge} />
            </Section>

            <Section title="פרטי נסיעה">
              <Row label="מדינה" value={r.countryCode} />
              <Row label="תאריך איסוף" value={`${r.pickupDate} ${r.pickupTime}`} />
              <Row label="תאריך החזרה" value={`${r.dropoffDate} ${r.dropoffTime}`} />
              <Row label="ימי השכרה" value={r.rentalDays} />
              <Row label="נקודת איסוף" value={r.pickupLocationName} />
              <Row label="נקודת החזרה" value={r.dropoffLocationName} />
              {r.flightNumber && <Row label="מספר טיסה" value={r.flightNumber} />}
            </Section>

            <Section title="תכנית">
              <Row label="ברוקר" value={r.broker} />
              <Row label="קוד ספק" value={r.supplierCode} />
              <Row
                label="מה כלול"
                value={
                  r.planInclusions?.length ? (
                    <ul className="list-inside list-disc space-y-0.5">
                      {r.planInclusions.map((inc) => (
                        <li key={inc}>{inc}</li>
                      ))}
                    </ul>
                  ) : (
                    "אין"
                  )
                }
              />
            </Section>

            <Section title="תשלום באיסוף">
              <PayAtPickupSection r={r} />
            </Section>

            <Section title="תמחור">
              <Row label="מטבע" value={r.currencyCode} />
              <Row label="שער חליפין" value={r.currencyRate.toFixed(4)} />
              <Row label="מע״מ" value={pct(r.vatPercentage)} />
              <Row label="מחיר רכישה" value={formatPriceFloat(r.purchasePrice, r.currencyCode)} />
              <Row label="מחיר רכישה בש״ח" value={formatPriceFloat(r.purchasePriceInILS, "ILS")} />
              <Row label="אחוז עמלה" value={pct(r.markupPercentage)} />
              <Row label="מחיר ERP ברוקר" value={formatPriceFloat(r.brokerErpPrice, r.currencyCode)} />
              <Row label='מחיר מכירה + ERP (מטבע)' value={formatPriceFloat(r.carSellPriceWithBrokerERP, r.currencyCode)} />
              <Row label='מחיר מכירה + ERP (ש״ח)' value={formatPriceFloat(r.carSellPriceWithBrokerERPInILS, "ILS")} />
              <Row label="מחיר ERP BT" value={formatPriceFloat(r.btErpPrice, r.currencyCode)} />
              <Row label="מחיר ERP BT בש״ח" value={formatPriceFloat(r.btErpPriceInILS, "ILS")} />
              <Row label="הנחה" value={pct(r.discountPercentage)} />
              <Row label="מחיר סופי" value={formatPriceFloat(r.totalPrice, r.currencyCode)} />
              <Row label="מחיר סופי בש״ח" value={formatPriceFloat(r.totalPriceInILS, "ILS")} />
              <Row label="רווח" value={formatPriceFloat(r.profit, r.currencyCode)} />
              <Row label="רווח בש״ח" value={formatPriceFloat(r.profitInILS, "ILS")} />
              <Row label="אחוז רווח" value={pct(r.profitPercentage)} />
            </Section>

            <Section title="זמנים">
              <Row label="נוצר ב" value={formatIsraeliDateTime(r.createdAt)} />
              <Row label="עודכן ב" value={formatIsraeliDateTime(r.updatedAt)} />
            </Section>

            {/* A fee is only ever charged on a reservation that was canceled. */}
            {r.reservationStatus === RESERVATION_STATUS_CANCELED && (
              <Section title="דמי ביטול / אי הגעה">
                <CreatePenaltyForm
                  reservationId={r.id}
                  currencyCode={r.currencyCode}
                />
              </Section>
            )}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
