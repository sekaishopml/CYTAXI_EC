import { cn } from "../utils/cn";
import { colors, radius, spacing, typography, fontFamily, transition, shadows, duration } from "@cytaxi/design-tokens";
import { useRef, useState } from "react";

type ButtonVariant = "primary" | "secondary" | "ghost" | "danger" | "outline";
type ButtonSize = "sm" | "md" | "lg";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  loading?: boolean;
  icon?: React.ReactNode;
  fullWidth?: boolean;
}

const variantStyles: Record<ButtonVariant, React.CSSProperties> = {
  primary: {
    background: colors.brand.green,
    color: colors.text.inverse,
    border: "none",
    boxShadow: shadows.buttonGreen,
  },
  secondary: {
    background: colors.surface.card,
    color: colors.text.primary,
    border: `1px solid ${colors.border.strong}`,
  },
  ghost: {
    background: "transparent",
    color: colors.text.primary,
    border: "none",
  },
  danger: {
    background: colors.status.errorBg,
    color: colors.status.error,
    border: "none",
  },
  outline: {
    background: "transparent",
    color: colors.brand.green,
    border: `1.5px solid ${colors.brand.green}`,
  },
};

const sizeStyles: Record<ButtonSize, React.CSSProperties> = {
  sm: { height: 34, padding: `0 ${spacing[4]}px`, fontSize: typography.bodySm.fontSize },
  md: { height: 44, padding: `0 ${spacing[6]}px`, fontSize: typography.body.fontSize },
  lg: { height: 52, padding: `0 ${spacing[8]}px`, fontSize: typography.body.fontSize },
};

export function Button({
  variant = "primary",
  size = "md",
  loading = false,
  icon,
  fullWidth,
  children,
  disabled,
  style,
  className,
  onClick,
  ...rest
}: ButtonProps) {
  const [scale, setScale] = useState(1);
  const timerRef = useRef<ReturnType<typeof setTimeout>>();

  const baseStyle: React.CSSProperties = {
    ...variantStyles[variant],
    ...sizeStyles[size],
    borderRadius: radius.md,
    fontWeight: typography.label.fontWeight,
    fontFamily: fontFamily.sans,
    cursor: disabled || loading ? "not-allowed" : "pointer",
    opacity: disabled ? 0.5 : loading ? 0.8 : 1,
    display: "inline-flex",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing[2],
    transition: `${transition.hover}, transform ${duration.hover}ms cubic-bezier(0.34, 1.56, 0.64, 1)`,
    transform: `scale(${scale})`,
    outline: "none",
    whiteSpace: "nowrap",
    width: fullWidth ? "100%" : undefined,
    ...style,
  };

  const handlePointerDown = () => {
    if (disabled || loading) return;
    setScale(0.96);
    clearTimeout(timerRef.current);
  };

  const handlePointerUp = () => {
    setScale(1);
    timerRef.current = setTimeout(() => setScale(1), 200);
  };

  return (
    <button
      className={cn(className)}
      style={baseStyle}
      disabled={disabled || loading}
      onPointerDown={handlePointerDown}
      onPointerUp={handlePointerUp}
      onPointerLeave={handlePointerUp}
      onClick={loading ? undefined : onClick}
      {...rest}
    >
      {loading ? <Spinner /> : icon}
      {children}
    </button>
  );
}

function Spinner() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" style={{ animation: "spin 0.6s linear infinite" }}>
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" opacity="0.2" />
      <path d="M12 2a10 10 0 019.95 9" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
    </svg>
  );
}

