import { broker } from "@/shared/client";
import { SupplierLogo } from "../SupplierLogo";
import Image from "next/image";
import { CarDetailsPills } from "../CarDetailsPills";
import { useTranslations } from "next-intl";
import { Suspense } from "react";
import { Loading } from "../../Loading";

export function SelectedCarHeader({
  carDetails,
}: {
  carDetails: broker.CarDetails;
}) {
  const t = useTranslations("booking.results");
  return (
    <>
      <div className="flex-col flex items-center">
        <div className="mb-12">
          <Suspense fallback={<Loading />}>
            <SupplierLogo supplierName={carDetails.supplierName} />
          </Suspense>
        </div>

        <Image
          src={carDetails.imageUrl}
          alt={carDetails.model}
          width={176}
          height={100}
          className="w-44 h-25 object-cover"
        />
      </div>
      <div className="flex gap-2 flex-col items-start">
        <h5 className="type-h5 text-navy">
          {carDetails.model} ({carDetails.acriss})
        </h5>
        <span className="type-paragraph text-navy">
          {t("carDetails.orSimilar")}
        </span>
      </div>
      <div className="mb-6">
        <CarDetailsPills carDetails={carDetails} />
      </div>
    </>
  );
}
