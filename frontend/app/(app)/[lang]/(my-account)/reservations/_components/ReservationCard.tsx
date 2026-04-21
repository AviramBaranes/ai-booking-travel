import { reservation } from "@/shared/client";
import { useTranslations } from "next-intl";
import {
  statusToBg,
  statusToColor,
} from "../[reservationId]/_components/ReservationSummary/HeaderSection";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useState } from "react";
import { DeleteReservationDialog } from "./DeleteReservationDialog";
import { VoucherReservationDialog } from "./VoucherReservationDialog";

export function ReservationCard({
  reservation,
  refetchReservations,
}: {
  reservation: reservation.ReservationSummary;
  refetchReservations: () => void;
}) {
  const { lang } = useParams();
  const tLabels = useTranslations("MyAccount.reservations.labels");
  const tButtons = useTranslations("MyAccount.reservations.buttons");
  const tStatuses = useTranslations("MyAccount.reservation.summary.status");

  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [voucherModalOpen, setVoucherModalOpen] = useState(false);

  function refetchAndClose() {
    refetchReservations();
    setDeleteModalOpen(false);
    setVoucherModalOpen(false);
  }

  return (
    <>
      <Link
        href={`/${lang}/reservations/${reservation.id}`}
        className="p-6 flex flex-col gap-4 rounded-xl bg-white shadow-card hover:shadow-card-hover hover:border hover:border-brand"
      >
        <ReservationCardLabelValue
          label={tLabels("pickupDate")}
          value={reservation.pickupDate}
        />
        <ReservationCardLabelValue
          label={tLabels("pickupLocation")}
          value={reservation.pickupLocationName}
        />
        <ReservationCardLabelValue
          label={tLabels("driverName")}
          value={`${reservation.driverTitle} ${reservation.driverFirstName} ${reservation.driverLastName}`}
        />
        <ReservationCardLabelValue
          valClassName="font-semibold"
          label={tLabels("bookingNumber")}
          value={reservation.brokerReservationId}
        />
        <div className="px-6 py-1 flex flex-col">
          <p className="text-xs text-muted">{tLabels("status")}</p>
          <p
            className={`rounded-md py-1 mt-2 px-2 w-fit text-sm ${statusToBg(reservation.status)} ${statusToColor(reservation.status)}`}
          >
            {tStatuses(reservation.status)}
          </p>
        </div>
        {reservation.status !== "canceled" && (
          <div className="flex justify-between px-4">
            <Button
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                setDeleteModalOpen(true);
              }}
              className="type-paragraph text-destructive"
              variant="ghost"
            >
              {tButtons("cancelOrder")}
            </Button>
            {reservation.status === "booked" ? (
              <Button
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  setVoucherModalOpen(true);
                }}
                className="type-paragraph text-navy"
                variant="ghost"
              >
                {tButtons("voucherOrder")}
              </Button>
            ) : (
              <div className="w-1/2" />
            )}
          </div>
        )}
      </Link>
      {reservation.status !== "canceled" && (
        <DeleteReservationDialog
          open={deleteModalOpen}
          setOpen={setDeleteModalOpen}
          reservationId={reservation.id}
          refetch={refetchAndClose}
        />
      )}
      {reservation.status === "booked" && (
        <VoucherReservationDialog
          open={voucherModalOpen}
          setOpen={setVoucherModalOpen}
          refetch={refetchAndClose}
          reservationId={reservation.id}
          bookingId={reservation.brokerReservationId}
        />
      )}
    </>
  );
}

function ReservationCardLabelValue({
  label,
  value,
  valClassName,
}: {
  label: string;
  value: string;
  valClassName?: string;
}) {
  return (
    <div className="px-6 py-1 flex flex-col">
      <p className="text-xs text-muted">{label}</p>
      <p className={`text-sm text-navy ${valClassName}`}>{value}</p>
    </div>
  );
}
