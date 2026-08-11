"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import type { penalties } from "@/shared/client";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import {
  createPenalty,
  validateFlexPaymentSummary,
  type RejectedSupplierPayment,
  type ValidateFlexPaymentSummaryResponse,
} from "@/shared/api/reservations-api";
import { formatPriceFloat } from "@/shared/utils/formatPrice";

const REJECTION_REASONS: Record<string, string> = {
  not_found: "לא נמצאה הזמנה במערכת",
  already_paid: "כבר שולמה לספק",
  canceled: "ההזמנה בוטלה — טרם נרשם קנס",
  invalid_price: "סכום לא תואם",
};

const rejectionReason = (reason: string) => REJECTION_REASONS[reason] ?? reason;

/** REASON_CANCELED marks the lines we can offer to record as a fee. */
const REASON_CANCELED = "canceled";

/**
 * Flex charges a no-show at a higher rate than a late cancellation, so the amount on the line is
 * what tells the two apart.
 */
const NO_SHOW_MIN_AMOUNT = 60;

const PENALTY_LABELS: Record<string, string> = {
  no_show: "אי-הגעה",
  cancellation: "ביטול מאוחר",
};

const suggestedPenaltyType = (balance: number) =>
  balance > NO_SHOW_MIN_AMOUNT ? "no_show" : "cancellation";

interface LoadPaymentSummaryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Selects the approved reservations and fees in the table behind the dialog. */
  onSelectApproved: (reservationIds: number[], penaltyIds: number[]) => void;
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
    onSelectApproved(
      result?.approved.map((r) => r.reservationId) ?? [],
      result?.approvedPenalties.map((p) => p.penaltyId) ?? [],
    );
    handleOpenChange(false);
  };

  // A fee recorded from a suggestion settles the line it was suggested on, so it moves out of the
  // rejected list and to the top of the approved one, where it is easy to spot.
  const onPenaltyCreated = (
    brokerReservationId: string,
    penalty: penalties.CreatePenaltyResponse,
  ) => {
    setResult((prev) =>
      prev
        ? {
            ...prev,
            rejected: prev.rejected.filter(
              (r) => r.reservationId !== penalty.reservationId,
            ),
            approvedPenalties: [
              {
                penaltyId: penalty.id,
                reservationId: penalty.reservationId,
                brokerReservationId,
                type: penalty.type,
                currencyCode: penalty.currencyCode,
                amount: penalty.amount,
              },
              ...prev.approvedPenalties,
            ],
          }
        : prev,
    );
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
            <ApprovedSection
              approved={result.approved}
              approvedPenalties={result.approvedPenalties}
            />
            <RejectedSection
              rejected={result.rejected}
              rejectedPenalties={result.rejectedPenalties}
              onPenaltyCreated={onPenaltyCreated}
            />

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
  approvedPenalties,
}: {
  approved: ValidateFlexPaymentSummaryResponse["approved"];
  approvedPenalties: ValidateFlexPaymentSummaryResponse["approvedPenalties"];
}) {
  const total = approved.length + approvedPenalties.length;

  return (
    <section className="flex flex-col gap-2">
      <h3 className="type-h6 text-navy ps-8">מאושרות לתשלום ({total})</h3>
      {total === 0 ? (
        <p className="type-paragraph text-text-secondary">
          אף שורה בקובץ לא תואמת הזמנה או קנס פתוחים.
        </p>
      ) : (
        <ul className="flex flex-col gap-1 max-h-48 overflow-y-auto rounded-xl border border-border-light/60 p-3">
          {/* Fees lead the list, so one just recorded from a suggestion is the first thing seen. */}
          {approvedPenalties.map((p) => (
            <li
              key={`penalty-${p.penaltyId}`}
              className="flex items-center justify-between gap-4 text-sm"
            >
              <span className="text-navy font-medium">
                {p.brokerReservationId}
                <span className="text-text-secondary font-normal">
                  {" "}
                  • קנס {PENALTY_LABELS[p.type] ?? p.type}
                </span>
              </span>
              <span className="text-text-secondary">
                {formatPriceFloat(p.amount, p.currencyCode)}
              </span>
            </li>
          ))}
          {approved.map((r) => (
            <li
              key={`reservation-${r.reservationId}`}
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
  rejectedPenalties,
  onPenaltyCreated,
}: {
  rejected: RejectedSupplierPayment[];
  rejectedPenalties: ValidateFlexPaymentSummaryResponse["rejectedPenalties"];
  onPenaltyCreated: (
    brokerReservationId: string,
    penalty: penalties.CreatePenaltyResponse,
  ) => void;
}) {
  const total = rejected.length + rejectedPenalties.length;
  if (total === 0) return null;

  return (
    <section className="flex flex-col gap-2">
      <h3 className="type-h6 text-navy">לא מאושרות ({total})</h3>
      <ul className="flex flex-col gap-3 max-h-64 overflow-y-auto rounded-xl border border-border-light/60 p-3">
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
            {r.reason === REASON_CANCELED && r.reservationId !== undefined && (
              <CreatePenaltySuggestion
                reservationId={r.reservationId}
                brokerReservationId={r.brokerReservationId}
                balance={r.balance}
                currencyCode={r.currencyCode}
                onCreated={onPenaltyCreated}
              />
            )}
          </li>
        ))}

        {rejectedPenalties.map((p) => (
          <li
            key={`penalty-${p.penaltyId}`}
            className="flex flex-col gap-0.5 text-sm"
          >
            <div className="flex items-center justify-between gap-4">
              <span className="text-navy font-medium">
                {p.brokerReservationId}
                <span className="text-text-secondary font-normal">
                  {" "}
                  • קנס {PENALTY_LABELS[p.type] ?? p.type}
                </span>
              </span>
              <span className="text-destructive">
                {rejectionReason(p.reason)}
              </span>
            </div>
            <span className="text-text-secondary">
              בקובץ {formatPriceFloat(p.balance, p.currencyCode)} • אצלנו{" "}
              {formatPriceFloat(
                p.expectedAmount,
                p.expectedCurrencyCode || p.currencyCode,
              )}
            </span>
          </li>
        ))}
      </ul>
    </section>
  );
}

/**
 * CreatePenaltySuggestion offers to record the fee behind a line whose reservation was canceled.
 * The amount is what the supplier charged, and the type is inferred from it.
 */
function CreatePenaltySuggestion({
  reservationId,
  brokerReservationId,
  balance,
  currencyCode,
  onCreated,
}: {
  reservationId: number;
  brokerReservationId: string;
  balance: number;
  currencyCode: string;
  onCreated: (
    brokerReservationId: string,
    penalty: penalties.CreatePenaltyResponse,
  ) => void;
}) {
  const queryClient = useQueryClient();
  const type = suggestedPenaltyType(balance);

  const { mutate, isPending, isError } = useMutation({
    mutationFn: () => createPenalty({ reservationId, type, amount: balance }),
    // Awaited so the table behind holds the new fee by the time the approved lines are marked.
    onSuccess: async (penalty) => {
      await queryClient.invalidateQueries({
        queryKey: ["unpaid-supplier-reservations"],
      });
      onCreated(brokerReservationId, penalty);
    },
  });

  return (
    <div className="flex flex-col gap-1 items-start">
      <Button
        type="button"
        variant="outline"
        className="h-8 px-3"
        loading={isPending}
        onClick={() => mutate()}
      >
        צור קנס {PENALTY_LABELS[type]} על סך{" "}
        {formatPriceFloat(balance, currencyCode)}
      </Button>
      {isError && <ErrorDisplay>יצירת הקנס נכשלה. נסה שוב.</ErrorDisplay>}
    </div>
  );
}
