"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { RentalPriceForDays } from "../../_components/RentalPriceForDays";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { booking } from "@/shared/client";
import { formatPrice } from "@/shared/utils/formatPrice";
import { CircleCheckBig, Circle, ShieldCheck, X, Crown } from "lucide-react";

interface OtherPlansButtonProps {
  plans: booking.Plan[];
  selectedPlan: number;
  onSelectPlan: (index: number) => void;
  currency: string;
  daysCount: number;
}

export function OtherPlansButton({
  plans,
  selectedPlan,
  onSelectPlan,
  currency,
  daysCount,
}: OtherPlansButtonProps) {
  const t = useTranslations("booking.plansDialog");
  const [open, setOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<"terms" | "inclusions">("terms");

  const handleSelectPlan = (index: number) => {
    onSelectPlan(index);
    setOpen(false);
  };

  const dialogWidthClass =
    plans.length >= 3 ? "!max-w-[1260px]" : "!max-w-[970px]";

  return (
    <>
      <Button
        variant="outline"
        className="px-8 py-6 type-paragraph text-navy rounded-lg"
        onClick={() => setOpen(true)}
      >
        {t("otherPlans")}
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          className={`${dialogWidthClass} p-0 bg-background border-none rounded-3xl shadow-[0px_4px_24px_0px_rgba(0,0,0,0.25)] overflow-hidden ring-0`}
          showCloseButton={false}
        >
          {/* Header — RTL order: title+icon on the right (first in DOM), X on the left (last in DOM) */}
          <div className="flex items-center justify-between px-12 pt-12 pb-0">
            <DialogTitle className="flex items-center gap-4">
              <ShieldCheck className="w-8 h-8 text-success" />
              <span className="type-h5 text-navy">{t("title")}</span>
            </DialogTitle>
            <button
              onClick={() => setOpen(false)}
              className="p-2 cursor-pointer"
            >
              <X className="w-6 h-6 text-navy" />
            </button>
          </div>

          {/* Divider + Tabs */}
          <div className="flex flex-col items-center gap-6 px-12">
            <div className="w-full border-t border-border-light" />
            <div className="flex gap-6 items-center text-[22px] leading-[30.8px] text-brand-blue">
              <button
                className={`cursor-pointer underline ${
                  activeTab === "inclusions" ? "font-bold" : "font-normal"
                }`}
                onClick={() => setActiveTab("inclusions")}
              >
                {t("whatsIncluded")}
              </button>
              <button
                className={`cursor-pointer underline ${
                  activeTab === "terms" ? "font-bold" : "font-normal"
                }`}
                onClick={() => setActiveTab("terms")}
              >
                {t("rentalTerms")}
              </button>
            </div>
          </div>

          {/* Plan Cards */}
          <div className="flex gap-6 items-stretch justify-center px-12 pb-6">
            {plans.map((plan, index) => {
              const isSelected = index === selectedPlan;
              const items =
                activeTab === "terms" ? plan.info : plan.planInclusions;

              return (
                <div
                  key={plan.planName}
                  className={`flex-1 min-w-0 bg-white rounded-3xl shadow-[0px_4px_12px_0px_rgba(63,63,63,0.1)] overflow-hidden flex flex-col ${
                    isSelected
                      ? "border-3 border-brand"
                      : "border border-border-light"
                  }`}
                >
                  <div className="flex flex-col justify-between flex-1 pt-6 px-6">
                    {/* Plan name + inclusions */}
                    <div className="flex flex-col gap-6">
                      <h5 className="type-h5 text-navy">
                        {index === plans.length - 1 && (
                          <Crown className="w-6 mb-3 h-6 text-brand inline-block m-2" />
                        )}
                        {t(`planNames.${plan.planName}`)}
                      </h5>

                      <div className="flex flex-col gap-3">
                        {items.map((item) => (
                          <div key={item} className="flex gap-2.5 items-start">
                            <CircleCheckBig className="w-4 h-4 text-success shrink-0 mt-1.5" />
                            <span className="type-paragraph text-text-secondary flex-1">
                              {item}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* Price section — RTL order: price on the right (first in DOM), label on the left (last in DOM) */}
                    <div className="flex items-center justify-between mt-8 mb-6">
                      <div className="type-label text-navy text-center flex flex-col items-start">
                        <p>{t("totalToPayLine1")}</p>
                        <p>{t("totalToPayLine2")}</p>
                      </div>
                      <div className="flex flex-col gap-1 items-end">
                        <span className="type-h4 text-brand">
                          {formatPrice(plan.price, currency)}
                        </span>
                        <RentalPriceForDays daysCount={daysCount} />
                      </div>
                    </div>
                  </div>

                  {/* Select button */}
                  <button
                    onClick={() => handleSelectPlan(index)}
                    className={`flex items-center justify-center gap-3 py-6 px-6 cursor-pointer type-paragraph font-bold w-full ${
                      isSelected
                        ? "bg-brand text-white"
                        : "bg-white text-navy shadow-[0px_-4px_8px_0px_rgba(0,0,0,0.08)]"
                    }`}
                  >
                    {isSelected ? (
                      <CircleCheckBig className="w-4 h-4" />
                    ) : (
                      <Circle className="w-4 h-4" />
                    )}
                    <span>
                      {isSelected ? t("selectedPlan") : t("selectThisPlan")}
                    </span>
                  </button>
                </div>
              );
            })}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
