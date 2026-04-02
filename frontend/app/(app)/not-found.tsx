import "../globals.css";
import { getPayload } from "payload";
import config from "@payload-config";
import { NotFoundContent } from "@/shared/components/NotFoundContent";
import { Footer } from "./[lang]/components/footer/Footer";
import { Navbar } from "./[lang]/components/navbar/Navbar";

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
    <html lang="he" dir="rtl" className="h-full antialiased">
      <Navbar lang="he" isAuthenticated={true} displayLangSwitcher={false} />
      <NotFoundContent
        title={notFoundData.title ?? ""}
        subtitle={notFoundData.subtitle ?? ""}
        buttonText={notFoundData.buttonText ?? ""}
        homepageUrl="/he"
      />
      <Footer lang="he" />
    </html>
  );
}
