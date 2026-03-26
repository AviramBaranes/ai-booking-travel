import { NextAuthOptions } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

import { login } from "../api/auth-api";
import Client, { Local } from "../client";
import { auth } from "../client";

async function refreshAccessToken(token: Record<string, unknown>) {
  try {
    // Call the backend directly, bypassing withErrorHandler to avoid
    // triggering getServerSession which would re-enter the JWT callback.
    const client = new Client(Local);
    const refreshed = await client.auth.RefreshTokens({
      RefreshToken: token.refreshToken as string,
    });

    if (!refreshed) {
      throw new Error("Failed to refresh token");
    }

    return {
      ...token,
      accessToken: refreshed.accessToken,
      refreshToken: refreshed.refreshToken,
      customExp: Math.floor(Date.now() / 1000) + 60 * 15, // 15 minutes
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
        username: { label: "שם משתמש", type: "text" },
        password: { label: "סיסמה", type: "password" },
      },
      async authorize(credentials) {
        const user = await login({
          username: credentials?.username ?? "",
          password: credentials?.password ?? "",
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
        const customExp = Math.floor(Date.now() / 1000) + 60 * 15; // 15 minutes
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
      session.user = token as unknown as auth.LoginResponse & {
        customExp: number;
      };
      return session;
    },
  },
};
