import React from "react";
import { Layout } from "@/components/layout/layout";
import { ProfileCard } from "@/components/ui/profile_card";
import { useAuth } from "@/contexts/auth_context";

export default function ProfilePage() {
  const { user } = useAuth();

  return (
    <Layout>
      <section className="space-y-6">
        <h1 className="text-2xl font-bold">Profile</h1>
        {user ? (
          <ProfileCard name={user.name} phone={user.phone} email={user.email} tripsCount={0} />
        ) : (
          <p className="text-muted">Please log in to view your profile.</p>
        )}
        <nav className="space-y-2">
          <a href="/settings" className="card block hover:shadow-md">Settings</a>
          <a href="/help" className="card block hover:shadow-md">Help & Support</a>
        </nav>
      </section>
    </Layout>
  );
}
