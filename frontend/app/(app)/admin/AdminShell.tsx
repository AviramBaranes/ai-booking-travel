"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import clsx from "clsx";
import {
  Home,
  Network,
  Building2,
  UserCog,
  Users,
  Contact,
  Receipt,
  Ticket,
  Coins,
  MapPin,
  CalendarCheck,
  BarChart3,
  Languages,
  Landmark,
} from "lucide-react";
import AdminNavbar from "@/shared/components/admin/AdminNavbar";

const navItems = [
  { label: "ראשי", href: "/admin", icon: Home },
  { label: "מנהלים", href: "/admin/admins", icon: UserCog },
  { label: "רשתות", href: "/admin/organizations", icon: Landmark },
  { label: "משרדים", href: "/admin/offices", icon: Building2 },
  { label: "אנשי קשר", href: "/admin/contacts", icon: Contact },
  { label: "סוכנים", href: "/admin/agents", icon: Users },
  { label: "מחירונים", href: "/admin/pricing", icon: Receipt },
  { label: "קופונים", href: "/admin/coupons", icon: Ticket },
  { label: "מטבעות", href: "/admin/currencies", icon: Coins },
  { label: "מיקומים", href: "/admin/locations", icon: MapPin },
  { label: "תרגומים", href: "/admin/translations", icon: Languages },
  { label: "הזמנות", href: "/admin/bookings", icon: CalendarCheck },
  { label: "דוחות", href: "/admin/reports", icon: BarChart3 },
];

export default function AdminShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  const isAsideActive = (href: string) => {
    if (href === "/admin") return pathname === "/admin";
    return pathname.startsWith(href);
  };

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <div className="flex flex-1 flex-col overflow-hidden">
        <AdminNavbar />
        <main className="flex-1 overflow-y-auto p-6 bg-background pr-60">
          {children}
        </main>
      </div>
      {/* Sidebar */}
      <aside className="w-56 shrink-0 bg-white border-l border-gray-200 fixed top-14 bottom-0 shadow-sm">
        <nav className="flex flex-col overflow-y-auto py-3">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isAsideActive(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={clsx(
                  "flex items-center gap-3 px-5 py-2.5 text-sm transition-colors",
                  active
                    ? "bg-blue-50 text-blue-700 font-semibold border-l-3 border-blue-600"
                    : "text-gray-700 hover:bg-gray-50 hover:text-gray-900",
                )}
              >
                <Icon
                  size={18}
                  className={clsx(active ? "text-blue-600" : "text-gray-400")}
                />
                {item.label}
              </Link>
            );
          })}
        </nav>
      </aside>
    </div>
  );
}
