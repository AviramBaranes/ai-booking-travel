"use client";

import { useTranslations } from "next-intl";
import { ArrowRight, ArrowLeft } from "lucide-react";
import clsx from "clsx";
import { useParams } from "next/navigation";
import { useDirection } from "@/shared/hooks/useDirection";

const BOOKING_STEPS = [
  { key: "results", path: "results" },
  { key: "plans", path: "plans" },
  { key: "ordering", path: "ordering" },
] as const;

interface BookingStepperProps {
  currentStep: (typeof BOOKING_STEPS)[number]["key"];
}

const ARROW_RIGHT = "\u2192";
const ARROW_LEFT = "\u2190";

export function BookingStepper({ currentStep }: BookingStepperProps) {
  const dir = useDirection();
  const t = useTranslations("booking.steps");

  const currentStepIndex = BOOKING_STEPS.findIndex(
    (step) => step.key === currentStep,
  );

  if (currentStepIndex === -1) return null;

  const handleBackNavigation = (targetStepIndex: number) => {
    const stepsBack = currentStepIndex - targetStepIndex;
    window.history.go(-stepsBack);
  };

  return (
    <div className="flex">
      {BOOKING_STEPS.map((step, i) => (
        <div key={step.key} className="flex items-center">
          {currentStepIndex > i ? (
            <button
              onClick={() => handleBackNavigation(i)}
              className="text-sm text-text-secondary cursor-pointer"
            >
              {`${t("step")} ${i + 1}: ${t(step.key)}`}
            </button>
          ) : (
            <span
              className={clsx(
                "text-sm type-label font-medium text-text-secondary",
                { "font-semibold text-navy!": currentStepIndex === i },
              )}
            >
              {`${t("step")} ${i + 1}: ${t(step.key)}`}
            </span>
          )}
          {i < BOOKING_STEPS.length - 1 &&
            (dir === "rtl" ? (
              <span className="mx-2">{ARROW_LEFT}</span>
            ) : (
              <span className="mx-2">{ARROW_RIGHT}</span>
            ))}
        </div>
      ))}
    </div>
  );
}
