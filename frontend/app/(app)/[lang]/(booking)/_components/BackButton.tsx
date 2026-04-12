"use client";

import { Button } from "@/components/ui/button";
import { useDirection } from "@/shared/hooks/useDirection";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";

const ARROW_RIGHT = "\u2192";
const ARROW_LEFT = "\u2190";
export function BackButton() {
  const t = useTranslations("booking.steps");
  const router = useRouter();
  const dir = useDirection();

  return (
    <Button
      variant="ghost"
      className="flex gap-2 cursor-pointer mt-8"
      onClick={() => router.back()}
    >
      <span className="text-link text-sm">
        {dir === "rtl" ? ARROW_RIGHT : ARROW_LEFT}
      </span>
      <span className="text-link text-sm">{t("backToResults")}</span>
    </Button>
  );
}
