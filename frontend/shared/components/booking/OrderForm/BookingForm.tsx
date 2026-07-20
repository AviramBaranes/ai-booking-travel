import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { ChevronDown } from "lucide-react";
import { OrderFormValues } from "./orderFormSchema";
import { useTranslations } from "next-intl";
import { Controller, useFormContext } from "react-hook-form";
import { OrderFormInput } from "./OrderFormInput";

export function BookingForm() {
  const t = useTranslations("booking.orderPage");
  const {
    control,
    register,
    formState: { errors },
  } = useFormContext<OrderFormValues>();

  const normalizeName = (value: unknown) =>
    String(value ?? "")
      .replace(/[^a-zA-Z\s]/g, "")
      .toUpperCase();

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
        <OrderFormInput
          placeholder={t("firstName")}
          aria-invalid={!!errors.driverFirstName}
          {...register("driverFirstName", {
            setValueAs: normalizeName,
          })}
        />
        <ErrorDisplay>{errors.driverFirstName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("lastName")}
          aria-invalid={!!errors.driverLastName}
          {...register("driverLastName", {
            setValueAs: normalizeName,
          })}
        />
        <ErrorDisplay>{errors.driverLastName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("flightNumber")}
          aria-invalid={!!errors.flightNumber}
          {...register("flightNumber")}
        />
        <ErrorDisplay>{errors.flightNumber?.message}</ErrorDisplay>
      </div>
    </>
  );
}
