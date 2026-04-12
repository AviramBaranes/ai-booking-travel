import { getPayload } from "payload";
import config from "@payload-config";
import { parseSearchQuery, toSearchRequest } from "../results/searchQuery";
import { redirect } from "next/navigation";
import { getLang } from "@/shared/lang/lang";
import { BookingStepper } from "../_components/BookingStepper";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BackButton } from "../_components/BackButton";

async function getAddOnsGallery() {
  const payload = await getPayload({ config });
  return payload.findGlobal({
    slug: "addonsGallery",
    draft: false,
  });
}

export default async function PlansPage({
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
  const addOnsGallery = await getAddOnsGallery();
  const searchRequest = toSearchRequest(query);

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NextIntlClientProvider locale={lang} messages={messages}>
        <BookingStepper currentStep="plans" />
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
          />
        </div>
        <BackButton />
      </NextIntlClientProvider>
    </main>
  );
}
