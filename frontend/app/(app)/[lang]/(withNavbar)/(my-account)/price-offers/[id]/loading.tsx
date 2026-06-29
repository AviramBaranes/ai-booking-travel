import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { SummarySkeleton } from "@/shared/components/booking/SummarySkeleton";

export default function Loading() {
  return (
    <main className="w-2/3 mx-auto pt-4 pb-6">
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
