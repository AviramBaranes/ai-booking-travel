import { Footer } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Image from "next/image";
import Link from "next/link";

interface FooterSecondFloorProps {
  footerData: Footer;
  lang: string;
}

export function FooterSecondFloor({
  footerData,
  lang,
}: FooterSecondFloorProps) {
  return (
    <div className="flex flex-col lg:flex-row justify-between items-center gap-6 flex-wrap border-b border-light-white bg-navy lg:pl-16 lg:pr-4">
      <div className="flex flex-1 flex-col items-center lg:items-start gap-2 py-7 px-5 mx-auto lg:pt-24 lg:pb-28">
        <Link href={`/${lang}`}>
          <Image
            src="/logo-dark.png"
            alt="AIBookingTravel"
            width={200}
            height={40}
            className="w-50 h-10"
          />
        </Link>
        <p className="text-white type-paragraph font-normal mt-4">
          {footerData.socialsTitle}
        </p>
        <div className="flex gap-2">
          {footerData.socialsLinks?.map((social) => (
            <Link
              key={social.id}
              className="border-light-white rounded-full bg-brand-blue/35 border w-13 h-13 flex items-center justify-center text-medium-white"
              target="_blank"
              href={social.link}
            >
              {social.label}
            </Link>
          ))}
        </div>
      </div>
      {/* Desktop */}
      <div className="hidden flex-5 items-center justify-around lg:flex lg:flex-row">
        {footerData.linkGroups?.map((group) => (
          <div className="flex flex-col" key={group.id}>
            <span className="h-0.5 w-5 border-none bg-brand" />

            <h6 className="mb-4 type-h6 text-white">{group.title}</h6>

            <div className="flex flex-col gap-2">
              {group.links?.map((link) => (
                <Link
                  key={link.id}
                  href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
                  className="type-paragraph text-white opacity-52"
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Mobile */}
      <div className="flex w-full flex-col lg:hidden">
        {footerData.linkGroups?.map((group) => (
          <details
            key={group.id}
            className="group border-b border-white/10 py-4 px-5"
          >
            <summary className="flex cursor-pointer list-none items-center justify-between [&::-webkit-details-marker]:hidden">
              <h6 className="type-h6 text-white">{group.title}</h6>

              <span className="text-2xl leading-none text-brand group-open:hidden">
                +
              </span>

              <span className="hidden text-2xl leading-none text-brand group-open:block">
                −
              </span>
            </summary>

            <div className="mt-4 flex flex-col gap-2">
              {group.links?.map((link) => (
                <Link
                  key={link.id}
                  href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
                  className="type-paragraph text-white opacity-52"
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </details>
        ))}
      </div>
    </div>
  );
}
