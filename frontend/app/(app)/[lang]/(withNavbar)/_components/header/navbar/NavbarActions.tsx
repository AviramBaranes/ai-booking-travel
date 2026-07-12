"use client";

import { useParams } from "next/navigation";
import { LoginModal } from "../login/LoginModal";
import { AuthenticatedDropdown } from "./AuthenticatedDropdown";
import { LangSwitcher } from "../login/LangSwitcher";
import useAuthStore from "@/shared/auth/authStore";

export function NavbarActions() {
  const { lang } = useParams();
  const user = useAuthStore((state) => state.user);
  const isAuthenticated = !!user && user.role !== "admin";

  return (
    <div className="flex items-center gap-3 lg:gap-4">
      <LangSwitcher lang={lang as string} />
      {isAuthenticated ? <AuthenticatedDropdown /> : <LoginModal />}
    </div>
  );
}
