import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useBookingSettings } from "@/shared/hooks/useBookingSettings";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";
import Image from "next/image";

interface ErpDialogProps {
  open: boolean;
  onApprove: () => void;
  onDecline: () => void;
  erpPrice: number;
  erpPriceCurrency: string;
}

export function ErpDialog({
  open,
  erpPrice,
  onApprove,
  onDecline,
  erpPriceCurrency,
}: ErpDialogProps) {
  const t = useTranslations("booking.plansPage");
  const { data: bookingSettings } = useBookingSettings();

  return (
    <Dialog open={open}>
      <DialogContent
        className="min-w-1/3 max-w-md py-6 px-10 flex flex-col gap-4 bg-white border-border-light/50 rounded-2xl shadow-modal"
        showCloseButton={false}
        onEscapeKeyDown={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
      >
        <DialogTitle className="flex items-center justify-between p-6 pb-15 border-b border-muted/50">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/icons/stamp.gif"
              alt="stamp"
              width={24}
              height={24}
              className="w-6 h-6"
            />
            <h5 className="type-h5 text-navy">
              {bookingSettings.erpPopupTitle}
            </h5>
          </div>
        </DialogTitle>
        <p className="type-h5 font-normal text-text-secondary">
          {bookingSettings.erpPopupContent.split("\n").map((line, index) => (
            <span key={index}>
              {line}
              <br />
            </span>
          ))}
        </p>
        <h3 className="type-h3 font-normal my-3">
          {t("erpTotal")} {formatPrice(erpPrice, erpPriceCurrency)}
        </h3>
        <div className="flex justify-center gap-6">
          <Button
            variant="outline"
            className="text-destructive border border-destructive rounded-md w-1/2 py-4"
            onClick={onDecline}
          >
            {bookingSettings.erpPopupDeclineButtonText}
          </Button>
          <Button
            variant="brand"
            className="font-semibold rounded-md w-1/2 py-4"
            onClick={onApprove}
          >
            {bookingSettings.erpPopupApproveButtonText}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
