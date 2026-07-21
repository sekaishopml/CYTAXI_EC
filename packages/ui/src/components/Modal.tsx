import { useEffect, useRef } from "react";
import { colors, radius, spacing, shadows, zIndex, duration, easing } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface ModalProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  closeOnOverlay?: boolean;
  style?: React.CSSProperties;
  className?: string;
}

export function Modal({
  open,
  onClose,
  title,
  children,
  closeOnOverlay = true,
  style,
  className,
}: ModalProps) {
  const overlayRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", handleKey);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", handleKey);
      document.body.style.overflow = "";
    };
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      ref={overlayRef}
      onClick={(e) => { if (closeOnOverlay && e.target === overlayRef.current) onClose(); }}
      style={{
        position: "fixed",
        inset: 0,
        zIndex: zIndex.modal,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: spacing[6],
        background: colors.surface.overlay,
        animation: `fadeIn ${duration.modal}ms ${easing.easeOut}`,
      }}
    >
      <div
        className={cn(className)}
        style={{
          background: colors.surface.white,
          borderRadius: radius.xl,
          boxShadow: shadows.modal,
          padding: spacing[8],
          maxWidth: 480,
          width: "100%",
          maxHeight: "90vh",
          overflowY: "auto",
          animation: `scaleIn ${duration.modal}ms ${easing.easeOut}`,
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
            <button
              onClick={onClose}
              style={{
                background: "none",
                border: "none",
                fontSize: 20,
                cursor: "pointer",
                padding: 4,
                lineHeight: 1,
                color: colors.text.muted,
              }}
            >
              ✕
            </button>
          </div>
        )}
        {children}
      </div>
    </div>
  );
}
