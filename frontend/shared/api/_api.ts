import Client, { isAPIError, Local, Environment, APIError } from "../client";
import { getLang } from "../lang/lang";
import { AppError } from "./AppError";
import { getValidToken } from "./apiClient";

export function getBaseURL(): string {
  const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (baseURL) return baseURL;
  const env = process.env.NEXT_PUBLIC_ENCORE_ENV;
  if (!env || env === "local") return Local;
  return Environment(env);
}

// Kept for backward-compat with AuthTokenProvider (no-op; auth is now Zustand).
type SessionStatus = "authenticated" | "unauthenticated" | "loading";
export function setClientSession(
  _accessToken: string | null,
  _status: SessionStatus,
) {}

export async function withErrorHandler<T>(
  apiCall: (client: Client) => Promise<T>,
  options?: {
    skipAuthRedirect?: boolean;
    onExpectedError?: { [status: number]: (error: APIError) => void };
  },
) {
  const lang = await getLang();
  let localClient = new Client(getBaseURL(), {
    requestInit: { credentials: "include" },
  });

  if (lang) {
    localClient = localClient.with({
      requestInit: { credentials: "include", headers: { "X-Lang": lang } },
    });
  }

  const token = await getValidToken();
  if (token) {
    localClient = localClient.with({ auth: token });
  }

  try {
    return await apiCall(localClient);
  } catch (error) {
    if (!isAPIError(error)) throw error;

    if (options?.onExpectedError?.[error.status]) {
      options.onExpectedError[error.status](error);
      return null as unknown as T; // Return a dummy value to satisfy the return type
    }

    if (
      error.status === 401 &&
      !options?.skipAuthRedirect &&
      typeof window !== "undefined"
    ) {
      // The token is resolved before the request is sent, so a 401 here means a
      // genuinely invalid or expired session.
      window.location.href = `/${lang}?login=open`;
    }

    if (process.env.NODE_ENV === "development") console.log({ error });
    if (error.details && typeof error.details.code === "string") {
      throw new AppError(error.details.code, error.details.field ?? null);
    }
    throw error;
  }
}
