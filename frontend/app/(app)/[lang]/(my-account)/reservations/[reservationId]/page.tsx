import { Loading } from "@/shared/components/Loading";
import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { Suspense } from "react";
import { SearchDataBannerWrapper } from "./_components/SearchDataBannerWrapper";
import { SearchDataBannerDisplaySkeleton } from "@/shared/components/booking/SearchDataBannerDisplaySkeleton";
import { ReservationCarCard } from "./_components/ReservationCarCard";
import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchSuppliersGallery } from "@/shared/server/cms";
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { BackButton } from "../../../../../../shared/components/booking/BackButton";
import { ReservationSummary } from "./_components/ReservationSummary/ReservationSummary";
import { ReservationSummarySkeleton } from "./_components/ReservationSummarySkeleton";

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
  await queryClient.fetchQuery({
    queryKey: suppliersGalleryKey,
    queryFn: fetchSuppliersGallery,
  });

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
            <Suspense fallback={<ReservationSummarySkeleton />}>
              <ReservationSummary reservationId={Number(reservationId)} />
            </Suspense>
          </div>
          <div className="w-1/4 print:w-full">
            <Suspense fallback={<SelectedCarCardSkeleton />}>
              <ReservationCarCard reservationId={Number(reservationId)} />
            </Suspense>
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
