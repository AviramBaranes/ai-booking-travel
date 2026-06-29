"use client";

import { signOut, useSession } from "next-auth/react";
import { useEffect } from "react";
import { setClientSession } from "@/shared/api/_api";

export function AuthTokenProvider({ children }: { children: React.ReactNode }) {
  const { data: session, status } = useSession();

  useEffect(() => {
    if (status === "loading") return;

    if (session?.user?.error === "RefreshTokenExpired") {
      setClientSession(null, "unauthenticated");
      signOut({ redirect: false });
      return;
    }

    setClientSession(session?.user?.accessToken ?? null, status);
  }, [session, status]);

  return <>{children}</>;
}