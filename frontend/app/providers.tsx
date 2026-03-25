"use client";

import { SessionProvider, useSession } from "next-auth/react";
import { useEffect, useState } from "react";

import {
  removeAuthorizationHeader,
  setAuthorizationHeader,
} from "@/shared/api/_api";

function AuthTokenProvider({ children }: { children: React.ReactNode }) {
  const { data: session, status, update } = useSession();
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    if (!session?.user?.customExp) {
      setAuthenticated(false);
      removeAuthorizationHeader();
      return;
    }

    const msUntilExpiry = session.user.customExp * 1000 - Date.now();
    // When the token expires, trigger a session update which will
    // invoke the JWT callback server-side and refresh the token.
    const timer = setTimeout(
      () => {
        update();
      },
      Math.max(msUntilExpiry, 0),
    );

    if (session.user?.accessToken) {
      setAuthorizationHeader(session.user.accessToken);
      setAuthenticated(true);
    } else {
      // Refresh failed — no accessToken means unauthenticated
      removeAuthorizationHeader();
      setAuthenticated(false);
    }

    return () => clearTimeout(timer);
  }, [session, update]);

  if (status === "loading") {
    return null;
  }

  if (status === "unauthenticated") {
    return <>{children}</>;
  }

  if (!authenticated) {
    return null;
  }

  return <>{children}</>;
}

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <AuthTokenProvider>{children}</AuthTokenProvider>
    </SessionProvider>
  );
}
