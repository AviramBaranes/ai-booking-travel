"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { booking } from "@/shared/client";
import { bookingKeys } from "@/shared/hooks/useAvailableCars";
import { searchRequestToParams } from "../results/searchQuery";

const REDIRECT_SECONDS = 5;
const SEARCH_TTL_MS = 15 * 60 * 1000; // 15 minutes

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
  const { lang } = useParams<{ lang: string }>();

  const [isExpired, setIsExpired] = useState(false);
  const [secondsLeft, setSecondsLeft] = useState(REDIRECT_SECONDS);
  const [isTabFocused, setIsTabFocused] = useState(true);
  const redirectedRef = useRef(false);
  const pendingRedirectRef = useRef(false);

  const redirectHref = useMemo(
    () => `/${lang}/results?${searchRequestToParams(searchRequest)}`,
    [lang, searchRequest],
  );

  const redirectToResults = useCallback(() => {
    if (redirectedRef.current) return;
    redirectedRef.current = true;
    location.href = redirectHref;
  }, [redirectHref]);

  const redirectWhenFocused = useCallback(() => {
    if (document.hidden || !document.hasFocus()) {
      pendingRedirectRef.current = true;
      return;
    }

    redirectToResults();
  }, [redirectToResults]);

  useEffect(() => {
    const handleVisibilityOrFocus = () => {
      const focused = !document.hidden && document.hasFocus();
      setIsTabFocused(focused);

      if (focused && pendingRedirectRef.current) {
        pendingRedirectRef.current = false;
        redirectToResults();
      }
    };

    handleVisibilityOrFocus();

    document.addEventListener("visibilitychange", handleVisibilityOrFocus);
    window.addEventListener("focus", handleVisibilityOrFocus);
    window.addEventListener("blur", handleVisibilityOrFocus);

    return () => {
      document.removeEventListener("visibilitychange", handleVisibilityOrFocus);
      window.removeEventListener("focus", handleVisibilityOrFocus);
      window.removeEventListener("blur", handleVisibilityOrFocus);
    };
  }, [redirectToResults]);

  // On mount: compute time remaining based on when data was originally fetched.
  // If there's no data at all, redirect immediately.
  useEffect(() => {
    const queryKey = bookingKeys.availability(searchRequest);
    const state = queryClient.getQueryState(queryKey);

    if (!state?.dataUpdatedAt) {
      redirectWhenFocused();
      return;
    }

    const msRemaining = SEARCH_TTL_MS - (Date.now() - state.dataUpdatedAt);

    if (msRemaining <= 0) {
      // Already expired — evict the stale data so we don't show wrong results,
      // then start the redirect countdown.
      queryClient.removeQueries({ queryKey });
      setIsExpired(true);
      return;
    }

    // Schedule expiry exactly when the TTL elapses.
    const expiryTimer = setTimeout(() => {
      queryClient.removeQueries({ queryKey });
      setIsExpired(true);
    }, msRemaining);

    return () => clearTimeout(expiryTimer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryClient, redirectWhenFocused, searchRequest]);

  // Start redirect countdown once expired.
  useEffect(() => {
    if (!isExpired || !isTabFocused) return;

    const interval = setInterval(() => {
      setSecondsLeft((previous) => {
        if (previous <= 1) {
          clearInterval(interval);
          redirectWhenFocused();
          return 0;
        }
        return previous - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [isExpired, isTabFocused, redirectWhenFocused]);

  if (!isExpired) {
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
