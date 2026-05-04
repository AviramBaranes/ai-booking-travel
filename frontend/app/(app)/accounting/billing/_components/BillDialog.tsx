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
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { bill } from "@/shared/api/bill-api";
import { cn } from "@/lib/utils";
import type { BillingEntity } from "./BillingEntityCombobox";

const schema = z.object({
  total_paid: z
    .number({ message: "שדה חובה" })
    .positive("הסכום חייב להיות גדול מאפס"),
  transfer_date: z.date({ message: "יש לבחור תאריך" }),
});

type FormValues = z.infer<typeof schema>;

interface BillDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  entity: BillingEntity;
  currencyCode: string;
  selectedIds: number[];
  onSuccess: () => void;
}

export function BillDialog({
  open,
  onOpenChange,
  entity,
  currencyCode,
  selectedIds,
  onSuccess,
}: BillDialogProps) {
  const queryClient = useQueryClient();
  const [submitError, setSubmitError] = useState<Error | null>(null);
  const [datePopoverOpen, setDatePopoverOpen] = useState(false);
  const [success, setSuccess] = useState(false);
  const translatedError = useTranslatedError(submitError);

  const {
    control,
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      total_paid: undefined as unknown as number,
      transfer_date: undefined as unknown as Date,
    },
  });

  const mutation = useMutation({
    mutationFn: (values: FormValues) =>
      bill({
        ids: selectedIds,
        total_paid: values.total_paid,
        transfer_date: format(values.transfer_date, "yyyy-MM-dd"),
        organization_id: entity.kind === "org" ? entity.id : undefined,
        office_id: entity.kind === "office" ? entity.id : undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["open-reservations", entity.kind, entity.id],
      });
      setSuccess(true);
      setTimeout(() => {
        onSuccess();
        onOpenChange(false);
      }, 1500);
    },
    onError: (err: Error) => setSubmitError(err),
  });

  const onSubmit = (values: FormValues) => {
    setSubmitError(null);
    mutation.mutate(values);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="min-w-96 max-w-md p-6 flex flex-col gap-6 bg-white border-border-light/50 rounded-2xl shadow-modal">
        <div className="flex flex-col gap-1">
          <DialogTitle className="type-h5 text-navy">
            הפקת חיוב
          </DialogTitle>
          <p className="type-paragraph text-text-secondary">
            {`${entity.name} • ${currencyCode} • ${selectedIds.length} הזמנות`}
          </p>
        </div>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-4"
        >
          {success ? (
            <SuccessBadge text="החיוב הופק בהצלחה" />
          ) : (
            <>
              <div className="flex flex-col gap-1.5">
            <Label htmlFor="total_paid" className="type-label text-navy">
              סכום ששולם ({currencyCode})
            </Label>
            <Input
              id="total_paid"
              type="number"
              step="0.01"
              min="0"
              dir="ltr"
              className="text-base"
              placeholder="0.00"
              {...register("total_paid", { valueAsNumber: true })}
            />
            {errors.total_paid && (
              <ErrorDisplay>{errors.total_paid.message}</ErrorDisplay>
            )}
          </div>

          <div className="flex flex-col gap-1.5">
            <Label className="type-label text-navy">תאריך העברה</Label>
            <Controller
              control={control}
              name="transfer_date"
              render={({ field }) => (
                <Popover
                  open={datePopoverOpen}
                  onOpenChange={setDatePopoverOpen}
                >
                  <PopoverTrigger asChild>
                    <Button
                      type="button"
                      variant="outline"
                      className={cn(
                        "w-full justify-between font-normal",
                        !field.value && "text-muted-foreground",
                      )}
                    >
                      {field.value
                        ? format(field.value, "d בMMMM yyyy", { locale: he })
                        : "בחר תאריך..."}
                      <CalendarIcon className="size-4 opacity-60" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      locale={he}
                      selected={field.value}
                      onSelect={(date) => {
                        if (!date) return;
                        field.onChange(date);
                        setDatePopoverOpen(false);
                      }}
                      autoFocus
                    />
                  </PopoverContent>
                </Popover>
              )}
            />
            {errors.transfer_date && (
              <ErrorDisplay>{errors.transfer_date.message}</ErrorDisplay>
            )}
          </div>

          {translatedError && <ErrorDisplay>{translatedError}</ErrorDisplay>}

          <div className="flex gap-3 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              className="flex-1"
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
              חייב
            </Button>
          </div>
            </>
          )}
        </form>
      </DialogContent>
    </Dialog>
  );
}
