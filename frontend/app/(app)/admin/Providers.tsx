"use client";

import { AuthTokenProvider } from "@/shared/components/providers/AuthTokenProvider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SessionProvider } from "next-auth/react";

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <SessionProvider>
        <AuthTokenProvider>{children}</AuthTokenProvider>
      </SessionProvider>
    </QueryClientProvider>
  );
}
