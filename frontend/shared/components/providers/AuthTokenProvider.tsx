"use client";

import { signOut, useSession } from "next-auth/react";
import { useEffect } from "react";
import {
  removeAuthorizationHeader,
  setAuthorizationHeader,
  setAuthLoadingState,
} from "@/shared/api/_api";

export function AuthTokenProvider({ children }: { children: React.ReactNode }) {
  const { data: session, status } = useSession();

  useEffect(() => {
    // Synchronize the loading state immediately with the API engine
    if (status === "loading") {
      setAuthLoadingState(true);
      return;
    }
    
    setAuthLoadingState(false);

    if (session?.user?.error === "RefreshTokenExpired") {
      removeAuthorizationHeader();
      signOut({ redirect: false });
      return;
    }

    if (session?.user?.accessToken) {
      setAuthorizationHeader(session.user.accessToken);
    } else {
      removeAuthorizationHeader();
    }
  }, [session, status]);

  return <>{children}</>;
}