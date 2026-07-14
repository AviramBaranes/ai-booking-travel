import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { getCustomerPaymentStatus } from "@/shared/api/bill-api";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";
import { usePaymentSuccess } from "@/shared/hooks/usePaymentSuccess";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { searchRequestToParams } from "../../results/searchQuery";
import { useSearchRequest } from "../../_hooks/useSearchRequest";
import { useSelectedVehicle } from "../../plans/_hooks/useSelectedVehicle";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";
import { Loading } from "@/shared/components/Loading";

interface PaymentDialogProps {
  pendingPayment: {
    url: string;
    token: string;
  } | null;
  setPendingPayment: (value: { url: string; token: string } | null) => void;
}

export function PaymentDialog({
  pendingPayment,
  setPendingPayment,
}: PaymentDialogProps) {
  const t = useTranslations("booking.orderPage");
  const { lang } = useParams();
  const router = useRouter();

  const { searchRequest } = useSearchRequest();
  const vehicle = useSelectedVehicle(searchRequest);
  const { data } = useAvailableCars(searchRequest, { fromCache: true });

  const user = useAuthStore((s) => s.user);
  const setSession = useAuthStore((s) => s.setSession);

  const [paymentLoading, setPaymentLoading] = useState(false);
  const [paymentError, setPaymentError] = useState(false);

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

  return (
    <Dialog
      open={!!pendingPayment}
      onOpenChange={() => {
        setPendingPayment(null);
      }}
    >
      <DialogContent className="w-full max-w-3xl! max-sm:mx-auto max-sm:w-11/12 max-h-[calc(100dvh-2rem)] overflow-y-auto">
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
  );
}
