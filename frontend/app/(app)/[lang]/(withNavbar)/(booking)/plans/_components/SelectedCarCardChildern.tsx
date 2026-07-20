import { Button } from "@/components/ui/button";
import { FreeCancellationBadge } from "@/shared/components/booking/FreeCancellationBadge";
import { useDirection } from "@/shared/hooks/useDirection";
import clsx from "clsx";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useSearchRequest } from "../../_hooks/useSearchRequest";
import useAuthStore from "@/shared/auth/authStore";
import { useTranslations } from "next-intl";

interface SelectedCarCardChildrenProps {
  isErpSelected: boolean;
  setIsErpDialogOpen: (open: boolean) => void;
  setIsPriceOfferDialogOpen: (open: boolean) => void;
}

export function SelectedCarCardChildren({
  isErpSelected,
  setIsErpDialogOpen,
  setIsPriceOfferDialogOpen,
}: SelectedCarCardChildrenProps) {
  const t = useTranslations("booking.plansPage");
  const { lang } = useParams();
  const dir = useDirection();
  const router = useRouter();
  const { searchRequest } = useSearchRequest();
  const currentSearchParams = useSearchParams();
  const user = useAuthStore((s) => s.user);
  const isAgent = user?.role === "agent";

  function handleContinue() {
    if (isErpSelected) {
      router.push(`/${lang}/order?${currentSearchParams.toString()}`);
    } else {
      setIsErpDialogOpen(true);
    }
  }

  function handleCreatePriceOffer() {
    setIsPriceOfferDialogOpen(true);
  }
  
  return (
    <>
    {/* Desktop */}
      <div className="hidden lg:contents">
        <FreeCancellationBadge
          pickupDate={searchRequest.PickupDate}
          pickupTime={searchRequest.PickupTime}
          text={t("freeCancellation")}
        />
        <Button
          variant="brand"
          className="mt-4 type-paragraph font-bold py-6 px-8 cursor-pointer"
          onClick={handleContinue}
        >
          {t("continueCta")}
        </Button>
        {isAgent && (
          <Button
            variant="brand"
            className="type-paragraph font-bold py-6 px-8 cursor-pointer bg-navy"
            onClick={handleCreatePriceOffer}
          >
            {t("createPriceOffer")}
          </Button>
        )}
      </div>
      
      {/* Mobile */}
      <div className="lg:hidden absolute bottom-0 w-full right-0 left-0">
        <div className="w-fit mx-auto text-center">
          <FreeCancellationBadge
            pickupDate={searchRequest.PickupDate}
            pickupTime={searchRequest.PickupTime}
            text={t("freeCancellation")}
          />
        </div>
        {isAgent ? (
          <div className="w-full">
            <Button
              variant="brand"
              className={clsx(
                "type-paragraph font-bold px-8 cursor-pointer border border-navy bg-navy w-1/2 rounded-none!",
                {
                  "rounded-br-md!": dir === "rtl",
                  "rounded-bl-md!": dir === "ltr",
                },
              )}
              onClick={handleCreatePriceOffer}
            >
              {t("createPriceOffer")}
            </Button>

            <Button
              variant="brand"
              className={clsx(
                "type-paragraph font-bold px-8 cursor-pointer border border-brand w-1/2 rounded-none!",
                {
                  "rounded-bl-md!": dir === "rtl",
                  "rounded-br-md!": dir === "ltr",
                },
              )}
              onClick={handleContinue}
            >
              {t("continueCta")}
            </Button>
          </div>
        ) : (
          <Button
            variant="brand"
            className="type-paragraph font-bold w-full px-8 cursor-pointer rounded-t-none border border-brand"
            onClick={handleContinue}
          >
            {t("continueCta")}
          </Button>
        )}
      </div>
    </>
  );
}
