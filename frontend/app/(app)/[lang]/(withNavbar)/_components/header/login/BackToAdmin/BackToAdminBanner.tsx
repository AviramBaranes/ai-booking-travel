"use client";

import { BackToAdminButton } from "./BackToAdminButton";
import { useSession } from "next-auth/react";
import { useTranslations } from "next-intl";

export function BackToAdminBanner() {
  const t = useTranslations("BackToAdmin");
  const { data: session } = useSession();

  if (!session?.user?.isAdminAsAgent) return null;

  return (
    <div className="bg-brand text-white text-sm py-2 px-4 flex items-center justify-center gap-3">
      <span>{t("message")}</span>
      <BackToAdminButton
        accessToken={session.user.accessToken}
        buttonText={t("button")}
      />
    </div>
  );
}
