import { colors, radius, spacing, shadows, transition } from "@cytaxi/design-tokens";
import { cn } from "../utils/cn";

interface CardProps {
  children: React.ReactNode;
  variant?: "default" | "glass" | "flat";
  padding?: keyof typeof spacing | number;
  hover?: boolean;
  onClick?: () => void;
  style?: React.CSSProperties;
  className?: string;
}

const variantStyles: Record<string, React.CSSProperties> = {
  default: {
    background: colors.surface.white,
    boxShadow: shadows.card,
    border: `1px solid ${colors.border.light}`,
  },
  glass: {
    background: colors.surface.glass,
    backdropFilter: "blur(24px) saturate(180%)",
    WebkitBackdropFilter: "blur(24px) saturate(180%)",
    border: `1px solid ${colors.border.light}`,
  },
  flat: {
    background: colors.surface.bg,
    border: `1px solid ${colors.border.medium}`,
  },
};

export function Card({
  children,
  variant = "default",
  padding = 6,
  hover,
  onClick,
  style,
  className,
}: CardProps) {
  const padValue = typeof padding === "number" ? padding : spacing[padding];
  const isClickable = !!onClick || hover;

  return (
    <div
      className={cn(className)}
      onClick={onClick}
      style={{
        borderRadius: radius.lg,
        padding: padValue,
        transition: `${transition.hover}, box-shadow ${hover ? "0.2s" : "0s"} ease`,
        cursor: isClickable ? "pointer" : undefined,
        ...variantStyles[variant],
        ...style,
      }}
    >
      {children}
    </div>
  );
}
