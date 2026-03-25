import { getServerSession } from "next-auth";
import Client, { APIError, isAPIError, Local, PreviewEnv } from "../client";
import { authOptions } from "../auth/authOptions";

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

export async function withErrorHandler<T>(
  apiCall: (client: Client) => Promise<T>,
  errorHandlers?: Record<number, ((e: APIError) => T | undefined) | undefined>,
  defaultErrorHandler?: () => T | undefined,
) {
  try {
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
    if (!errorHandlers || !(error.status in errorHandlers)) {
      if (!defaultErrorHandler) return null;
      return defaultErrorHandler();
    }
    return errorHandlers[error.status]?.(error);
  }
}
