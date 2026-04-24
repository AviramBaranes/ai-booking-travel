"use client";

import { REGEXP_ONLY_DIGITS } from "input-otp";
import { AlertCircle } from "lucide-react";
import { signIn } from "next-auth/react";
import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { Loading } from "@/shared/components/Loading";
import { useMutation } from "@tanstack/react-query";
import { sendOTP } from "@/shared/api/accounts-api";

const RESEND_COUNTDOWN = 45;

function formatTimer(seconds: number) {
  const m = Math.floor(seconds / 60)
    .toString()
    .padStart(2, "0");
  const s = (seconds % 60).toString().padStart(2, "0");
  return `${m}:${s}`;
}

interface Props {
  phone: string;
  onSuccess: () => void;
}

export function CustomerOtpForm({ phone, onSuccess }: Props) {
  const t = useTranslations("Login");
  const tError = useTranslations("ApiErrors");
  const [otp, setOtp] = useState("");
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [resendTimer, setResendTimer] = useState(RESEND_COUNTDOWN);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const startResendTimer = () => {
    setResendTimer(RESEND_COUNTDOWN);
    if (timerRef.current) clearInterval(timerRef.current);
    timerRef.current = setInterval(() => {
      setResendTimer((prev) => {
        if (prev <= 1) {
          clearInterval(timerRef.current!);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  useEffect(() => {
    startResendTimer();
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSubmit = async () => {
    setError(null);
    setIsPending(true);
    try {
      const phoneNumber = phone.replace(/[\s-]/g, "");
      const result = await signIn("customer-login", {
        redirect: false,
        phoneNumber,
        otp,
      });
      const res = result as { error?: string } | undefined;
      if (res?.error) throw new Error(res.error ?? "unknown_error");
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown_error");
    } finally {
      setIsPending(false);
    }
  };

  const {
    mutate: resendOtp,
    isPending: resendPending,
  } = useMutation({
    mutationFn: async () => sendOTP({ phoneNumber: phone }),
    onSuccess: () => {
      startResendTimer();
    },
  });

  return (
    <div className="flex flex-col gap-6 items-center w-full">
      <InputOTP
        maxLength={6}
        dir="ltr"
        value={otp}
        onChange={setOtp}
        pattern={REGEXP_ONLY_DIGITS}
      >
        <InputOTPGroup className="gap-2" dir="ltr">
          {Array.from({ length: 6 }).map((_, i) => (
            <InputOTPSlot
              key={i}
              index={i}
              className="size-11.5 rounded-lg border-2 border-border-light bg-white shadow-subtle type-h6 text-foreground data-[active=true]:border-brand/80"
            />
          ))}
        </InputOTPGroup>
      </InputOTP>

      {resendTimer > 0 ? (
        <p className="type-paragraph text-text-secondary text-end underline whitespace-nowrap">
          {t("customer.otpResend", { time: formatTimer(resendTimer) })}
        </p>
      ) : (
        <Button
          variant="ghost"
          loading={resendPending}
          onClick={() => resendOtp()}
          className="type-paragraph text-navy text-end underline whitespace-nowrap cursor-pointer"
        >
          {t("agent.otpResend")}
        </Button>
      )}

      {error && (
        <div role="alert" className="flex items-center gap-1 w-full">
          <AlertCircle className="size-3.5 text-destructive shrink-0" />
          <span className="type-paragraph text-destructive">
            {tError(error)}
          </span>
        </div>
      )}

      <Button
        variant="brand"
        className="w-full py-3.5 h-auto"
        disabled={isPending || otp.length < 6}
        onClick={handleSubmit}
      >
        {isPending && <Loading />}
        {t("customer.confirmCode")}
      </Button>
    </div>
  );
}
