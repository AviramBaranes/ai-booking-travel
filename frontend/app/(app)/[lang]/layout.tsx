import "@/app/globals.css";
import localFont from "next/font/local";
import { getServerSession } from "next-auth/next";
import { redirect } from "next/dist/client/components/navigation";
import { notFound } from "next/navigation";
import { authOptions } from "@/shared/auth/authOptions";
import { Navbar } from "./_components/header/navbar/Navbar";
import { Footer } from "./_components/footer/Footer";
import { BackToAdminBanner } from "./_components/header/login/BackToAdmin/BackToAdminBanner";

const polin = localFont({
  src: [
    { path: "../../fonts/Polin-Regular.otf", weight: "400" },
    { path: "../../fonts/Polin-SemiBold.otf", weight: "600" },
    { path: "../../fonts/Polin-Bold.otf", weight: "700" },
    { path: "../../fonts/Polin-Black.otf", weight: "900" },
  ],
  variable: "--font-polin",
  display: "swap",
});

export default async function AppRootLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const session = await getServerSession(authOptions);
  const isAdmin = session?.user?.role === "admin";
  if (isAdmin) {
    redirect("/admin/");
  }
  const { lang } = await params;

  if (!["he", "en"].includes(lang)) {
    notFound();
  }

  return (
    <html
      lang={lang}
      dir={lang === "he" ? "rtl" : "ltr"}
      className={`h-full antialiased ${polin.variable}`}
    >
      <body className={`min-h-full flex flex-col`}>
        <BackToAdminBanner />
        <Navbar lang={lang} />
        {children}
        <Footer lang={lang} />
      </body>
    </html>
  );
}
