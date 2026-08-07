"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns/format";
import { he } from "date-fns/locale/he";
import { CalendarIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { paySupplierReservations } from "@/shared/api/reservations-api";
import { formatDate } from "@/shared/utils/formatDate";
import { cn } from "@/lib/utils";
import type { supplier_payments } from "@/shared/client";

const FAILURE_REASONS: Record<string, string> = {
  unsupported_broker: "ספק לא נתמך",
  expense_creation_failed: "יצירת ההוצאה נכשלה",
  mark_paid_failed: "ההוצאה נוצרה אך ההזמנה לא סומנה כשולמה",
};

const failureReason = (reason: string) => FAILURE_REASONS[reason] ?? reason;

interface MarkAsPaidDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  brokerCode: string;
  reservationIds: number[];
  /** Clears the table selection once the payment went through. */
  onPaid: () => void;
}

export function MarkAsPaidDialog({
  open,
  onOpenChange,
  brokerCode,
  reservationIds,
  onPaid,
}: MarkAsPaidDialogProps) {
  const queryClient = useQueryClient();
  const [paymentDate, setPaymentDate] = useState<Date>(new Date());
  const [datePopoverOpen, setDatePopoverOpen] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [result, setResult] = useState<supplier_payments.PaySupplierReservationsResponse>();

  const mutation = useMutation({
    mutationFn: (date: Date) =>
      paySupplierReservations({
        reservationIds,
        paymentDate: formatDate(date),
      }),
    onSuccess: (response) => {
      queryClient.invalidateQueries({
        queryKey: ["unpaid-supplier-reservations", brokerCode],
      });
      onPaid();

      // Everything settled — nothing left to report, so get out of the way.
      if (response.failed.length === 0 && response.skipped.length === 0) {
        close();
        return;
      }
      setResult(response);
    },
    onError: () =>
      setSubmitError("סימון התשלום נכשל. נסה שוב או פנה לתמיכה."),
  });

  const close = () => {
    setPaymentDate(new Date());
    setSubmitError(null);
    setResult(undefined);
    mutation.reset();
    onOpenChange(false);
  };

  const handleOpenChange = (next: boolean) => {
    if (next) onOpenChange(true);
    else close();
  };

  const onSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitError(null);
    mutation.mutate(paymentDate);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="min-w-96 max-w-md p-6 flex flex-col gap-6 bg-white border-border-light/50 rounded-2xl shadow-modal">
        {/* ps-8 keeps the heading clear of the close button in the top start corner. */}
        <div className="flex flex-col gap-1 ps-8">
          <DialogTitle className="type-h5 text-navy">סמן כשולם</DialogTitle>
          <p className="type-paragraph text-text-secondary">
            {reservationIds.length} הזמנות ייווצרו כהוצאה ויסומנו כשולמו לספק.
          </p>
        </div>

        {result ? (
          <div className="flex flex-col gap-4">
            <p className="type-paragraph text-navy">
              {result.paid.length} הזמנות סומנו כשולמו.
            </p>

            {result.skipped.length > 0 && (
              <p className="type-paragraph text-text-secondary">
                {result.skipped.length} הזמנות דולגו מכיוון שכבר שולמו או בוטלו.
              </p>
            )}

            {result.failed.length > 0 && (
              <section className="flex flex-col gap-2">
                <h3 className="type-h6 text-navy">
                  נכשלו ({result.failed.length})
                </h3>
                <ul className="flex flex-col gap-1 max-h-48 overflow-y-auto rounded-xl border border-border-light/60 p-3">
                  {result.failed.map((f) => (
                    <li
                      key={f.reservationId}
                      className="flex items-center justify-between gap-4 text-sm"
                    >
                      <span className="text-navy font-medium">
                        {f.reservationId}
                      </span>
                      <span className="text-destructive">
                        {failureReason(f.reason)}
                        {f.expenseId && ` (הוצאה ${f.expenseId})`}
                      </span>
                    </li>
                  ))}
                </ul>
              </section>
            )}

            <Button
              type="button"
              variant="brand"
              className="h-10"
              onClick={close}
            >
              סגור
            </Button>
          </div>
        ) : (
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <Label className="type-label text-navy">תאריך תשלום</Label>
              <Popover open={datePopoverOpen} onOpenChange={setDatePopoverOpen}>
                <PopoverTrigger asChild>
                  <Button
                    type="button"
                    variant="outline"
                    className={cn("w-full justify-between font-normal")}
                  >
                    {format(paymentDate, "d בMMMM yyyy", { locale: he })}
                    <CalendarIcon className="size-4 opacity-60" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    locale={he}
                    selected={paymentDate}
                    onSelect={(date) => {
                      if (!date) return;
                      setPaymentDate(date);
                      setDatePopoverOpen(false);
                    }}
                    autoFocus
                  />
                </PopoverContent>
              </Popover>
            </div>

            {submitError && <ErrorDisplay>{submitError}</ErrorDisplay>}

            <div className="flex gap-3 pt-2">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                onClick={close}
                disabled={mutation.isPending}
              >
                ביטול
              </Button>
              <Button
                type="submit"
                variant="brand"
                className="flex-1 h-10"
                loading={mutation.isPending}
              >
                סמן כשולם
              </Button>
            </div>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}
