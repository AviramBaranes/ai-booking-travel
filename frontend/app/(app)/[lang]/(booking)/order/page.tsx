import { getLang } from "@/shared/lang/lang";
import { parseSearchQuery, toSearchRequest } from "../results/searchQuery";
import { redirect } from "next/navigation";
import { BookingStepper } from "../_components/BookingStepper";
import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BackButton } from "../../../../../shared/components/booking/BackButton";
import { OrderPageContent } from "./_components/OrderPageContent";
import { ExpiredSearchGate } from "../_components/ExpiredSearchGate";

export default async function OrderPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string>>;
}) {
  const lang = await getLang();
  const resolvedParams = await searchParams;
  const query = parseSearchQuery(new URLSearchParams(resolvedParams));

  if (!query) {
    redirect(`/${lang}`);
  }

  const searchRequest = toSearchRequest(query);

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <BookingStepper currentStep="ordering" />
      <div className="my-4">
        <SearchDataBanner
          pickUpLocationId={query.pickupLocationId}
          dropOffLocationId={query.returnLocationId}
          pickUpTime={query.pickupTime}
          dropOffTime={query.returnTime}
          pickUpDate={query.pickupDate}
          dropOffDate={query.returnDate}
          driverAge={query.driverAge}
          couponCode={query.couponCode}
          searchRequest={searchRequest}
          showButton
          fromCache
        />
      </div>
      <BackButton translationKey="backToPlans" />
      <ExpiredSearchGate searchRequest={searchRequest}>
        <OrderPageContent searchRequest={searchRequest} />
      </ExpiredSearchGate>
    </main>
  );
}
