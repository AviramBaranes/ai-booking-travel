import { getTranslations } from "next-intl/server";
import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { getLang } from "@/shared/lang/lang";
import clsx from "clsx";

const BOOKING_STEPS = [
  { key: "results", path: "results" },
  { key: "addOns", path: "add-ons" },
  { key: "ordering", path: "ordering" },
] as const;

interface BookingStepperProps {
  currentStep: (typeof BOOKING_STEPS)[number]["key"];
}

export async function BookingStepper({ currentStep }: BookingStepperProps) {
  const t = await getTranslations("booking.steps");
  const lang = await getLang();

  const currentStepIndex = BOOKING_STEPS.findIndex(
    (step) => step.key === currentStep,
  );

  if (currentStepIndex === -1) return null;

  return (
    <div className="flex">
      {BOOKING_STEPS.map((step, i) => (
        <div key={step.key} className="flex items-center">
          {currentStepIndex > i ? (
            <Link
              href={`/${step.path}`}
              className="text-sm text-text-secondary"
            >
              {`${t("step")} ${i + 1}: ${t(step.key)}`}
            </Link>
          ) : (
            <span
              className={clsx(
                "text-sm type-label font-medium text-text-secondary",
                {
                  "font-semibold text-navy!": currentStepIndex === i,
                },
              )}
            >
              {`${t("step")} ${i + 1}: ${t(step.key)}`}
            </span>
          )}
          {i < BOOKING_STEPS.length - 1 &&
            (lang === "he" ? (
              <ArrowLeft className="mx-2" />
            ) : (
              <ArrowRight className="mx-2" />
            ))}
        </div>
      ))}
    </div>
  );
}
