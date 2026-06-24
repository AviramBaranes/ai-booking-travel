"use client";

import { Button } from "@/components/ui/button";
import { useDirection } from "@/shared/hooks/useDirection";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";

const ARROW_RIGHT = "\u2192";
const ARROW_LEFT = "\u2190";
export function BackButton({
  translationKey,
  href,
}: {
  translationKey: string;
  href?: string;
}) {
  const t = useTranslations("booking.steps");
  const router = useRouter();
  const dir = useDirection();

  return (
      <Button
        // variant="ghost"
        variant="outline"
        className="flex gap-2 cursor-pointer mt-8 font-bold my-4 p-6 border-navy rounded-md"
        onClick={() => {
          if (href) {
            router.push(href);
          } else {
            router.back();
          }
        }}
      >
        <span className="text-navy text-sm">
          {dir === "rtl" ? ARROW_RIGHT : ARROW_LEFT}
        </span>
        <span className="text-navy text-sm">{t(translationKey)}</span>
      </Button>
  );
}
