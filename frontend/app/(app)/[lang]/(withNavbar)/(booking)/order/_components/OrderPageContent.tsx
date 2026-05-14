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
import { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { bookCar } from "@/shared/api/booking-api";
import { useForm, Controller, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ErpCheckbox } from "../../plans/_components/ErpCheckbox";
import Link from "next/link";
import { searchRequestToParams } from "../../results/searchQuery";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import {
  orderFormSchema,
  OrderFormValues,
} from "@/shared/components/booking/OrderForm/orderFormSchema";
import { BookingForm } from "@/shared/components/booking/OrderForm/BookingForm";
import { useMutation } from "@tanstack/react-query";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";

interface OrderPageContentProps {
  searchRequest: booking.SearchAvailabilityRequest;
}

export function OrderPageContent({ searchRequest }: OrderPageContentProps) {
  const t = useTranslations("booking.orderPage");
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

  const schema = orderFormSchema(t);
  const formMethods = useForm<OrderFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      driverTitle: "" as unknown as "Mr",
      driverFirstName: "",
      driverLastName: "",
      flightNumber: "",
      termsAccepted: false as unknown as true,
    },
  });

  const { control, handleSubmit, watch } = formMethods;

  const termsAccepted = watch("termsAccepted");

  if (!vehicle || !data) {
    return <Loading />;
  }

  const selectedPlan = vehicle.plans[selectedPlanIndex];

  const { mutate, isPending, error } = useMutation({
    mutationFn: bookCar,
    onSuccess: ({ reservationId }) => {
      router.push(`/${lang}/reservations/${reservationId}`);
    },
  });

  const translatedError = useTranslatedError(error);

  return (
    <form
      className="flex gap-4 mt-4"
      onSubmit={handleSubmit((formData) =>
        mutate({
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
        }),
      )}
    >
      <div className="w-3/4">
        <h2 className="type-h4 text-navy mb-6">{t("driverDetails")}</h2>
        <FormProvider {...formMethods}>
          <div className="flex gap-4">
            <BookingForm />
          </div>
        </FormProvider>

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
            <p className="mt-4 text-destructive type-paragraph">
              {translatedError}
            </p>
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
              disabled={isPending || !termsAccepted}
              loading={isPending}
              className="w-full py-6 type-paragraph font-bold"
            >
              {t("confirmCta")}
            </Button>
          </>
        </SelectedCarCard>
      </div>
    </form>
  );
}
