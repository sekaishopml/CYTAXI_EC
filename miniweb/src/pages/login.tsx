import React, { useState } from "react";
import { Layout } from "@/components/layout/layout";
import { useAuth } from "@/contexts/auth_context";

export default function LoginPage() {
  const { login } = useAuth();
  const [phone, setPhone] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login(phone);
  };

  return (
    <Layout>
      <div className="max-w-sm mx-auto py-12">
        <h1 className="text-2xl font-bold text-center mb-6">Welcome to CYTAXI</h1>
        <form onSubmit={handleSubmit} className="card space-y-4" role="form">
          <input
            className="input"
            type="tel"
            placeholder="Phone number"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            required
            aria-label="Phone number"
          />
          <button type="submit" className="btn-primary w-full">
            Continue with Phone
          </button>
          <p className="text-xs text-muted text-center">
            We&apos;ll send a verification code via WhatsApp
          </p>
        </form>
      </div>
    </Layout>
  );
}
