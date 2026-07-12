"use client";

import { AuthBootstrap } from "@/shared/components/providers/AuthBootstrap";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthBootstrap>
        {children}
      </AuthBootstrap>
    </QueryClientProvider>
  );
}
