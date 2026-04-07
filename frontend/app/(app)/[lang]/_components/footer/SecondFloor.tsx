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
    <div className="flex justify-between items-center gap-6 flex-wrap border-b border-light-white bg-navy pt-24 pb-28 pl-16 pr-4">
      <div className="flex flex-1 flex-col gap-2">
        <Link href={`/${lang}`}>
          <Image
            src="/logo-dark.png"
            alt="AIBookingTravel"
            width={200}
            height={50}
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
      <div className="flex justify-around items-center flex-5">
        {footerData.linkGroups?.map((group) => (
          <div className="flex flex-col" key={group.id}>
            <span className="bg-brand w-5 h-0.5 border-none"></span>
            <h6 className="text-white type-h6 mb-4">{group.title}</h6>
            <div className="flex flex-col gap-2">
              {group.links?.map((link) => (
                <Link
                  key={link.id}
                  href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
                  className="text-white opacity-52 type-paragraph"
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
