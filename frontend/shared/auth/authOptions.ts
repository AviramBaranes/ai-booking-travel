import { NextAuthOptions } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

import { login, loginWithOTP } from "../api/accounts-api";
import Client, { BaseURL, Local } from "../client";
import { accounts } from "../client";
import { JWT } from "next-auth/jwt";

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
    const client = new Client(Local as BaseURL, {});
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
      customExp: Math.floor(Date.now() / 1000) + 60 * 14, // 14 minutes
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

        const client = new Client(Local as BaseURL, {
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

        const client = new Client(Local as BaseURL, {
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
  ],
  callbacks: {
    async jwt({ token, user, trigger }) {
      // Initial sign in
      if (trigger === "signIn" && user) {
        const customExp = Math.floor(Date.now() / 1000) + 60 * 14; // 14 minutes
        return { ...token, ...user, customExp };
      }

      // Return previous token if it hasn't expired yet
      if (
        typeof token.customExp === "number" &&
        Date.now() / 1000 < token.customExp
      ) {
        return token;
      }

      // Access token has expired, try to refresh
      const res = await refreshAccessToken(token);
      return res;
    },
    async session({ session, token }) {
      session.user = token as unknown as accounts.LoginResponse & {
        customExp: number;
        isAdminAsAgent?: boolean;
      };
      return session;
    },
  },
};
