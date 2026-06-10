import { redirect, useParams, useSearchParams } from "next/navigation";
import { parseSearchQuery, toSearchRequest } from "../results/searchQuery";

export function useSearchRequest() {
  const { lang } = useParams();
  const searchParams = useSearchParams();
  const query = parseSearchQuery(new URLSearchParams(searchParams));
  if (!query) {
    redirect(`/${lang}`);
  }

  const searchRequest = toSearchRequest(query);

  return {searchRequest, query};
}
