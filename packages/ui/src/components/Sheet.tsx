import { useEffect, useRef } from "react";
import { colors, radius, spacing, shadows, zIndex, duration, easing, typography } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface SheetProps {
  open: boolean;
  onClose?: () => void;
  title?: string;
  snapPoints?: number[];
  children: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
}

export function Sheet({
  open,
  onClose,
  title,
  children,
  style,
  className,
}: SheetProps) {
  const sheetRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose?.();
    };
    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        zIndex: zIndex.sheet,
        animation: `fadeIn ${duration.modal}ms ${easing.easeOut}`,
      }}
    >
      <div
        onClick={onClose}
        style={{
          position: "absolute", inset: 0,
          background: colors.surface.overlay,
        }}
      />
      <div
        ref={sheetRef}
        className={cn(className)}
        style={{
          position: "absolute",
          bottom: 0,
          left: 0,
          right: 0,
          background: colors.surface.white,
          borderRadius: `${radius.xl}px ${radius.xl}px 0 0`,
          boxShadow: shadows.sheet,
          padding: spacing[8],
          paddingBottom: `calc(${spacing[8]}px + env(safe-area-inset-bottom, 0px))`,
          maxHeight: "85vh",
          overflowY: "auto",
          animation: `slideUp ${duration.sheet}ms ${easing.easeOut}`,
          ...style,
        }}
      >
        {title && (
          <div style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: spacing[6],
          }}>
            <h2 style={{ margin: 0, fontSize: 18, fontWeight: 600 }}>{title}</h2>
            {onClose && (
              <button
                onClick={onClose}
                style={{
                  background: colors.surface.bg,
                  border: "none",
                  width: 32, height: 32,
                  borderRadius: "50%",
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: 16,
                  color: colors.text.muted,
                }}
              >
                ✕
              </button>
            )}
          </div>
        )}
        {children}
      </div>
    </div>
  );
}
