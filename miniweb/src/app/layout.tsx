import { Providers } from "@/app/providers";
import "@/styles/globals.css";
import { Layout } from "@/components/layout/layout";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <title>CYTAXI - Conversational Mobility</title>
        <meta name="description" content="Conversation-First Mobility Platform" />
        <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
        <meta name="theme-color" content="#000000" />
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body>
        <Providers>
          <Layout>{children}</Layout>
        </Providers>
      </body>
    </html>
  );
}
