"use client";

import AdminNavbar from "@/shared/components/admin/AdminNavbar";

export default function AccountingShell({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <AdminNavbar hideLinks />
      <main className="flex-1 overflow-y-auto p-6 bg-background">{children}</main>
    </div>
  );
}
