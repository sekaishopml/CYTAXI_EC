import React from "react";
import Link from "next/link";

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur">
      <div className="container mx-auto flex items-center justify-between h-16 px-4 max-w-3xl">
        <Link href="/" className="text-xl font-bold">CYTAXI</Link>
        <nav className="flex items-center gap-4" aria-label="Main navigation">
          <Link href="/trip" className="text-sm hover:underline">Trip</Link>
          <Link href="/profile" className="text-sm hover:underline">Profile</Link>
          <Link href="/notifications" className="text-sm hover:underline">Notifications</Link>
        </nav>
      </div>
    </header>
  );
}
