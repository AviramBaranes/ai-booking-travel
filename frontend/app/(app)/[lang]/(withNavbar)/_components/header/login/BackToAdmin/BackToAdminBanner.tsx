"use client";

import { BackToAdminButton } from "./BackToAdminButton";
import { useSession } from "next-auth/react";
import { useTranslations } from "next-intl";
import Link from "next/link";

export function BackToAdminBanner() {
  const t = useTranslations("BackToAdmin");
  const { data: session } = useSession();

  const isAdmin = session?.user?.role === "admin";
  const isAdminAsAgent = session?.user?.isAdminAsAgent;
  if (!isAdminAsAgent && !isAdmin) return null;

  return (
    <div className="bg-brand text-white text-sm py-2 px-4 flex items-center justify-center gap-3">
      <span>{isAdminAsAgent ? t("adminAsAgentMsg") : t("adminMsg")}</span>
      {isAdminAsAgent ? (
        <BackToAdminButton
          accessToken={session.user.accessToken}
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
