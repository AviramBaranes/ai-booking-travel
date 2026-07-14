import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { SummarySkeleton } from "@/shared/components/booking/SummarySkeleton";

export default function Loading() {
  return (
    <main className="lg:w-2/3 mx-5 lg:mx-auto lg:pt-15 pb-6">
      <div className="flex flex-col-reverse lg:flex-row gap-2 mt-6">
        <div className="lg:w-3/4">
          <SummarySkeleton />
        </div>
        <div className="lg:w-1/4">
          <SelectedCarCardSkeleton />
        </div>
      </div>
    </main>
  );
}
