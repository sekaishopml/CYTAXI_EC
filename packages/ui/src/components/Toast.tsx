import { useEffect, useState } from "react";
import { colors, radius, spacing, shadows, zIndex, duration, easing, typography, fontFamily, fontWeight } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

type ToastVariant = "success" | "error" | "info" | "warning";

interface ToastData {
  id: string;
  message: string;
  variant?: ToastVariant;
  duration?: number;
}

interface ToastContainerProps {
  toasts: ToastData[];
  onDismiss: (id: string) => void;
  className?: string;
}

const variantColors: Record<ToastVariant, { bg: string; color: string; icon: string }> = {
  success: { bg: colors.brand.green, color: colors.text.inverse, icon: "✓" },
  error: { bg: colors.status.error, color: colors.text.inverse, icon: "✕" },
  info: { bg: colors.status.info, color: colors.text.inverse, icon: "ℹ" },
  warning: { bg: colors.status.warning, color: colors.text.primary, icon: "⚠" },
};

export function ToastItem({ toast, onDismiss }: { toast: ToastData; onDismiss: (id: string) => void }) {
  const [visible, setVisible] = useState(true);
  const v = variantColors[toast.variant || "info"];

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false);
      setTimeout(() => onDismiss(toast.id), 300);
    }, toast.duration || 4000);
    return () => clearTimeout(timer);
  }, [toast, onDismiss]);

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: spacing[4],
        padding: `${spacing[4]}px ${spacing[6]}px`,
        borderRadius: radius.md,
        background: v.bg,
        color: v.color,
        boxShadow: shadows.toast,
        fontSize: typography.body.fontSize,
        fontWeight: fontWeight.medium,
        fontFamily: fontFamily.sans,
        animation: visible
          ? `slideDown ${duration.toast}ms ${easing.spring}`
          : `fadeOut ${duration.fade}ms ${easing.easeIn}`,
        pointerEvents: "auto",
        maxWidth: 360,
      }}
    >
      <span style={{ fontWeight: 700, fontSize: 14 }}>{v.icon}</span>
      <span style={{ flex: 1 }}>{toast.message}</span>
      <button
        onClick={() => { setVisible(false); setTimeout(() => onDismiss(toast.id), 300); }}
        style={{
          background: "none",
          border: "none",
          color: "inherit",
          opacity: 0.7,
          cursor: "pointer",
          fontSize: 16,
          padding: 0,
          lineHeight: 1,
        }}
      >
        ✕
      </button>
    </div>
  );
}

export function ToastContainer({ toasts, onDismiss, className }: ToastContainerProps) {
  return (
    <div
      className={cn(className)}
      style={{
        position: "fixed",
        top: spacing[8],
        right: spacing[8],
        zIndex: zIndex.toast,
        display: "flex",
        flexDirection: "column",
        gap: spacing[3],
        pointerEvents: "none",
      }}
    >
      {toasts.map(t => (
        <ToastItem key={t.id} toast={t} onDismiss={onDismiss} />
      ))}
    </div>
  );
}
