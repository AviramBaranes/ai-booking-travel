"use client";

import { useEffect, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import { useParams } from "next/navigation";
import { useSearchParams, useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";

interface FixedBottomButtonsProps {
  isAgent: boolean;
  isErpSelected: boolean;
  setIsErpDialogOpen: (open: boolean) => void;
  setIsPriceOfferDialogOpen: (open: boolean) => void;
  /** The element to observe — buttons are shown only when it's out of view. */
  watchRef: React.RefObject<HTMLElement | null>;
}

export function FixedBottomButtons({
  isAgent,
  isErpSelected,
  setIsErpDialogOpen,
  setIsPriceOfferDialogOpen,
  watchRef,
}: FixedBottomButtonsProps) {
  const t = useTranslations("booking.plansPage");
  const { lang } = useParams();
  const router = useRouter(); 
  const currentSearchParams = useSearchParams();

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

  const handleContinue = () => {
    if (isErpSelected) {
      router.push(`/${lang}/order?${currentSearchParams.toString()}`);
    } else {
      setIsErpDialogOpen(true);
    }
  };

  return (
    <div
      id="fixed-bottom-buttons"
      className={`fixed bottom-0 w-full lg:hidden left-0 right-0 ${isInView ? "hidden!" : ""}`}
    >
      {isAgent ? (
        <div className="w-full">
          <Button
            variant="brand"
            className="type-paragraph font-bold px-8 py-8 cursor-pointer border-navy bg-navy w-1/2 rounded-none!"
            onClick={() => setIsPriceOfferDialogOpen(true)}
          >
            {t("createPriceOffer")}
          </Button>
          <Button
            variant="brand"
            className="type-paragraph font-bold p-8 cursor-pointer border-brand w-1/2 rounded-none!"
            onClick={handleContinue}
          >
            {t("continueCta")}
          </Button>
        </div>
      ) : (
        <Button
          variant="brand"
          className="type-paragraph font-bold w-full p-8 cursor-pointer rounded-t-none border border-brand"
          onClick={handleContinue}
        >
          {t("continueCta")}
        </Button>
      )}
    </div>
  );
}
