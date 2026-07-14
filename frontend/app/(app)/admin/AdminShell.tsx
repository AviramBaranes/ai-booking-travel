"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import clsx from "clsx";
import {
  Home,
  Calculator,
  Building2,
  UserCog,
  Users,
  Contact,
  Receipt,
  Ticket,
  MapPin,
  CalendarCheck,
  TrendingUp,
  Banknote,
  Languages,
  Landmark,
  Menu,
} from "lucide-react";
import AdminNavbar from "@/shared/components/admin/AdminNavbar";
import {
  Sheet,
  SheetContent,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import useAuthStore from "@/shared/auth/authStore";
import { useEffect, useState } from "react";

const navItems = [
  { label: "ראשי", href: "/admin", icon: Home },
  { label: "מנהלים", href: "/admin/admins", icon: UserCog },
  { label: "רואי חשבון", href: "/admin/accountants", icon: Calculator },
  { label: "רשתות", href: "/admin/organizations", icon: Landmark },
  { label: "משרדים", href: "/admin/offices", icon: Building2 },
  { label: "אנשי קשר", href: "/admin/contacts", icon: Contact },
  { label: "סוכנים", href: "/admin/agents", icon: Users },
  { label: "לקוחות פרטיים", href: "/admin/customers", icon: Users },
  { label: "מחירונים", href: "/admin/pricing", icon: Receipt },
  { label: "קופונים", href: "/admin/coupons", icon: Ticket },
  // { label: "מטבעות", href: "/admin/currencies", icon: Coins },
  { label: "מיקומים", href: "/admin/locations", icon: MapPin },
  { label: "תרגומים", href: "/admin/translations", icon: Languages },
  { label: "דוח הזמנות עסקי", href: "/admin/reports/reservations", icon: CalendarCheck },
  { label: "דוח רווחיות", href: "/admin/reports/profitability", icon: TrendingUp },
  { label: "דוח גבייה", href: "/admin/reports/collections", icon: Banknote },
];

export default function AdminShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const status = useAuthStore((state) => state.status);
  const [mobileNavOpen, setMobileNavOpen] = useState(false);
  const isAuthorized =
    status !== "loading" && status !== "idle" && user?.role === "admin";

  useEffect(() => {
    if (status === "loading" || status === "idle") {
      return;
    }

    if (!isAuthorized) {
      router.replace("/he/");
    }
  }, [isAuthorized, router, status]);

  const isAsideActive = (href: string) => {
    if (href === "/admin") return pathname === "/admin";
    return pathname.startsWith(href);
  };

  if (!isAuthorized) {
    return null;
  }

  const renderNavLinks = (onNavigate?: () => void) =>
    navItems.map((item) => {
      const Icon = item.icon;
      const active = isAsideActive(item.href);
      return (
        <Link
          key={item.href}
          href={item.href}
          onClick={onNavigate}
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
    });

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <div className="flex flex-1 flex-col overflow-hidden">
        <AdminNavbar />
        {/* Mobile toolbar with menu trigger */}
        <div className="md:hidden flex items-center border-b border-gray-200 bg-white px-4 py-2">
          <Sheet open={mobileNavOpen} onOpenChange={setMobileNavOpen}>
            <SheetTrigger
              className="flex items-center gap-2 text-sm font-medium text-gray-700"
              aria-label="פתח תפריט"
            >
              <Menu size={20} />
              תפריט
            </SheetTrigger>
            <SheetContent side="right" className="w-64 p-0">
              <SheetTitle className="px-5 pt-5 pb-2">תפריט ניהול</SheetTitle>
              <nav className="flex flex-col overflow-y-auto pb-3">
                {renderNavLinks(() => setMobileNavOpen(false))}
              </nav>
            </SheetContent>
          </Sheet>
        </div>
        <main className="flex-1 overflow-y-auto p-6 bg-background md:pr-60">
          {children}
        </main>
      </div>
      {/* Sidebar (desktop only) */}
      <aside className="hidden md:block w-56 shrink-0 bg-white border-l border-gray-200 fixed top-14 bottom-0 shadow-sm">
        <nav className="flex flex-col overflow-y-auto py-3">
          {renderNavLinks()}
        </nav>
      </aside>
    </div>
  );
}
