"use client";

import { useState } from "react";
import { LogIn } from "lucide-react";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";
import { loginAsAgent } from "@/shared/api/accounts-api";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

export default function LoginAsAgentButton({ agentId }: { agentId: number }) {
  const store = useAuthStore();
  const router = useRouter();

  const { mutate, isPending, isSuccess } = useMutation({
    mutationFn: () => loginAsAgent(agentId),
    onSuccess: (result) => {
      // Update auth store with new session
      store.setSession(result.accessToken, result.accessTokenExpiresAt, {
        id: result.id,
        email: result.email,
        firstName: result.firstName,
        lastName: result.lastName,
        role: result.role as UserRole,
        phoneNumber: result.phoneNumber,
        officeId: result.officeId,
        isAdminAsAgent: true,
      });
      router.push("/he");
    },
  });

  return (
    <Button
      variant="ghost"
      onClick={() => mutate()}
      loading={isPending}
      disabled={isSuccess}
      className="cursor-pointer inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 disabled:opacity-50"
    >
      התחבר כסוכן
      <LogIn size={14} />
    </Button>
  );
}
