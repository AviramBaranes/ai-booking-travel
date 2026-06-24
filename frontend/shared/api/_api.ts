import { getServerSession } from "next-auth";
import Client, { isAPIError, Local, Environment } from "../client";
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

let globalClientSideToken = "";
let isAuthLoading = true; // Tracks Next-Auth initialization state in browser

export function setAuthorizationHeader(token: string) {
  if (typeof window !== "undefined") {
    globalClientSideToken = token;
  }
}

export function removeAuthorizationHeader() {
  if (typeof window !== "undefined") {
    globalClientSideToken = "";
  }
}

export function setAuthLoadingState(loading: boolean) {
  if (typeof window !== "undefined") {
    isAuthLoading = loading;
  }
}

export async function withErrorHandler<T>(
  apiCall: (client: Client) => Promise<T>,
  options?: { skipAuthRedirect?: boolean },
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
    if (globalClientSideToken) {
      localClient = localClient.with({ auth: globalClientSideToken });
    }
  }

  try {
    return await apiCall(localClient);
  } catch (error) {
    if (!isAPIError(error)) throw error;

    if (error.status === 401 && !options?.skipAuthRedirect) {
      if (typeof window !== "undefined") {
        // CRITICAL SAFETY VALVE: If Next-Auth is still loading its session,
        // do NOT hard reload the browser. Throw the error and let React-Query retry.
        if (isAuthLoading) {
          throw error;
        }

        removeAuthorizationHeader();
        const targetUrl = `/${lang}?login=open`;
        window.location.href = targetUrl;
      }
    }

    if (process.env.NODE_ENV === "development") console.log({ error });
    if (error.details && typeof error.details.code === "string") {
      throw new AppError(error.details.code, error.details.field ?? null);
    }
    throw error;
  }
}