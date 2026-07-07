import { NextAuthOptions } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

import { login, loginWithOTP } from "../api/accounts-api";
import Client from "../client";
import { auth } from "../client";
import { JWT } from "next-auth/jwt";
import { getBaseURL } from "../api/_api";

// Deduplicates concurrent refresh calls for the same refresh token,
// preventing a race condition where multiple requests all see an expired
// token simultaneously and each invalidate the same refresh token.
const inflightRefreshes = new Map<string, Promise<JWT>>();

async function refreshAccessToken(token: JWT): Promise<JWT> {
  const refreshToken = token.refreshToken as string;

  const inflight = inflightRefreshes.get(refreshToken);
  if (inflight) return inflight;

  const promise = doRefreshAccessToken(token).finally(() => {
    inflightRefreshes.delete(refreshToken);
  });

  inflightRefreshes.set(refreshToken, promise);
  return promise;
}

async function doRefreshAccessToken(token: JWT): Promise<JWT> {
  try {
    // Call the backend directly, bypassing withErrorHandler to avoid
    // triggering getServerSession which would re-enter the JWT callback.
    const client = new Client(getBaseURL());
    const refreshed = await client.accounts.RefreshTokens({
      RefreshToken: token.refreshToken as string,
    });

    if (!refreshed) {
      throw new Error("Failed to refresh token");
    }

    return {
      ...token,
      accessToken: refreshed.accessToken,
      refreshToken: refreshed.refreshToken,
      customExp: refreshed.accessTokenExpiresAt,
    };
  } catch (error) {
    return { ...token, error: "RefreshTokenExpired" };
  }
}

export const authOptions: NextAuthOptions = {
  providers: [
    CredentialsProvider({
      name: "Credentials",
      type: "credentials",
      credentials: {
        email: { label: 'דוא"ל', type: "email" },
        password: { label: "סיסמה", type: "password" },
      },
      async authorize(credentials) {
        const user = await login({
          email: credentials?.email ?? "",
          password: credentials?.password ?? "",
        });

        if (user) return { ...user, id: String(user.id) };
        return null;
      },
    }),
    CredentialsProvider({
      id: "agent-login",
      name: "Agent Login",
      type: "credentials",
      credentials: {
        agentId: { type: "text" },
        accessToken: { type: "text" },
      },
      async authorize(credentials) {
        if (!credentials?.agentId || !credentials?.accessToken) return null;

        const client = new Client(getBaseURL(), {
          auth: credentials.accessToken,
        });
        const user = await client.accounts.LoginAsAgent({
          agentId: Number(credentials.agentId),
        });

        if (user) return { ...user, id: String(user.id), isAdminAsAgent: true };
        return null;
      },
    }),
    CredentialsProvider({
      id: "admin-login-back",
      name: "Admin Login Back",
      type: "credentials",
      credentials: {
        accessToken: { type: "text" },
      },
      async authorize(credentials) {
        if (!credentials?.accessToken) return null;

        const client = new Client(getBaseURL(), {
          auth: credentials.accessToken,
        });
        const user = await client.accounts.LoginBackToAdmin();

        if (user) return { ...user, id: String(user.id) };
        return null;
      },
    }),
    CredentialsProvider({
      id: "customer-login",
      name: "Customer Login",
      type: "credentials",
      credentials: {
        phoneNumber: { type: "text" },
        otp: { type: "text" },
      },
      async authorize(credentials) {
        const user = await loginWithOTP({
          otp: credentials?.otp ?? "",
          phoneNumber: credentials?.phoneNumber ?? "",
        });

        if (user) return { ...user, id: String(user.id) };
        return null;
      },
    }),
    CredentialsProvider({
      id: "customer-login-after-payment",
      name: "Customer Login After Payment",
      type: "credentials",
      credentials: {
        login: { type: "object" },
      },
      async authorize(credentials) {
        if (!credentials?.login) return null;
        const user = JSON.parse(credentials.login) as auth.LoginResponse;

        return { ...user, id: String(user.id) };
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user, trigger }) {
      // Initial sign in
      if (trigger === "signIn" && user) {
        return {
          ...token,
          ...user,
          customExp: (user as unknown as auth.LoginResponse)
            .accessTokenExpiresAt,
        };
      }
      const nowInSeconds = Math.floor(Date.now() / 1000);
      // Safety buffer: If the token expires in less than 6 minutes, refresh it now
      const refreshBuffer = 360;

      if (
        typeof token.customExp === "number" &&
        nowInSeconds < token.customExp - refreshBuffer
      ) {
        return token;
      }

      // Access token has expired or is nearing expiration, invoke refresh loop
      return await refreshAccessToken(token);
    },
    async session({ session, token }) {
      session.user = token as unknown as auth.LoginResponse & {
        customExp: number;
        isAdminAsAgent?: boolean;
      };
      return session;
    },
  },
};
