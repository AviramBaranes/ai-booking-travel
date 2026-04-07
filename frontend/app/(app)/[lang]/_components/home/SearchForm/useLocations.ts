import { searchLocations } from "@/shared/api/locations-api";
import { useQuery } from "@tanstack/react-query";

export function useLocations(search: string) {
  const { data, isLoading, error } = useQuery({
    queryKey: ["locations", search],
    enabled: search.length >= 3,
    queryFn: async () => searchLocations(search),
  });

  return { locations: data?.locations ?? [], isLoading, error };
}
