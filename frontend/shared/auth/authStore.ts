import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';

export type UserRole = 'admin' | 'agent' | 'accountant' | 'customer';
export interface User {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  isAdminAsAgent: boolean;
  phoneNumber?: string;
  officeId?: number;
}

export type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'unauthenticated' | 'error';

export interface AuthState {
  // State
  accessToken: string | null;
  expiresAt: number | null; // UNIX timestamp in seconds
  user: User | null;
  status: AuthStatus;
  error: string | null;

  // Actions
  setSession: (accessToken: string, expiresAt: number, user: User) => void;
  setUser: (user: User) => void;
  setAccessToken: (token: string, expiresAt: number) => void;
  setStatus: (status: AuthStatus) => void;
  setError: (error: string | null) => void;
  logout: () => void;
  isTokenExpiringSoon: (bufferSeconds?: number) => boolean;
  hasValidToken: () => boolean;
}

const useAuthStore = create<AuthState>()(
  subscribeWithSelector((set, get) => ({
    accessToken: null,
    expiresAt: null,
    user: null,
    status: 'idle',
    error: null,

    setSession: (accessToken: string, expiresAt: number, user: User) => {
      set({
        accessToken,
        expiresAt,
        user,
        status: 'authenticated',
        error: null,
      });
    },

    setUser: (user: User) => {
      set({ user });
    },

    setAccessToken: (token: string, expiresAt: number) => {
      set({
        accessToken: token,
        expiresAt,
      });
    },

    setStatus: (status: AuthStatus) => {
      set({ status });
    },

    setError: (error: string | null) => {
      set({ error });
    },

    logout: () => {
      set({
        accessToken: null,
        expiresAt: null,
        user: null,
        status: 'unauthenticated',
        error: null,
      });
    },

    isTokenExpiringSoon: (bufferSeconds: number = 300) => {
      const { expiresAt } = get();
      if (!expiresAt) return true;

      const now = Math.floor(Date.now() / 1000);
      return now + bufferSeconds >= expiresAt;
    },

    hasValidToken: () => {
      const { accessToken, expiresAt } = get();
      if (!accessToken || !expiresAt) return false;

      const now = Math.floor(Date.now() / 1000);
      return now < expiresAt;
    },
  }))
);

export default useAuthStore;
