import { withErrorHandler } from "./_api";

export function listMonthlyReports() {
  return withErrorHandler((client) => client.billing.ListMonthlyReports());
}

/**
 * Returns a short-lived signed URL to one archived monthly report. The link points straight at the
 * private bucket, so it must be used right after it is fetched.
 */
export function getMonthlyReportUrl(
  entityType: string,
  entityId: number,
  period: string,
) {
  return withErrorHandler((client) =>
    client.billing.GetMonthlyReportURL(entityType, entityId, period),
  );
}
