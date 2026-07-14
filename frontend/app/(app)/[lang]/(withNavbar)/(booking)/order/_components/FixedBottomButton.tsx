"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";

interface FixedBottomButtonProps {
  isDisabled: boolean;
  loading: boolean;
  watchRef: React.RefObject<HTMLElement | null>;
}

export function FixedBottomButton({
  loading,
  watchRef,
  isDisabled,
}: FixedBottomButtonProps) {
  const t = useTranslations("booking.plansPage");
  const [isInView, setIsInView] = useState(true);

  useEffect(() => {
    const el = watchRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => setIsInView(entry.isIntersecting),
      { threshold: 0 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [watchRef]);

  return (
    <div
      id="fixed-bottom-buttons"
      className={`fixed z-999 bottom-0 w-full lg:hidden left-0 right-0 ${isInView ? "hidden!" : ""}`}
    >
      <Button
        type="submit"
        variant="brand"
        className="type-paragraph font-bold w-full p-8 cursor-pointer rounded-t-none border border-brand"
        disabled={isDisabled}
        loading={loading}
      >
        {t("continueCta")}
      </Button>
    </div>
  );
}
