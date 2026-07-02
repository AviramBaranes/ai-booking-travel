import { he } from "react-day-picker/locale";

export function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export function formatRangeDate(lang: string, date: Date) {
  const locale = lang === "he" ? he : undefined;
  const parts = new Intl.DateTimeFormat(locale?.code ?? "en-GB", {
    weekday: "short",
    day: "numeric",
    month: "short",
  }).formatToParts(date);

  const get = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((p) => p.type === type)?.value ?? "";

  return `${get("weekday")} ${get("day")} ${get("month")}`;
}
