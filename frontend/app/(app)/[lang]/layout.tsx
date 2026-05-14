import "@/app/globals.css";
import localFont from "next/font/local";
import { getServerSession } from "next-auth/next";
import { redirect } from "next/dist/client/components/navigation";
import { notFound } from "next/navigation";
import { authOptions } from "@/shared/auth/authOptions";
import { Toaster } from "@/components/ui/sonner";

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

export default async function LangRootLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const session = await getServerSession(authOptions);
  const role = session?.user?.role;
  if (role === "admin") {
    redirect("/admin/");
  }
  if (role === "accountant") {
    redirect("/accounting/billing");
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
        {children}
        <Toaster />
      </body>
    </html>
  );
}
