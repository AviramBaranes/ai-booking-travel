"use client";

import AdminNavbar from "@/shared/components/admin/AdminNavbar";
import useAuthStore from "@/shared/auth/authStore";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function AccountingShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const status = useAuthStore((state) => state.status);
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    if (status === "loading" || status === "idle") {
      // Still loading auth
      return;
    }

    if (!user || !["accountant", "admin"].includes(user.role)) {
      // Not authorized, redirect
      router.replace("/he/");
      return;
    }

    setIsAuthorized(true);
  }, [user, status, router]);

  // Don't render content until we verify authorization
  if (!isAuthorized) {
    return null;
  }

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <AdminNavbar hideLinks />
      <main className="flex-1 overflow-y-auto p-6 bg-background">{children}</main>
    </div>
  );
}
