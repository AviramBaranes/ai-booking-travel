import { BookingStepper } from "../_components/BookingStepper";
import { OrderPageContent } from "./_components/OrderPageContent";
import { ExpiredSearchGate } from "../_components/ExpiredSearchGate";
import { BackButton } from "@/shared/components/booking/BackButton";
import { SearchDataBannerWithQuery } from "../_components/SearchDataBannerWithQuery";

export default function OrderPage() {
  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <BookingStepper currentStep="ordering" />
      <div className="my-4">
        <SearchDataBannerWithQuery />
      </div>
      <BackButton translationKey="backToPlans" />
      <ExpiredSearchGate>
        <OrderPageContent />
      </ExpiredSearchGate>
    </main>
  );
}
