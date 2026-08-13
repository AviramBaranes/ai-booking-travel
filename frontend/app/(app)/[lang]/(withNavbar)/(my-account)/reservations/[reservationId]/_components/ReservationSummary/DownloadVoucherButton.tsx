import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { Download } from "lucide-react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { downloadVoucher } from "@/shared/api/reservations-api";
import { getErrorKey } from "@/shared/hooks/useTranslatedError";

export function DownloadVoucherButton({
  reservationId,
}: {
  reservationId: number;
}) {
  const t = useTranslations("MyAccount.reservation.summary.downloadVoucher");
  const tErrors = useTranslations("ApiErrors");

  const { mutate, isPending } = useMutation({
    mutationFn: () => downloadVoucher(reservationId),
    onSuccess: ({ blob, filename }) => saveFile(blob, filename),
    onError: (err) => {
      toast.error(tErrors(getErrorKey(err)), {
        duration: 5000,
        position: "top-center",
      });
    },
  });

  return (
    <Button
      variant="ghost"
      className="py-6 px-6 text-border-muted font-semibold flex gap-4 print:hidden"
      onClick={() => mutate()}
      loading={isPending}
    >
      <Download className="w-6 h-6" aria-hidden />
      {t("button")}
    </Button>
  );
}

// saveFile hands the PDF to the browser. The voucher arrives as a blob rather than a linkable URL,
// so an object URL is what lets the download carry the filename the backend chose.
function saveFile(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}
