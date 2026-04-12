"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useRouter, useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { booking } from "@/shared/client";
import { bookingKeys } from "@/shared/hooks/useAvailableCars";
import { searchRequestToParams } from "../results/searchQuery";

const REDIRECT_SECONDS = 5;

interface ExpiredSearchGateProps {
  children: React.ReactNode;
  searchRequest: booking.SearchAvailabilityRequest;
}

export function ExpiredSearchGate({
  children,
  searchRequest,
}: ExpiredSearchGateProps) {
  const t = useTranslations("booking.expiredSearch");
  const queryClient = useQueryClient();
  const router = useRouter();
  const params = useParams<{ lang: string }>();

  const langParam = params.lang;
  const lang = Array.isArray(langParam) ? langParam[0] : langParam;

  const [secondsLeft, setSecondsLeft] = useState(REDIRECT_SECONDS);
  const redirectedRef = useRef(false);

  const hasAvailabilityData =
    queryClient.getQueryData(bookingKeys.availability(searchRequest)) !==
    undefined;

  const redirectHref = useMemo(
    () => `/${lang}/results?${searchRequestToParams(searchRequest)}`,
    [lang, searchRequest],
  );

  const redirectToResults = useCallback(() => {
    if (redirectedRef.current) return;
    redirectedRef.current = true;
    router.push(redirectHref);
  }, [router, redirectHref]);

  useEffect(() => {
    if (hasAvailabilityData) return;

    const interval = setInterval(() => {
      setSecondsLeft((previous) => {
        if (previous <= 1) {
          clearInterval(interval);
          redirectToResults();
          return 0;
        }
        return previous - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [hasAvailabilityData, redirectToResults]);

  if (hasAvailabilityData) {
    return <>{children}</>;
  }

  return (
    <>
      {children}
      <Dialog open>
        <DialogContent
          className="min-w-96 max-w-md p-6 flex flex-col gap-4 bg-white border-border-light/50 rounded-2xl shadow-modal"
          showCloseButton={false}
          onEscapeKeyDown={(event) => event.preventDefault()}
          onPointerDownOutside={(event) => event.preventDefault()}
        >
          <DialogTitle className="type-h5 text-navy">{t("title")}</DialogTitle>
          <p className="type-paragraph text-text-secondary">{t("message")}</p>
          <p className="type-paragraph text-text-secondary">
            {t("redirecting", { seconds: secondsLeft })}
          </p>

          <Button
            onClick={redirectToResults}
            className="bg-navy text-white hover:bg-navy/90 rounded-xl"
          >
            {t("redirectNow")}
          </Button>
        </DialogContent>
      </Dialog>
    </>
  );
}
