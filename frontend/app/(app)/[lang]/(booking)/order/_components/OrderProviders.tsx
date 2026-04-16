"use client";

import { AuthTokenProvider } from "@/shared/components/providers/AuthTokenProvider";
import { SessionProvider } from "next-auth/react";

export function OrderProviders({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <AuthTokenProvider>{children}</AuthTokenProvider>
    </SessionProvider>
  );
}
