import "../globals.css";
import { getPayload } from "payload";
import config from "@payload-config";
import { NotFoundContent } from "@/shared/components/NotFoundContent";
import { Footer } from "./[lang]/_components/footer/Footer";
import { Navbar } from "./[lang]/_components/header/navbar/Navbar";
import localFont from "next/font/local";

const polin = localFont({
  src: [
    { path: "../fonts/Polin-Regular.otf", weight: "400" },
    { path: "../fonts/Polin-SemiBold.otf", weight: "600" },
    { path: "../fonts/Polin-Bold.otf", weight: "700" },
    { path: "../fonts/Polin-Black.otf", weight: "900" },
  ],
  variable: "--font-polin",
  display: "swap",
});

export async function getNotFoundData(lang: string) {
  const payload = await getPayload({ config });
  const notFoundData = await payload.findGlobal({
    slug: "not-found",
    locale: lang as "he" | "en",
    draft: false,
  });

  return notFoundData;
}

export default async function NotFound() {
  const notFoundData = await getNotFoundData("he");

  return (
    <html
      lang="he"
      dir="rtl"
      className={`h-full antialiased ${polin.variable}`}
    >
      <body>
        <Navbar lang="he" isAuthenticated={true} isRootLayout />
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
