import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useTranslations } from "next-intl";
import { VoucherForm } from "../[reservationId]/_components/VoucherForm";
import { X } from "lucide-react";

interface VoucherReservationDialogProps {
  open: boolean;
  setOpen: (open: boolean) => void;
  reservationId: number;
  bookingId: string;
  refetch?: () => void;
}
export function VoucherReservationDialog({
  open,
  setOpen,
  reservationId,
  bookingId,
  refetch,
}: VoucherReservationDialogProps) {
  const t = useTranslations("MyAccount.reservations");

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent
        className="min-w-1/4 max-w-md p-12 flex flex-col gap-4 bg-white border-border-light/50 rounded-2xl shadow-modal"
        showCloseButton={false}
      >
        <DialogTitle className="type-h5 text-navy flex items-center justify-between p-6 pb-15 border-b border-muted/50">
          {t("voucherDialogTitle", { bookingId })}
          <X
            className="w-6 h-6 cursor-pointer"
            onClick={() => setOpen(false)}
          />
        </DialogTitle>
        <VoucherForm reservationId={reservationId} refetch={refetch} />
      </DialogContent>
    </Dialog>
  );
}
