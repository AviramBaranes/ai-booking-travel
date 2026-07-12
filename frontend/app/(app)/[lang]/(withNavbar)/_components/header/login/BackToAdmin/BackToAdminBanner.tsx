"use client";

import { BackToAdminButton } from "./BackToAdminButton";
import { useTranslations } from "next-intl";
import Link from "next/link";
import useAuthStore from "@/shared/auth/authStore";

export function BackToAdminBanner() {
  const t = useTranslations("BackToAdmin");
  const user = useAuthStore((state) => state.user);

  if (!user) return null;

  const isAdmin = user.role === "admin";
  const isAdminAsAgent = user.isAdminAsAgent

  if (!isAdminAsAgent && !isAdmin) return null;

  return (
    <div className="bg-brand text-white text-sm py-2 px-4 flex items-center justify-center gap-3">
      <span>{isAdminAsAgent ? t("adminAsAgentMsg") : t("adminMsg")}</span>
      {isAdminAsAgent ? (
        <BackToAdminButton
          buttonText={t("adminAsAgentBtn")}
        />
      ) : (
        <Link
          className="cursor-pointer inline-flex items-center gap-1 rounded bg-white px-3 py-1 text-xs font-semibold text-brand hover:bg-orange-50"
          href={"/admin"}
        >
          {t("adminBtn")}
        </Link>
      )}
    </div>
  );
}
