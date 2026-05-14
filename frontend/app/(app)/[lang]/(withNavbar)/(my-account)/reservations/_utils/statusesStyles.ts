/**
 * @param status - reservation status or payment status
 * @returns the tailwind classes for the text color based on the reservation status or payment status
 */
export function statusToColor(status: string) {
  switch (status) {
    case "vouchered":
    case "paid":
      return "text-success font-semibold";
    case "canceled":
      return "text-destructive font-semibold";
    case "booked":
    case "refund_pending":
    case "unpaid":
      return "text-brand font-semibold";
    case "refunded":
      return "text-brand-blue font-semibold";
    default:
      return "text-navy font-semibold";
  }
}

/**
 * @param status - reservation status or payment status
 * @returns the tailwind classes for the background color based on the reservation status or payment status
 */
export function statusToBg(status: string) {
  switch (status) {
    case "vouchered":
    case "paid":
      return "bg-success/10";
    case "canceled":
      return "bg-destructive/10";
    case "booked":
    case "refund_pending":
    case "unpaid":
      return "bg-brand/10";
    case "refunded":
      return "bg-brand-blue/10";
    default:
      return "bg-navy/10";
  }
}
