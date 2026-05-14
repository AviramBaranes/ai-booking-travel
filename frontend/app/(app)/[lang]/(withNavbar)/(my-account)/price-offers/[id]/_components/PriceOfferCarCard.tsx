"use client";

import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { SelectedCarCardWrapper } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardWrapper";
import { SelectedCarHeader } from "@/shared/components/booking/SelectedCarCard/SelectedCarHeader";
import { useTranslations } from "next-intl";
import { usePriceOffer } from "../_hooks/usePriceOffer";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { toast } from "sonner";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { ShieldCheck, X } from "lucide-react";
import { BookingForm } from "@/shared/components/booking/OrderForm/BookingForm";
import {
  orderFormSchema,
  OrderFormValues,
} from "@/shared/components/booking/OrderForm/orderFormSchema";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Checkbox } from "@/components/ui/checkbox";
import Link from "next/link";
import { useMutation } from "@tanstack/react-query";
import { bookPriceOffer } from "@/shared/api/price-offers-api";
import { useParams, useRouter } from "next/navigation";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { Page } from "@/payload-types";

export function PriceOfferCarCard({ priceOfferId }: { priceOfferId: number }) {
  const router = useRouter();
  const { lang } = useParams();

  const { data: priceOffer } = usePriceOffer(priceOfferId);
  const { data: bookingSettings } = useBookingSettings();

  const t = useTranslations("MyAccount.priceOffer");
  const tOrder = useTranslations("booking.orderPage");
  const [showForm, setShowForm] = useState(false);

  function handleOrderClick() {
    const { renewedAt } = priceOffer;
    const needRenew =
      new Date().getTime() - new Date(renewedAt).getTime() > 15 * 60 * 1000;

    if (needRenew) {
      toast.error("צריך לחדש את הצעת המחיר לפני ביצוע הזמנה", {
        duration: 3000,
        position: "top-center",
      });
    } else {
      setShowForm(true);
    }
  }

  const schema = orderFormSchema(tOrder);
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

  const { handleSubmit, control, watch } = formMethods;
  const termsAccepted = watch("termsAccepted");

  const { mutate, isPending, error } = useMutation({
    mutationFn: bookPriceOffer,
    onSuccess: ({ reservationId }) => {
      router.push(`/${lang}/reservations/${reservationId}`);
    },
  });

  const translatedError = useTranslatedError(error);

  return (
    <div className="sticky top-24">
      <SelectedCarCardWrapper>
        <SelectedCarHeader carDetails={priceOffer.carDetails} />
        <FreeCancellationBadge
          pickupDate={priceOffer.pickupDate}
          pickupTime={priceOffer.pickupTime}
          text={t("freeCancellation")}
        />
        {priceOffer.status === "open" && (
          <Button
            className="mt-4 type-paragraph font-bold py-6 px-8"
            variant="brand"
            onClick={handleOrderClick}
          >
            {t("orderCTA")}
          </Button>
        )}
      </SelectedCarCardWrapper>
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent
          className="min-w-1/3 max-w-md py-6 px-10 flex flex-col gap-4 bg-background border-border-light/50 rounded-2xl shadow-modal"
          showCloseButton={false}
        >
          <div className="flex items-center justify-between p-3 pb-0">
            <DialogTitle className="flex items-center gap-4">
              <ShieldCheck className="w-8 h-8 text-success" />
              <span className="type-h5 text-navy">{t("orderDialogTitle")}</span>
            </DialogTitle>
            <button
              onClick={() => setShowForm(false)}
              className="p-2 cursor-pointer"
            >
              <X className="w-6 h-6 text-navy" />
            </button>
          </div>
          <hr />
          <FormProvider {...formMethods}>
            <form
              onSubmit={handleSubmit((data) =>
                mutate({ priceOfferId, ...data }),
              )}
            >
              <div className="grid gap-2 grid-cols-2 my-6">
                <BookingForm />
              </div>
              <div className="flex flex-col items-center gap-3 mx-auto w-1/2">
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
                        {tOrder("termsCheckbox")}{" "}
                        <Link
                          target="_blank"
                          href={
                            typeof bookingSettings.orderTermsLink === "object"
                              ? `/${lang}/${(bookingSettings.orderTermsLink as Page).slug}`
                              : "#"
                          }
                          className="text-link underline type-label"
                        >
                          {tOrder("termsLink")}
                        </Link>
                      </span>
                    </label>
                  )}
                />
                <Button
                  type="submit"
                  variant="brand"
                  loading={isPending}
                  disabled={!termsAccepted || isPending}
                  className="w-full py-6 type-paragraph font-bold"
                >
                  {t("orderCTA")}
                </Button>
                <ErrorDisplay>{translatedError}</ErrorDisplay>
              </div>
            </form>
          </FormProvider>
        </DialogContent>
      </Dialog>
    </div>
  );
}
