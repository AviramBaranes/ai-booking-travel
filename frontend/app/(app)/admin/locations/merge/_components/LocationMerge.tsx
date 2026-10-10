"use client";

import { useState } from "react";
import { location } from "@/shared/client";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import MergeLocationSearch from "./MergeLocationSearch";
import MergeComparison from "./MergeComparison";
import MergeSuggestions from "./MergeSuggestions";

export default function LocationMerge() {
  const [a, setA] = useState<location.MergeLocation | null>(null);
  const [b, setB] = useState<location.MergeLocation | null>(null);
  const [merged, setMerged] = useState<location.MergeLocation | null>(null);

  const select = (
    setter: (l: location.MergeLocation | null) => void,
    l: location.MergeLocation | null,
  ) => {
    setMerged(null);
    setter(l);
  };

  return (
    <div className="space-y-6">
      <section className="bg-white rounded-lg border border-gray-200 p-4 space-y-4">
        <h2 className="text-lg font-semibold text-gray-700">בחירת מיקומים</h2>
        <div className="grid gap-4 md:grid-cols-2">
          <MergeLocationSearch
            label="מיקום א"
            value={a}
            onChange={(l) => select(setA, l)}
          />
          <MergeLocationSearch
            label="מיקום ב"
            value={b}
            onChange={(l) => select(setB, l)}
          />
        </div>

        {a && b && a.id === b.id && (
          <p className="text-sm text-destructive">בחרו שני מיקומים שונים.</p>
        )}

        {a && b && a.id !== b.id && (
          <MergeComparison
            key={`${a.id}-${b.id}`}
            a={a}
            b={b}
            onMerged={(l) => {
              setA(null);
              setB(null);
              setMerged(l);
            }}
          />
        )}

        {merged && (
          <SuccessBadge>
            המיקומים מוזגו למיקום #{merged.id} ({merged.name}).
          </SuccessBadge>
        )}
      </section>

      <MergeSuggestions
        onPick={(first, second) => {
          setMerged(null);
          setA(first);
          setB(second);
        }}
        onPickA={(l) => select(setA, l)}
        onPickB={(l) => select(setB, l)}
      />
    </div>
  );
}
