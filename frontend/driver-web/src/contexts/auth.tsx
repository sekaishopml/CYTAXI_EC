import React, { createContext, useContext, useState } from "react";

interface Driver {
  id: string; name: string; phone: string; status: string; rating: number;
}

interface AuthCtx {
  driver: Driver | null; isAuthenticated: boolean;
  login: (phone: string) => Promise<void>; logout: () => void;
}
const AuthContext = createContext<AuthCtx | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [driver, setDriver] = useState<Driver | null>(null);
  const login = async (phone: string) => { setDriver({ id: "drv_1000", name: "Driver", phone, status: "online", rating: 4.8 }); };
  const logout = () => setDriver(null);
  return <AuthContext.Provider value={{ driver, isAuthenticated: !!driver, login, logout }}>{children}</AuthContext.Provider>;
}
export function useAuth() { const c = useContext(AuthContext); if (!c) throw new Error("useAuth"); return c; }
