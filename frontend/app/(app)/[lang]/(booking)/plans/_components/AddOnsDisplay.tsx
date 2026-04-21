import { Button } from "@/components/ui/button";
import { Media } from "@/payload-types";
import { broker } from "@/shared/client";
import { Populated } from "@/shared/types/payload";
import { formatPrice } from "@/shared/utils/formatPrice";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { useParams } from "next/navigation";
import { useAddonsGallery } from "@/shared/hooks/useAddonsGallery";
import { Box } from "lucide-react";

interface AddOnsDisplayProps {
  addons: broker.AddOn[];
  selectedAddons: broker.SelectAddOn[];
  setSelectedAddons: (addons: broker.SelectAddOn[]) => void;
}

export function AddOnsDisplay({
  addons,
  selectedAddons,
  setSelectedAddons,
}: AddOnsDisplayProps) {
  const t = useTranslations("booking.addOns");
  const { lang } = useParams();
  const { data: addOnsGallery } = useAddonsGallery();

  function addQuantity(addOnId: number) {
    const existing = selectedAddons.find((a) => a.id === addOnId);
    if (existing) {
      setSelectedAddons(
        selectedAddons.map((a) =>
          a.id === addOnId ? { ...a, quantity: a.quantity + 1 } : a,
        ),
      );
    } else {
      setSelectedAddons([...selectedAddons, { id: addOnId, quantity: 1 }]);
    }
  }

  function subtractQuantity(addOnId: number) {
    const existing = selectedAddons.find((a) => a.id === addOnId);
    if (!existing) return;
    if (existing.quantity === 1) {
      setSelectedAddons(selectedAddons.filter((a) => a.id !== addOnId));
    } else {
      setSelectedAddons(
        selectedAddons.map((a) =>
          a.id === addOnId ? { ...a, quantity: a.quantity - 1 } : a,
        ),
      );
    }
  }

  return (
    <div className="pb-30">
      <h5 className="type-h5 text-navy">{t("title")}</h5>
      <div className="grid grid-cols-3 gap-4 mt-6 ">
        {addons.map((addOn) => {
          const addOnDetails = addOnsGallery.addons?.find(
            (item) => item.addonId === addOn.id.toString(),
          );
          const name =
            (lang === "he"
              ? addOnDetails?.hebrewName
              : addOnDetails?.englishName) ?? addOn.name;

          const image = addOnDetails?.image;
          const media =
            image && typeof image === "object" ? (image as Media) : null;

          const selectedQuantity =
            selectedAddons.find((a) => a.id === addOn.id)?.quantity ?? 0;
          return (
            <div
              key={addOn.id}
              className="bg-white border-border-muted rounded-2xl"
            >
              <div className="bg-white rounded-t-2xl shadow-card">
                {media?.url ? (
                  <Image
                    src={media.url}
                    alt={media?.alt ?? name ?? "add-on image"}
                    width={200}
                    height={200}
                    className="w-50 h-50 mx-auto"
                  />
                ) : (
                  <Box className="w-50 h-50 mx-auto text-muted" />
                )}
              </div>
              <div className="m-6 border-b border-border-muted bg-white pb-6">
                <p className="type-paragraph font-bold text-navy w-2/3">
                  {name}
                </p>
                <p className="type-paragraph font-bold text-navy">
                  {formatPrice(addOn.price, addOn.currency)} {t(addOn.period)}
                </p>
                <p className="type-paragraph text-navy">{t("payAtPickup")}</p>
              </div>
              <div className="flex m-6 justify-between">
                <p className="type-paragraph text-muted">{t("quantity")}</p>
                <div className="border-cars-border border rounded-md">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="px-2"
                    onClick={() => subtractQuantity(addOn.id)}
                    disabled={selectedQuantity === 0}
                  >
                    -
                  </Button>
                  <span className="px-2"> {selectedQuantity} </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="px-2"
                    onClick={() => addQuantity(addOn.id)}
                    disabled={selectedQuantity >= addOn.allowedQuantity}
                  >
                    +
                  </Button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
