import { Homepage } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Image from "next/image";
import { SearchForm } from "./SearchForm/SearchForm";

interface Props {
  title: string;
  subtitle: string;
  image: Populated<Homepage["featuredImage"]>;
}
export function Hero({ title, subtitle, image }: Props) {
  if (!image?.url) return null;

  return (
    <section className="relative">
      <Image
        src={image.url}
        alt={image.alt}
        width={image.width ?? 1200}
        height={image.height ?? 630}
        style={{ width: "100%", height: "auto" }}
        priority
      />
      <div className="w-full absolute top-38 -translate-x-1/2 left-1/2">
        <h1 className="text-[55px] text-center type-h1 text-white">{title}</h1>
        <h6 className="text-center mt-2 type-h6 text-white">{subtitle}</h6>
        <SearchForm />
      </div>
    </section>
  );
}
