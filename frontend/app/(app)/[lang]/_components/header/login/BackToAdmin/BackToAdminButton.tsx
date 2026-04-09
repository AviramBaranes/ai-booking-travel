"use client";
import { Loading } from "@/shared/components/Loading";
import { ArrowLeftRight } from "lucide-react";
import { signIn } from "next-auth/react";
import { useState } from "react";

interface BackToAdminButtonProps {
  accessToken: string;
  buttonText: string;
}

export function BackToAdminButton({
  accessToken,
  buttonText,
}: BackToAdminButtonProps) {
  const [loading, setLoading] = useState(false);

  async function handleClick() {
    if (!accessToken) return;
    setLoading(true);

    try {
      const result = await signIn("admin-login-back", {
        redirect: false,
        accessToken,
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
    <button
      onClick={handleClick}
      disabled={loading}
      className="cursor-pointer inline-flex items-center gap-1 rounded bg-white px-3 py-1 text-xs font-semibold text-brand hover:bg-orange-50 disabled:opacity-50"
    >
      <ArrowLeftRight size={14} />
      {loading ? <Loading /> : buttonText}
    </button>
  );
}
