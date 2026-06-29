import { getServerSession } from "next-auth";
import Client, { isAPIError, Local, Environment, APIError } from "../client";
import { authOptions } from "../auth/authOptions";
import { getLang } from "../lang/lang";
import { AppError } from "./AppError";

export function getBaseURL(): string {
  const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (baseURL) return baseURL;
  const env = process.env.NEXT_PUBLIC_ENCORE_ENV;
  if (!env || env === "local") return Local;
  return Environment(env);
}

// --- Client-side access token (in-memory, fed by AuthTokenProvider) ---
// The token is exposed as an awaitable value so requests made during Next-Auth
// initialization wait for the first session resolution instead of racing it
// (which previously caused token-less 401s and a spurious redirect to login).
let clientToken: string | null = null;
let sessionSettled = false;
let resolveSessionReady: () => void;
const sessionReady = new Promise<void>((resolve) => {
  resolveSessionReady = resolve;
});

type SessionStatus = "authenticated" | "unauthenticated" | "loading";

/**
 * Pushes the latest Next-Auth session token into the API engine. Called by
 * AuthTokenProvider. The transient "loading" status is ignored so that
 * getClientToken() only resumes once the session has actually resolved.
 */
export function setClientSession(
  accessToken: string | null,
  status: SessionStatus,
) {
  if (typeof window === "undefined" || status === "loading") return;
  clientToken = accessToken;
  if (!sessionSettled) {
    sessionSettled = true;
    resolveSessionReady();
  }
}

/**
 * Resolves the current client access token, waiting for Next-Auth to finish
 * initializing on the first call. Returns null when the user is unauthenticated.
 */
async function getClientToken(): Promise<string | null> {
  if (!sessionSettled) await sessionReady;
  return clientToken;
}

export async function withErrorHandler<T>(
  apiCall: (client: Client) => Promise<T>,
  options?: {
    skipAuthRedirect?: boolean;
    onExpectedError?: { [status: number]: (error: APIError) => void };
  },
) {
  const isServer = typeof window === "undefined";
  const lang = await getLang();
  let localClient = new Client(getBaseURL());

  if (lang) {
    localClient = localClient.with({
      requestInit: { headers: { "X-Lang": lang } },
    });
  }

  if (isServer) {
    const session = await getServerSession(authOptions);
    if (session?.user?.accessToken) {
      localClient = localClient.with({ auth: session.user.accessToken });
    }
  } else {
    // Resolve the token, waiting for Next-Auth to finish initializing so the
    // request never goes out token-less during hydration.
    const token = await getClientToken();
    if (token) {
      localClient = localClient.with({ auth: token });
    }
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
      // genuinely invalid or expired session rather than an initialization race.
      clientToken = null;
      window.location.href = `/${lang}?login=open`;
    }

    if (process.env.NODE_ENV === "development") console.log({ error });
    if (error.details && typeof error.details.code === "string") {
      throw new AppError(error.details.code, error.details.field ?? null);
    }
    throw error;
  }
}
