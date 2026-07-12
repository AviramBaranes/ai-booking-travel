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
import { useMutation, useQuery } from "@tanstack/react-query";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { useSearchRequest } from "../../_hooks/useSearchRequest";
import { CustomerForm } from "@/shared/components/booking/OrderForm/CustomerForm";
import {
  getCustomerPaymentIframe,
  getCustomerPaymentStatus,
} from "@/shared/api/bill-api";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { usePaymentSuccess } from "@/shared/hooks/usePaymentSuccess";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";

export function OrderPageContent() {
  const t = useTranslations("booking.orderPage");
  const { lang } = useParams();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const setSession = useAuthStore((s) => s.setSession);
  const isAgent = user?.role === "agent";
  const { data: bookingSettings } = useBookingSettings();

  const { searchRequest } = useSearchRequest();
  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const selectedPlanIndex = useBookingSessionStore((s) => s.selectedPlanIndex);
  const isErpSelected = useBookingSessionStore((s) => s.isErpSelected);
  const selectedAddons = useBookingSessionStore((s) => s.selectedAddons);
  const setIsErpSelected = useBookingSessionStore((s) => s.setIsErpSelected);

  const [showErp] = useState(!isErpSelected);
  const [paymentLoading, setPaymentLoading] = useState(false);
  const [paymentError, setPaymentError] = useState(false);
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

  const { control, handleSubmit, watch } = formMethods;

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

  const { refetch, data: paymentStatus } = useQuery({
    queryKey: ["bookingSettings", pendingPayment?.token],
    enabled: !!pendingPayment?.token,
    queryFn: () => getCustomerPaymentStatus(pendingPayment?.token ?? ""),
  });

  usePaymentSuccess({
    onSuccess: () => {
      if (!user) {
        if (!paymentStatus?.login) {
          console.warn("No login data found in payment status");
          return;
        }

        setSession(
          paymentStatus.login.accessToken,
          paymentStatus.login.accessTokenExpiresAt,
          {
            id: paymentStatus.login.id,
            email: paymentStatus.login.email,
            firstName: paymentStatus.login.firstName,
            lastName: paymentStatus.login.lastName,
            role: paymentStatus.login.role as UserRole,
            phoneNumber: paymentStatus.login.phoneNumber,
            officeId: paymentStatus.login.officeId,
            isAdminAsAgent: false,
          },
        );
      }
      router.push(`/${lang}/reservations/${paymentStatus?.reservationId}`);
    },
    onFailure: async () => {
      if (!user) {
        if (!paymentStatus?.login) {
          console.warn("No login data found in payment status");
          return;
        }

        setSession(
          paymentStatus.login.accessToken,
          paymentStatus.login.accessTokenExpiresAt,
          {
            id: paymentStatus.login.id,
            email: paymentStatus.login.email,
            firstName: paymentStatus.login.firstName,
            lastName: paymentStatus.login.lastName,
            role: paymentStatus.login.role as UserRole,
            phoneNumber: paymentStatus.login.phoneNumber,
            officeId: paymentStatus.login.officeId,
            isAdminAsAgent: false,
          },
        );
      }
      setPaymentLoading(false);
      setPaymentError(true);
    },
    onLoadingStart: () => {
      setPaymentLoading(true);
    },
    refetch,
    isSuccess: paymentStatus?.paymentStatus === "completed",
    isFailure: paymentStatus?.paymentStatus === "failed",
  });

  if (!vehicle || !data) {
    return <Loading />;
  }

  const selectedPlan = vehicle.plans[selectedPlanIndex];

  return (
    <>
      <form className="flex gap-4 mt-4" onSubmit={handleSubmit(submitHandler)}>
        <div className="w-3/4">
          <FormProvider {...formMethods}>
            {!isAgent && (
              <>
                <h5 className="type-h5 text-navy mb-6">
                  {t("customerDetails")}
                </h5>
                <div className="flex gap-4 my-4">
                  <CustomerForm />
                </div>
              </>
            )}
            <h5 className="type-h5 text-navy mb-6">{t("driverDetails")}</h5>
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
                disabled={
                  isBookingPending || isPaymentIframePending || !termsAccepted
                }
                loading={isBookingPending || isPaymentIframePending}
                className="w-full py-6 type-paragraph font-bold"
              >
                {t("confirmCta")}
              </Button>
            </>
          </SelectedCarCard>
        </div>
      </form>
      <Dialog
        open={!!pendingPayment}
        onOpenChange={() => {
          setPendingPayment(null);
        }}
      >
        <DialogContent className="w-full max-w-3xl!">
          {paymentError ? (
            <>
              <DialogTitle>
                <p className="type-h5 text-navy mt-4">
                  {t("paymentFailedTitle")}
                </p>
              </DialogTitle>

              <div className="space-y-4">
                <p className="type-paragraph text-navy">
                  {t("paymentFailedMessage")}
                </p>

                <div className="flex gap-3">
                  <Button
                    variant="brand"
                    onClick={() => {
                      router.push(`/${lang}/צור-קשר`);
                    }}
                  >
                    {t("contactSupport")}
                  </Button>

                  <Button
                    variant="outline"
                    onClick={() => {
                      router.push(
                        `/${lang}/results?${searchRequestToParams(searchRequest).toString()}`,
                      );
                    }}
                  >
                    {t("backToSearchResults")}
                  </Button>
                </div>
              </div>
            </>
          ) : (
            <>
              <DialogTitle>
                <p className="type-h5 text-navy mt-4">
                  {t(paymentLoading ? "paymentLoading" : "paymentTitle")}
                </p>
                <p className="type-paragraph">
                  {t("paymentSubtitle", {
                    carModel: vehicle.carDetails.model,
                    pickupLocation: data.pickupLocationName,
                    pickupDate: searchRequest.PickupDate,
                    pickupTime: searchRequest.PickupTime,
                  })}
                </p>
              </DialogTitle>

              <iframe
                src={pendingPayment?.url ?? undefined}
                className="w-full h-160"
              />
            </>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}
