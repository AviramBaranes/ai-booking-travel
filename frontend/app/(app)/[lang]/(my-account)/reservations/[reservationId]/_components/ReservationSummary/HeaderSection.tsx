import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { useParams } from "next/navigation";
import { reservation } from "@/shared/client";
import { OrderSummaryRow } from "./OrderSummaryRow";
import { DeleteReservationButton } from "./DeleteReservationButton";

export function statusToColor(status: string) {
  switch (status) {
    case "vouchered":
    case "paid":
      return "text-success font-semibold";
    case "canceled":
      return "text-destructive font-semibold";
    case "booked":
      return "text-brand font-semibold";
    default:
      return "text-navy font-semibold";
  }
}

export function statusToBg(status: string) {
  switch (status) {
    case "vouchered":
    case "paid":
      return "bg-success/10";
    case "canceled":
      return "bg-destructive/10";
    case "booked":
      return "bg-brand/10";
    default:
      return "bg-navy/10";
  }
}

export function HeaderSection({
  reservation: res,
}: {
  reservation: reservation.GetReservationResponse;
}) {
  const { lang } = useParams();
  const t = useTranslations("MyAccount.reservation.summary");

  return (
    <>
      <div className="flex items-center justify-between">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <div className="flex gap-1 items-center w-1/4 justify-end">
          {res.status !== "canceled" && (
            <>
              <DeleteReservationButton reservationId={res.id} />
              <div className="border-l border-cars-border h-5"></div>
            </>
          )}
          <Button
            variant="ghost"
            className="py-6 px-6 text-border-muted font-semibold flex gap-4 print:hidden"
            onClick={() => window.print()}
          >
            <Image
              src="/assets/icons/printer.svg"
              alt={t("print")}
              width={24}
              height={24}
              className="w-6 h-6"
            />
            {t("print")}
          </Button>
        </div>
      </div>
      <hr />
      <OrderSummaryRow
        label={t("labels.driverName")}
        value={`${res.driverFirstName} ${res.driverLastName}`}
      />
      <OrderSummaryRow
        label={t("labels.reservationNumber")}
        value={res.brokerReservationId}
      />
      <OrderSummaryRow
        label={t("labels.status")}
        value={t(`status.${res.status}`)}
        valClassName={statusToColor(res.status)}
      />
      <OrderSummaryRow
        label={t("labels.createdAt")}
        value={new Date(res.createdAt).toLocaleDateString(lang)}
      />
    </>
  );
}
