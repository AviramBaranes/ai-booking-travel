import { useTranslations } from "next-intl";
import { Suspense, useCallback, useEffect, useRef, useState } from "react";
import { User, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { AgentLoginForm } from "./AgentLoginForm";
import { AgentSuccessScreen } from "./AgentSuccessScreen";
import { CustomerPhoneForm } from "./CustomerPhoneForm";
import { CustomerOtpForm } from "./CustomerOtpForm";
import { useDialogOpenFromQuery } from "./useDialogOpenFromQuery";
import { ForgotPasswordForm } from "./ForgotPasswordForm";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import useAuthStore from "@/shared/auth/authStore";

type LoginMode = "agent" | "customer";
type AgentStep = "credentials" | "success" | "passwordReset";
type CustomerStep = "phone" | "otp";

type LoginQuerySyncProps = {
  open: () => void;
  registerClearQueryFlag: (clearQueryFlag: () => void) => void;
};

// Keeps useSearchParams behind Suspense without hiding the login button.
function LoginQuerySync({ open, registerClearQueryFlag }: LoginQuerySyncProps) {
  const { clearQueryFlag } = useDialogOpenFromQuery({ open });

  useEffect(() => {
    registerClearQueryFlag(clearQueryFlag);
    return () => registerClearQueryFlag(() => {});
  }, [clearQueryFlag, registerClearQueryFlag]);

  return null;
}

interface LoginModalProps {
  trigger?: React.ReactNode;
}

export function LoginModal({ trigger }: LoginModalProps = {}) {
  const t = useTranslations("Login");

  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<LoginMode>("agent");
  const [agentStep, setAgentStep] = useState<AgentStep>("credentials");
  const [customerStep, setCustomerStep] = useState<CustomerStep>("phone");
  const [customerPhone, setCustomerPhone] = useState("");
  const [passwordResetSuccess, setPasswordResetSuccess] = useState("");
  // The query listener is isolated behind Suspense, so closing the modal calls
  // the latest registered cleanup through this ref instead of reading URL state here.
  const clearQueryFlagRef = useRef<() => void>(() => {});
  const { user } = useAuthStore();

  const openDialog = useCallback(() => {
    setOpen(true);
  }, []);

  const registerClearQueryFlag = useCallback((clearQueryFlag: () => void) => {
    clearQueryFlagRef.current = clearQueryFlag;
  }, []);

  const handleOpenChange = (next: boolean) => {
    if (!next) {
      setMode("agent");
      setAgentStep("credentials");
      setCustomerStep("phone");
      setCustomerPhone("");
      clearQueryFlagRef.current();
    }
    setOpen(next);
  };

  const handleModeSwitch = (newMode: LoginMode) => {
    setMode(newMode);
    setAgentStep("credentials");
    setCustomerStep("phone");
  };

  const handleAgentSuccess = () => {
    setAgentStep("success");
  };

  const handleContinueToSite = async () => {
    handleOpenChange(false);
    if (user?.role === "admin") {
      window.location.href = "/admin";
    }
  };

  const handleCustomerPhoneSubmit = (phone: string) => {
    setCustomerPhone(phone);
    setCustomerStep("otp");
  };

  const headerTitle = () => {
    if (mode === "agent") return t("agent.title");
    if (customerStep === "otp") return t("customer.otpTitle");
    return t("customer.title");
  };

  const headerSubtitle = () => {
    if (mode === "agent") {
      return agentStep === "credentials" ? t("agent.subtitle") : null;
    }
    if (customerStep === "phone") return t("customer.subtitle");
    return t("customer.otpSubtitle", { phone: customerPhone });
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <Suspense fallback={null}>
        <LoginQuerySync
          open={openDialog}
          registerClearQueryFlag={registerClearQueryFlag}
        />
      </Suspense>

      <DialogTrigger asChild>
        {trigger ?? (
          <div>
            <Button size="outline" variant="outline" className="hidden lg:flex">
              <User className="size-5" />
              {t("openModal")}
            </Button>
            <Button variant="ghost" className="lg:hidden px-0">
              <User className="size-5" />
            </Button>
          </div>
        )}
      </DialogTrigger>

      <DialogContent
        className="lg:min-w-96 w-80 max-w-md lg:p-6 p-3 flex flex-col lg:gap-6 gap-3 bg-white border-border-light/50 rounded-2xl shadow-modal"
        showCloseButton={false}
      >
        {/* Header — title on inline-start, close on inline-end */}
        <div className="flex items-start justify-between w-full gap-4">
          <div className="flex flex-col gap-1">
            <DialogTitle className="type-h5 text-navy">
              {headerTitle()}
            </DialogTitle>
            {headerSubtitle() && (
              <p className="type-paragraph text-text-secondary">
                {headerSubtitle()}
              </p>
            )}
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => handleOpenChange(false)}
            className="text-text-secondary -me-1.5 -mt-1.5 size-8 shrink-0"
            aria-label="Close"
          >
            <X size={18} />
          </Button>
        </div>

        {/* Success state */}
        {mode === "agent" && agentStep === "success" ? (
          <AgentSuccessScreen onContinue={handleContinueToSite} />
        ) : (
          <>
            {/* Tab switcher — agent on inline-start (right in RTL) */}
            <div className="flex gap-4 items-center w-full">
              <Button
                onClick={() => handleModeSwitch("agent")}
                className={cn(
                  "flex-1 lg:py-4 lg:px-9 py-3 px-5 rounded-xl type-paragraph font-bold h-auto transition-colors",
                  mode === "agent"
                    ? "bg-navy text-white hover:bg-navy/90"
                    : "bg-background border border-navy text-navy hover:bg-navy/5",
                )}
              >
                {t("tab.agent")}
              </Button>
              <Button
                onClick={() => handleModeSwitch("customer")}
                className={cn(
                  "flex-1 lg:py-4 lg:px-9 py-3 px-5 rounded-xl type-paragraph font-bold h-auto transition-colors",
                  mode === "customer"
                    ? "bg-navy text-white hover:bg-navy/90"
                    : "bg-background border border-navy text-navy hover:bg-navy/5",
                )}
              >
                {t("tab.customer")}
              </Button>
            </div>

            <div className="h-px w-full bg-border-light/50" />

            {mode === "agent" && agentStep === "credentials" && (
              <>
                {passwordResetSuccess && (
                  <SuccessBadge>{passwordResetSuccess}</SuccessBadge>
                )}
                <AgentLoginForm
                  onSuccess={handleAgentSuccess}
                  onForgotPassword={() => setAgentStep("passwordReset")}
                />
              </>
            )}

            {agentStep === "passwordReset" && (
              <ForgotPasswordForm
                onBackToLogin={() => setAgentStep("credentials")}
                onSuccess={() => {
                  setPasswordResetSuccess(t("agent.passwordResetSuccess"));
                  setAgentStep("credentials");
                }}
              />
            )}

            {mode === "customer" && customerStep === "phone" && (
              <CustomerPhoneForm onSubmit={handleCustomerPhoneSubmit} />
            )}

            {mode === "customer" && customerStep === "otp" && (
              <CustomerOtpForm
                phone={customerPhone}
                onSuccess={handleContinueToSite}
              />
            )}
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
