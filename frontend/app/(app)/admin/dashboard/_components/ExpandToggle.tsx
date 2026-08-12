"use client";

import { ChevronsDownUp, ChevronsUpDown } from "lucide-react";

/**
 * Categories are folded into "אחר" past a limit because colour and label space stop
 * separating them, not to keep cards the same height. This lets you override that per card
 * when you actually need the full list.
 */
export function ExpandToggle({
  expanded,
  onToggle,
  hiddenCount,
  limit,
}: {
  expanded: boolean;
  onToggle: () => void;
  hiddenCount: number;
  limit: number;
}) {
  if (hiddenCount <= 0) return null;

  return (
    <button
      type="button"
      onClick={onToggle}
      aria-pressed={expanded}
      className="flex items-center gap-1 rounded border border-gray-200 px-2 py-1 text-xs text-gray-600 transition-colors hover:bg-gray-50"
      title={
        expanded
          ? `הצגת ${limit} המובילים בלבד`
          : `הצגת כל הערכים (${hiddenCount} מקובצים תחת "אחר")`
      }
    >
      {expanded ? (
        <ChevronsDownUp className="size-3.5" />
      ) : (
        <ChevronsUpDown className="size-3.5" />
      )}
      {expanded ? `${limit} מובילים` : `הצג הכל (${hiddenCount}+)`}
    </button>
  );
}
