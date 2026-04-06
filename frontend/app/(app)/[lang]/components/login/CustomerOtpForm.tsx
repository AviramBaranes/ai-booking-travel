"use client";

import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";

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
}

export function CustomerOtpForm({ phone }: Props) {
  const t = useTranslations("Login");
  const [otp, setOtp] = useState("");
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

  return (
    <div className="flex flex-col gap-6 items-center w-full">
      <InputOTP maxLength={6} dir="ltr" value={otp} onChange={setOtp}>
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
        <button
          onClick={startResendTimer}
          className="type-paragraph text-navy text-end underline whitespace-nowrap cursor-pointer"
        >
          {t("agent.otpResend")}
        </button>
      )}

      <Button
        variant="brand"
        className="w-full py-3.5 h-auto"
        disabled={otp.length < 6}
      >
        {t("customer.confirmCode")}
      </Button>
    </div>
  );
}
