"use client";

import { Button } from "@/components/ui/button";
import { Dialog, DialogTitle, DialogContent } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { insertLocationAlias } from "@/shared/api/locations-api";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { zodResolver } from "@hookform/resolvers/zod/dist/zod.js";
import { useMutation } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const addLocationAliasSchema = z.object({
  locationId: z.number().min(1, "שדה חובה"),
  alias: z.string().min(1, "שדה חובה"),
});
type AddLocationAliasFormData = z.infer<typeof addLocationAliasSchema>;

export function AddLocationAliasButton() {
  const [isOpen, setIsOpen] = useState(false);

  const {
    handleSubmit,
    register,
    formState: { errors },
    reset,
  } = useForm<AddLocationAliasFormData>({
    resolver: zodResolver(addLocationAliasSchema),
  });

  const onSubmit = (data: AddLocationAliasFormData) => {
    mutate(data);
  };

  const { mutate, isPending, error, isSuccess } = useMutation({
    mutationFn: (data: AddLocationAliasFormData) =>
      insertLocationAlias({
        aliases: [{ value: data.alias }],
        locationId: data.locationId,
      }),
    onSuccess: () => {
      reset();
      setTimeout(() => {
        setIsOpen(false);
      }, 1000);
    },
  });

  const tError = useTranslatedError(error);

  return (
    <>
      <Button
        onClick={() => setIsOpen(true)}
        className="gap-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
      >
        <Plus className="w-4 h-4" />
        הוספת שם למקום
      </Button>
      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent className="max-w-md rounded-lg">
          <div className="pt-2">
            <DialogTitle className="text-center text-lg font-semibold text-gray-900">
              הוספת שם חלופי למקום
            </DialogTitle>
            <p className="text-center text-sm text-gray-500 mt-1">
              הוסף שם חלופי לכל מקום בטווח הזמנות
            </p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5 mt-6">
            <div className="space-y-2">
              <label
                htmlFor="locationId"
                className="block text-sm font-medium text-gray-700"
              >
                מזהה המקום
              </label>
              <Input
                id="locationId"
                {...register("locationId", { valueAsNumber: true })}
                type="number"
                placeholder="לדוגמה: 12345"
                className="w-full"
              />
              {errors.locationId && (
                <ErrorDisplay>{errors.locationId.message}</ErrorDisplay>
              )}
            </div>

            <div className="space-y-2">
              <label
                htmlFor="alias"
                className="block text-sm font-medium text-gray-700"
              >
                שם חלופי
              </label>
              <Input
                id="alias"
                {...register("alias")}
                type="text"
                placeholder="לדוגמה: מקום בעיר"
                className="w-full"
              />
              {errors.alias && (
                <ErrorDisplay>{errors.alias.message}</ErrorDisplay>
              )}
            </div>

            <div className="flex gap-3 pt-4 border-t border-gray-200">
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsOpen(false)}
                className="flex-1"
                loading={isPending}
              >
                ביטול
              </Button>
              <Button
                type="submit"
                className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
              >
                הוסף
              </Button>
            </div>
            <div className="block text-center">
              <ErrorDisplay>{tError}</ErrorDisplay>
            </div>
            {isSuccess && (
              <div className="block text-center">
                <SuccessBadge>השם נוסף בהצלחה</SuccessBadge>
              </div>
            )}
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}
