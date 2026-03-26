import { getServerSession } from "next-auth";
import Client, { isAPIError, Local } from "../client";
import { authOptions } from "../auth/authOptions";
import { getLang } from "../lang/lang";
import { AppError } from "./AppError";

let client = new Client(Local);

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
) {
  try {
    const lang = await getLang();
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
    if (process.env.NODE_ENV === "development") console.error({ error });
    if (error.details && typeof error.details.code === "string") {
      throw new AppError(error.details.code, error.details.field ?? null);
    }
    throw error;
  }
}
