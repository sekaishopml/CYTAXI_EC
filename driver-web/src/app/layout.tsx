import React from "react";
import { Providers } from "@/app/providers";
import "@/styles/globals.css";
import { Sidebar } from "@/components/layout/sidebar";
import { Header } from "@/components/layout/header";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <title>CYTAXI Driver Portal</title>
        <meta name="description" content="Driver Portal" />
        <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
        <meta name="theme-color" content="#000000" />
      </head>
      <body>
        <Providers>
          <div className="flex min-h-screen">
            <Sidebar />
            <div className="flex-1 flex flex-col lg:ml-[var(--sidebar-width)]">
              <Header />
              <main className="flex-1 p-4 lg:p-6 max-w-4xl">{children}</main>
            </div>
          </div>
        </Providers>
      </body>
    </html>
  );
}
