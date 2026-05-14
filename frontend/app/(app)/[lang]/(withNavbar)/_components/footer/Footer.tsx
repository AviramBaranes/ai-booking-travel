import { getPayload } from "payload";
import config from "@payload-config";
import { FooterFirstFloor } from "./FirstFloor";
import { FooterSecondFloor } from "./SecondFloor";
import { FooterThirdFloor } from "./ThirdFloor";

interface FooterProps {
  lang: string;
}

async function getFooterContent(lang: string) {
  const payload = await getPayload({ config });
  const footerContent = await payload.findGlobal({
    slug: "footer",
    locale: lang as "he" | "en",
    draft: false,
  });

  return footerContent;
}

export async function Footer({ lang }: FooterProps) {
  const footerContent = await getFooterContent(lang);

  return (
    <>
      <FooterFirstFloor links={footerContent.firstFloorLinks} lang={lang} />
      <FooterSecondFloor footerData={footerContent} lang={lang} />
      <FooterThirdFloor
        lang={lang}
        links={footerContent.thirdFloorLinks}
        rightsText={footerContent.rights ?? ""}
      />
    </>
  );
}
