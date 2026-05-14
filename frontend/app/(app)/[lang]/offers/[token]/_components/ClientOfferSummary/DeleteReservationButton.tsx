import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { useState } from "react";
import { useReservation } from "../../_hooks/useReservation";
import { DeleteReservationDialog } from "../../../_components/DeleteReservationDialog";

export function DeleteReservationButton({
  reservationId,
}: {
  reservationId: number;
}) {
  const t = useTranslations("MyAccount.reservation.summary.cancel");
  const [open, setOpen] = useState(false);
  const { refetch } = useReservation(reservationId);

  return (
    <>
      <Button
        variant="ghost"
        className="py-6 px-6 text-border-muted font-semibold flex gap-4 print:hidden"
        onClick={() => setOpen(true)}
      >
        <Image
          src="/assets/icons/trash.svg"
          alt={t("button")}
          width={24}
          height={24}
          className="w-6 h-6"
        />
        {t("button")}
      </Button>
      <DeleteReservationDialog
        open={open}
        reservationId={reservationId}
        setOpen={setOpen}
        refetch={refetch}
      />
    </>
  );
}
