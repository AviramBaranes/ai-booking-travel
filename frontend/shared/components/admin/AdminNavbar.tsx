"use client";
/*
We use px in styling to ensure consistency across (payload) and (app) sections, 
as rem units can be affected by global styles in (payload) that we don't control. 
 */

import "./admin-tailwind.css";
import clsx from "clsx";
import { FileText, LayoutDashboard, LogOut } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { signOut } from "next-auth/react";
import Image from "next/image";

export default function AdminNavbar({ hideLinks = false }: { hideLinks?: boolean }) {
  const pathname = usePathname();

  const isHeaderActive = (href: string) =>
    pathname === href || pathname.startsWith(href);

  return (
    <header className="h-[56px] shrink-0 bg-white border-b border-gray-200 flex items-center justify-between px-[24px] shadow-sm text-[20px] m-0">
      <div
        className="flex items-center gap-[16px] m-0 p-0"
        style={{ boxSizing: "border-box" }}
      >
        <div
          className="pl-[38px] border-l-2 py-[16px] border-b border-gray-200"
          style={{ boxSizing: "border-box" }}
        >
          <Image
            src="/logo.png"
            alt="AIBookingTravel"
            width={160}
            height={40}
          />
        </div>
        {!hideLinks && [
          {
            href: "/admin",
            label: "ניהול מערכת",
            icon: LayoutDashboard,
          },
          {
            href: "/cms",
            label: "ניהול תוכן",
            icon: FileText,
          },
        ].map((item) => {
          const active = isHeaderActive(item.href);
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              style={{ boxSizing: "border-box" }}
              className={clsx(
                "flex items-center gap-[8px] text-[16px] font-semibold transition-colors px-[8px] py-[4px] rounded-md no-underline m-0",
                active
                  ? "bg-blue-50 text-blue-700 border-l-[3px] border-blue-600"
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
      </div>
      <button
        onClick={() => signOut({ callbackUrl: "/he/" })}
        style={{ boxSizing: "border-box" }}
        className="appearance-none bg-transparent border-none p-0 m-0 flex items-center gap-[8px] text-[14px] text-gray-600 hover:text-red-600 transition-colors cursor-pointer"
      >
        <LogOut size={16} />
        התנתק
      </button>
    </header>
  );
}
