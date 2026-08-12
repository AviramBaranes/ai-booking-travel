import { formatPrice } from "@/shared/utils/formatPrice";

/** Money on the dashboard is always ILS — the backend converts every row before sending it. */
export function ils(value: number): string {
  return formatPrice(value, "ILS");
}

/** Axis ticks and dense labels, where the full number would not fit. */
export function ilsCompact(value: number): string {
  const abs = Math.abs(value);
  if (abs >= 1_000_000) return `₪${(value / 1_000_000).toFixed(1)}M`;
  if (abs >= 1_000) return `₪${Math.round(value / 1_000)}K`;
  return `₪${Math.round(value)}`;
}

/**
 * Recharts hands label formatters a loosely typed value, so wrap the numeric formatters
 * rather than casting at every call site.
 */
export function labelFormatter(format: (value: number) => string) {
  return (value: unknown) => (value == null ? "" : format(Number(value)));
}

export function percent(value: number, digits = 1): string {
  return `${value.toFixed(digits)}%`;
}

export function count(value: number): string {
  return new Intl.NumberFormat("he-IL").format(value);
}

export function decimal(value: number, digits = 1): string {
  return new Intl.NumberFormat("he-IL", {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  }).format(value);
}

/**
 * The percentage change between two periods. Returns null when the previous period was
 * empty — "up from zero" is not a meaningful percentage.
 */
export function delta(current: number, previous: number): number | null {
  if (previous === 0) return null;
  return ((current - previous) / Math.abs(previous)) * 100;
}

export const RESERVATION_STATUS_LABELS: Record<string, string> = {
  booked: "הוזמן",
  vouchered: "הופק שובר",
  canceled: "בוטל",
};

export const PAYMENT_STATUS_LABELS: Record<string, string> = {
  unpaid: "לא שולם",
  paid: "שולם",
  refund_pending: "ממתין לזיכוי",
  refunded: "זוכה",
};

export const GEAR_TYPE_LABELS: Record<string, string> = {
  auto: "אוטומט",
  manual: "ידני",
  electric: "חשמלי",
};

export const BROKER_LABELS: Record<string, string> = {
  flex: "Flex",
  hertz: "Hertz",
};

/** Status colours are reserved for state and are always paired with the label above. */
export const RESERVATION_STATUS_COLORS: Record<string, string> = {
  booked: "var(--chart-warning)",
  vouchered: "var(--chart-good)",
  canceled: "var(--chart-critical)",
};

export const PAYMENT_STATUS_COLORS: Record<string, string> = {
  unpaid: "var(--chart-warning)",
  paid: "var(--chart-good)",
  refund_pending: "var(--chart-serious)",
  refunded: "var(--chart-critical)",
};

export const SERIES = {
  revenue: "var(--chart-1)",
  cost: "var(--chart-2)",
  profit: "var(--chart-3)",
  business: "var(--chart-4)",
  private: "var(--chart-5)",
  primary: "var(--chart-1)",
} as const;

/** Categorical slots in their validated order. Assign in order; never generate a 9th. */
export const CATEGORICAL = [
  "var(--chart-1)",
  "var(--chart-2)",
  "var(--chart-3)",
  "var(--chart-4)",
  "var(--chart-5)",
  "var(--chart-6)",
  "var(--chart-7)",
  "var(--chart-8)",
] as const;

const COUNTRY_NAMES: Record<string, string> = {
  IL: "ישראל",
  US: "ארה״ב",
  GB: "בריטניה",
  GR: "יוון",
  IT: "איטליה",
  ES: "ספרד",
  FR: "צרפת",
  DE: "גרמניה",
  PT: "פורטוגל",
  CY: "קפריסין",
  NL: "הולנד",
  BG: "בולגריה",
  GE: "גאורגיה",
  TR: "טורקיה",
  AE: "איחוד האמירויות",
  TH: "תאילנד",
  CA: "קנדה",
  AT: "אוסטריה",
  CH: "שווייץ",
  BE: "בלגיה",
  CZ: "צ׳כיה",
  HU: "הונגריה",
  PL: "פולין",
  RO: "רומניה",
  HR: "קרואטיה",
  SI: "סלובניה",
  ME: "מונטנגרו",
  AL: "אלבניה",
  RS: "סרביה",
  IE: "אירלנד",
  DK: "דנמרק",
  SE: "שוודיה",
  NO: "נורווגיה",
  FI: "פינלנד",
  MX: "מקסיקו",
  AU: "אוסטרליה",
  NZ: "ניו זילנד",
  ZA: "דרום אפריקה",
  JP: "יפן",
  MA: "מרוקו",
};

export function countryLabel(code: string): string {
  if (!code) return "לא ידוע";
  return COUNTRY_NAMES[code] ? `${COUNTRY_NAMES[code]} (${code})` : code;
}

export function orEmptyLabel(value: string, fallback = "לא ידוע"): string {
  return value || fallback;
}
