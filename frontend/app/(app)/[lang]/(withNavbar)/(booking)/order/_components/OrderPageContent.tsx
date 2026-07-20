"use client";

import { Page } from "@/payload-types";
import { useSelectedVehicle } from "../../plans/_hooks/useSelectedVehicle";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { Loading } from "@/shared/components/Loading";
import { SelectedCarCard } from "@/shared/components/booking/SelectedCarCard/SelectedCarCard";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { useRef, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { bookCar } from "@/shared/api/booking-api";
import { useForm, FormProvider } from "react-hook-form";
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
import { useSearchRequest } from "../../_hooks/useSearchRequest";
import { CustomerForm } from "@/shared/components/booking/OrderForm/CustomerForm";
import { getCustomerPaymentIframe } from "@/shared/api/bill-api";
import useAuthStore from "@/shared/auth/authStore";
import { FixedBottomButton } from "./FixedBottomButton";
import { TermsCheckbox } from "./TermsCheckbox";
import { PaymentDialog } from "./PaymentDialog";

export function OrderPageContent() {
  const t = useTranslations("booking.orderPage");
  const { lang } = useParams();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const status = useAuthStore((s) => s.status);
  const isAgent = user?.role === "agent";
  const { data: bookingSettings } = useBookingSettings();

  const selectedCarCardRef = useRef<HTMLDivElement>(null);

  const { searchRequest } = useSearchRequest();
  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const selectedPlanIndex = useBookingSessionStore((s) => s.selectedPlanIndex);
  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedAddons = useBookingSessionStore((s) => s.selectedAddons);
  const setIsErpSelected = useBookingSessionStore((s) => s.setIsErpSelected);

  const [showErp] = useState(!isErpSelected);
  const [pendingPayment, setPendingPayment] = useState<{
    url: string;
    token: string;
  } | null>(null);

  const schema = orderFormSchema(t, !isAgent);
  const formMethods = useForm<OrderFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      driverTitle: "" as unknown as "Mr",
      driverFirstName: "",
      driverLastName: "",
      customerFirstName: isAgent ? "" : (user?.firstName ?? ""),
      customerLastName: isAgent ? "" : (user?.lastName ?? ""),
      customerPhone: isAgent ? "" : (user?.phoneNumber ?? ""),
      customerEmail: isAgent ? "" : (user?.email ?? ""),
      termsAccepted: false as unknown as true,
    },
  });

  const { handleSubmit, watch } = formMethods;

  const termsAccepted = watch("termsAccepted");

  const {
    mutate: doBooking,
    isPending: isBookingPending,
    error: bookingError,
  } = useMutation({
    mutationFn: bookCar,
    onSuccess: ({ reservationId }) => {
      router.push(`/${lang}/reservations/${reservationId}`);
    },
  });

  const {
    mutate: getPaymentIframe,
    isPending: isPaymentIframePending,
    error: paymentIframeError,
  } = useMutation({
    mutationFn: getCustomerPaymentIframe,
    onSuccess: ({ url, pendingPaymentToken }) => {
      setPendingPayment({ url, token: pendingPaymentToken });
    },
  });

  const translatedError = useTranslatedError(
    bookingError || paymentIframeError,
  );

  function submitHandler(formData: OrderFormValues) {
    const basePayload = {
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
    };

    if (isAgent) return doBooking(basePayload);

    getPaymentIframe({
      ...basePayload,
      phone: formData.customerPhone ?? "",
      email: formData.customerEmail ?? "",
      firstName: formData.customerFirstName ?? "",
      lastName: formData.customerLastName ?? "",
    });
  }

  if (!vehicle || !data) {
    return <Loading />;
  }

  const selectedPlan = vehicle.plans[selectedPlanIndex];

  const btnDisabled =
    isBookingPending || isPaymentIframePending || !termsAccepted;
  const btnLoading = isBookingPending || isPaymentIframePending;
  const termsHref =
    typeof bookingSettings?.orderTermsLink === "object"
      ? `/${lang}/${(bookingSettings.orderTermsLink as Page).slug}`
      : "#";

  return (
    <>
      <FormProvider {...formMethods}>
        <form
          className="flex flex-col-reverse lg:flex-row gap-4 mt-4 max-sm:mx-5"
          onSubmit={handleSubmit(submitHandler)}
        >
          <div className="lg:w-3/4">
            <>
              {!isAgent && (
                <>
                  <h5 className="type-h5 text-navy mb-6">
                    {t("customerDetails")}
                  </h5>
                  <div className="flex flex-col lg:flex-row gap-4 my-4">
                    <CustomerForm isReadOnly={status === "authenticated"} />
                  </div>
                </>
              )}
              <h5 className="type-h5 text-navy mb-6">{t("driverDetails")}</h5>
              <div className="flex flex-col lg:flex-row gap-4">
                <BookingForm />
                <div className="lg:hidden">
                  <TermsCheckbox href={termsHref} />
                </div>
              </div>
            </>

            {showErp && (
              <ErpCheckbox
                isSelected={isErpSelected}
                setSelected={setIsErpSelected}
                vehicle={vehicle}
                selectedPlan={selectedPlanIndex}
                daysCount={data.daysCount}
              />
            )}

            {translatedError && (
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

          <div className="lg:w-1/4 relative">
            <SelectedCarCard
              isErpSelected={isErpSelected}
              daysCount={data.daysCount}
              vehicle={vehicle}
              selectedPlanIndex={selectedPlanIndex}
              headerClassName="max-sm:mb-0!"
            >
              <>
                <div className="w-fit mx-auto mb-2 text-center">
                  <FreeCancellationBadge
                    pickupDate={searchRequest.PickupDate}
                    pickupTime={searchRequest.PickupTime}
                    text={t("freeCancellation")}
                  />
                </div>
                <div className="hidden lg:contents">
                  <TermsCheckbox href={termsHref} />
                  <Button
                    type="submit"
                    variant="brand"
                    disabled={btnDisabled}
                    loading={btnLoading}
                    className="w-full py-6 type-paragraph font-bold"
                  >
                    {t("confirmCta")}
                  </Button>
                </div>
                <div className="lg:hidden absolute bottom-0 w-full right-0 left-0">
                  <Button
                    type="submit"
                    variant="brand"
                    disabled={btnDisabled}
                    loading={btnLoading}
                    className="type-paragraph font-bold w-full px-8 cursor-pointer rounded-t-none border border-brand"
                  >
                    {t("confirmCta")}
                  </Button>
                </div>
              </>
            </SelectedCarCard>
            <div className="absolute bottom-20" ref={selectedCarCardRef} />
          </div>
          {!pendingPayment && (
            <FixedBottomButton
              isDisabled={false}
              loading={btnLoading}
              watchRef={selectedCarCardRef}
            />
          )}
        </form>
      </FormProvider>
      <PaymentDialog
        pendingPayment={pendingPayment}
        setPendingPayment={setPendingPayment}
      />
    </>
  );
}
