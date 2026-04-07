"use client";

import { useSession } from "next-auth/react";
import { signIn } from "next-auth/react";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { ArrowLeftRight } from "lucide-react";

export function BackToAdminBanner() {
  const { data: session } = useSession();
  const t = useTranslations("BackToAdmin");
  const [loading, setLoading] = useState(false);

  if (!session?.user?.isAdminAsAgent) return null;

  async function handleClick() {
    if (!session?.user?.accessToken) return;
    setLoading(true);

    try {
      const result = await signIn("admin-login-back", {
        redirect: false,
        accessToken: session.user.accessToken,
      });

      if (result?.error) {
        console.error("Login back to admin failed:", result.error);
        return;
      }

      window.location.href = "/admin/agents";
    } catch (error) {
      console.error("Login back to admin failed:", error);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="bg-brand text-white text-sm py-2 px-4 flex items-center justify-center gap-3">
      <span>{t("message")}</span>
      <button
        onClick={handleClick}
        disabled={loading}
        className="cursor-pointer inline-flex items-center gap-1 rounded bg-white px-3 py-1 text-xs font-semibold text-brand hover:bg-orange-50 disabled:opacity-50"
      >
        <ArrowLeftRight size={14} />
        {loading ? "..." : t("button")}
      </button>
    </div>
  );
}
