"use client";

import { Button } from "@/components/ui/button";
import { LogOut } from "lucide-react";
import { signOut } from "next-auth/react";
import { useTranslations } from "next-intl";

export function LogoutButton() {
  const t = useTranslations("Logout");

  return (
    <Button variant="ghost" onClick={() => signOut({ callbackUrl: "/he/" })}>
      <LogOut size={16} />
      {t("Logout")}
    </Button>
  );
}
