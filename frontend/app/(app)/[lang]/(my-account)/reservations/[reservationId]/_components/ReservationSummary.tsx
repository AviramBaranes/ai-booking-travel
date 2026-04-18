"use client";

import { Button } from "@/components/ui/button";
import { useReservation } from "../_hooks/useReservation";

import { useTranslations } from "next-intl";
import Image from "next/image";
import { useParams } from "next/navigation";
import { formatPrice } from "@/shared/utils/formatPrice";

export function ReservationSummary({
  reservationId,
}: {
  reservationId: number;
}) {
  const { lang } = useParams();
  const t = useTranslations("MyAccount.reservation.summary");
  const { data: reservation } = useReservation(reservationId);

  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      <div className="flex items-center justify-between">
        <h5 className="type-h5 text-navy">{t("title")}</h5>
        <Button
          variant="outline"
          className="border py-6 px-6 text-border-muted font-semibold flex gap-4"
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
      <hr />
      <OrderSummaryRow
        label={t("labels.driverName")}
        value={`${reservation.driverFirstName} ${reservation.driverLastName}`}
      />
      <OrderSummaryRow
        label={t("labels.reservationNumber")}
        value={reservation.brokerReservationId}
      />
      <OrderSummaryRow
        label={t("labels.status")}
        value={t(`status.${reservation.status}`)}
      />
      <OrderSummaryRow
        label={t("labels.createdAt")}
        value={new Date(reservation.createdAt).toLocaleDateString(lang)}
      />
      <OrderSummarySubTitle title={t("sections.costBreakdown")} />
      <OrderSummaryRow
        label={t("labels.rentalPrice")}
        value={formatPrice(reservation.priceBefDesc, reservation.currencyCode)}
      />
      <OrderSummaryRow
        label={t("labels.couponDiscount")}
        value={formatPrice(
          reservation.discountAmount,
          reservation.currencyCode,
        )}
      />
      {reservation.erpPrice > 0 && (
        <OrderSummaryRow
          label={t("labels.fullCoverage")}
          value={formatPrice(reservation.erpPrice, reservation.currencyCode)}
        />
      )}
      <div className="text-white bg-brand-blue py-3 5 px-5 flex justify-between items-center rounded-xl mt-8">
        <span className="type-paragraph">{t("labels.totalToPay")}</span>
        <h4 className="type-h4">
          {formatPrice(reservation.totalPrice, reservation.currencyCode)}
        </h4>
      </div>
      <OrderSummarySubTitle title={t("sections.carDetails")} />
      <OrderSummaryRow
        label={t("labels.rentalDays")}
        value={reservation.rentalDays.toString()}
      />
      <OrderSummaryRow
        label={t("labels.carType")}
        value={reservation.carDetails.carType}
      />
      <OrderSummaryRow
        label={t("labels.brand")}
        value={reservation.carDetails.supplierName}
      />
      <OrderSummarySubTitle title={t("sections.included")} />
      <ul>
        {reservation.planInclusions.map((inclusion) => (
          <li
            key={inclusion}
            className="type-paragraph text-navy mx-4 my-2 flex"
          >
            <Image
              src="/assets/icons/V.svg"
              alt={t("includedIconAlt")}
              width={28}
              height={28}
              className="inline w-7 h-7"
            />
            {inclusion}
          </li>
        ))}
      </ul>
    </div>
  );
}

function OrderSummaryRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between">
      <span className="type-paragraph text-text-secondary">{label}</span>
      <span className="type-paragraph text-value">{value}</span>
    </div>
  );
}

function OrderSummarySubTitle({ title }: { title: string }) {
  return (
    <div className="mt-8">
      <h5 className="type-h5 text-navy mb-2">{title}</h5>
      <hr />
    </div>
  );
}
