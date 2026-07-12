import { Loading } from "@/shared/components/Loading";
import { ArrowLeftRight } from "lucide-react";
import { useState } from "react";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";
import { loginBackToAdmin } from "@/shared/api/accounts-api";
import { Button } from "@/components/ui/button";

interface BackToAdminButtonProps {
  buttonText: string;
}

export function BackToAdminButton({ buttonText }: BackToAdminButtonProps) {
  const store = useAuthStore();
  const [loading, setLoading] = useState(false);

  async function handleClick() {
    setLoading(true);

    try {
      const result = await loginBackToAdmin();

      // Update auth store with new session
      store.setSession(result.accessToken, result.accessTokenExpiresAt, {
        id: result.id,
        email: result.email,
        firstName: result.firstName,
        lastName: result.lastName,
        role: result.role as UserRole,
        phoneNumber: result.phoneNumber,
        officeId: result.officeId,
        isAdminAsAgent: false,
      });

      window.location.href = "/admin/agents";
    } catch (error) {
      console.error("Login back to admin failed:", error);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Button
      variant="brand"
      onClick={handleClick}
      disabled={loading}
      className="cursor-pointer inline-flex items-center gap-1 rounded bg-white px-3 py-1 text-xs font-semibold text-brand hover:bg-orange-50 disabled:opacity-50"
    >
      <ArrowLeftRight size={14} />
      {loading ? <Loading /> : buttonText}
    </Button>
  );
}
