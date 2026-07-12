"use client";

import useAuthStore from "@/shared/auth/authStore";
import { useEffect, ReactNode, Suspense } from "react";
import { getValidToken } from "@/shared/api/apiClient";

/**
 * Inner component that loads auth state and suspends if needed.
 */
function AuthBootstrapInner() {
  const store = useAuthStore();

  useEffect(() => {
    const initAuth = async () => {
      try {
        // Check if refresh_token cookie exists (gr_session hint)
        const hasSession = document.cookie.includes("gr_session=1");

        if (hasSession) {
          const token = await getValidToken();

          if (token) {
            store.setStatus("authenticated");
          } else {
            store.setStatus("unauthenticated");
          }
        } else {
          // No session cookie, user is not logged in
          store.setStatus("unauthenticated");
        }
      } catch (error) {
        console.error("Failed to initialize auth:", error);
        store.setStatus("error");
        store.setError(
          error instanceof Error ? error.message : "Unknown error",
        );
      }
    };

    initAuth();
  }, []);

  // Return null to allow suspense boundary to show fallback
  // The component that uses useAuthStore will check isInitialized
  return null;
}

/**
 * AuthBootstrap initializes the auth state on app boot.
 * It reads the refresh_token cookie and attempts to validate the session.
 * Wrapped in Suspense with a fallback to prevent hydration issues.
 */
export function AuthBootstrap({ children }: { children: ReactNode }) {
  return (
    <Suspense fallback={null}>
      <AuthBootstrapInner />
      {children}
    </Suspense>
  );
}
