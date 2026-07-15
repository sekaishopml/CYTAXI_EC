import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";

const links = [
  { href: "/dashboard", label: "Dashboard", icon: "■" },
  { href: "/trips", label: "Trip Queue", icon: "▶" },
  { href: "/trip/current", label: "Current Trip", icon: "●" },
  { href: "/trip/history", label: "History", icon: "◷" },
  { href: "/vehicle", label: "Vehicle", icon: "◆" },
  { href: "/documents", label: "Documents", icon: "◈" },
  { href: "/notifications", label: "Notifications", icon: "◉" },
  { href: "/settings", label: "Settings", icon: "⚙" },
];

export function Sidebar() {
  const router = useRouter();
  return (
    <aside className="hidden lg:flex flex-col fixed left-0 top-0 bottom-0 w-[var(--sidebar-width)] bg-surface border-r border-border p-4 z-40">
      <Link href="/dashboard" className="text-xl font-bold mb-8 px-3">CYTAXI Driver</Link>
      <nav className="flex-1 space-y-1" role="navigation" aria-label="Main">
        {links.map(link => (
          <Link key={link.href} href={link.href}
            className={router.pathname === link.href ? "sidebar-link-active" : "sidebar-link"}>
            <span className="w-5 text-center">{link.icon}</span>
            {link.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
