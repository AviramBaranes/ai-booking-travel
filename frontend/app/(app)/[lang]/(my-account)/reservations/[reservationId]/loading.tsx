import { BackButton } from "@/shared/components/booking/BackButton";
import { SearchDataBannerDisplaySkeleton } from "@/shared/components/booking/SearchDataBannerDisplaySkeleton";
import { ReservationSummarySkeleton } from "./_components/ReservationSummarySkeleton";
import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { getLang } from "@/shared/lang/lang";

export default async function ReservationLoading() {
  const lang = await getLang();
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <SearchDataBannerDisplaySkeleton dir={lang === "he" ? "rtl" : "ltr"} />
      <BackButton
        translationKey="backToReservations"
        href={`/${lang}/reservations`}
      />
      <div className="flex gap-2 mt-6">
        <div className="w-3/4">
          <ReservationSummarySkeleton />
        </div>
        <div className="w-1/4">
          <SelectedCarCardSkeleton />
        </div>
      </div>
    </main>
  );
}
