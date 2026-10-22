"use client";

import { useMutation } from "@tanstack/react-query";
import { LogIn } from "lucide-react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { loginAsUser } from "@/shared/api/accounts-api";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";

type LoginAsUserButtonProps = {
  userId: number;
  label: string;
  compact?: boolean;
};

export default function LoginAsUserButton({
  userId,
  label,
  compact = false,
}: LoginAsUserButtonProps) {
  const store = useAuthStore();
  const router = useRouter();

  const { mutate, isPending, isSuccess } = useMutation({
    mutationFn: () => loginAsUser(userId),
    onSuccess: (result) => {
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

  if (compact) {
    return (
      <Button
        variant="ghost"
        size="icon"
        onClick={() => mutate()}
        loading={isPending}
        disabled={isSuccess}
        title={label}
        aria-label={label}
        className="size-7 cursor-pointer text-blue-600 hover:bg-blue-50 disabled:opacity-50"
      >
        <LogIn size={14} />
      </Button>
    );
  }

  return (
    <Button
      variant="ghost"
      onClick={() => mutate()}
      loading={isPending}
      disabled={isSuccess}
      className="cursor-pointer inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 disabled:opacity-50"
    >
      {label}
      <LogIn size={14} />
    </Button>
  );
}