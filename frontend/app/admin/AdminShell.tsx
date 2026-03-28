"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { signOut } from "next-auth/react";
import {
  Home,
  Network,
  Building2,
  Users,
  Receipt,
  Ticket,
  Coins,
  MapPin,
  CalendarCheck,
  BarChart3,
  Languages,
  LogOut,
} from "lucide-react";

const navItems = [
  { label: "ראשי", href: "/admin", icon: Home },
  { label: "רשתות", href: "/admin/networks", icon: Network },
  { label: "משרדים", href: "/admin/offices", icon: Building2 },
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

  const isActive = (href: string) => {
    if (href === "/admin") return pathname === "/admin";
    return pathname.startsWith(href);
  };

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <aside className="w-56 shrink-0 bg-white border-l border-gray-200 flex flex-col shadow-sm">
        <div className="px-5 py-4 border-b border-gray-200">
          <Image
            src="/logo.png"
            alt="AIBookingTravel"
            width={160}
            height={40}
          />
        </div>

        <nav className="flex-1 overflow-y-auto py-3">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 px-5 py-2.5 text-sm transition-colors ${
                  active
                    ? "bg-blue-50 text-blue-700 font-semibold border-l-3 border-blue-600"
                    : "text-gray-700 hover:bg-gray-50 hover:text-gray-900"
                }`}
              >
                <Icon
                  size={18}
                  className={active ? "text-blue-600" : "text-gray-400"}
                />
                {item.label}
              </Link>
            );
          })}
        </nav>
      </aside>

      {/* Main area */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Header */}
        <header className="h-14 shrink-0 bg-white border-b border-gray-200 flex items-center justify-between px-6 shadow-sm">
          <h1 className="text-lg font-semibold text-gray-800">ניהול מערכת</h1>
          <button
            onClick={() => signOut({ callbackUrl: "/he/" })}
            className="flex items-center gap-2 text-sm text-gray-600 hover:text-red-600 transition-colors cursor-pointer"
          >
            <LogOut size={16} />
            התנתק
          </button>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto p-6 bg-background">
          {children}
        </main>
      </div>
    </div>
  );
}
