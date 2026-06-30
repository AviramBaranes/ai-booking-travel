import "../globals.css";
import { NotFoundContent } from "@/shared/components/NotFoundContent";
import localFont from "next/font/local";
import { Footer } from "./[lang]/(withNavbar)/_components/footer/Footer";
import { NavbarLinks } from "./[lang]/(withNavbar)/_components/header/navbar/NavbarLinks";
import { getCachedPayload } from "@/shared/server/cms";
import { getHeaderData } from "./[lang]/(withNavbar)/_components/header/navbar/Navbar";

const polin = localFont({
  src: [
    { path: "../fonts/Polin-Regular.otf", weight: "400" },
    { path: "../fonts/Polin-Semibold.otf", weight: "600" },
    { path: "../fonts/Polin-Bold.otf", weight: "700" },
    { path: "../fonts/Polin-Black.otf", weight: "900" },
  ],
  variable: "--font-polin",
  display: "swap",
});

export async function getNotFoundData(lang: string) {
  const payload = await getCachedPayload();
  const notFoundData = await payload.findGlobal({
    slug: "not-found",
    locale: lang as "he" | "en",
    draft: false,
  });

  return notFoundData;
}

export default async function NotFound() {
  const notFoundData = await getNotFoundData("he");
  const headerData = await getHeaderData("he");

  return (
    <html
      lang="he"
      dir="rtl"
      className={`h-full antialiased ${polin.variable}`}
    >
      <body>
        <NavbarLinks lang="he" headerData={headerData}/>
        <NotFoundContent
          title={notFoundData.title ?? ""}
          subtitle={notFoundData.subtitle ?? ""}
          buttonText={notFoundData.buttonText ?? ""}
          homepageUrl="/he"
        />
        <Footer lang="he" />
      </body>
    </html>
  );
}
