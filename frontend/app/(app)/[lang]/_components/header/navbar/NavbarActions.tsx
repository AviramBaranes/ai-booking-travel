"use client";

import { useParams } from "next/navigation";
import { LoginModal } from "../login/LoginModal";
import { AuthenticatedDropdown } from "./AuthenticatedDropdown";
import { useSession } from "next-auth/react";
import { LangSwitcher } from "../login/LangSwitcher";

export function NavbarActions() {
  const { lang } = useParams();
  const session = useSession();
  const isAuthenticated =
    !!session.data?.user && session.data.user.role !== "admin";

  return (
    <div className="flex items-center gap-4">
      <LangSwitcher lang={lang as string} />
      {isAuthenticated ? <AuthenticatedDropdown /> : <LoginModal />}
    </div>
  );
}
