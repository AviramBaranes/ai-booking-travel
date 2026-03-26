import "../globals.css";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

import { authOptions } from "@/shared/auth/authOptions";
import Providers from "@/app/providers";
import AdminShell from "./AdminShell";

export default async function AdminRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect("/he/");
  }

  const isAdmin = session?.user?.role === "admin";
  if (!isAdmin) {
    redirect("/he/");
  }

  return (
    <html lang="he" dir="rtl" className="h-full antialiased">
      <body className="h-full">
        <Providers>
          <AdminShell>{children}</AdminShell>
        </Providers>
      </body>
    </html>
  );
}
