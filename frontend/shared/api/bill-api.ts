import { billing } from "../client";
import { withErrorHandler } from "./_api";

export function bill(params: billing.BillRequestParams) {
    return withErrorHandler((client) => client.billing.Bill(params));
}