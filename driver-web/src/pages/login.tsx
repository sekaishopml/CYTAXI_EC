import React, { useState } from "react";
import { useAuth } from "@/contexts/auth";

export default function LoginPage() {
  const { login } = useAuth();
  const [phone, setPhone] = useState("");
  const handler = (e: React.FormEvent) => { e.preventDefault(); login(phone); };
  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <div className="w-full max-w-sm space-y-6">
        <div className="text-center">
          <h1 className="text-2xl font-bold">CYTAXI</h1>
          <p className="text-muted-foreground mt-1">Driver Portal</p>
        </div>
        <form onSubmit={handler} className="card space-y-4">
          <input className="input" type="tel" placeholder="Phone number" value={phone} onChange={e => setPhone(e.target.value)} required aria-label="Phone number" />
          <button className="btn-primary w-full">Sign In</button>
          <p className="text-xs text-muted-foreground text-center">We&apos;ll verify your identity via WhatsApp</p>
        </form>
      </div>
    </div>
  );
}
