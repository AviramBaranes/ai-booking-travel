"use client";

import { useState } from "react";
import clsx from "clsx";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AlertTriangle, Merge } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
} from "@/components/ui/dialog";
import { location } from "@/shared/client";
import { mergeLocations } from "@/shared/api/locations-api";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { BrokerCodeBadge } from "./LocationSummary";

type Side = "a" | "b";

type Field =
  | "id"
  | "name"
  | "country"
  | "country_code"
  | "city"
  | "iata"
  | "is_airport";

const fields: { key: Field; label: string }[] = [
  { key: "id", label: "מזהה שנשמר" },
  { key: "name", label: "שם" },
  { key: "country", label: "מדינה" },
  { key: "country_code", label: "קוד מדינה" },
  { key: "city", label: "עיר" },
  { key: "iata", label: "IATA" },
  { key: "is_airport", label: "שדה תעופה" },
];

function isEmpty(v: unknown) {
  return v === null || v === undefined || v === "" || v === false;
}

function display(v: unknown) {
  if (v === true) return "כן";
  if (v === false) return "לא";
  if (v === null || v === undefined || v === "") {
    return <span className="text-muted-foreground">—</span>;
  }
  return String(v);
}

// A field defaults to location A, unless A has nothing there and B does.
function defaultPicks(
  a: location.MergeLocation,
  b: location.MergeLocation,
): Record<Field, Side> {
  const picks = {} as Record<Field, Side>;
  for (const { key } of fields) {
    picks[key] = key !== "id" && isEmpty(a[key]) && !isEmpty(b[key]) ? "b" : "a";
  }
  return picks;
}

function uniqueIgnoreCase(values: string[]) {
  const seen = new Set<string>();
  return values.filter((v) => {
    const k = v.trim().toLowerCase();
    if (seen.has(k)) return false;
    seen.add(k);
    return true;
  });
}

interface MergeComparisonProps {
  a: location.MergeLocation;
  b: location.MergeLocation;
  onMerged: (l: location.MergeLocation) => void;
}

export default function MergeComparison({
  a,
  b,
  onMerged,
}: MergeComparisonProps) {
  const queryClient = useQueryClient();
  const [picks, setPicks] = useState(() => defaultPicks(a, b));
  const [confirmOpen, setConfirmOpen] = useState(false);

  const sides = { a, b };
  const valueOf = (key: Field) => sides[picks[key]][key];

  const kept = sides[picks.id];
  const removed = picks.id === "a" ? b : a;
  const name = String(valueOf("name"));

  const brokersOfA = new Set(a.broker_codes.map((c) => c.broker));
  const sharedBrokers = uniqueIgnoreCase(
    b.broker_codes.map((c) => c.broker).filter((br) => brokersOfA.has(br)),
  );

  const resultAliases = uniqueIgnoreCase([
    ...a.aliases,
    ...b.aliases,
    ...[a.name, b.name].filter(
      (n) => n.trim().toLowerCase() !== name.trim().toLowerCase(),
    ),
  ]);

  const { mutate, isPending, error } = useMutation({
    mutationFn: () =>
      mergeLocations({
        keep_id: kept.id,
        remove_id: removed.id,
        name,
        country: String(valueOf("country")),
        country_code: String(valueOf("country_code")),
        city: (valueOf("city") as string | null) ?? "",
        iata: (valueOf("iata") as string | null) ?? "",
        is_airport: Boolean(valueOf("is_airport")),
      }),
    onSuccess: (resp) => {
      setConfirmOpen(false);
      queryClient.invalidateQueries({
        queryKey: ["location-merge-suggestions"],
      });
      queryClient.invalidateQueries({ queryKey: ["location-merge-search"] });
      if (resp) onMerged(resp.location);
    },
  });
  const tError = useTranslatedError(error);

  const pickCell = (key: Field, side: Side) => (
    <td className="p-1">
      <button
        type="button"
        onClick={() => setPicks((p) => ({ ...p, [key]: side }))}
        aria-pressed={picks[key] === side}
        className={clsx(
          "w-full rounded-md border px-3 py-2 text-start transition-colors",
          picks[key] === side
            ? "border-primary bg-primary/10 font-medium"
            : "border-transparent hover:bg-muted/50",
        )}
      >
        {key === "id" ? `#${sides[side].id}` : display(sides[side][key])}
      </button>
    </td>
  );

  return (
    <div className="space-y-4">
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-200 text-gray-500">
              <th className="p-2 text-start font-medium w-32">שדה</th>
              <th className="p-2 text-start font-medium">מיקום א</th>
              <th className="p-2 text-start font-medium">מיקום ב</th>
              <th className="p-2 text-start font-medium bg-muted/30">
                תוצאה
              </th>
            </tr>
          </thead>
          <tbody>
            {fields.map(({ key, label }) => (
              <tr key={key} className="border-b border-gray-100">
                <td className="p-2 text-gray-500">{label}</td>
                {pickCell(key, "a")}
                {pickCell(key, "b")}
                <td className="p-2 px-4 bg-muted/30 font-medium">
                  {key === "id" ? `#${kept.id}` : display(valueOf(key))}
                </td>
              </tr>
            ))}
            <tr className="border-b border-gray-100 align-top">
              <td className="p-2 text-gray-500">ספקים</td>
              {[a, b].map((l) => (
                <td key={l.id} className="p-2 px-4">
                  <div className="flex flex-wrap gap-1">
                    {l.broker_codes.map((c) => (
                      <BrokerCodeBadge key={c.id} code={c} />
                    ))}
                  </div>
                </td>
              ))}
              <td className="p-2 px-4 bg-muted/30">
                <div className="flex flex-wrap gap-1">
                  {[...a.broker_codes, ...b.broker_codes].map((c) => (
                    <BrokerCodeBadge key={c.id} code={c} />
                  ))}
                </div>
              </td>
            </tr>
            <tr className="align-top">
              <td className="p-2 text-gray-500">שמות חלופיים</td>
              {[a, b].map((l) => (
                <td key={l.id} className="p-2 px-4">
                  {l.aliases.length ? l.aliases.join(", ") : display(null)}
                </td>
              ))}
              <td className="p-2 px-4 bg-muted/30">
                {resultAliases.length
                  ? resultAliases.join(", ")
                  : display(null)}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      {sharedBrokers.length > 0 && (
        <div className="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/5 p-3 text-sm text-destructive">
          <AlertTriangle className="size-4 shrink-0" />
          לשני המיקומים יש קוד של אותו ספק ({sharedBrokers.join(", ")}), ולכן
          אי אפשר למזג אותם.
        </div>
      )}

      <div className="flex items-center gap-3">
        <Button
          onClick={() => setConfirmOpen(true)}
          disabled={sharedBrokers.length > 0}
          className="gap-2 bg-blue-600 hover:bg-blue-700 text-white"
        >
          <Merge className="size-4" />
          מיזוג
        </Button>
        {!confirmOpen && <ErrorDisplay>{tError}</ErrorDisplay>}
      </div>

      <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <DialogContent className="max-w-md">
          <DialogTitle>אישור מיזוג</DialogTitle>
          <DialogDescription>
            מיקום #{removed.id} ({removed.name}) יימחק ויאוחד לתוך מיקום #
            {kept.id}, שייקרא &quot;{name}&quot;. אי אפשר לבטל את הפעולה.
          </DialogDescription>
          <ErrorDisplay>{tError}</ErrorDisplay>
          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setConfirmOpen(false)}>
              ביטול
            </Button>
            <Button
              onClick={() => mutate()}
              loading={isPending}
              className="bg-blue-600 hover:bg-blue-700 text-white"
            >
              מיזוג
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
