import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BookingStepper } from "../_components/BookingStepper";
import { getMessages } from "next-intl/server";
import { getLang } from "@/shared/lang/lang";
import { NextIntlClientProvider } from "next-intl";
import { redirect } from "next/navigation";
import { parseSearchQuery } from "./searchQuery";

export default async function ResultsPage({
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

  return (
    <main className="w-2/3 mx-auto pt-15 pb-300">
      <BookingStepper currentStep="results" />
      <NextIntlClientProvider locale={lang} messages={messages}>
        <div className="my-4">
          <SearchDataBanner
            pickUpLocation={{
              id: query.pickupLocationId,
              name: "Holland, Amsterdam - Schiphol Airport",
            }}
            dropOffLocation={{
              id: query.returnLocationId,
              name: "Holland, Amsterdam - Schiphol Airport",
            }}
            pickUpTime={query.pickupTime}
            dropOffTime={query.returnTime}
            pickUpDate={query.pickupDate}
            dropOffDate={query.returnDate}
            driverAge={query.driverAge}
            couponCode={query.couponCode}
            showButton
          />
        </div>
      </NextIntlClientProvider>
    </main>
  );
}
