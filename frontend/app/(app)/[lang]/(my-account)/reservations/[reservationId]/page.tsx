import { Loading } from "@/shared/components/Loading";
import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { Suspense } from "react";
import { SearchDataBannerWrapper } from "./_components/SearchDataBannerWrapper";
import { SearchDataBannerDisplaySkeleton } from "@/shared/components/booking/SearchDataBannerDisplaySkeleton";

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

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <Suspense
        fallback={
          <SearchDataBannerDisplaySkeleton
            dir={lang === "he" ? "rtl" : "ltr"}
          />
        }
      >
        <SearchDataBannerWrapper reservationId={Number(reservationId)} />
      </Suspense>
    </main>
  );
}
