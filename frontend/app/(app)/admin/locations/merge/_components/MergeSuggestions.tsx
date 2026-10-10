"use client";

import { useQuery } from "@tanstack/react-query";
import { AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { location } from "@/shared/client";
import { listLocationMergeSuggestions } from "@/shared/api/locations-api";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { BrokerCodeBadge } from "./LocationSummary";

interface MergeSuggestionsProps {
  onPick: (a: location.MergeLocation, b: location.MergeLocation) => void;
  onPickA: (l: location.MergeLocation) => void;
  onPickB: (l: location.MergeLocation) => void;
}

export default function MergeSuggestions({
  onPick,
  onPickA,
  onPickB,
}: MergeSuggestionsProps) {
  const { data, isLoading, error } = useQuery({
    queryKey: ["location-merge-suggestions"],
    queryFn: listLocationMergeSuggestions,
  });
  const tError = useTranslatedError(error);
  const groups = data?.groups ?? [];

  return (
    <section className="bg-white rounded-lg border border-gray-200 p-4 space-y-4">
      <div>
        <h2 className="text-lg font-semibold text-gray-700">הצעות למיזוג</h2>
        <p className="text-sm text-gray-500">
          מיקומים באותה מדינה ששמותיהם זהים אחרי הסרת רווחים, סימנים והבדלי
          אותיות גדולות וקטנות.
        </p>
      </div>

      {isLoading && <p className="text-sm text-gray-500">טוען...</p>}
      <ErrorDisplay>{tError}</ErrorDisplay>
      {!isLoading && !error && groups.length === 0 && (
        <p className="text-sm text-gray-500">אין הצעות למיזוג.</p>
      )}

      <div className="grid gap-3 lg:grid-cols-2">
        {groups.map((g) => (
          <div
            key={`${g.country_code}-${g.normalized_name}`}
            className="rounded-md border border-gray-200 p-3 space-y-2"
          >
            <div className="flex items-center justify-between gap-2">
              <span className="text-xs font-medium text-gray-500">
                {g.country_code}
              </span>
              <div className="flex items-center gap-2">
                {g.shared_broker && (
                  <span className="flex items-center gap-1 text-xs text-destructive">
                    <AlertTriangle className="size-3.5" />
                    ספק משותף
                  </span>
                )}
                {g.locations.length === 2 && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => onPick(g.locations[0], g.locations[1])}
                  >
                    השווה
                  </Button>
                )}
              </div>
            </div>
            <ul className="space-y-2">
              {g.locations.map((l) => (
                <li
                  key={l.id}
                  className="flex items-center justify-between gap-2 text-sm"
                >
                  <div className="min-w-0 space-y-1">
                    <div className="truncate">
                      <span className="text-gray-400">#{l.id}</span> {l.name}
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {l.broker_codes.map((c) => (
                        <BrokerCodeBadge key={c.id} code={c} />
                      ))}
                    </div>
                  </div>
                  <div className="flex shrink-0 gap-1">
                    <Button size="sm" variant="ghost" onClick={() => onPickA(l)}>
                      צד א
                    </Button>
                    <Button size="sm" variant="ghost" onClick={() => onPickB(l)}>
                      צד ב
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </section>
  );
}
