"use client";

import { booking } from "@/shared/client";
import { Page } from "@/payload-types";
import { useSelectedVehicle } from "../../plans/_hooks/useSelectedVehicle";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { Loading } from "@/shared/components/Loading";
import { SelectedCarCard } from "@/shared/components/booking/SelectedCarCard/SelectedCarCard";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { bookCar } from "@/shared/api/booking-api";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { orderFormSchema, type OrderFormValues } from "./orderFormSchema";
import { isAppError } from "@/shared/api/AppError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { ErpCheckbox } from "../../plans/_components/ErpCheckbox";
import Image from "next/image";
import { ChevronDown } from "lucide-react";
import Link from "next/link";
import { searchRequestToParams } from "../../results/searchQuery";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";

interface OrderPageContentProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function OrderPageContent({ searchRequest }: OrderPageContentProps) {
  const t = useTranslations("booking.orderPage");
  const tError = useTranslations("ApiErrors");
  const { lang } = useParams();
  const router = useRouter();
  const { data: bookingSettings } = useBookingSettings();

  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const selectedPlanIndex = useBookingSessionStore((s) => s.selectedPlanIndex);
  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedAddons = useBookingSessionStore((s) => s.selectedAddons);
  const setIsErpSelected = useBookingSessionStore((s) => s.setIsErpSelected);

  const [showErp] = useState(!isErpSelected);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const schema = orderFormSchema(t);
  const {
    control,
    handleSubmit,
    watch,
    register,
    formState: { errors },
  } = useForm<OrderFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      driverTitle: "" as unknown as "Mr",
      driverFirstName: "",
      driverLastName: "",
      flightNumber: "",
      termsAccepted: false as unknown as true,
    },
  });

  const termsAccepted = watch("termsAccepted");

  if (!vehicle || !data) {
    return <Loading />;
  }

  const selectedPlan = vehicle.plans[selectedPlanIndex];

  async function onSubmit(formData: OrderFormValues) {
    setError(null);
    setIsSubmitting(true);
    try {
      const result = await bookCar({
        snapshotId: data!.snapshotId,
        rateQualifier: selectedPlan.rateQualifier,
        supplierCode: selectedPlan.supplierCode,
        planId: String(selectedPlan.planId),
        includeERP: isErpSelected,
        selectedAddOns: selectedAddons,
        driverTitle: formData.driverTitle,
        driverFirstName: formData.driverFirstName,
        driverLastName: formData.driverLastName,
        flightNumber: formData.flightNumber,
      });
      router.push(`/${lang}/reservations/${result.reservationId}`);
    } catch (err) {
      setError(isAppError(err) ? tError(err.code) : t("submitError"));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="flex gap-4 mt-4" onSubmit={handleSubmit(onSubmit)}>
      <div className="w-3/4">
        <h2 className="type-h4 text-navy mb-6">{t("driverDetails")}</h2>
        <div className="flex gap-4">
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
                      <span>
                        {field.value
                          ? t(`title${field.value as "Mr" | "Ms"}`)
                          : t("title")}
                      </span>
                      <ChevronDown className="w-4 h-4 text-muted shrink-0" />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    align="start"
                    className="w-(--radix-dropdown-menu-trigger-width)"
                  >
                    <DropdownMenuItem onClick={() => field.onChange("Mr")}>
                      {t("titleMr")}
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => field.onChange("Ms")}>
                      {t("titleMs")}
                    </DropdownMenuItem>
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
        </div>

        {showErp && (
          <ErpCheckbox
            isSelected={isErpSelected}
            setSelected={setIsErpSelected}
            vehicle={vehicle}
            selectedPlan={selectedPlanIndex}
            daysCount={data.daysCount}
          />
        )}

        {error && (
          <>
            <p className="mt-4 text-destructive type-paragraph">{error}</p>
            <Link
              href={`/${lang}/results?${searchRequestToParams(searchRequest).toString()}`}
              className="text-link underline"
            >
              {t("reSearch")}
            </Link>
          </>
        )}
      </div>

      <div className="w-1/4">
        <SelectedCarCard
          isErpSelected={isErpSelected}
          daysCount={data.daysCount}
          vehicle={vehicle}
          selectedPlanIndex={selectedPlanIndex}
        >
          <>
            <FreeCancellationBadge
              pickupDate={searchRequest.PickupDate}
              pickupTime={searchRequest.PickupTime}
              text={t("freeCancellation")}
            />
            <Controller
              name="termsAccepted"
              control={control}
              render={({ field }) => (
                <label className="flex items-center gap-2 cursor-pointer text-navy mx-auto">
                  <Checkbox
                    checked={field.value}
                    onCheckedChange={field.onChange}
                    className="border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand"
                  />
                  <span className="type-paragraph text-navy">
                    {t("termsCheckbox")}{" "}
                    <Link
                      target="_blank"
                      href={
                        typeof bookingSettings.orderTermsLink === "object"
                          ? `/${lang}/${(bookingSettings.orderTermsLink as Page).slug}`
                          : "#"
                      }
                      className="text-link underline type-label"
                    >
                      {t("termsLink")}
                    </Link>
                  </span>
                </label>
              )}
            />
            <Button
              type="submit"
              variant="brand"
              disabled={isSubmitting || !termsAccepted}
              className="w-full py-6 type-paragraph font-bold"
            >
              {isSubmitting ? t("submitting") : t("confirmCta")}
            </Button>
          </>
        </SelectedCarCard>
      </div>
    </form>
  );
}
