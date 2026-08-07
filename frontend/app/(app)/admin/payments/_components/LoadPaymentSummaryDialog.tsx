"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import {
  validateFlexPaymentSummary,
  type RejectedSupplierPayment,
  type ValidateFlexPaymentSummaryResponse,
} from "@/shared/api/reservations-api";
import { formatPriceFloat } from "@/shared/utils/formatPrice";

const REJECTION_REASONS: Record<string, string> = {
  not_found_or_already_paid: "לא נמצאה הזמנה פתוחה / כבר שולמה",
  invalid_price: "סכום לא תואם",
};

const rejectionReason = (reason: string) => REJECTION_REASONS[reason] ?? reason;

interface LoadPaymentSummaryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Selects the approved reservations in the table behind the dialog. */
  onSelectApproved: (reservationIds: number[]) => void;
}

export function LoadPaymentSummaryDialog({
  open,
  onOpenChange,
  onSelectApproved,
}: LoadPaymentSummaryDialogProps) {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<ValidateFlexPaymentSummaryResponse>();
  const [submitError, setSubmitError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: (selectedFile: File) =>
      validateFlexPaymentSummary(selectedFile),
    onSuccess: setResult,
    onError: () =>
      setSubmitError("טעינת הקובץ נכשלה. ודא שנבחר קובץ חיובים תקין ונסה שוב."),
  });

  const reset = () => {
    setFile(null);
    setResult(undefined);
    setSubmitError(null);
    mutation.reset();
  };

  const handleOpenChange = (next: boolean) => {
    if (!next) reset();
    onOpenChange(next);
  };

  const onSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (!file) return;
    setSubmitError(null);
    mutation.mutate(file);
  };

  const markSelected = () => {
    onSelectApproved(result?.approved.map((r) => r.reservationId) ?? []);
    handleOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="min-w-96 max-w-2xl p-6 flex flex-col gap-6 bg-white border-border-light/50 rounded-2xl shadow-modal">
        {/* ps-8 keeps the heading clear of the close button, which sits in the top start corner. */}
        {result ? (
          <DialogTitle className="sr-only">טען קובץ חיובים</DialogTitle>
        ) : (
          <div className="flex flex-col gap-1 ps-8">
            <DialogTitle className="type-h5 text-navy">
              טען קובץ חיובים
            </DialogTitle>
            <p className="type-paragraph text-text-secondary">
              העלה את קובץ החיובים שהתקבל מהספק כדי להשוות אותו להזמנות הפתוחות.
            </p>
          </div>
        )}

        {result ? (
          <div className="flex flex-col gap-5">
            <ApprovedSection approved={result.approved} />
            <RejectedSection rejected={result.rejected} />

            <div className="flex gap-3 pt-2">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                onClick={() => handleOpenChange(false)}
              >
                סגור
              </Button>
              <Button
                type="button"
                variant="brand"
                className="flex-1 h-10"
                disabled={result.approved.length === 0}
                onClick={markSelected}
              >
                סמן נבחרים
              </Button>
            </div>
          </div>
        ) : (
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="summary-file" className="type-label text-navy">
                קובץ חיובים (אקסל)
              </Label>
              <Input
                id="summary-file"
                type="file"
                accept=".xlsx"
                dir="ltr"
                className="text-base"
                onChange={(e) => setFile(e.target.files?.[0] ?? null)}
              />
            </div>

            {submitError && <ErrorDisplay>{submitError}</ErrorDisplay>}

            <div className="flex gap-3 pt-2">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                onClick={() => handleOpenChange(false)}
                disabled={mutation.isPending}
              >
                ביטול
              </Button>
              <Button
                type="submit"
                variant="brand"
                className="flex-1 h-10"
                disabled={!file}
                loading={mutation.isPending}
              >
                שלח
              </Button>
            </div>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}

function ApprovedSection({
  approved,
}: {
  approved: ValidateFlexPaymentSummaryResponse["approved"];
}) {
  return (
    <section className="flex flex-col gap-2">
      <h3 className="type-h6 text-navy ps-8">
        מאושרות לתשלום ({approved.length})
      </h3>
      {approved.length === 0 ? (
        <p className="type-paragraph text-text-secondary">
          אף שורה בקובץ לא תואמת הזמנה פתוחה.
        </p>
      ) : (
        <ul className="flex flex-col gap-1 max-h-48 overflow-y-auto rounded-xl border border-border-light/60 p-3">
          {approved.map((r) => (
            <li
              key={r.reservationId}
              className="flex items-center justify-between gap-4 text-sm"
            >
              <span className="text-navy font-medium">
                {r.brokerReservationId}
              </span>
              <span className="text-text-secondary">
                {formatPriceFloat(r.amount, r.currencyCode)}
              </span>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function RejectedSection({
  rejected,
}: {
  rejected: RejectedSupplierPayment[];
}) {
  if (rejected.length === 0) return null;

  return (
    <section className="flex flex-col gap-2">
      <h3 className="type-h6 text-navy">לא מאושרות ({rejected.length})</h3>
      <ul className="flex flex-col gap-2 max-h-48 overflow-y-auto rounded-xl border border-border-light/60 p-3">
        {rejected.map((r, index) => (
          <li
            key={`${r.brokerReservationId}-${index}`}
            className="flex flex-col gap-0.5 text-sm"
          >
            <div className="flex items-center justify-between gap-4">
              <span className="text-navy font-medium">
                {r.brokerReservationId}
              </span>
              <span className="text-destructive">
                {rejectionReason(r.reason)}
              </span>
            </div>
            {r.expectedAmount !== undefined && (
              <span className="text-text-secondary">
                בקובץ {formatPriceFloat(r.balance, r.currencyCode)} • אצלנו{" "}
                {formatPriceFloat(
                  r.expectedAmount,
                  r.expectedCurrencyCode || r.currencyCode,
                )}
              </span>
            )}
          </li>
        ))}
      </ul>
    </section>
  );
}
