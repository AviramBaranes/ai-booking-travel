/**
 * Returns the current lang segment from the URL path.
 * - Server-side: reads the X-Lang header set by middleware.
 * - Client-side: parses window.location.pathname.
 */
export async function getLang(): Promise<string> {
  if (typeof window !== "undefined") {
    return window.location.pathname.split("/")[1] || "he";
  }

  const { headers } = await import("next/headers");
  const h = await headers();
  return h.get("x-lang") || "he";
}
