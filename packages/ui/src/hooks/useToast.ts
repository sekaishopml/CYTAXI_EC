import { useState, useCallback } from "react";

type ToastVariant = "success" | "error" | "info" | "warning";

interface ToastData {
  id: string;
  message: string;
  variant: ToastVariant;
  duration?: number;
}

export function useToast() {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const addToast = useCallback((message: string, variant: ToastVariant = "info", duration?: number) => {
    const id = `toast-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;
    setToasts(prev => [...prev, { id, message, variant, duration }]);
    return id;
  }, []);

  const dismissToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(t => t.id !== id));
  }, []);

  return { toasts, addToast, dismissToast };
}
