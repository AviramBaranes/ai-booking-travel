import useAuthStore, { User } from "../auth/authStore";
import { Local } from "../client";

interface ApiError {
  code?: string;
  message?: string;
  details?: Record<string, unknown>;
}

interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
  status: number;
}

// Broadcast channel for cross-tab sync
let broadcastChannel: BroadcastChannel | null = null;

try {
  broadcastChannel = new BroadcastChannel("auth-channel");
  broadcastChannel.onmessage = (event) => {
    if (event.data.type === "token-refreshed") {
      onTokenRefreshed(event.data);
    } else if (event.data.type === "logout") {
      useAuthStore.setState({
        accessToken: null,
        expiresAt: null,
        user: null,
        status: "unauthenticated",
      });
    }
  };
} catch (e) {
  // BroadcastChannel not available in some environments
  console.warn("BroadcastChannel not available");
}

function onTokenRefreshed({
  accessToken,
  expiresAt,
  user,
}: {
  accessToken: string;
  expiresAt: number;
  user?: User;
}) {
  if (accessToken && expiresAt && user) {
    useAuthStore.setState({
      accessToken,
      expiresAt,
      user,
    });
  }
}

/**
 * Refresh token using Web Locks for single-flight.
 * If another tab is already refreshing, wait for its result via broadcast.
 */
async function getValidToken(): Promise<string | null> {
  const isServer = typeof window === "undefined";
  if (isServer) return null;
  
  const store = useAuthStore.getState();

  // Check if token is still valid
  if (store.hasValidToken() && !store.isTokenExpiringSoon()) {
    return store.accessToken;
  }

  // Need to refresh - use Web Locks to ensure single-flight in this tab
  if (!("locks" in navigator)) {
    // Web Locks not available, fall back to refresh without coordination
    return await refreshToken();
  }

  try {
    const token = await new Promise<string | null>((resolve) => {
      navigator.locks.request("auth-refresh", async () => {
        // Re-check token after acquiring lock (another tab may have refreshed)
        const currentState = useAuthStore.getState();
        if (
          currentState.hasValidToken() &&
          !currentState.isTokenExpiringSoon()
        ) {
          resolve(currentState.accessToken);
          return;
        }

        // Perform refresh
        const newToken = await refreshToken();
        resolve(newToken);
      });
    });

    return token;
  } catch (e) {
    console.error("Error acquiring lock for token refresh:", e);
    return await refreshToken();
  }
}

/**
 * Call the refresh endpoint to get a new access token.
 * The backend will read the refresh_token from the cookie automatically.
 */
async function refreshToken(): Promise<string | null> {
  const store = useAuthStore.getState();

  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL || Local}/refresh`,
      {
        method: "POST",
        credentials: "include", // Include cookies (refresh_token)
        headers: {
          "Content-Type": "application/json",
        },
        body: "{}", // Empty body
      },
    );

    if (!response.ok) {
      if (response.status === 401) {
        // Refresh token invalid or expired
        store.logout();
        broadcastChannel?.postMessage({ type: "logout" });
      }
      return null;
    }

    const data = await response.json();
    const { accessToken, accessTokenExpiresAt, ...user } = data;

    if (accessToken && accessTokenExpiresAt) {
      store.setAccessToken(accessToken, accessTokenExpiresAt);

      // Notify other tabs
      const data = {
        type: "token-refreshed",
        accessToken,
        expiresAt: accessTokenExpiresAt,
        user,
      };
      broadcastChannel?.postMessage(data);
      onTokenRefreshed(data);

      return accessToken;
    }

    return null;
  } catch (error) {
    console.error("Token refresh failed:", error);
    return null;
  }
}

/**
 * Fetch wrapper that handles token refresh automatically.
 * Ensures token is valid before making the request, and attaches it to the Authorization header.
 */
export async function apiFetch<T = unknown>(
  path: string,
  options: RequestInit = {},
): Promise<ApiResponse<T>> {
  const store = useAuthStore.getState();

  // Ensure token is valid before making the request
  const token = await getValidToken();

  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  // Check for gr_session hint cookie to determine if refresh is available
  const hasSessionCookie = document.cookie.includes("gr_session=1");
  if (!token && !hasSessionCookie) {
    // No token and no session cookie - definitely unauthenticated
    store.setStatus("unauthenticated");
  }

  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || "";
  const url = path.startsWith("http") ? path : `${baseUrl}${path}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers,
      credentials: "include", // Include cookies
    });

    if (response.status === 401) {
      // Unauthorized - token may have expired or been revoked
      store.logout();
      broadcastChannel?.postMessage({ type: "logout" });
    }

    const data = await response.json();

    return {
      data: data,
      status: response.status,
    };
  } catch (error) {
    console.error("API request failed:", error);
    return {
      error: {
        message: error instanceof Error ? error.message : "Unknown error",
      },
      status: 0,
    };
  }
}

export { getValidToken, refreshToken };
