"use client";

import { CheckCircle2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

interface Props {
  onContinue: () => void;
}

export function AgentSuccessScreen({ onContinue }: Props) {
  const t = useTranslations("Login");

  return (
    <div className="flex flex-col items-center gap-6 pb-2 mx-auto">
      <CheckCircle2 className="size-16 text-green-500" strokeWidth={1.5} />
      <div className="flex flex-col items-center gap-2 text-center">
        <h3 className="type-h5 text-navy">{t("agent.successTitle")}</h3>
        <p className="type-paragraph text-text-secondary">
          {t("agent.successMessage")}
        </p>
      </div>
      <Button
        variant="brand"
        className="w-full py-3.5 h-auto"
        onClick={onContinue}
      >
        {t("agent.continueToSite")}
      </Button>
    </div>
  );
}
