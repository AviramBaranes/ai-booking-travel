import { getLang } from "@/shared/lang/lang";
import { parseSearchQuery, toSearchRequest } from "../results/searchQuery";
import { redirect } from "next/navigation";
import { getMessages } from "next-intl/server";
import { NextIntlClientProvider } from "next-intl";
import { BookingStepper } from "../_components/BookingStepper";
import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BackButton } from "../_components/BackButton";
import { ExpiredSearchGate } from "../_components/ExpiredSearchGate";
import { OrderPageContent } from "./_components/OrderPageContent";
import { OrderProviders } from "./_components/OrderProviders";

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

  const messages = await getMessages({ locale: lang });

  const searchRequest = toSearchRequest(query);

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NextIntlClientProvider locale={lang} messages={messages}>
        <ExpiredSearchGate searchRequest={searchRequest}>
          <OrderProviders>
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
            <BackButton />
            <OrderPageContent searchRequest={searchRequest} />
          </OrderProviders>
        </ExpiredSearchGate>
      </NextIntlClientProvider>
    </main>
  );
}
