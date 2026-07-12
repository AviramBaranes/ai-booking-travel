"use client";

import { LogOut } from "lucide-react";
import useAuthStore from "@/shared/auth/authStore";
import { logout } from "@/shared/api/accounts-api";

interface LogoutButtonProps {
  buttonText: string;
  onLogout?: () => void;
}
export function LogoutButton({ buttonText, onLogout }: LogoutButtonProps) {
  const store = useAuthStore();

  const handleLogout = async () => {
    onLogout?.();
    try {
      await logout();
    } catch (error) {
      console.error("Logout error:", error);
    }
    store.logout();
    window.location.href = "/he/";
  };

  return (
    <button
      className="flex w-full cursor-pointer items-center gap-2 px-3 py-3 text-[15px] font-medium text-navy transition-colors hover:bg-brand/30! md:min-h-18 md:px-4 md:py-0 md:text-[16px]"
      onClick={handleLogout}
    >
      <LogOut className="size-4 lg:size-6 text-brand shrink-0" />
      <span>{buttonText}</span>
    </button>
  );
}
