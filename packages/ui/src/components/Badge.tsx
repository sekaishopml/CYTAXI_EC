import { colors, radius, typography } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

type BadgeVariant = "default" | "success" | "warning" | "error" | "info" | "brand";

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  dot?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

const variantStyles: Record<BadgeVariant, { bg: string; color: string }> = {
  default: { bg: colors.surface.bg, color: colors.text.muted },
  success: { bg: colors.brand.greenBg, color: colors.brand.green },
  warning: { bg: "rgba(255,193,7,0.12)", color: "#b8860b" },
  error: { bg: colors.status.errorBg, color: colors.status.error },
  info: { bg: "rgba(68,138,255,0.1)", color: colors.status.info },
  brand: { bg: colors.brand.green, color: colors.text.inverse },
};

export function Badge({ children, variant = "default", dot, className, style }: BadgeProps) {
  const v = variantStyles[variant];
  return (
    <span
      className={cn(className)}
      style={{
        display: "inline-flex",
        alignItems: "center",
        gap: 5,
        padding: "3px 10px",
        borderRadius: radius.full,
        fontSize: typography.caption.fontSize,
        fontWeight: typography.label.fontWeight,
        background: v.bg,
        color: v.color,
        letterSpacing: "0.02em",
        whiteSpace: "nowrap",
        ...style,
      }}
    >
      {dot && (
        <span
          style={{
            width: 6, height: 6,
            borderRadius: "50%",
            background: v.color,
            flexShrink: 0,
          }}
        />
      )}
      {children}
    </span>
  );
}
