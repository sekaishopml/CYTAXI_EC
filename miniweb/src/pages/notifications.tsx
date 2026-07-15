import React from "react";
import { Layout } from "@/components/layout/layout";

export default function NotificationsPage() {
  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Notifications</h1>
        <p className="text-muted text-center py-8">No notifications yet</p>
      </section>
    </Layout>
  );
}
