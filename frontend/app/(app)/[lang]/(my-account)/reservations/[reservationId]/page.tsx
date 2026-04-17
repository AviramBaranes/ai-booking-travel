import { Loading } from "@/shared/components/Loading";
import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { Suspense } from "react";
import { SearchDataBannerWrapper } from "./_components/SearchDataBannerWrapper";
import { SearchDataBannerDisplaySkeleton } from "@/shared/components/booking/SearchDataBannerDisplaySkeleton";
import { ReservationCarCard } from "./_components/ReservationCarCard";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchSuppliersGallery } from "@/shared/server/cms";
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";

export default async function ReservationDetailsPage({
  params,
}: {
  params: Promise<{ reservationId: string }>;
}) {
  const lang = await getLang();
  const { reservationId } = await params;
  console.log("reservationId", reservationId);

  if (!reservationId || isNaN(Number(reservationId))) {
    redirect(`/${lang}/reservations`);
  }

  const queryClient = getQueryClient();
  await queryClient.fetchQuery({
    queryKey: suppliersGalleryKey,
    queryFn: fetchSuppliersGallery,
  });

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <Suspense
          fallback={
            <SearchDataBannerDisplaySkeleton
              dir={lang === "he" ? "rtl" : "ltr"}
            />
          }
        >
          <SearchDataBannerWrapper reservationId={Number(reservationId)} />
        </Suspense>
        <div className="flex gap-2 mt-6">
          <div className="w-3/4"></div>
          <div className="w-1/4">
            <Suspense fallback={<Loading />}>
              <ReservationCarCard reservationId={Number(reservationId)} />
            </Suspense>
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
