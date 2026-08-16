import Image from "next/image";
import { getTranslations } from "next-intl/server";

const WHATSAPP_NUMBER = "97235555999";
const WHATSAPP_URL = `https://wa.me/${WHATSAPP_NUMBER}`;

export async function WhatsAppButton({ lang }: { lang: string }) {
  const t = await getTranslations({ locale: lang, namespace: "WhatsApp" });
  const label = t("tooltip");

  return (
    <a
      href={WHATSAPP_URL}
      target="_blank"
      rel="noopener noreferrer"
      aria-label={label}
      className="group fixed bottom-5 inset-e-5 z-50 flex h-14 w-14 items-center justify-center rounded-full bg-white shadow-lg transition-transform hover:scale-105 print:hidden"
    >
      {/* Served unoptimized at its native 115px into a 56px box: the source is
          already tiny, and re-encoding it to a 56px WebP softens the edges. */}
      <Image
        src="/assets/icons/Whatsapp.png"
        alt=""
        width={115}
        height={115}
        unoptimized
        className="h-14 w-14 rounded-full"
      />
      <span
        role="tooltip"
        className="pointer-events-none absolute end-full top-1/2 me-3 -translate-y-1/2 whitespace-nowrap rounded-md bg-[#25D366] px-3 py-1.5 text-sm font-medium text-white opacity-0 shadow-md transition-opacity duration-150 group-hover:opacity-100 group-focus-visible:opacity-100"
      >
        {label}
      </span>
    </a>
  );
}
