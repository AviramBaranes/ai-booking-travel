import { ReactNode } from "react";

import { cn } from "@/lib/utils";
import { reservation } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";

import { ReservationReportFilters } from "./ReservationsReportFilterBar";

export type ReportColumn<T> = {
  key: string;
  label: string;
  className: string;
  headerClassName?: string;
  render: (row: T) => ReactNode;
};

const statusLabels: Record<string, string> = {
  booked: "הוזמן",
  vouchered: "הופק שובר",
  canceled: "בוטל",
};

const statusClasses: Record<string, string> = {
  booked: "bg-navy/5 text-navy ring-navy/15",
  vouchered: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  canceled: "bg-rose-50 text-rose-700 ring-rose-200",
};

export function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={cn(
        "inline-flex rounded-full px-2.5 py-1 text-xs font-semibold ring-1",
        statusClasses[status] ?? "bg-gray-50 text-gray-700 ring-gray-200",
      )}
    >
      {statusLabels[status] ?? status}
    </span>
  );
}

export function MoneyCell({
  value,
  currency,
  strong = false,
}: {
  value: number;
  currency: string;
  strong?: boolean;
}) {
  return (
    <span
      className={cn(
        "tabular-nums",
        strong ? "font-bold text-emerald-700" : "text-gray-800",
      )}
    >
      {formatPrice(value, currency)}
    </span>
  );
}

export function buildPageNumbers(current: number, total: number): (number | "...")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  const pages: (number | "...")[] = [1];
  if (current > 3) pages.push("...");
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let page = start; page <= end; page++) pages.push(page);
  if (current < total - 2) pages.push("...");
  pages.push(total);
  return pages;
}

export function formatDate(value?: string): string {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString("he-IL");
}

export function formatDateTime(value?: string): string {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("he-IL", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function toNumber(value: string): number | undefined {
  return value ? Number(value) : undefined;
}

export function buildRequest(
  page: number,
  pageSize: number,
  filters: ReservationReportFilters,
): reservation.ReportParams {
  return {
    Page: page,
    PageSize: pageSize,
    PickupDateFrom: filters.pickupFrom || undefined,
    PickupDateTo: filters.pickupTo || undefined,
    CreatedDateFrom: filters.createdFrom || undefined,
    CreatedDateTo: filters.createdTo || undefined,
    VoucheredAtFrom: filters.voucheredFrom || undefined,
    VoucheredAtTo: filters.voucheredTo || undefined,
    Status: filters.status || undefined,
    Broker: filters.broker || undefined,
    Supplier: filters.supplier || undefined,
    OrganizationID: toNumber(filters.organizationId),
    OfficeID: toNumber(filters.officeId),
    AgentID: toNumber(filters.agentId),
    IsBusiness: true,
  };
}

export function makeBaseColumns<
  T extends reservation.BusinessReservationReportRow,
>(): ReportColumn<T>[] {
  return [
    {
      key: "reservationId",
      label: "מזהה",
      className: "min-w-24 text-gray-500",
      render: (row) => row.reservationId,
    },
    {
      key: "brokerReservationId",
      label: "מספר הזמנה",
      className: "min-w-44",
      render: (row) => row.brokerReservationId,
    },
    {
      key: "status",
      label: "סטטוס",
      className: "min-w-32",
      render: (row) => <StatusBadge status={row.status} />,
    },
    {
      key: "organizationName",
      label: "רשת",
      className: "min-w-56 font-medium text-gray-900",
      render: (row) => row.organizationName || "-",
    },
    {
      key: "officeName",
      label: "משרד",
      className: "min-w-52",
      render: (row) => row.officeName || "-",
    },
    {
      key: "agentName",
      label: "סוכן",
      className: "min-w-44",
      render: (row) => row.agentName || "-",
    },
    {
      key: "adminName",
      label: "אדמין",
      className: "min-w-40",
      render: (row) => row.adminName || "-",
    },
    {
      key: "brokerName",
      label: "ברוקר",
      className: "min-w-28",
      render: (row) => (
        <span className="rounded bg-slate-100 px-2 py-1 text-xs font-bold uppercase text-slate-700">
          {row.brokerName}
        </span>
      ),
    },
    {
      key: "supplierName",
      label: "ספק",
      className: "min-w-44",
      render: (row) => row.supplierName || "-",
    },
    {
      key: "countryCode",
      label: "מדינה",
      className: "min-w-24 text-center font-semibold text-gray-600",
      render: (row) => row.countryCode || "-",
    },
    {
      key: "pickupDate",
      label: "איסוף",
      className: "min-w-32 font-medium",
      render: (row) => formatDate(row.pickupDate),
    },
    {
      key: "dropoffDate",
      label: "החזרה",
      className: "min-w-32",
      render: (row) => formatDate(row.dropoffDate),
    },
    {
      key: "rentalDays",
      label: "ימים",
      className: "min-w-20 text-center font-semibold",
      render: (row) => row.rentalDays,
    },
    {
      key: "driverName",
      label: "נהג",
      className: "min-w-52 font-medium text-gray-900",
      render: (row) => row.driverName || "-",
    },
    {
      key: "currencyCode",
      label: "מטבע",
      className: "min-w-24 text-center",
      render: (row) => row.currencyCode,
    },
    {
      key: "currencyRate",
      label: "שער",
      className: "min-w-24 tabular-nums",
      render: (row) => row.currencyRate,
    },
    {
      key: "carSellPriceWithBrokerERP",
      label: "מחיר רכב",
      className: "min-w-36 tabular-nums",
      render: (row) => formatPrice(row.carSellPriceWithBrokerERP, row.currencyCode),
    },
    {
      key: "carSellPriceWithBrokerERPInILS",
      label: 'מחיר רכב בש"ח',
      className: "min-w-44",
      render: (row) => (
        <MoneyCell value={row.carSellPriceWithBrokerERPInILS} currency="ILS" />
      ),
    },
    {
      key: "btERPPrice",
      label: "BT ERP",
      className: "min-w-32 tabular-nums",
      render: (row) => formatPrice(row.btERPPrice, row.currencyCode),
    },
    {
      key: "btERPPriceInILS",
      label: 'BT ERP בש"ח',
      className: "min-w-40",
      render: (row) => <MoneyCell value={row.btERPPriceInILS} currency="ILS" />,
    },
    {
      key: "totalPrice",
      label: 'סה"כ',
      className: "min-w-36 font-semibold tabular-nums",
      render: (row) => formatPrice(row.totalPrice, row.currencyCode),
    },
    {
      key: "totalPriceInILS",
      label: 'סה"כ בש"ח',
      className: "min-w-44 bg-emerald-50/70",
      headerClassName: "bg-emerald-100/80 text-emerald-900",
      render: (row) => <MoneyCell value={row.totalPriceInILS} currency="ILS" strong />,
    },
    {
      key: "voucherNumber",
      label: "שובר",
      className: "min-w-36 text-gray-600",
      render: (row) => row.voucherNumber || "-",
    },
    {
      key: "voucheredAt",
      label: "כורטס בתאריך",
      className: "min-w-44 text-gray-600",
      render: (row) => formatDateTime(row.voucheredAt),
    },
    {
      key: "createdAt",
      label: "נוצר בתאריך",
      className: "min-w-44 text-gray-600",
      render: (row) => formatDateTime(row.createdAt),
    },
  ];
}
