"use client";

import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns/format";
import { he } from "date-fns/locale/he";
import { CalendarIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { bill, isBillingFailedError } from "@/shared/api/bill-api";
import { cn } from "@/lib/utils";
import type { BillingEntity } from "./BillingEntityCombobox";
import { toast } from "sonner";

interface PromptBillDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  entity: BillingEntity;
  selectedIds: number[];
  onSuccess: () => void;
  onContinueToInvoiceCreation: () => void;
}

export function PromptBillDialog({
  open,
  onOpenChange,
  entity,
  selectedIds,
  onSuccess,
  onContinueToInvoiceCreation,
}: PromptBillDialogProps) {
  const queryClient = useQueryClient();
  const [isSuccess, setIsSuccess] = useState(false);
  const [submitError, setSubmitError] = useState<Error | null>(null);
  const translatedError = useTranslatedError(
    submitError && isBillingFailedError(submitError) ? null : submitError,
  );

  const errorMessage =
    submitError && isBillingFailedError(submitError)
      ? submitError.message
      : translatedError;

  const { isPending, mutate } = useMutation({
    mutationFn: () =>
      bill({
        ids: selectedIds,
        // This screen does not select fees yet, so none are billed.
        penalty_ids: [],
        total_paid: 1, //skip validation
        transfer_date: "2006-01-02", //skip validation
        organization_id: entity.kind === "org" ? entity.id : undefined,
        office_id: entity.kind === "office" ? entity.id : undefined,
        skip_invoice_creation: true,
      }),
    onSuccess: (res) => {
      setIsSuccess(true);
      queryClient.invalidateQueries({
        queryKey: ["open-reservations", entity.kind, entity.id],
      });
      setTimeout(() => {
        onSuccess();
        setIsSuccess(false);
        onOpenChange(false);
      }, 2500);
    },
    onError: (err: Error) => setSubmitError(err),
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="min-w-96 max-w-lg! p-6 flex flex-col gap-6 bg-white border-border-light/50 rounded-2xl shadow-modal">
        <div className="flex flex-col gap-1">
          <DialogTitle className="type-h5 text-navy">
            האם תרצה להפיק חשבונית?
          </DialogTitle>
          <p>⚠️ יש לבחור 'לא' אך ורק אם כבר הופקה חשבונית מחוץ למערכת.</p>
        </div>
        {isSuccess && <SuccessBadge>החיוב סומן כשולם בהצלחה!</SuccessBadge>}
        <div className="flex gap-2">
          <Button
            className="font-bold py-6 rounded-lg"
            variant="brand"
            onClick={onContinueToInvoiceCreation}
          >
            כן, בוא נתקדם להפקת החשבונית
          </Button>
          <Button
            variant="destructive"
            className="font-bold py-6 rounded-lg bg-transparent text-destructive border-destructive hover:bg-destructive/10"
            onClick={() => mutate()}
            loading={isPending}
          >
            לא, רק לסמן את החשבוניות כשולמו
          </Button>
        </div>

        {errorMessage && <ErrorDisplay>{errorMessage}</ErrorDisplay>}
      </DialogContent>
    </Dialog>
  );
}
