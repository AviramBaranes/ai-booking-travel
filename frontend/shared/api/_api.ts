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

let client = new Client(getBaseURL());

export function setAuthorizationHeader(token: string) {
  client = client.with({
    auth: token,
  });
}

export function removeAuthorizationHeader() {
  client = client.with({
    auth: "",
  });
}

function setLangHeader(lang: string) {
  if (!lang) return;
  client = client.with({
    requestInit: {
      headers: {
        "X-Lang": lang,
      },
    },
  });
}

export async function withErrorHandler<T>(
  apiCall: (client: Client) => Promise<T>,
  options?: { skipAuthRedirect?: boolean },
) {
  const lang = await getLang();
  try {
    setLangHeader(lang);
    if (typeof window === "undefined") {
      const session = await getServerSession(authOptions);
      if (session) {
        setAuthorizationHeader(session.user.accessToken);
      }
    }

    return await apiCall(client);
  } catch (error) {
    if (!isAPIError(error)) throw error;
    if (error.status === 401 && !options?.skipAuthRedirect) {
      removeAuthorizationHeader();
      if (typeof window !== "undefined") {
        window.location.href = `/${lang}?login=open`;
      }
    }
    if (process.env.NODE_ENV === "development") console.log({ error });
    if (error.details && typeof error.details.code === "string") {
      throw new AppError(error.details.code, error.details.field ?? null);
    }
    throw error;
  }
}
