import { Homepage } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Image from "next/image";
import { SearchForm } from "./SearchForm/SearchForm";
import { AppProviders } from "../../../_components/providers/AppProviders";
import { getMessages } from "next-intl/server";
import { QueryProvider } from "../../../_components/providers/QueryProvider";
import { NextIntlClientProvider } from "next-intl";

interface Props {
  lang: string;
  title: string;
  subtitle: string;
  image: Populated<Homepage["featuredImage"]>;
}
export async function Hero({ lang, title, subtitle, image }: Props) {
  if (!image?.url) return null;
  const messages = await getMessages({ locale: lang });

  return (
    <section className="relative">
      <Image
        src={image.url}
        alt={image.alt}
        width={image.width ?? 1200}
        height={image.height ?? 630}
        sizes="(max-width: 768px) 390px, 100vw"
        className="w-full h-177.5 object-cover object-top lg:object-center md:h-auto"
        priority
      />
      <div className="w-full absolute top-14 lg:top-38 -translate-x-1/2 left-1/2">
        <h1 className="text-[28px] lg:text-[55px] text-center type-h1 text-white">
          {title}
        </h1>
        <h6 className="text-center mt-2 type-h6 text-white">{subtitle}</h6>

        <QueryProvider showDevtools={false}>
          <NextIntlClientProvider locale={lang} messages={messages}>
            <SearchForm />
          </NextIntlClientProvider>
        </QueryProvider>
      </div>
    </section>
  );
}
