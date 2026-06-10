import { BookingStepper } from "../_components/BookingStepper";
import { BackButton } from "@/shared/components/booking/BackButton";
import { PlansPageContent } from "./_components/PlansPageContent";
import { ExpiredSearchGate } from "../_components/ExpiredSearchGate";
import { SearchDataBannerWithQuery } from "../_components/SearchDataBannerWithQuery";

export default function PlansPage() {
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <BookingStepper currentStep="plans" />
      <div className="my-4">
        <SearchDataBannerWithQuery />
      </div>
      <BackButton translationKey="backToResults" />
      <ExpiredSearchGate>
        <PlansPageContent />
      </ExpiredSearchGate>
    </main>
  );
}
