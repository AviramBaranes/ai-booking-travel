"use client";

import { useSession } from "next-auth/react";
import { signIn } from "next-auth/react";
import { useState } from "react";
import { LogIn } from "lucide-react";

export default function LoginAsAgentButton({ agentId }: { agentId: number }) {
  const { data: session } = useSession();
  const [loading, setLoading] = useState(false);

  async function handleClick() {
    if (!session?.user?.accessToken) return;
    setLoading(true);

    try {
      const result = await signIn("agent-login", {
        redirect: false,
        agentId: String(agentId),
        accessToken: session.user.accessToken,
      });

      if (result?.error) {
        console.error("Login as agent failed:", result.error);
        return;
      }

      window.location.href = "/he";
    } catch (error) {
      console.error("Login as agent failed:", error);
    } finally {
      setLoading(false);
    }
  }

  return (
    <button
      onClick={handleClick}
      disabled={loading}
      className="cursor-pointer inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 disabled:opacity-50"
    >
      {loading ? "מתחבר..." : "התחבר כסוכן"}
      <LogIn size={14} />
    </button>
  );
}
