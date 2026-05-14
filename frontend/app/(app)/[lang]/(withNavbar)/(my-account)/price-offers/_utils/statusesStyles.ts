
/**
 * @param status - price offer status
 * @returns the tailwind classes for the text color based on the price offer status
 */
export function statusToColor(status: string) {
  switch (status) {
    case "booked":
      return "text-success font-semibold";
    case "declined":
    case "unavailable":
      return "text-destructive font-semibold";
    case "open":
      return "text-brand font-semibold";
    default:
      return "text-navy font-semibold";
  }
}

/**
 * @param status - price offer status
 * @returns the tailwind classes for the background color based on the price offer status
 */
export function statusToBg(status: string) {
  switch (status) {
    case "booked":
      return "bg-success/10";
    case "declined":
      return "bg-destructive/10";
    case "open":
      return "bg-brand/10";
    default:
      return "bg-navy/10";
  }
}
