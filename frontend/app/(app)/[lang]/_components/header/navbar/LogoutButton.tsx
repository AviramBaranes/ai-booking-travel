"use client";

import { LogOut } from "lucide-react";
import { signOut } from "next-auth/react";

interface LogoutButtonProps {
  buttonText: string;
  onLogout?: () => void;
}
export function LogoutButton({ buttonText, onLogout }: LogoutButtonProps) {
  return (
    <button
      className="flex items-center justify-end gap-2 px-4 min-h-18 cursor-pointer w-full font-medium text-[16px] text-navy transition-colors hover:bg-brand/30!"
      onClick={() => {
        onLogout?.();
        signOut({ callbackUrl: "/he/" });
      }}
    >
      <span>{buttonText}</span>
      <LogOut className="size-6 text-brand shrink-0" />
    </button>
  );
}
