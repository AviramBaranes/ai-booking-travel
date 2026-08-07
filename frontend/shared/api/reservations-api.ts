import { queries, reports, supplier_payments } from "../client";
import { withErrorHandler } from "./_api";

export function listReservations(params: queries.ListReservationsParams) {
  return withErrorHandler((client) =>
    client.reservation.ListReservations(params),
  );
}

export function getReservationById(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.GetReservation(reservationId),
  );
}

export function applyVoucher(reservationId: number, voucherCode: string) {
  return withErrorHandler((client) =>
    client.reservation.ApplyVoucher(reservationId, {
      voucher: voucherCode,
    }),
  );
}

export function cancelReservation(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.CancelReservation(reservationId),
  );
}

export function listOpenReservations(
  params: queries.ListOpenReservationsByBillingEntityParams,
) {
  return withErrorHandler((client) =>
    client.reservation.ListOpenReservationsByBillingEntity(params),
  );
}

export function businessReservationReport(params: reports.ReportParams) {
  return withErrorHandler((client) =>
    client.reservation.GetBusinessReport(params),
  );
}

export function profitabilityReport(params: reports.ReportParams) {
  return withErrorHandler((client) =>
    client.reservation.GetProfitReport(params),
  );
}

export function collectionsReport() {
  return withErrorHandler((client) =>
    client.reservation.GetBusinessesBalancesReport(),
  );
}

export function listUnpaidSupplierReservations(
  params: supplier_payments.ListUnpaidSupplierReservationsParams,
) {
  return withErrorHandler((client) =>
    client.reservation.ListUnpaidSupplierReservations(params),
  );
}

/**
 * ValidateFlexPaymentSummary is a raw endpoint, so the generated client returns an untyped
 * Response. These shapes mirror the handler's response and must be kept in sync with
 * services/reservation/handlers/supplier_payments/validate_flex_payment_summary.go.
 */
export interface ApprovedSupplierPayment {
  reservationId: number;
  brokerReservationId: string;
  currencyCode: string;
  amount: number;
}

export interface RejectedSupplierPayment {
  /** Absent when no reservation matched the summary line. */
  reservationId?: number;
  brokerReservationId: string;
  reason: string;
  currencyCode: string;
  balance: number;
  expectedAmount?: number;
  expectedCurrencyCode?: string;
}

export interface ValidateFlexPaymentSummaryResponse {
  approved: ApprovedSupplierPayment[];
  rejected: RejectedSupplierPayment[];
}

export function validateFlexPaymentSummary(file: File) {
  const formData = new FormData();
  formData.append("file", file);

  return withErrorHandler(async (client) => {
    const response = await client.reservation.ValidateFlexPaymentSummary(
      "POST",
      formData,
    );
    return (await response.json()) as ValidateFlexPaymentSummaryResponse;
  });
}

export function paySupplierReservations(
  params: supplier_payments.PaySupplierReservationsParams,
) {
  return withErrorHandler((client) =>
    client.reservation.PaySupplierReservations(params),
  );
}

export function getFullReservation(reservationId: number) {
  return withErrorHandler((client) =>
    client.reservation.GetFullReservation(reservationId),
  );
}
