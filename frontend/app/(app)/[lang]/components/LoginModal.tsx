"use client";

import { getSession } from "next-auth/react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { User, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { AgentLoginForm } from "./login/AgentLoginForm";
import { AgentSuccessScreen } from "./login/AgentSuccessScreen";
import { CustomerPhoneForm } from "./login/CustomerPhoneForm";
import { CustomerOtpForm } from "./login/CustomerOtpForm";

type LoginMode = "agent" | "customer";
type AgentStep = "credentials" | "success";
type CustomerStep = "phone" | "otp";

export function LoginModal() {
  const t = useTranslations("Login");
  const router = useRouter();

  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<LoginMode>("agent");
  const [agentStep, setAgentStep] = useState<AgentStep>("credentials");
  const [customerStep, setCustomerStep] = useState<CustomerStep>("phone");
  const [customerPhone, setCustomerPhone] = useState("");

  // Restore success screen if next-auth's session update caused a remount
  useEffect(() => {
    if (sessionStorage.getItem("agentLoginSuccess")) {
      sessionStorage.removeItem("agentLoginSuccess");
      setOpen(true);
      setAgentStep("success");
    }
  }, []);

  const handleOpenChange = (next: boolean) => {
    if (!next) {
      setMode("agent");
      setAgentStep("credentials");
      setCustomerStep("phone");
      setCustomerPhone("");
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
    const session = await getSession();
    handleOpenChange(false);
    if (session?.user?.role === "admin") {
      router.push("/admin");
    } else {
      router.refresh();
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
      <DialogTrigger asChild>
        <Button size="outline" variant="outline">
          <User className="size-5" />
          {t("openModal")}
        </Button>
      </DialogTrigger>

      <DialogContent
        className="min-w-96 max-w-md p-6 flex flex-col gap-6 bg-white border-border-light/50 rounded-2xl shadow-modal"
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
                  "flex-1 py-4 px-9 rounded-xl type-paragraph font-bold h-auto transition-colors",
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
                  "flex-1 py-4 px-9 rounded-xl type-paragraph font-bold h-auto transition-colors",
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
              <AgentLoginForm onSuccess={handleAgentSuccess} />
            )}

            {mode === "customer" && customerStep === "phone" && (
              <CustomerPhoneForm onSubmit={handleCustomerPhoneSubmit} />
            )}

            {mode === "customer" && customerStep === "otp" && (
              <CustomerOtpForm phone={customerPhone} />
            )}
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
