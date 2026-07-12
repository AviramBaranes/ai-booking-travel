import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { ReservationCarCard } from "./_components/ReservationCarCard";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchAddonsGallery, fetchSuppliersGallery } from "@/shared/server/cms";
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { BackButton } from "@/shared/components/booking/BackButton";
import { ReservationSummary } from "./_components/ReservationSummary/ReservationSummary";
import { addonsGalleryKey } from "@/shared/hooks/useAddonsGallery";

export default async function ReservationDetailsPage({
  params,
}: {
  params: Promise<{ reservationId: string }>;
}) {
  const lang = await getLang();
  const { reservationId } = await params;

  if (!reservationId || isNaN(Number(reservationId))) {
    redirect(`/${lang}/reservations`);
  }

  const queryClient = getQueryClient();
  await Promise.all([
    queryClient.fetchQuery({
      queryKey: suppliersGalleryKey,
      queryFn: fetchSuppliersGallery,
    }),
    queryClient.fetchQuery({
      queryKey: addonsGalleryKey,
      queryFn: fetchAddonsGallery,
    }),
  ]);

  return (
    <main className="w-2/3 mx-auto pt-4 pb-6 print:w-full">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <div className="print:hidden">
          <BackButton
            translationKey="backToReservations"
            href={`/${lang}/reservations`}
          />
        </div>
        <div className="flex gap-2 mt-6 print:flex-col print:gap-6">
          <div className="w-3/4 print:w-full">
            <ReservationSummary reservationId={Number(reservationId)} />
          </div>
          <div className="w-1/4 print:w-full">
            <ReservationCarCard reservationId={Number(reservationId)} />
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
