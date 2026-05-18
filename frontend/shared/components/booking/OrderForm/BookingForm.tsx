import { availability } from "@/shared/client";
import { Page } from "@/payload-types";
import { SelectedCarCard } from "@/shared/components/booking/SelectedCarCard/SelectedCarCard";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { isAppError } from "@/shared/api/AppError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { ChevronDown } from "lucide-react";
import Link from "next/link";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { orderFormSchema, OrderFormValues } from "./orderFormSchema";
import { useTranslations } from "next-intl";
import { Controller, useForm, useFormContext } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { useParams } from "next/navigation";

export function BookingForm() {
  const t = useTranslations("booking.orderPage");
  const {
    control,
    register,
    formState: { errors },
  } = useFormContext<OrderFormValues>();

  return (
    <>
      <div className="flex-1">
        <Controller
          name="driverTitle"
          control={control}
          render={({ field }) => (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  className={`w-full flex items-center justify-between bg-white border rounded-lg px-4 h-12 type-paragraph text-text-secondary cursor-pointer ${errors.driverTitle ? "border-destructive" : "border-cars-border"}`}
                >
                  <span>{field.value ? field.value : t("title")}</span>
                  <ChevronDown className="w-4 h-4 text-muted shrink-0" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align="start"
                className="w-(--radix-dropdown-menu-trigger-width)"
              >
                {["Mr", "Mrs", "Ms", "Miss", "Dr"].map((title) => (
                  <DropdownMenuItem
                    key={title}
                    onClick={() => field.onChange(title)}
                  >
                    {title}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        />
        <ErrorDisplay>{errors.driverTitle?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <Input
          className="bg-white border border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary w-full"
          placeholder={t("firstName")}
          aria-invalid={!!errors.driverFirstName}
          {...register("driverFirstName", {
            onChange: (e) => {
              e.target.value = e.target.value
                .replace(/[^a-zA-Z\s]/g, "")
                .toUpperCase();
            },
          })}
        />
        <ErrorDisplay>{errors.driverFirstName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <Input
          className="bg-white border border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary w-full"
          placeholder={t("lastName")}
          aria-invalid={!!errors.driverLastName}
          {...register("driverLastName", {
            onChange: (e) => {
              e.target.value = e.target.value
                .replace(/[^a-zA-Z\s]/g, "")
                .toUpperCase();
            },
          })}
        />
        <ErrorDisplay>{errors.driverLastName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <Input
          className="bg-white border border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary  w-full"
          placeholder={t("flightNumber")}
          aria-invalid={!!errors.flightNumber}
          {...register("flightNumber")}
        />
        <ErrorDisplay>{errors.flightNumber?.message}</ErrorDisplay>
      </div>
    </>
  );
}
