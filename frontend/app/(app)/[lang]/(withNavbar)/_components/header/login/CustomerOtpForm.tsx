"use client";

import { REGEXP_ONLY_DIGITS } from "input-otp";
import { AlertCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { useMutation } from "@tanstack/react-query";
import { sendOTP, validateOTP } from "@/shared/api/accounts-api";
import useAuthStore, { UserRole } from "@/shared/auth/authStore";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";

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
  const store = useAuthStore();
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

  const {
    mutate: login,
    isPending,
    error,
  } = useMutation({
    mutationFn: () => validateOTP(phone, otp),
    onSuccess: (response) => {
      store.setSession(response.accessToken, response.accessTokenExpiresAt, {
        id: response.id,
        email: response.email,
        firstName: response.firstName,
        lastName: response.lastName,
        role: response.role as UserRole,
        phoneNumber: response.phoneNumber,
        officeId: response.officeId,
        isAdminAsAgent: false,
      });
      onSuccess();
    },
  });

  const translatedError = useTranslatedError(error)

  const { mutate: resendOtp, isPending: resendPending } = useMutation({
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
            {translatedError}
          </span>
        </div>
      )}

      <Button
        variant="brand"
        className="w-full py-3.5 h-auto"
        disabled={isPending || otp.length < 6}
        loading={isPending}
        onClick={() => login()}
      >
        {t("customer.confirmCode")}
      </Button>
    </div>
  );
}
