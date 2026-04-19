import { useTranslations } from "next-intl";
import Image from "next/image";
import { OrderSummarySubTitle } from "./OrderSummarySubTitle";

export function IncludedSection({
  planInclusions,
}: {
  planInclusions: string[];
}) {
  const t = useTranslations("MyAccount.reservation.summary");

  return (
    <>
      <OrderSummarySubTitle title={t("sections.included")} />
      <ul>
        {planInclusions.map((inclusion) => (
          <li
            key={inclusion}
            className="type-paragraph text-navy mx-4 my-2 flex"
          >
            <Image
              src="/assets/icons/V.svg"
              alt={t("includedIconAlt")}
              width={28}
              height={28}
              className="inline w-7 h-7"
            />
            {inclusion}
          </li>
        ))}
      </ul>
    </>
  );
}
