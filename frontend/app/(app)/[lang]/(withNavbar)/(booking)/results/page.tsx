import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/navigation";
import { parseSearchQuery } from "./searchQuery";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import {
  fetchSuppliersGallery,
  fetchAddonsGallery,
  fetchBookingSettings,
} from "@/shared/server/cms";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { addonsGalleryKey } from "@/shared/hooks/useAddonsGallery";
import { bookingSettingsKey } from "@/shared/hooks/useBookingSettings";
import ErrorResultPageContent from "./_components/ErrorPage";
import { ResultsPageContent } from "./_components/ResultsPageContent";

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
  const queryClient = getQueryClient();

  try {
    await Promise.all([
      queryClient.fetchQuery({
        queryKey: suppliersGalleryKey,
        queryFn: fetchSuppliersGallery,
      }),
      queryClient.fetchQuery({
        queryKey: addonsGalleryKey,
        queryFn: fetchAddonsGallery,
      }),
      queryClient.fetchQuery({
        queryKey: bookingSettingsKey,
        queryFn: fetchBookingSettings,
      }),
    ]);
  } catch {
    return <ErrorResultPageContent />;
  }

  return (
    <main>
      <HydrationBoundary state={dehydrate(queryClient)}>
        <ResultsPageContent query={query} />
      </HydrationBoundary>
    </main>
  );
}
