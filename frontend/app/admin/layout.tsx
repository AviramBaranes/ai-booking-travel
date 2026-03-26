import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

import { authOptions } from "@/shared/auth/authOptions";

export default async function AdminRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession(authOptions);
  console.log("logging the session from the admin layout:");
  console.log({ session });
  if (!session) {
    redirect("/he/");
  }

  const isAdmin = session?.user?.role === "admin";
  if (!isAdmin) {
    redirect("/he/");
  }

  return (
    <html lang="he" dir="rtl" className="h-full antialiased">
      <body>{children}</body>
    </html>
  );
}
