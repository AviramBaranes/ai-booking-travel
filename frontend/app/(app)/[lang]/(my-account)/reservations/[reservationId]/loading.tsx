import { BackButton } from "@/shared/components/booking/BackButton";
import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { getLang } from "@/shared/lang/lang";
import { SummarySkeleton } from "../../_components/SummarySkeleton";

export default async function ReservationLoading() {
  const lang = await getLang();
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <BackButton
        translationKey="backToReservations"
        href={`/${lang}/reservations`}
      />
      <div className="flex gap-2 mt-6">
        <div className="w-3/4">
          <SummarySkeleton />
        </div>
        <div className="w-1/4">
          <SelectedCarCardSkeleton />
        </div>
      </div>
    </main>
  );
}
